package watcher

import (
	log "github.com/sirupsen/logrus"
	dv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"time"
)

func DeploymentChannel(client kubernetes.Interface) <- chan *dv1.DeploymentList {
	updates := make(chan *dv1.DeploymentList)
	go func() {
		for {
			log.Debug("Checking deployments...")
			di := client.AppsV1().Deployments("")
			dlist, err := di.List(v1.ListOptions{})
			if err != nil {
				log.Warn(err)
			}

			if len(dlist.Items) > 0 {
				log.Debug("Found deployments...")
				updates <- dlist
			}
			time.Sleep(time.Second * 2)
		}
	}()
	return updates
}
