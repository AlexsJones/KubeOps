package kubernetes

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"strings"
)

// GetKubeClient creates a Kubernetes config and client for a given kubeconfig context.
func GetKubeClient(context string) (*rest.Config, kubernetes.Interface, error) {

	var config *rest.Config
	if context == "" {
		c, err := rest.InClusterConfig()
		if err != nil {
			return nil, nil, err
		}
		config = c
	} else {
		cf, err := configForContext(context)
		if err != nil {
			return nil, nil, err
		}
		config = cf
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, nil, fmt.Errorf("could not get Kubernetes client: %s", err)
	}
	return config, client, nil
}

// configForContext creates a Kubernetes REST client configuration for a given kubeconfig context.
func configForContext(context string) (*rest.Config, error) {
	config, err := getConfig(context).ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("could not get Kubernetes config for context %q: %s", context, err)
	}
	return config, nil
}

func findKubeConfig() (string, error) {
	env := os.Getenv("KUBECONFIG")
	if env != "" {
		return env, nil
	}
	path, err := homedir.Expand("~/.kube/config")
	if err != nil {
		return "", err
	}
	return path, nil
}

// getConfig returns a Kubernetes client config for a given context.
func getConfig(c string) clientcmd.ClientConfig {
	rules := clientcmd.NewDefaultClientConfigLoadingRules()
	rules.DefaultClientConfig = &clientcmd.DefaultClientConfig

	overrides := &clientcmd.ConfigOverrides{ClusterDefaults: clientcmd.ClusterDefaults}

	kubeConfigPath, err := findKubeConfig()
	if err != nil {
		log.Fatal(err)
	}

	kubeConfig, err := clientcmd.LoadFromFile(kubeConfigPath)
	if err != nil {
		log.Fatal(err)
	}

	if c != "" {
		for mapContext, _ := range kubeConfig.Contexts {
			if strings.Contains(mapContext, c) {
				overrides.CurrentContext = mapContext
				continue
			}
		}
	}
	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(rules, overrides)
}
