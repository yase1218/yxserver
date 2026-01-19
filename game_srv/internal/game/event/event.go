package event

type EventFunc func(event IEvent)

type IEvent interface {
	RouteID() uint64
	CallBack() EventFunc
}
