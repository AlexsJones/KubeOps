/*
Copyright © 2020 alexsimonjones@gmail.com

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package main

import (
	"context"
	"flag"
	"github.com/AlexsJones/KubeOps/lib/runtime"
	"github.com/AlexsJones/KubeOps/lib/subscription"
	"github.com/AlexsJones/KubeOps/lib/watcher"
	"github.com/AlexsJones/KubeOps/subscriptions"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog"
	"os"
	"os/signal"
	"time"
)

var (
	masterURL  string
	kubeconfig string
)


func main() {

	klog.InitFlags(nil)
	flag.Parse()

	start := time.Now()
	klog.Infof("Starting @ %s", start.String())
	klog.Info("Got watcher client...")

	cfg, err := clientcmd.BuildConfigFromFlags(masterURL,kubeconfig)
	if err != nil {
		klog.Fatalf("Error building kubeconfig: %s", err.Error())
	}

	kubeClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		klog.Fatalf("Error building watcher clientset: %s", err.Error())
	}

	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	defer func() {
		signal.Stop(c)
		cancel()
	}()
	go func() {
		select {
		case <-c:
			cancel()
		case <-ctx.Done():
		}
	}()

	err = runtime.EventBuffer(ctx, kubeClient, &subscription.Registry{
		Subscriptions: []subscription.ISubscription{
			// Subscribe to these built-in type events
			subscriptions.ExamplePodOperator{},
			subscriptions.ExampleFooCRDOperator{},
			subscriptions.ExampleDeploymentOperator{},
		},
	},[]watcher.IObject{
		kubeClient.CoreV1().Pods(""),
		kubeClient.AppsV1().Deployments(""),
		kubeClient.CoreV1().ConfigMaps(""),

	})
	if err != nil {
		klog.Error(err)
	}
}

func init() {
	flag.StringVar(&kubeconfig, "kubeconfig", "", "Path to a kubeconfig. Only required if out-of-cluster.")
	flag.StringVar(&masterURL, "master", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
}