package runtime

import (
	"github.com/AlexsJones/kubeops/lib/kubernetes"
	"github.com/AlexsJones/kubeops/lib/subscription"
	"github.com/AlexsJones/kubeops/lib/watcher"
	log "github.com/sirupsen/logrus"
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
	watchers, err := watcher.GenerateWatchers(client)
	if err != nil {
		log.Fatal(err)
	}
	for {
		for _, w := range watchers {
			select {
			case update, hasUpdate := <-w:
				if hasUpdate {
					registry.OnEvent(subscription.Message{
						Event:update,
						Client:client,
					})
				}
				break
			}

		}
		time.Sleep(time.Second)
	}

}
