package main

import (
	"context"
	"fmt"
	"time"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/util/wait"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/retry"

	installerv1alpha1 "github.com/operator-framework/api/pkg/operators/v1alpha1"
)

func main() {
	// Load the Kubernetes configuration from the default location or provide your own path.
	config, err := rest.InClusterConfig()
	if err != nil {
		// Handle error
		panic(err.Error())
	}

	// Create a Kubernetes clientset.
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		// Handle error
		panic(err.Error())
	}

	// Namespace where the InstallerSets are located.
	namespace := "openshift-pipelines"

	// List InstallerSets in the specified namespace.
	installerSetList, err := clientset.OperatorsV1alpha1().InstallerSets(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		// Handle error
		panic(err.Error())
	}

	// Iterate over each InstallerSet to check its readiness.
	for _, installerSet := range installerSetList.Items {
		if isInstallerSetReady(installerSet) {
			fmt.Printf("InstallerSet %s is ready\n", installerSet.Name)
		} else {
			fmt.Printf("InstallerSet %s is not ready\n", installerSet.Name)
		}
	}
}

func isInstallerSetReady(installerSet installerv1alpha1.InstallerSet) bool {
	// Check the InstallerSet conditions to determine if it's ready.
	for _, condition := range installerSet.Status.Conditions {
		if condition.Type == installerv1alpha1.InstallerSetConditionTypeReady && condition.Status == installerv1alpha1.ConditionStatusTrue {
			return true
		}
	}
	return false
}
