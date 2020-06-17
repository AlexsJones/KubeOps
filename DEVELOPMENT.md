# Development

You will see in the code `./operators` defines how we filter and act on events

Here is a simple example that adds a label to the kubeops pod when it is installed in the cluster as a deployment (through the helm chart).

```go 

type ExamplePodOperator struct{}

func (ExamplePodOperator) WithElectedResource() interface{} {

	return &v1.Pod{}
}

func (ExamplePodOperator) WithEventType() []watch.EventType {

	return []watch.EventType {watch.Added, watch.Modified}
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

```
This example struct adheres to the `ISubscription` interface.

Once you've created the operator add it into the main.go as a subscription.

```go
  registry := &subscription.Registry{
    Subscriptions: []subscription.ISubscription{
      operators.ExamplePodOperator{},
      operators.ExampleDeploymentOperator{},

    },
  }
```

_Anything you have registered in the watch channel you can filter for in the ISubscription_

```
type ISubscription interface {
	OnEvent(msg Message)
	WithFilter() interface{}
    WithEventType() []watch.EventType
}   

```

- The WithFilter is important as it allows the runtime to selectively send messages to your handler `OnEvent`
   You'll also get the kubernetes client interface passed through to you to act on any incoming messages.
- The WithEvent function is important as it allows selective event types to be filtered out of the incoming operator handler.


```go

  registry := &subscription.Registry{
    Subscriptions: []subscription.ISubscription{
      operators.ExamplePodOperator{},

    },
  }
```

That's it - now you will get filtered watch events for your type (Providing they are in the `lib/watcher` GenerateWatchers channel setup.)


## Adding watches for different Kubernetes resources

Out of the box you get some channels registered in `lib/watcher/deployment`

```go
//main.go
  runtime.EventBuffer(context,client, registry,[]kubernetes.IObject{
    client.CoreV1().Pods(""),
    client.AppsV1().Deployments(""),
    client.CoreV1().ConfigMaps(""),
  })

```

These will automatically get picked up in the runtime. 
