package kube

import (
	"fmt"
	"github.com/applauseoss/metronomikon/config"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Get a list of Kubernetes namespaces, taking into account the configured blacklist/whitelist
func GetNamespaces() ([]string, error) {
	ret := []string{}
	cfg := config.GetConfig()
	namespaces, err := client.CoreV1().Namespaces().List(metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("could not list namespaces: %s", err)
	}
	// Whitelist takes precedence over blacklist
	if len(cfg.Kubernetes.Namespace.Whitelist) > 0 {
		for _, namespace := range namespaces.Items {
			for _, tmp_namespace := range cfg.Kubernetes.Namespace.Whitelist {
				if tmp_namespace == namespace.ObjectMeta.Name {
					ret = append(ret, namespace.ObjectMeta.Name)
					break
				}
			}
		}
	} else {
		for _, namespace := range namespaces.Items {
			found := false
			for _, tmp_namespace := range cfg.Kubernetes.Namespace.Blacklist {
				if tmp_namespace == namespace.ObjectMeta.Name {
					found = true
					break
				}
			}
			if !found {
				ret = append(ret, namespace.ObjectMeta.Name)
			}
		}
	}
	return ret, nil
}
