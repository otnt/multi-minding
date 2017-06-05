package taskstore

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"reflect"

	"github.com/golang/protobuf/proto"
	commonpb "github.com/otnt/multi-minding/common/proto"
)

func TestInMemoryTaskStoreCommit_Update(t *testing.T) {
	taskSetUps := []struct {
		existedTasks []*commonpb.Task
		updateTasks  []*commonpb.Task
		remainTasks  taskIDToTaskMap
	}{
		{
			existedTasks: []*commonpb.Task{},
			updateTasks:  []*commonpb.Task{},
			remainTasks:  map[string]*commonpb.Task{},
		},
		{
			existedTasks: []*commonpb.Task{},
			updateTasks:  []*commonpb.Task{&commonpb.Task{Id: proto.String("1")}},
			remainTasks: map[string]*commonpb.Task{
				"1": &commonpb.Task{Id: proto.String("1")},
			},
		},
		{
			existedTasks: []*commonpb.Task{},
			updateTasks: []*commonpb.Task{
				&commonpb.Task{Id: proto.String("1")},
				&commonpb.Task{Id: proto.String("2")},
			},
			remainTasks: map[string]*commonpb.Task{
				"1": &commonpb.Task{Id: proto.String("1")},
				"2": &commonpb.Task{Id: proto.String("2")},
			},
		},
		{
			existedTasks: []*commonpb.Task{&commonpb.Task{Id: proto.String("1")}},
			updateTasks:  []*commonpb.Task{&commonpb.Task{Id: proto.String("2")}},
			remainTasks: map[string]*commonpb.Task{
				"1": &commonpb.Task{Id: proto.String("1")},
				"2": &commonpb.Task{Id: proto.String("2")},
			},
		},
	}

	user := commonpb.User{Id: proto.Int64(1)}

	for _, taskSetUp := range taskSetUps {
		taskStore := MakeInMemoryTaskStore().WithTasks(&user, taskSetUp.existedTasks)
		commitGroup := CommitGroup{upsertTasks: taskSetUp.updateTasks}

		if err := taskStore.Commit(&user, &commitGroup); err != nil {
			t.Errorf("Commit(%+v, %+v) got error %v, wanted nil.", user, commitGroup, err)
			continue
		}

		actualRemainTasks := taskStore.userIDToTasksMap[user.GetId()]
		if !reflect.DeepEqual(actualRemainTasks, taskSetUp.remainTasks) {
			t.Errorf("Remain tasks got %#v, wanted %#v", actualRemainTasks, taskSetUp.remainTasks)
		}
	}
}

func TestInMemoryTaskStoreQuery(t *testing.T) {
	taskSetUps := [][]*commonpb.Task{
		[]*commonpb.Task{
			{},
		},
		[]*commonpb.Task{
			&commonpb.Task{Id: proto.String("1")},
		},
		[]*commonpb.Task{
			&commonpb.Task{Id: proto.String("1")},
			&commonpb.Task{Id: proto.String("2")},
		},
	}

	user := commonpb.User{Id: proto.Int64(1)}

	for _, taskSetUp := range taskSetUps {
		taskStore := MakeInMemoryTaskStore().WithTasks(&user, taskSetUp)

		queryResult, err := taskStore.Query(&user)
		if err != nil {
			t.Errorf("Query(%+v) got error %v, wanted nil.", user, err)
			continue
		}

		if !reflect.DeepEqual(queryResult, taskSetUp) {
			t.Errorf("Query tasks got %#v, wanted %#v", queryResult, taskSetUp)
		}
	}
}

func TestInMemoryTaskStoreQuery_UserNotExists(t *testing.T) {
	user := commonpb.User{Id: proto.Int64(1)}
	taskStore := MakeInMemoryTaskStore()

	queryResult, err := taskStore.Query(&user)
	if err != nil {
		t.Errorf("Query(%+v) got error %v, wanted nil.", user, err)
		return
	}

	wantedResult := []*commonpb.Task{}
	if !reflect.DeepEqual(queryResult, wantedResult) {
		t.Errorf("Query tasks got %#v, wanted %#v", queryResult, wantedResult)
	}
}

func TestInMemoryTaskStoreCommit_InvalidUserArguments(t *testing.T) {
	taskStore := MakeInMemoryTaskStore()
	user := &commonpb.User{}
	commitGroup := &CommitGroup{}

	if err := taskStore.Commit(user, commitGroup); err == nil {
		t.Errorf("Commit(%+v, %+v) got %v, wanted %v", user, commitGroup, err, grpc.Errorf(codes.InvalidArgument, ""))
	}
}

func TestInMemoryTaskStoreCommit_InvalidTasksArguments(t *testing.T) {
	user := &commonpb.User{Id: proto.Int64(1)}
	commitGroups := []CommitGroup{
		{
			upsertTasks: []*commonpb.Task{},
			deleteTasks: []*commonpb.Task{},
		},
		{
			upsertTasks: []*commonpb.Task{{}},
			deleteTasks: []*commonpb.Task{},
		},
		{
			upsertTasks: []*commonpb.Task{},
			deleteTasks: []*commonpb.Task{{}},
		},
		{
			upsertTasks: []*commonpb.Task{{Id: proto.String("1")}, {Id: proto.String("1")}},
			deleteTasks: []*commonpb.Task{},
		},
		{
			upsertTasks: []*commonpb.Task{},
			deleteTasks: []*commonpb.Task{{Id: proto.String("1")}, {Id: proto.String("1")}},
		},
		{
			upsertTasks: []*commonpb.Task{{Id: proto.String("1")}},
			deleteTasks: []*commonpb.Task{{Id: proto.String("1")}},
		},
	}

	for i, commitGroup := range commitGroups {
		taskStore := MakeInMemoryTaskStore()
		if err := taskStore.Commit(user, &commitGroup); i == 0 {
			if err != nil {
				t.Fatalf("Commit(%+v, %+v) got %v, wanted nil", user, commitGroup, err)
			}
		} else if err == nil {
			t.Errorf("Commit(%+v, %+v) got %v, wanted %v", user, commitGroup, err, grpc.Errorf(codes.InvalidArgument, ""))
		}
	}
}

func TestInMemoryTaskStoreCommit_Delete(t *testing.T) {
	taskSetUps := []struct {
		existedTasks []*commonpb.Task
		deleteTasks  []*commonpb.Task
		remainTasks  taskIDToTaskMap
	}{
		// Nothing exists, then delete nothing.
		{
			existedTasks: []*commonpb.Task{},
			deleteTasks:  []*commonpb.Task{},
			remainTasks:  taskIDToTaskMap{},
		},

		// One exists, then delete one.
		{
			existedTasks: []*commonpb.Task{{Id: proto.String("1")}},
			deleteTasks:  []*commonpb.Task{{Id: proto.String("1")}},
			remainTasks:  taskIDToTaskMap{},
		},

		// Two exist, then delete one.
		{
			existedTasks: []*commonpb.Task{{Id: proto.String("1")}, {Id: proto.String("2")}},
			deleteTasks:  []*commonpb.Task{{Id: proto.String("1")}},
			remainTasks:  taskIDToTaskMap{"2": &commonpb.Task{Id: proto.String("2")}},
		},

		// Two exist, then delete two.
		{
			existedTasks: []*commonpb.Task{{Id: proto.String("1")}, {Id: proto.String("2")}},
			deleteTasks:  []*commonpb.Task{{Id: proto.String("1")}, {Id: proto.String("2")}},
			remainTasks:  taskIDToTaskMap{},
		},

		// Three exist, then delete one.
		{
			existedTasks: []*commonpb.Task{{Id: proto.String("1")}, {Id: proto.String("2")}, {Id: proto.String("3")}},
			deleteTasks:  []*commonpb.Task{{Id: proto.String("1")}},
			remainTasks: taskIDToTaskMap{
				"2": &commonpb.Task{Id: proto.String("2")},
				"3": &commonpb.Task{Id: proto.String("3")},
			},
		},
	}

	user := commonpb.User{Id: proto.Int64(1)}

	for _, taskSetUp := range taskSetUps {
		taskStore := MakeInMemoryTaskStore().WithTasks(&user, taskSetUp.existedTasks)

		commitGroup := CommitGroup{deleteTasks: taskSetUp.deleteTasks}

		if err := taskStore.Commit(&user, &commitGroup); err != nil {
			t.Errorf("Commit(%+v, %+v) got error %v, wanted nil.", user, commitGroup, err)
			continue
		}

		actualRemainTasks := taskStore.userIDToTasksMap[user.GetId()]
		if !reflect.DeepEqual(actualRemainTasks, taskSetUp.remainTasks) {
			t.Errorf("Remain tasks got %#v, wanted %#v", actualRemainTasks, taskSetUp.remainTasks)
		}
	}
}

func TestInMemoryTaskStoreCommit_Delete_TaskNotExists(t *testing.T) {
	user := commonpb.User{Id: proto.Int64(1)}
	taskStore := MakeInMemoryTaskStore().WithTasks(&user, []*commonpb.Task{})
	commitGroup := CommitGroup{deleteTasks: []*commonpb.Task{{Id: proto.String("1")}}}
	if err := taskStore.Commit(&user, &commitGroup); err == nil {
		t.Errorf("Commit(%+v, %+v) got nil, wanted %+v.", user, commitGroup, grpc.Errorf(codes.InvalidArgument, ""))
	}
}

func TestInMemoryTaskStoreCommit_Delete_UserNotExists(t *testing.T) {
	taskStore := MakeInMemoryTaskStore()
	user := commonpb.User{Id: proto.Int64(1)}
	commitGroup := CommitGroup{deleteTasks: []*commonpb.Task{{Id: proto.String("1")}}}
	if err := taskStore.Commit(&user, &commitGroup); err == nil {
		t.Errorf("Commit(%+v, %+v) got nil, wanted %+v.", user, commitGroup, grpc.Errorf(codes.InvalidArgument, ""))
	}
}

func TestInMemoryTaskStore_MultithreadUpsert(t *testing.T) {
	const threadNum int = 1000

	taskStore := MakeInMemoryTaskStore()

	user := commonpb.User{Id: proto.Int64(1)}

	resultTaskMap := taskIDToTaskMap{}
	for i := 0; i < threadNum; i++ {
		id := strconv.FormatInt(int64(i), 10)
		resultTaskMap[id] = &commonpb.Task{Id: proto.String(id)}
	}

	// Spawn a bunch of goroutine to upsert value to task store.
	var wg sync.WaitGroup
	for i := 0; i < threadNum; i++ {
		iLocal := i
		wg.Add(1)
		go func() {
			defer wg.Done()

			id := strconv.FormatInt(int64(iLocal), 10)
			commitGroup := CommitGroup{
				upsertTasks: []*commonpb.Task{
					{Id: proto.String(id)},
				},
			}
			if err := taskStore.Commit(&user, &commitGroup); err != nil {
				panic(err)
			}
		}()
	}
	wg.Wait()

	// Assert result.
	actualTaskMap := taskStore.userIDToTasksMap[user.GetId()]
	if !reflect.DeepEqual(actualTaskMap, resultTaskMap) {
		t.Errorf("Got %+v.", actualTaskMap)
		t.Errorf("Wanted %+v.", resultTaskMap)
	}
}

func TestInMemoryTaskStore_MultithreadDelete(t *testing.T) {
	const threadNum int = 1000

	user := commonpb.User{Id: proto.Int64(1)}

	tasks := make([]*commonpb.Task, threadNum)
	for i := 0; i < threadNum; i++ {
		id := strconv.FormatInt(int64(i), 10)
		tasks[i] = &commonpb.Task{Id: proto.String(id)}
	}

	taskStore := MakeInMemoryTaskStore().WithTasks(&user, tasks)

	// Spawn a bunch of goroutine to delete value from task store.
	var wg sync.WaitGroup
	for i := 0; i < threadNum; i++ {
		iLocal := i
		wg.Add(1)
		go func() {
			defer wg.Done()

			id := strconv.FormatInt(int64(iLocal), 10)
			commitGroup := CommitGroup{
				deleteTasks: []*commonpb.Task{
					{Id: proto.String(id)},
				},
			}
			if err := taskStore.Commit(&user, &commitGroup); err != nil {
				panic(err)
			}
		}()
	}
	wg.Wait()

	// Assert result.
	actualTasks := taskStore.userIDToTasksMap[user.GetId()]
	wantedTasks := make(taskIDToTaskMap)
	if !reflect.DeepEqual(actualTasks, wantedTasks) {
		t.Errorf("Got %+v.", actualTasks)
		t.Errorf("Wanted %+v.", wantedTasks)
	}
}

func TestInMemoryTaskStore_MultithreadUpsertDeleteQuery(t *testing.T) {
	const threadNum int = 6000
	const deleteTaskIndexFrom int = 0
	const deleteTaskIndexTo int = 2000
	const upsertTaskIndexFrom int = 2000
	const upsertTaskIndexTo int = 4000

	user := commonpb.User{Id: proto.Int64(1)}

	existingTasks := make([]*commonpb.Task, deleteTaskIndexTo-deleteTaskIndexFrom)
	for i := deleteTaskIndexFrom; i < deleteTaskIndexTo; i++ {
		id := strconv.FormatInt(int64(i), 10)
		existingTasks[i-deleteTaskIndexFrom] = &commonpb.Task{Id: proto.String(id)}
	}

	remainingTasks := make(taskIDToTaskMap)
	for i := upsertTaskIndexFrom; i < upsertTaskIndexTo; i++ {
		id := strconv.FormatInt(int64(i), 10)
		remainingTasks[id] = &commonpb.Task{Id: proto.String(id)}
	}

	taskStore := MakeInMemoryTaskStore().WithTasks(&user, existingTasks)

	// Create shuffled indices so as to mix up delete/update/query.
	indices := make([]int, threadNum)
	for i := 0; i < threadNum; i++ {
		indices[i] = i
		j := rand.Intn(i + 1)
		indices[i], indices[j] = indices[j], indices[i]
	}

	// Spawn a bunch of goroutine to delete/update/query value from task store.
	var wg sync.WaitGroup
	for i := 0; i < threadNum; i++ {
		iLocal := indices[i]
		wg.Add(1)
		go func() {
			defer wg.Done()

			if deleteTaskIndexFrom <= iLocal && iLocal < deleteTaskIndexTo {
				// Do delete.
				id := strconv.FormatInt(int64(iLocal), 10)
				commitGroup := CommitGroup{
					deleteTasks: []*commonpb.Task{
						{Id: proto.String(id)},
					},
				}
				if err := taskStore.Commit(&user, &commitGroup); err != nil {
					panic(err)
				}
			} else if upsertTaskIndexFrom <= iLocal && iLocal < upsertTaskIndexTo {
				// Do upsert.
				id := strconv.FormatInt(int64(iLocal), 10)
				commitGroup := CommitGroup{
					upsertTasks: []*commonpb.Task{
						{Id: proto.String(id)},
					},
				}
				if err := taskStore.Commit(&user, &commitGroup); err != nil {
					panic(err)
				}
			} else {
				// Do query.
				if _, err := taskStore.Query(&user); err != nil {
					panic(err)
				}
			}
		}()
	}
	wg.Wait()

	// Assert result.
	actualTasks := taskStore.userIDToTasksMap[user.GetId()]
	if !reflect.DeepEqual(actualTasks, remainingTasks) {
		t.Errorf("Got %+v.", actualTasks)
		t.Errorf("Wanted %+v.", remainingTasks)
	}
}
