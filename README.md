<img src="image/logo_transparent.png" data-canonical-src="image/logo_transparent.png" width="300" />


A simple programmatic Kubernetes Operator template.

Use this to create your own Kubernetes operators with golang.

- Build with KIND (Kubernetes in Docker).
- Express custom behaviours in code - no DSL to learn.
- Generate your CRD's with controller-runtime and import them here.
- Works with built-in kubernetes resource types and custom resource definitions ( example included ).


<img src="image/local.gif" width="800" />

## Run the example...

Run builtin examples locally:

- `make up` to start a K.I.N.D cluster on Docker.
- `go run main.go --kubeconfig=/Users/<name>/.kube/config`

This creates some resource interfaces and subscribes to them with some basic subscriptions `./subscriptions`

_Resources to watch_

```go
  runtime.EventBuffer(ctx, kubeClient, registry,[]watcher.IObject{

    // Buffer events for these built-in types
    kubeClient.CoreV1().Pods(""),
    kubeClient.AppsV1().Deployments(""),
    kubeClient.CoreV1().ConfigMaps(""),
    //Example CRD imported into the runtime-----------------------------------------------------------
    //exampleClient.SamplecontrollerV1alpha1().Foos(""),
    // -----------------------------------------------------------------------------------------------
  })
```

_Subscriptions on the watched resources_

```go
  registry := &subscription.Registry{
    Subscriptions: []subscription.ISubscription{
      // Subscribe to these built-in type events
      subscriptions.ExamplePodOperator{},
      subscriptions.ExampleFooCRDOperator{},
      subscriptions.ExampleDeploymentOperator{},
    },
  }
```


Build docker image and install into cluster locally:

- `make`


_Please see [REQUIREMENTS.md](REQUIREMENTS.md) for installation requirements_

_Please see [LICENCE.md](LICENCE.md) for licence enquiries_

<img src="image/kubeops.png" width="800" />

### Development perks

- A simple golang based implementation of an Operator with the boiler plate done.
- Has a pre-made helm chart, so you can build an image of this code and push it into a cluster with your changes.
- Example of using k8s golang API
- Example of using helm


## Commands

|   |   |
|---|---|
| make up  | Creates a kind cluster   |
| make down | Deletes the kind cluster  |
| make | Builds the project, dockerfile, side loads then installs into the cluster |
| make delete | Deletes all currently installed helm releases  |


## Suggested development workflow

1. `<write some code> `

2. Test with the above steps.

3. `make`

4. View your changes in the cluster

5. If you like it, push the docker image with the kubeops image and use the helm chart to install it `cd helm && helm install . --generate-name`

## How does it work in a nutshell?

- Kubeops uses the go-client for Kubernetes and leverages the watch capability.

- The value of this project is wrapping those calls in interfaces and creating some utility functionality for cluster connection.

- More information on how to develop your own operator watchers/subscriptions can be found [here](DEVELOPMENT.md).
