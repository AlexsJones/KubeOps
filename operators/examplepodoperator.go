package operators

import (
	"github.com/AlexsJones/kubeops/lib/subscription"
	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/watch"
)

type ExamplePodOperator struct{}

func (ExamplePodOperator) WithFilter() interface{} {

	return &v1.Pod{}
}

func (ExamplePodOperator) OnEvent(msg subscription.Message) {

	log.Debug("Pod event ----> %v", msg.Event)

	switch msg.Event.Type {
	case watch.Added:
		break
	case watch.Modified:
		break
	case watch.Deleted:
		break
	case watch.Error:
		break
	case watch.Bookmark:

	}
}
