# KubeOps


<img src="image/SPACEGIRL_GOPHER.png" data-canonical-src="image/SPACEGIRL_GOPHER.png" width="300" />

- A kubernetes Operator built in go and leveraging `Kind` for speedy local development.
- Utilises golang to express your operator requirements not some DSL or CRD system.

### Development perks

- Uses a local registry to build a golang project then push into the k8s cluster directly
- Example of using k8s golang API
- Example of using helm

## Requirements

|   |
|---|
| Kind  |   
| Golang |  
| Docker |
| Helm |  


## Install something and get it running

`make up`

`make`

These commands should be enough to create the kind cluster and install `kubeops`


## Commands

|   |   |
|---|---|
| make up  | Creates a kind cluster   |
| make down | Deletes the kind cluster  |
| make | Builds the project, dockerfile, side loads then installs into the cluster |
| make delete | Deletes all currently installed helm releases  |

## Development without having to push an image

`go run main.go --context kind-kind `

Allows you to connect using the local kubeconfig to the cluster and operate externally.

## Development and pushing to the cluster

1. `<write some code> `

2. Test with the above (`go run main.go --context kind-kind `)

3. `make`

4. View your changes in the cluster

# Development

You will see in the code `operators` defines how we filter and act on events

```go
func (ExamplePodOperator) WithFilter() interface{} {

	return &v1.Pod{}
}

func (ExamplePodOperator) OnEvent(msg subscription.Message) {

	log.Debug("Pod event ----> %v", msg.Event)
}
```
This example struct adheres to the `ISubscription` interface.

Once you've created the operator add it into the main.go

```go

  registry := &subscription.Registry{
    Subscriptions: []subscription.ISubscription{
      operators.ExamplePodOperator{},

    },
  }
```

That's it - now you will get filtered watch events for your type (Providing they are in the `lib/watcher` GenerateWatchers channel setup.)


## Err can you explain in more detail...?

Out of the box you get two channels registered in `lib/watcher/deployment`

```go
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
```

These will automatically get picked up in the runtime. Add more if you like.

_Anything you have registered in the watch channel you can filter for in the ISubscription_

```
type ISubscription interface {
	OnEvent(msg Message)
	WithFilter() interface{}
}

```

The WithFilter is important as it allows the runtime to selectively send messages to your handler `OnEvent`
You'll also get the kubernetes client interface passed through to you to act on any incoming messages.

_For example_

```go
switch msg.Event.Type {
	case watch.Added:
        //Basic on some precondition...
		pod := msg.Event.Object.(*v1.Pod)
		err := msg.Client.CoreV1().Pods(pod.Namespace).Delete(pod.Name,&metav1.DeleteOptions{})
		if err != nil {
			log.Error(err)
		}

	}
```