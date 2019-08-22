// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: model.proto

package proto

import (
	fmt "fmt"
	io "io"
	math "math"
	reflect "reflect"
	strconv "strconv"
	strings "strings"

	proto "github.com/gogo/protobuf/proto"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

type Command int32

const (
	Unknown    Command = 0
	Ping       Command = 1
	Pong       Command = 2
	Connect    Command = 3
	DisConnect Command = 4
	Packet     Command = 5
)

var Command_name = map[int32]string{
	0: "Unknown",
	1: "Ping",
	2: "Pong",
	3: "Connect",
	4: "DisConnect",
	5: "Packet",
}

var Command_value = map[string]int32{
	"Unknown":    0,
	"Ping":       1,
	"Pong":       2,
	"Connect":    3,
	"DisConnect": 4,
	"Packet":     5,
}

func (Command) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_4c16552f9fdb66d8, []int{0}
}

type Message struct {
	Cmd Command `protobuf:"varint,1,opt,name=cmd,proto3,enum=proto.Command" json:"cmd,omitempty"`
}

func (m *Message) Reset()      { *m = Message{} }
func (*Message) ProtoMessage() {}
func (*Message) Descriptor() ([]byte, []int) {
	return fileDescriptor_4c16552f9fdb66d8, []int{0}
}
func (m *Message) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Message) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Message.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Message) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Message.Merge(m, src)
}
func (m *Message) XXX_Size() int {
	return m.Size()
}
func (m *Message) XXX_DiscardUnknown() {
	xxx_messageInfo_Message.DiscardUnknown(m)
}

var xxx_messageInfo_Message proto.InternalMessageInfo

func (m *Message) GetCmd() Command {
	if m != nil {
		return m.Cmd
	}
	return Unknown
}

func init() {
	proto.RegisterEnum("proto.Command", Command_name, Command_value)
	proto.RegisterType((*Message)(nil), "proto.Message")
}

func init() { proto.RegisterFile("model.proto", fileDescriptor_4c16552f9fdb66d8) }

var fileDescriptor_4c16552f9fdb66d8 = []byte{
	// 210 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0xce, 0xcd, 0x4f, 0x49,
	0xcd, 0xd1, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x05, 0x53, 0x4a, 0xda, 0x5c, 0xec, 0xbe,
	0xa9, 0xc5, 0xc5, 0x89, 0xe9, 0xa9, 0x42, 0x0a, 0x5c, 0xcc, 0xc9, 0xb9, 0x29, 0x12, 0x8c, 0x0a,
	0x8c, 0x1a, 0x7c, 0x46, 0x7c, 0x10, 0x65, 0x7a, 0xce, 0xf9, 0xb9, 0xb9, 0x89, 0x79, 0x29, 0x41,
	0x20, 0x29, 0xad, 0x60, 0x2e, 0x76, 0x28, 0x5f, 0x88, 0x9b, 0x8b, 0x3d, 0x34, 0x2f, 0x3b, 0x2f,
	0xbf, 0x3c, 0x4f, 0x80, 0x41, 0x88, 0x83, 0x8b, 0x25, 0x20, 0x33, 0x2f, 0x5d, 0x80, 0x11, 0xcc,
	0xca, 0xcf, 0x4b, 0x17, 0x60, 0x02, 0x29, 0x70, 0xce, 0xcf, 0xcb, 0x4b, 0x4d, 0x2e, 0x11, 0x60,
	0x16, 0xe2, 0xe3, 0xe2, 0x72, 0xc9, 0x2c, 0x86, 0xf1, 0x59, 0x84, 0xb8, 0xb8, 0xd8, 0x02, 0x12,
	0x93, 0xb3, 0x53, 0x4b, 0x04, 0x58, 0x9d, 0x4c, 0x2e, 0x3c, 0x94, 0x63, 0xb8, 0xf1, 0x50, 0x8e,
	0xe1, 0xc3, 0x43, 0x39, 0xc6, 0x86, 0x47, 0x72, 0x8c, 0x2b, 0x1e, 0xc9, 0x31, 0x9e, 0x78, 0x24,
	0xc7, 0x78, 0xe1, 0x91, 0x1c, 0xe3, 0x83, 0x47, 0x72, 0x8c, 0x2f, 0x1e, 0xc9, 0x31, 0x7c, 0x78,
	0x24, 0xc7, 0x38, 0xe1, 0xb1, 0x1c, 0xc3, 0x85, 0xc7, 0x72, 0x0c, 0x37, 0x1e, 0xcb, 0x31, 0x24,
	0xb1, 0x81, 0x9d, 0x67, 0x0c, 0x08, 0x00, 0x00, 0xff, 0xff, 0x98, 0x81, 0x83, 0xd1, 0xd4, 0x00,
	0x00, 0x00,
}

func (x Command) String() string {
	s, ok := Command_name[int32(x)]
	if ok {
		return s
	}
	return strconv.Itoa(int(x))
}
func (this *Message) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Message)
	if !ok {
		that2, ok := that.(Message)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if this.Cmd != that1.Cmd {
		return false
	}
	return true
}
func (this *Message) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 5)
	s = append(s, "&proto.Message{")
	s = append(s, "Cmd: "+fmt.Sprintf("%#v", this.Cmd)+",\n")
	s = append(s, "}")
	return strings.Join(s, "")
}
func valueToGoStringModel(v interface{}, typ string) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("func(v %v) *%v { return &v } ( %#v )", typ, typ, pv)
}
func (m *Message) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Message) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.Cmd != 0 {
		dAtA[i] = 0x8
		i++
		i = encodeVarintModel(dAtA, i, uint64(m.Cmd))
	}
	return i, nil
}

func encodeVarintModel(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *Message) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Cmd != 0 {
		n += 1 + sovModel(uint64(m.Cmd))
	}
	return n
}

func sovModel(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozModel(x uint64) (n int) {
	return sovModel(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (this *Message) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&Message{`,
		`Cmd:` + fmt.Sprintf("%v", this.Cmd) + `,`,
		`}`,
	}, "")
	return s
}
func valueToStringModel(v interface{}) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("*%v", pv)
}
func (m *Message) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowModel
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Message: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Message: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Cmd", wireType)
			}
			m.Cmd = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowModel
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Cmd |= Command(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipModel(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthModel
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthModel
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipModel(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowModel
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowModel
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
			return iNdEx, nil
		case 1:
			iNdEx += 8
			return iNdEx, nil
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowModel
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthModel
			}
			iNdEx += length
			if iNdEx < 0 {
				return 0, ErrInvalidLengthModel
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowModel
					}
					if iNdEx >= l {
						return 0, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					innerWire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				innerWireType := int(innerWire & 0x7)
				if innerWireType == 4 {
					break
				}
				next, err := skipModel(dAtA[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
				if iNdEx < 0 {
					return 0, ErrInvalidLengthModel
				}
			}
			return iNdEx, nil
		case 4:
			return iNdEx, nil
		case 5:
			iNdEx += 4
			return iNdEx, nil
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
	}
	panic("unreachable")
}

var (
	ErrInvalidLengthModel = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowModel   = fmt.Errorf("proto: integer overflow")
)
