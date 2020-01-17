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
  "time"

  "github.com/AlexsJones/kubeops/lib/kubernetes"
  "github.com/AlexsJones/kubeops/lib/runtime"
  "github.com/AlexsJones/kubeops/lib/subscription"
  "github.com/AlexsJones/kubeops/operators"
  log "github.com/sirupsen/logrus"
)

var (
  c string
)

func main() {

  log.SetLevel(log.DebugLevel)

  flag.StringVar(&c,"context","",
    "Kubernetes context")
  flag.Parse()

  registry := &subscription.Registry{
    Subscriptions: []subscription.ISubscription{
      operators.ExamplePodOperator{},
      operators.ExampleDeploymentOperator{},

    },
  }

  start := time.Now()
  log.Infof("Starting @ %s", start.String())
  log.Info("Got kubernetes client...")

  _, client, err := kubernetes.GetKubeClient(c)
  if err != nil {
    log.Fatal(err)
  }
  log.Info("Started event buffer...")

  ctx, _ := context.WithCancel(context.Background())


  runtime.EventBuffer(ctx, client, registry,[]kubernetes.IObject{
    client.CoreV1().Pods(""),
    client.AppsV1().Deployments(""),
    client.CoreV1().ConfigMaps(""),
  })

}