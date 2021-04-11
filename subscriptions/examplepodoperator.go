package subscriptions

import (
	"context"
	"KubeOps/app/lib/subscription"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/klog"
)

type ExamplePodOperator struct{}

func (ExamplePodOperator) WithElectedResource() interface{} {

	return &v1.Pod{}
}

func (ExamplePodOperator) WithEventType() []watch.EventType {

	return []watch.EventType {watch.Added, watch.Deleted, watch.Modified}
}

func (ExamplePodOperator) OnEvent(msg subscription.Message) {
	// This is a silly example that checks the kubeops pod
	// When we add, delete or modify the pod
	// Then we add an example label to it.
	pod := msg.Event.Object.(*v1.Pod)
	klog.Infof("Incoming pod event from %s",pod.Name)
	if pod.Labels["app.watcher.io/name"] == "kubeops" {

		klog.Infof("%v",pod)
		existingLabels := pod.Labels

		if _, ok := existingLabels["sneaky-label"]; !ok {
			//Let's add a label...
			existingLabels["sneaky-label"] = "is-an-example"
			pod.SetLabels(existingLabels)
			// Invoke a new pod client interface
			pi := msg.Client.CoreV1().Pods(pod.Namespace)
			if _, err := pi.Update(context.TODO(),pod, metav1.UpdateOptions{}); err != nil {
				klog.Error(err)
			}else {
				klog.Infof("Added a new label...")
			}
		}
	}
}
