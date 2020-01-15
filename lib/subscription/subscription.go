package subscription

import "k8s.io/apimachinery/pkg/watch"

type ISubscription interface {
	OnEvent(msg Message)
	WithElectedResource() interface{}
	// Return an empty array for all event types or selective set an event type to filter for
	//Added    EventType = "ADDED"
	//Modified EventType = "MODIFIED"
	//Deleted  EventType = "DELETED"
	//Bookmark EventType = "BOOKMARK"
	//Error    EventType = "ERROR"
	WithEventType() []watch.EventType
}
