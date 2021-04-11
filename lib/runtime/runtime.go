package runtime

import (
	"context"
	"errors"
	"io"
	"math/rand"
	"sync"
	"time"

	"github.com/AlexsJones/KubeOps/lib/metrics"
	"github.com/AlexsJones/KubeOps/lib/subscription"
	"github.com/AlexsJones/KubeOps/lib/watcher"
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	k "k8s.io/client-go/kubernetes"
)

var (
	minWatchTimeout = 5 * time.Minute
	timeoutSeconds  = int64(minWatchTimeout.Seconds() * (rand.Float64() + 1.0))
)

func EventBuffer(context context.Context, client k.Interface,
	registry *subscription.Registry, obj []watcher.IObject) error {

	if len(obj) == 0 {
		return errors.New("no watchers selected, exiting.")
	}
	var watchers []<-chan watch.Event
	for _, o := range obj {
		funcObj := o
		w, err := funcObj.Watch(context, metav1.ListOptions{
			TimeoutSeconds:      &timeoutSeconds,
			AllowWatchBookmarks: true})
		defer w.Stop()
		if err != nil {
			switch {
			case err == io.EOF:
				// watch closed normally
			case err == io.ErrUnexpectedEOF:
				log.Infof("closed with unexpected EOF")
			}
		}
		watchers = append(watchers, w.ResultChan())
	}
	log.Debugf("%+v", watchers)

	var wg sync.WaitGroup
	wg.Add(len(watchers))
	for x, o := range watchers {
		go func(t int, c <-chan watch.Event) error {
			defer wg.Done()
			counter := 0
			for {
				select {
				case update, hasUpdate := <-c:
					if hasUpdate {
						err := registry.OnEvent(subscription.Message{
							Event:  update,
							Client: client,
						})
						if err != nil {
							return err
						}
						metrics.TotalEventOps.Inc()
					}
				}
				counter++
			}

		}(x, o)
	}
	wg.Wait()

	return nil
}
