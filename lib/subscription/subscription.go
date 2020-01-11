package subscription

type ISubscription interface {
	OnEvent(msg Message)
	WithFilter() interface{}
}
