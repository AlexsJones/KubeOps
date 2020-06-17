package watcher

import (
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

type IObject interface {
	Watch(context context.Context, opts metav1.ListOptions) (watch.Interface, error)

}
