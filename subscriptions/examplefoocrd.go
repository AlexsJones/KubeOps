package subscriptions

import (
	"github.com/AlexsJones/KubeOps/lib/subscription"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/klog"
	// In this example we have imported the CRD type from the Kubernetes sample-controller project
	"k8s.io/sample-controller/pkg/apis/samplecontroller/v1alpha1"
)

type ExampleFooCRDOperator struct{}

func (ExampleFooCRDOperator) WithElectedResource() interface{} {

	return &v1alpha1.Foo{}
}

func (ExampleFooCRDOperator) WithEventType() []watch.EventType {

	return []watch.EventType {watch.Added, watch.Deleted, watch.Modified}
}

func (ExampleFooCRDOperator) OnEvent(msg subscription.Message) {

	_ = msg.Event.Object.(*v1alpha1.Foo)

	klog.Infof("Found our FOO CRD!")

}