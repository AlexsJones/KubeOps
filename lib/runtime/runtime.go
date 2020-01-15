package runtime

import (
	"github.com/AlexsJones/kubeops/lib/kubernetes"
	"github.com/AlexsJones/kubeops/lib/subscription"
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"time"
)

func EventBuffer(context string, registry *subscription.Registry) {
	// ----------------------------------------------------------------------
	start := time.Now()
	log.Infof("Starting @ %s", start.String())
	log.Info("Got kubernetes client...")

	_, client, err := kubernetes.GetKubeClient(context)
	if err != nil {
		log.Fatal(err)
	}
	log.Info("Started event buffer...")
	// ----------------------------------------------------------------------
	processChan := func(c <-chan watch.Event) {
		select {
		case update, hasUpdate := <-c:
			log.Debug("Channel trigger...")
			if hasUpdate {
				err := registry.OnEvent(subscription.Message{
					Event:  update,
					Client: client,
				})
				if err != nil {
					log.Error(err)
				}
			}
		}
		log.Debug("Finished watch cycle...")
	}

	go func() {
		pi := client.CoreV1().Pods("")
		w, err := pi.Watch(metav1.ListOptions{Watch:true})
		wChan := w.ResultChan()
		if err != nil {
			panic(err)
		}
		for {
			processChan(wChan)
		}
	}()

	go func() {
		di := client.AppsV1().Deployments("")
		w, err := di.Watch(metav1.ListOptions{Watch:true})
		wChan := w.ResultChan()
		if err != nil {
			panic(err)
		}
		for {
			processChan(wChan)
		}
	}()

	go func() {
		ci := client.CoreV1().ConfigMaps("")
		w, err := ci.Watch(metav1.ListOptions{ Watch:true})
		wChan := w.ResultChan()
		if err != nil {
			panic(err)
		}
		for {
			processChan(wChan)
		}
	}()

	for {
		time.Sleep(time.Second * 2)
	}
}
