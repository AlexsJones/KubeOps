package subscription

import (
	"k8s.io/apimachinery/pkg/runtime"
	"reflect"
)

type Registry struct {
	Subscriptions []ISubscription
}


func (r *Registry) Add(subscription ISubscription) error {

	r.Subscriptions = append(r.Subscriptions, subscription)
	return nil
}

func (r *Registry) OnEvent(msg Message) error {

	unstructuredObj, err := runtime.DefaultUnstructuredConverter.ToUnstructured(msg.Event.Object)
	if err != nil {
		return err
	}
	newObj := reflect.New(reflect.TypeOf(msg.Event.Object).Elem()).Interface().(runtime.Object)
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(unstructuredObj, newObj)
	if err != nil {
		return err
	}

	for _, subscription := range r.Subscriptions {

		//log.Debugf("Comparing subscription %v with event %v", reflect.TypeOf(subscription.WithElectedResource()),reflect.TypeOf(newObj))
		if reflect.TypeOf(subscription.WithElectedResource()) == reflect.TypeOf(newObj) {

			evts := subscription.WithEventType()
			if len(evts) == 0 {
				subscription.OnEvent(msg)
			}
			for _, evt := range evts {
				if evt == msg.Event.Type {
					subscription.OnEvent(msg)
				}
			}
		}
	}
	return nil
}