package taskstore

import (
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"reflect"

	"github.com/golang/protobuf/proto"
	commonpb "github.com/otnt/multi-minding/common/proto"
)

func TestInMemoryTaskStoreCommit_Update(t *testing.T) {
	taskLists := [][]*commonpb.Task{
		// Empty task.
		{},
		// One task.
		{&commonpb.Task{Id: proto.String("TaskId-1")}},
		// Two tasks.
		{
			&commonpb.Task{Id: proto.String("TaskId-1")},
			&commonpb.Task{Id: proto.String("TaskId-2")},
		},
	}

	userLists := [][]commonpb.User{
		// Single user.
		{commonpb.User{Id: proto.Int64(1)}},
		// Two users.
		{
			commonpb.User{Id: proto.Int64(1)},
			commonpb.User{Id: proto.Int64(2)},
		},
	}

	// For each users, tasks combination.
	for _, tasks := range taskLists {
		for _, users := range userLists {
			taskStore := MakeInMemoryTaskStore()
			commitGroup := CommitGroup{upsertTasks: tasks}

			// Upsert tasks for user.
			for _, user := range users {
				if err := taskStore.Commit(&user, &commitGroup); err != nil {
					t.Errorf("Commit(%+v, %+v) got error %v, wanted nil.", user, commitGroup, err)
					continue
				}
			}

			// Query user tasks.
			for _, user := range users {
				queryTasks, err := taskStore.Query(&user)
				if err != nil {
					t.Errorf("Query(%+v) got error %v, wanted nil", user, err)
					continue
				}
				if !reflect.DeepEqual(tasks, queryTasks) {
					t.Errorf("Query(%+v) got tasks %+v, wanted %+v", user, queryTasks, tasks)
				}
			}
		}
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
