package operators

import (
	"github.com/AlexsJones/kubeops/lib/subscription"
	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/watch"
)

type ExampleDeploymentOperator struct{}

func (ExampleDeploymentOperator) WithElectedResource() interface{} {

	return &v1.Deployment{}
}

func (ExampleDeploymentOperator) WithEventType() []watch.EventType {

	return []watch.EventType {watch.Added, watch.Deleted, watch.Modified}
}

func (ExampleDeploymentOperator) OnEvent(msg subscription.Message) {

	d := msg.Event.Object.(*v1.Deployment)

	log.Debugf("Deployment %s has %d Available replicas",d.Name,d.Status.AvailableReplicas)

}
