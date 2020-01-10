package runtime

import (
	"github.com/AlexsJones/kubeops/lib/kubernetes"
	"github.com/AlexsJones/kubeops/lib/watcher"
	log "github.com/sirupsen/logrus"
	"time"
)

func EventBuffer(context string) {
	// ----------------------------------------------------------------------
	start := time.Now()
	log.Info("Starting @ %s", start.String())
	log.Info("Got kubernetes client...")

	_, client, err := kubernetes.GetKubeClient(context)
	if err != nil {
		log.Fatal(err)
	}
	log.Info("Started event buffer...")
	// ----------------------------------------------------------------------
	deploymentUpdates := watcher.DeploymentChannel(client)
	select {
		case update,hasUpdate := <- deploymentUpdates:
			if hasUpdate {
				for _, i := range update.Items {
					log.Info(i.Name)
				}
			}
			break

	}
	for {
		time.Sleep(time.Second)
	}
}
