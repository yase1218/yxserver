package processor

import (
	"encoding/binary"
	"errors"
	"fmt"
	"reflect"

	"github.com/v587-zyf/gc/log"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

type Processor struct {
	msgID   map[reflect.Type]uint32
	msgType map[uint32]reflect.Type
}

func NewProcessor() *Processor {
	return &Processor{
		msgID:   make(map[reflect.Type]uint32),
		msgType: make(map[uint32]reflect.Type),
	}
}

func (p *Processor) RegisterWithMsgId(msgId uint32, msg proto.Message) uint32 {
	msgType := reflect.TypeOf(msg)
	if msgType == nil || msgType.Kind() != reflect.Ptr {
		log.Panic("protobuf message pointer required")
	}
	if _, ok := p.msgID[msgType]; ok {
		log.Panic("msg already registered", zap.Reflect("msgType", msgType))
	}

	p.msgID[msgType] = msgId
	p.msgType[msgId] = msgType
	return msgId
}

func (p *Processor) GetMsgId(msg proto.Message) uint32 {
	msgType := reflect.TypeOf(msg)
	if msgType == nil || msgType.Kind() != reflect.Ptr {
		log.Error("protobuf message pointer required")
		return 0
	}

	if id, ok := p.msgID[msgType]; ok {
		return id
	}
	log.Error("msg id not exists", zap.Reflect("msgType", msgType))
	return 0
}

func (p *Processor) Unmarshal(data []byte) (uint32, uint32, proto.Message, error) {
	if len(data) < 8 {
		return 0, 0, nil, errors.New("protobuf data too short")
	}

	packetId := binary.LittleEndian.Uint32(data[4:])
	id := binary.LittleEndian.Uint32(data[8:])

	if t, ok := p.msgType[id]; ok {
		msg := reflect.New(t.Elem()).Interface()
		pb_msg, pb_ok := msg.(proto.Message)
		if !pb_ok {
			return 0, 0, nil, fmt.Errorf("msgId:%v reflect failed", id)
		}
		return packetId, id, pb_msg, proto.Unmarshal(data[12:], msg.(proto.Message))
	}
	return 0, 0, nil, fmt.Errorf("msgId:%v msgType not found", id)
}

func (p *Processor) UnmarshalUnlen(data []byte) (uint32, uint32, proto.Message, error) {
	if len(data) < 8 {
		return 0, 0, nil, errors.New("protobuf data too short")
	}

	packetId := binary.LittleEndian.Uint32(data)
	id := binary.LittleEndian.Uint32(data[4:])

	if t, ok := p.msgType[id]; ok {
		msg := reflect.New(t.Elem()).Interface()
		pb_msg, pb_ok := msg.(proto.Message)
		if !pb_ok {
			return 0, 0, nil, fmt.Errorf("msgId:%v reflect failed", id)
		}
		return packetId, id, pb_msg, proto.Unmarshal(data[8:], msg.(proto.Message))
	}
	return 0, 0, nil, fmt.Errorf("msgId:%v msgType not found", id)
}

func (p *Processor) UnmarshlUnHead(id uint32, data []byte) (proto.Message, error) {
	if t, ok := p.msgType[id]; ok {
		msg := reflect.New(t.Elem()).Interface()
		pb_msg, pb_ok := msg.(proto.Message)
		if !pb_ok {
			return nil, fmt.Errorf("msgId:%v reflect failed", id)
		}
		return pb_msg, proto.Unmarshal(data, pb_msg)
	}
	return nil, fmt.Errorf("msgId:%v msgType not found", id)
}

func (p *Processor) Marshal(packetID uint32, message proto.Message) ([]byte, error) {
	messageType := reflect.TypeOf(message)

	// id
	id, ok := p.msgID[messageType]
	if !ok {
		return nil, fmt.Errorf("message %s not registered", messageType)
	}

	// data
	data, err := proto.Marshal(message)
	if err != nil {
		return nil, err
	}

	buffer := make([]byte, 12+len(data))
	binary.LittleEndian.PutUint32(buffer, uint32(len(data)+12))
	binary.LittleEndian.PutUint32(buffer[4:], packetID)
	binary.LittleEndian.PutUint32(buffer[8:], id)
	copy(buffer[12:], data) // use copy instead of append to avoid creating a slice
	return buffer, nil
}
