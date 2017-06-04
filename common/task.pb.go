// Code generated by protoc-gen-go. DO NOT EDIT.
// source: task.proto

/*
Package common is a generated protocol buffer package.

It is generated from these files:
	task.proto

It has these top-level messages:
	Snapshot
	History
	Priority
	Task
*/
package common

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import google_protobuf "github.com/golang/protobuf/ptypes/timestamp"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// Type of task. This should only be framework recognized type. For
// example, a REMINDER type task could be pushed to user. Examples
// that should not be put here is NEWS, SHOPPING etc.
type Type int32

const (
	// REMINDER task could be queired by reminder and pushed to
	// users. A reminder task must has reminding time.
	Type_REMINDER Type = 0
	// REPEATED task contains information of how to repeat a task
	// based on different task state.
	Type_REPEATED Type = 1
)

var Type_name = map[int32]string{
	0: "REMINDER",
	1: "REPEATED",
}
var Type_value = map[string]int32{
	"REMINDER": 0,
	"REPEATED": 1,
}

func (x Type) Enum() *Type {
	p := new(Type)
	*p = x
	return p
}
func (x Type) String() string {
	return proto.EnumName(Type_name, int32(x))
}
func (x *Type) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(Type_value, data, "Type")
	if err != nil {
		return err
	}
	*x = Type(value)
	return nil
}
func (Type) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

// State of task.
type State int32

const (
	// A task has just been created. Initial state of a task.
	State_CREATED State = 0
	// A task has been dispatched to user and user is working on
	// this task.
	State_WORK_IN_PROGRESS State = 1
	// A task has been dispatched to user but user reject to work
	// on this task.
	State_REJECTED State = 2
	// A task has been dispatched to user, user worked on this task
	// but is now blocked and cannot continue working on this task.
	State_BLOCKED State = 3
	// A task has been canceled by user.
	State_CANCELED State = 4
	// A task has been delayed by user.
	State_DELAYED State = 5
	// A task has been pending too long and is considered as
	// expired.
	State_EXPIRED State = 6
)

var State_name = map[int32]string{
	0: "CREATED",
	1: "WORK_IN_PROGRESS",
	2: "REJECTED",
	3: "BLOCKED",
	4: "CANCELED",
	5: "DELAYED",
	6: "EXPIRED",
}
var State_value = map[string]int32{
	"CREATED":          0,
	"WORK_IN_PROGRESS": 1,
	"REJECTED":         2,
	"BLOCKED":          3,
	"CANCELED":         4,
	"DELAYED":          5,
	"EXPIRED":          6,
}

func (x State) Enum() *State {
	p := new(State)
	*p = x
	return p
}
func (x State) String() string {
	return proto.EnumName(State_name, int32(x))
}
func (x *State) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(State_value, data, "State")
	if err != nil {
		return err
	}
	*x = State(value)
	return nil
}
func (State) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

type Snapshot struct {
	// Same as id in task.
	Id *string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	// The time when this snapshot is took.
	Time *google_protobuf.Timestamp `protobuf:"bytes,2,opt,name=time" json:"time,omitempty"`
	// The state when this snapshot is took.
	State *State `protobuf:"varint,3,opt,name=state,enum=common.State" json:"state,omitempty"`
	// The reason that transite to this state.
	Reason           *string `protobuf:"bytes,4,opt,name=reason" json:"reason,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *Snapshot) Reset()                    { *m = Snapshot{} }
func (m *Snapshot) String() string            { return proto.CompactTextString(m) }
func (*Snapshot) ProtoMessage()               {}
func (*Snapshot) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Snapshot) GetId() string {
	if m != nil && m.Id != nil {
		return *m.Id
	}
	return ""
}

func (m *Snapshot) GetTime() *google_protobuf.Timestamp {
	if m != nil {
		return m.Time
	}
	return nil
}

func (m *Snapshot) GetState() State {
	if m != nil && m.State != nil {
		return *m.State
	}
	return State_CREATED
}

func (m *Snapshot) GetReason() string {
	if m != nil && m.Reason != nil {
		return *m.Reason
	}
	return ""
}

type History struct {
	Snapshot         []*Snapshot `protobuf:"bytes,1,rep,name=snapshot" json:"snapshot,omitempty"`
	XXX_unrecognized []byte      `json:"-"`
}

func (m *History) Reset()                    { *m = History{} }
func (m *History) String() string            { return proto.CompactTextString(m) }
func (*History) ProtoMessage()               {}
func (*History) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *History) GetSnapshot() []*Snapshot {
	if m != nil {
		return m.Snapshot
	}
	return nil
}

type Priority struct {
	Priority         *int64 `protobuf:"varint,1,opt,name=priority" json:"priority,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (m *Priority) Reset()                    { *m = Priority{} }
func (m *Priority) String() string            { return proto.CompactTextString(m) }
func (*Priority) ProtoMessage()               {}
func (*Priority) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *Priority) GetPriority() int64 {
	if m != nil && m.Priority != nil {
		return *m.Priority
	}
	return 0
}

type Task struct {
	// Id has format: ServerIP:ProcessID:Timestamp, where timestamp has
	// millisecond resolution.
	Id *string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	// A task could contain multiple types. For example, a task could be
	// a repeated reminder task.
	Type []Type `protobuf:"varint,2,rep,name=type,enum=common.Type" json:"type,omitempty"`
	// History of all task snapshots.
	History *History `protobuf:"bytes,7,opt,name=history" json:"history,omitempty"`
	// When the task is created.
	CreationTime *google_protobuf.Timestamp `protobuf:"bytes,4,opt,name=creation_time,json=creationTime" json:"creation_time,omitempty"`
	// Most recent time the task is modified.
	LastModifiedTime *google_protobuf.Timestamp `protobuf:"bytes,5,opt,name=last_modified_time,json=lastModifiedTime" json:"last_modified_time,omitempty"`
	// Priority of the task. Used to compare with other tasks to decide
	// what task to do next.
	Priority                     *Priority `protobuf:"bytes,6,opt,name=priority" json:"priority,omitempty"`
	proto.XXX_InternalExtensions `json:"-"`
	XXX_unrecognized             []byte `json:"-"`
}

func (m *Task) Reset()                    { *m = Task{} }
func (m *Task) String() string            { return proto.CompactTextString(m) }
func (*Task) ProtoMessage()               {}
func (*Task) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

var extRange_Task = []proto.ExtensionRange{
	{10000, 536870911},
}

func (*Task) ExtensionRangeArray() []proto.ExtensionRange {
	return extRange_Task
}

func (m *Task) GetId() string {
	if m != nil && m.Id != nil {
		return *m.Id
	}
	return ""
}

func (m *Task) GetType() []Type {
	if m != nil {
		return m.Type
	}
	return nil
}

func (m *Task) GetHistory() *History {
	if m != nil {
		return m.History
	}
	return nil
}

func (m *Task) GetCreationTime() *google_protobuf.Timestamp {
	if m != nil {
		return m.CreationTime
	}
	return nil
}

func (m *Task) GetLastModifiedTime() *google_protobuf.Timestamp {
	if m != nil {
		return m.LastModifiedTime
	}
	return nil
}

func (m *Task) GetPriority() *Priority {
	if m != nil {
		return m.Priority
	}
	return nil
}

func init() {
	proto.RegisterType((*Snapshot)(nil), "common.Snapshot")
	proto.RegisterType((*History)(nil), "common.History")
	proto.RegisterType((*Priority)(nil), "common.Priority")
	proto.RegisterType((*Task)(nil), "common.Task")
	proto.RegisterEnum("common.Type", Type_name, Type_value)
	proto.RegisterEnum("common.State", State_name, State_value)
}

func init() { proto.RegisterFile("task.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 447 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x51, 0x4f, 0x6f, 0xd3, 0x30,
	0x14, 0x5f, 0xfe, 0xb4, 0xcd, 0x5e, 0xbb, 0x62, 0x59, 0x08, 0x45, 0xbd, 0x10, 0x05, 0x09, 0x85,
	0x6a, 0xca, 0xa4, 0x5e, 0x38, 0xa2, 0x92, 0x58, 0xac, 0xac, 0x6b, 0x2b, 0xb7, 0x12, 0x70, 0xaa,
	0xc2, 0x9a, 0x6d, 0xd6, 0x96, 0x3a, 0x8a, 0xcd, 0xa1, 0xb7, 0xde, 0xb8, 0xf2, 0x2d, 0xf8, 0x9a,
	0xc8, 0x76, 0x12, 0x90, 0x38, 0xf4, 0x96, 0xf7, 0xfb, 0xe7, 0xf7, 0x7e, 0x01, 0x90, 0x99, 0x78,
	0x8a, 0xcb, 0x8a, 0x4b, 0x8e, 0xbb, 0x77, 0xbc, 0x28, 0xf8, 0x7e, 0xf4, 0xfa, 0x81, 0xf3, 0x87,
	0xe7, 0xfc, 0x4a, 0xa3, 0xdf, 0x7f, 0xdc, 0x5f, 0x49, 0x56, 0xe4, 0x42, 0x66, 0x45, 0x69, 0x84,
	0xe1, 0x4f, 0x0b, 0xbc, 0xf5, 0x3e, 0x2b, 0xc5, 0x23, 0x97, 0x78, 0x08, 0x36, 0xdb, 0xf9, 0x56,
	0x60, 0x45, 0xe7, 0xd4, 0x66, 0x3b, 0x1c, 0x83, 0xab, 0xf4, 0xbe, 0x1d, 0x58, 0x51, 0x7f, 0x32,
	0x8a, 0x4d, 0x58, 0xdc, 0x84, 0xc5, 0x9b, 0x26, 0x8c, 0x6a, 0x1d, 0x7e, 0x03, 0x1d, 0x21, 0x33,
	0x99, 0xfb, 0x4e, 0x60, 0x45, 0xc3, 0xc9, 0x45, 0x6c, 0xb6, 0x88, 0xd7, 0x0a, 0xa4, 0x86, 0xc3,
	0xaf, 0xa0, 0x5b, 0xe5, 0x99, 0xe0, 0x7b, 0xdf, 0xd5, 0x0f, 0xd5, 0x53, 0xf8, 0x1e, 0x7a, 0xd7,
	0x4c, 0x48, 0x5e, 0x1d, 0xf0, 0x25, 0x78, 0xa2, 0xde, 0xc9, 0xb7, 0x02, 0x27, 0xea, 0x4f, 0x50,
	0x1b, 0x55, 0xe3, 0xb4, 0x55, 0x84, 0x6f, 0xc1, 0x5b, 0x55, 0x8c, 0x57, 0x4c, 0x1e, 0xf0, 0x08,
	0xbc, 0xb2, 0xfe, 0xd6, 0x77, 0x38, 0xb4, 0x9d, 0xc3, 0xdf, 0x36, 0xb8, 0x9b, 0x4c, 0x3c, 0xfd,
	0x77, 0x66, 0x00, 0xae, 0x3c, 0x94, 0xea, 0x4c, 0x27, 0x1a, 0x4e, 0x06, 0xcd, 0x53, 0x9b, 0x43,
	0x99, 0x53, 0xcd, 0xe0, 0x77, 0xd0, 0x7b, 0x34, 0xbb, 0xf9, 0x3d, 0xdd, 0xc5, 0x8b, 0x46, 0x54,
	0xaf, 0x4c, 0x1b, 0x1e, 0x7f, 0x80, 0x8b, 0xbb, 0x2a, 0xcf, 0x24, 0xe3, 0xfb, 0xad, 0x2e, 0xcf,
	0x3d, 0x59, 0xde, 0xa0, 0x31, 0x28, 0x08, 0x5f, 0x03, 0x7e, 0xce, 0x84, 0xdc, 0x16, 0x7c, 0xc7,
	0xee, 0x59, 0xbe, 0x33, 0x29, 0x9d, 0x93, 0x29, 0x48, 0xb9, 0x6e, 0x6b, 0x93, 0x4e, 0xba, 0xfc,
	0xa7, 0x8c, 0xae, 0xf6, 0xb7, 0x35, 0x36, 0x85, 0xfd, 0xad, 0x67, 0x7c, 0xee, 0xfd, 0x5a, 0xa0,
	0xe3, 0xf1, 0x78, 0xb4, 0xc7, 0x21, 0xb8, 0xea, 0x78, 0x3c, 0x00, 0x8f, 0x92, 0xdb, 0xd9, 0x22,
	0x25, 0x14, 0x9d, 0x99, 0x69, 0x45, 0xa6, 0x1b, 0x92, 0x22, 0x6b, 0x5c, 0x40, 0x47, 0xff, 0x56,
	0xdc, 0x87, 0x5e, 0x42, 0x0d, 0x7a, 0x86, 0x5f, 0x02, 0xfa, 0xb2, 0xa4, 0x37, 0xdb, 0xd9, 0x62,
	0xbb, 0xa2, 0xcb, 0x4f, 0x94, 0xac, 0xd7, 0xc8, 0x32, 0xce, 0xcf, 0x24, 0x51, 0x1a, 0x5b, 0x19,
	0x3e, 0xce, 0x97, 0xc9, 0x0d, 0x49, 0x91, 0xa3, 0xa8, 0x64, 0xba, 0x48, 0xc8, 0x9c, 0xa4, 0xc8,
	0x55, 0x54, 0x4a, 0xe6, 0xd3, 0x6f, 0x24, 0x45, 0x1d, 0x35, 0x90, 0xaf, 0xab, 0x19, 0x25, 0x29,
	0xea, 0xfe, 0x09, 0x00, 0x00, 0xff, 0xff, 0x63, 0x25, 0x0d, 0xdd, 0xdd, 0x02, 0x00, 0x00,
}