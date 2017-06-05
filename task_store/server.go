package taskstore

import (
	commonpb "github.com/otnt/multi-minding/common/proto"

	"log"

	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// CommitGroup contains list of tasks to insert/update, and list of tasks to delete.
// The entire commit is done in an atomic transaction.
type CommitGroup struct {
	// List of tasks to insert or update.
	upsertTasks []*commonpb.Task

	// List of tasks to delete.
	deleteTasks []*commonpb.Task
}

// TaskStore interface defines how client could interact with task store backend.
type TaskStore interface {
	// Commit the "commitGroup" for "user" in an atomic transaction.
	//
	// Within "commitGroup", if task in "updates" does not exist for "user", then the
	// task is inserted, if task in "updates" already exists, then the task is updated (replaced
	// with new value).
	//
	// Preconditions
	// - all tasks have different task ids
	// - any task in "deletes" must already exists
	Commit(user *commonpb.User, commitGroup *CommitGroup) error

	// Query all tasks for the given "user".
	Query(user *commonpb.User) ([]*commonpb.Task, error)
}

type (
	// All tasks for one user. Map of task ID to task.
	taskIDToTaskMap map[string]*commonpb.Task

	// InMemoryTaskStore is an implementation of TaskStore mainly for testing.
	//
	// It only saves state in memory, therefore suffers error situation.
	InMemoryTaskStore struct {
		// Map of user id to all tasks for this user.
		userIDToTasksMap map[int64]taskIDToTaskMap

		// Protect multi-thread map read/write.
		lock sync.RWMutex
	}
)

// Commit update/delete tasks in "commitGroup" for "user" to in memory storage.
func (taskStore *InMemoryTaskStore) Commit(user *commonpb.User, commitGroup *CommitGroup) error {
	if err := assertUserIDExists(user); err != nil {
		return err
	}
	if err := assertTasksIDsExist(commitGroup.upsertTasks); err != nil {
		return err
	}
	if err := assertTasksIDsExist(commitGroup.deleteTasks); err != nil {
		return err
	}
	if err := assertTasksIDsUnique(commitGroup.upsertTasks, commitGroup.deleteTasks); err != nil {
		return err
	}

	taskStore.lock.Lock()
	defer taskStore.lock.Unlock()

	if err := assertDeleteTasksExist(taskStore.userIDToTasksMap, user, commitGroup.deleteTasks); err != nil {
		return err
	}

	taskStore.createUserIfNotPresent(user.GetId())

	taskStore.upsertTasks(user.GetId(), commitGroup.upsertTasks)

	taskStore.deleteTasks(user.GetId(), commitGroup.deleteTasks)

	return nil
}

func (taskStore *InMemoryTaskStore) createUserIfNotPresent(userID int64) {
	if _, ok := taskStore.userIDToTasksMap[userID]; !ok {
		taskStore.userIDToTasksMap[userID] = make(map[string]*commonpb.Task)
	}
}

func (taskStore *InMemoryTaskStore) upsertTasks(userID int64, upsertTasks []*commonpb.Task) {
	taskMap := taskStore.userIDToTasksMap[userID]
	for _, upsertTask := range upsertTasks {
		taskMap[upsertTask.GetId()] = upsertTask
	}
}

func (taskStore *InMemoryTaskStore) deleteTasks(userID int64, deleteTasks []*commonpb.Task) {
	taskMap := taskStore.userIDToTasksMap[userID]
	for _, deleteTask := range deleteTasks {
		delete(taskMap, deleteTask.GetId())
	}
}

// Query all tasks for the given "user" from in memory storage. If user does not exist,
// return empty slice.
func (taskStore *InMemoryTaskStore) Query(user *commonpb.User) ([]commonpb.Task, error) {
	if err := assertUserIDExists(user); err != nil {
		return nil, err
	}

	taskStore.lock.RLock()
	defer taskStore.lock.RUnlock()

	taskMap, ok := taskStore.userIDToTasksMap[user.GetId()]
	if !ok {
		return []commonpb.Task{}, nil
	}

	tasks := make([]commonpb.Task, 0)
	for _, task := range taskMap {
		tasks = append(tasks, *task)
	}
	return tasks, nil
}

// MakeInMemoryTaskStore creates a new InMemoryTaskStore.
func MakeInMemoryTaskStore() *InMemoryTaskStore {
	return &InMemoryTaskStore{userIDToTasksMap: make(map[int64]taskIDToTaskMap)}
}

// Only for testing.

// WithTasks add "tasks" to task store for "user". "user" must be absent from task store.
func (taskStore *InMemoryTaskStore) WithTasks(user *commonpb.User, tasks []*commonpb.Task) *InMemoryTaskStore {
	if err := assertUserIDExists(user); err != nil {
		log.Fatal(err)
	}
	if err := assertTasksIDsExist(tasks); err != nil {
		log.Fatal(err)
	}
	if _, ok := taskStore.userIDToTasksMap[user.GetId()]; ok {
		log.Fatalf("User %+v already exist.", user)
	}

	taskStore.createUserIfNotPresent(user.GetId())
	taskStore.upsertTasks(user.GetId(), tasks)
	return taskStore
}

// Internal helper functions.

func assertUserIDExists(user *commonpb.User) error {
	if user.Id == nil {
		return grpc.Errorf(codes.InvalidArgument, "Missing user id: %+v", user)
	}
	return nil
}

func assertTasksIDsExist(tasks []*commonpb.Task) error {
	for _, task := range tasks {
		if task.Id == nil {
			return grpc.Errorf(codes.InvalidArgument, "Missing task id: %+v", task)
		}
	}
	return nil
}

// Assert all IDs in all provided tasks are unique. This function assumes all IDs existed.
func assertTasksIDsUnique(first []*commonpb.Task, others ...[]*commonpb.Task) error {
	allTasks := [][]*commonpb.Task{first}
	allTasks = append(allTasks, others...)

	ids := make(map[string]bool)

	for _, tasks := range allTasks {
		for _, task := range tasks {
			if _, ok := ids[task.GetId()]; ok {
				return grpc.Errorf(codes.InvalidArgument, "Duplicate IDs: %s", task.GetId())
			}
			ids[task.GetId()] = true
		}
	}

	return nil
}

func assertDeleteTasksExist(userIDToTasksMap map[int64]taskIDToTaskMap, user *commonpb.User, deleteTasks []*commonpb.Task) error {
	if len(deleteTasks) == 0 {
		return nil
	}

	tasksMap, ok := userIDToTasksMap[user.GetId()]
	if !ok {
		return grpc.Errorf(codes.InvalidArgument, "User %+v does not exist.", user)
	}
	for _, deleteTask := range deleteTasks {
		if _, ok = tasksMap[deleteTask.GetId()]; !ok {
			return grpc.Errorf(codes.InvalidArgument, "Task %+v does not exist for user %+v.", deleteTask, user)
		}
	}
	return nil
}
