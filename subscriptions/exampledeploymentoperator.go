package subscriptions

import (
	"github.com/AlexsJones/KubeOps/lib/subscription"
	v1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/klog"
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

	klog.Infof("Deployment %s has %d Available replicas",d.Name,d.Status.AvailableReplicas)

	
}
