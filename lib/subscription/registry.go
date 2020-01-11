package subscription

import (
	log "github.com/sirupsen/logrus"
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

func (r *Registry) OnEvent(msg Message) {

	unstructuredObj, err := runtime.DefaultUnstructuredConverter.ToUnstructured(msg.Event.Object)
	if err != nil {
		return
	}
	newObj := reflect.New(reflect.TypeOf(msg.Event.Object).Elem()).Interface().(runtime.Object)
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(unstructuredObj, newObj)
	if err != nil {
		return
	}

	for _, subscription := range r.Subscriptions {

		log.Debugf("Comparing subscription %v with event %v", reflect.TypeOf(subscription.WithFilter()),reflect.TypeOf(newObj))
		if reflect.TypeOf(subscription.WithFilter()) == reflect.TypeOf(newObj) {
			subscription.OnEvent(msg)
		}
	}
}