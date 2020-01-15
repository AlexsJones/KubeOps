package operators

import (
	"github.com/AlexsJones/kubeops/lib/subscription"
	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/watch"
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
	if pod.Labels["app.kubernetes.io/name"] == "kubeops" {

		log.Debugf("%v",pod)
		existingLabels := pod.Labels

		if _, ok := existingLabels["sneaky-label"]; !ok {
			//Let's add a label...
			existingLabels["sneaky-label"] = "is-an-example"
			pod.SetLabels(existingLabels)
			// Invoke a new pod client interface
			pi := msg.Client.CoreV1().Pods(pod.Namespace)
			if _, err := pi.Update(pod); err != nil {
				log.Error(err)
			}else {
				log.Debug("Added a new label...")
			}
		}
	}
}
