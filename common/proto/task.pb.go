// Code generated by protoc-gen-go. DO NOT EDIT.
// source: task.proto

/*
Package proto is a generated protocol buffer package.

It is generated from these files:
	task.proto
	user.proto

It has these top-level messages:
	Snapshot
	History
	Priority
	Task
	User
*/
package proto

import proto1 "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import google_protobuf "github.com/golang/protobuf/ptypes/timestamp"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto1.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto1.ProtoPackageIsVersion2 // please upgrade the proto package

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
	return proto1.EnumName(Type_name, int32(x))
}
func (x *Type) UnmarshalJSON(data []byte) error {
	value, err := proto1.UnmarshalJSONEnum(Type_value, data, "Type")
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
	return proto1.EnumName(State_name, int32(x))
}
func (x *State) UnmarshalJSON(data []byte) error {
	value, err := proto1.UnmarshalJSONEnum(State_value, data, "State")
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
	State *State `protobuf:"varint,3,opt,name=state,enum=proto.State" json:"state,omitempty"`
	// The reason that transite to this state.
	Reason           *string `protobuf:"bytes,4,opt,name=reason" json:"reason,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *Snapshot) Reset()                    { *m = Snapshot{} }
func (m *Snapshot) String() string            { return proto1.CompactTextString(m) }
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
func (m *History) String() string            { return proto1.CompactTextString(m) }
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
func (m *Priority) String() string            { return proto1.CompactTextString(m) }
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
	Type []Type `protobuf:"varint,2,rep,name=type,enum=proto.Type" json:"type,omitempty"`
	// History of all task snapshots.
	History *History `protobuf:"bytes,7,opt,name=history" json:"history,omitempty"`
	// When the task is created.
	CreationTime *google_protobuf.Timestamp `protobuf:"bytes,4,opt,name=creation_time,json=creationTime" json:"creation_time,omitempty"`
	// Most recent time the task is modified.
	LastModifiedTime *google_protobuf.Timestamp `protobuf:"bytes,5,opt,name=last_modified_time,json=lastModifiedTime" json:"last_modified_time,omitempty"`
	// Priority of the task. Used to compare with other tasks to decide
	// what task to do next.
	Priority                      *Priority `protobuf:"bytes,6,opt,name=priority" json:"priority,omitempty"`
	proto1.XXX_InternalExtensions `json:"-"`
	XXX_unrecognized              []byte `json:"-"`
}

func (m *Task) Reset()                    { *m = Task{} }
func (m *Task) String() string            { return proto1.CompactTextString(m) }
func (*Task) ProtoMessage()               {}
func (*Task) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

var extRange_Task = []proto1.ExtensionRange{
	{10000, 536870911},
}

func (*Task) ExtensionRangeArray() []proto1.ExtensionRange {
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
	proto1.RegisterType((*Snapshot)(nil), "proto.Snapshot")
	proto1.RegisterType((*History)(nil), "proto.History")
	proto1.RegisterType((*Priority)(nil), "proto.Priority")
	proto1.RegisterType((*Task)(nil), "proto.Task")
	proto1.RegisterEnum("proto.Type", Type_name, Type_value)
	proto1.RegisterEnum("proto.State", State_name, State_value)
}

func init() { proto1.RegisterFile("task.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 442 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x51, 0xcb, 0x6e, 0xda, 0x40,
	0x14, 0x8d, 0x5f, 0xe0, 0x5c, 0x28, 0x1d, 0x8d, 0xaa, 0xca, 0x62, 0x13, 0xe4, 0x45, 0x65, 0x11,
	0xc9, 0x91, 0x58, 0x74, 0x5b, 0x51, 0x7b, 0xd4, 0xd0, 0x10, 0x40, 0x03, 0x52, 0xdb, 0x15, 0x72,
	0x8b, 0x93, 0x8c, 0x12, 0x18, 0xcb, 0x33, 0x5d, 0xb0, 0xf3, 0xaa, 0xeb, 0x7e, 0x44, 0x3f, 0xb4,
	0x9a, 0x87, 0xd3, 0x4a, 0x59, 0xb0, 0x82, 0x73, 0xcf, 0x63, 0xee, 0x3d, 0x06, 0x90, 0x85, 0x78,
	0x4c, 0xab, 0x9a, 0x4b, 0x8e, 0x03, 0xfd, 0x33, 0xbc, 0xb8, 0xe7, 0xfc, 0xfe, 0xa9, 0xbc, 0xd2,
	0xe8, 0xfb, 0xcf, 0xbb, 0x2b, 0xc9, 0xf6, 0xa5, 0x90, 0xc5, 0xbe, 0x32, 0xba, 0xf8, 0x97, 0x03,
	0xe1, 0xfa, 0x50, 0x54, 0xe2, 0x81, 0x4b, 0x3c, 0x00, 0x97, 0xed, 0x22, 0x67, 0xe4, 0x24, 0xe7,
	0xd4, 0x65, 0x3b, 0x9c, 0x82, 0xaf, 0xf4, 0x91, 0x3b, 0x72, 0x92, 0xde, 0x64, 0x98, 0x9a, 0xb0,
	0xb4, 0x0d, 0x4b, 0x37, 0x6d, 0x18, 0xd5, 0x3a, 0x1c, 0x43, 0x20, 0x64, 0x21, 0xcb, 0xc8, 0x1b,
	0x39, 0xc9, 0x60, 0xd2, 0x37, 0xca, 0x74, 0xad, 0x66, 0xd4, 0x50, 0xf8, 0x2d, 0x74, 0xea, 0xb2,
	0x10, 0xfc, 0x10, 0xf9, 0xfa, 0x1d, 0x8b, 0xe2, 0xf7, 0xd0, 0xbd, 0x66, 0x42, 0xf2, 0xfa, 0x88,
	0x2f, 0x21, 0x14, 0x76, 0xa5, 0xc8, 0x19, 0x79, 0x49, 0x6f, 0xf2, 0xba, 0x4d, 0xb2, 0x63, 0xfa,
	0x2c, 0x88, 0xdf, 0x41, 0xb8, 0xaa, 0x19, 0xaf, 0x99, 0x3c, 0xe2, 0x21, 0x84, 0x95, 0xfd, 0xaf,
	0xaf, 0xf0, 0xe8, 0x33, 0x8e, 0xff, 0xb8, 0xe0, 0x6f, 0x0a, 0xf1, 0xf8, 0xe2, 0xc8, 0x0b, 0xf0,
	0xe5, 0xb1, 0x52, 0x47, 0x7a, 0xc9, 0x60, 0xd2, 0xb3, 0x2f, 0x6d, 0x8e, 0x55, 0x49, 0x35, 0x81,
	0x13, 0xe8, 0x3e, 0x98, 0xcd, 0xa2, 0xae, 0x2e, 0x62, 0x60, 0x35, 0x76, 0x5f, 0xda, 0xd2, 0xf8,
	0x03, 0xbc, 0xfa, 0x51, 0x97, 0x85, 0x64, 0xfc, 0xb0, 0xd5, 0xc5, 0xf9, 0x27, 0x8b, 0xeb, 0xb7,
	0x06, 0x35, 0xc2, 0xd7, 0x80, 0x9f, 0x0a, 0x21, 0xb7, 0x7b, 0xbe, 0x63, 0x77, 0xac, 0xdc, 0x99,
	0x94, 0xe0, 0x64, 0x0a, 0x52, 0xae, 0x5b, 0x6b, 0xd2, 0x49, 0x97, 0xff, 0x55, 0xd1, 0xd1, 0xfe,
	0xb6, 0xc3, 0xb6, 0xad, 0x7f, 0xdd, 0x8c, 0xcf, 0xc3, 0xdf, 0x0b, 0xd4, 0x34, 0x4d, 0xe3, 0x8e,
	0x63, 0xf0, 0xd5, 0xe9, 0xb8, 0x0f, 0x21, 0x25, 0xb7, 0xb3, 0x45, 0x4e, 0x28, 0x3a, 0x33, 0x68,
	0x45, 0xa6, 0x1b, 0x92, 0x23, 0x67, 0xbc, 0x87, 0x40, 0x7f, 0x52, 0xdc, 0x83, 0x6e, 0x46, 0xcd,
	0xf4, 0x0c, 0xbf, 0x01, 0xf4, 0x65, 0x49, 0x6f, 0xb6, 0xb3, 0xc5, 0x76, 0x45, 0x97, 0x9f, 0x28,
	0x59, 0xaf, 0x91, 0x63, 0x9c, 0x9f, 0x49, 0xa6, 0x34, 0xae, 0x32, 0x7c, 0x9c, 0x2f, 0xb3, 0x1b,
	0x92, 0x23, 0x4f, 0x51, 0xd9, 0x74, 0x91, 0x91, 0x39, 0xc9, 0x91, 0xaf, 0xa8, 0x9c, 0xcc, 0xa7,
	0xdf, 0x48, 0x8e, 0x02, 0x05, 0xc8, 0xd7, 0xd5, 0x8c, 0x92, 0x1c, 0x75, 0xfe, 0x06, 0x00, 0x00,
	0xff, 0xff, 0xa4, 0x55, 0x47, 0x82, 0xd7, 0x02, 0x00, 0x00,
}
