package watcher

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
)

func GenerateWatchers(client kubernetes.Interface) ([]<-chan watch.Event, error){

	var watchers []<-chan watch.Event

	// ---------------------------------------------------------------------------------------
	di := client.AppsV1().Deployments("")
	wi, err := di.Watch(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	watchers = append(watchers,wi.ResultChan())
	// ----------------------------------------------------------------------------------------
	pi := client.CoreV1().Pods("")
	wpi, err := pi.Watch(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	watchers = append(watchers,wpi.ResultChan())
	// ----------------------------------------------------------------------------------------
	return watchers, nil
}