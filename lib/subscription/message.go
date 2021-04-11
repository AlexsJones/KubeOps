package subscription

import (
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
)

type Message struct {
	Client kubernetes.Interface
	Event watch.Event
}
