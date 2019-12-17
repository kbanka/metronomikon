package kube

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"sync"
)

// Globals for client object
var clientOnce sync.Once
var client *kubernetes.Clientset

func GetClient() (*kubernetes.Clientset, error) {
	var outerErr error

	// Create the client only once
	clientOnce.Do(func() {
		// The in-cluster config will automatically find the API endpoint
		// and service account credentials
		config, err := rest.InClusterConfig()
		if err != nil {
			outerErr = err
			return
		}
		// Create the clientset
		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			outerErr = err
			return
		}
		client = clientset
	})

	if outerErr != nil {
		return nil, outerErr
	}

	return client, nil
}

func TestClientConnection() error {
	c, err := GetClient()
	if err != nil {
		return err
	}
	// Retrieve server version to test connectivity to API
	// This should work regardless of service account permissions
	if _, err = c.Discovery().ServerVersion(); err != nil {
		return err
	}
	return nil
}
