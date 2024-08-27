package main

import (
	"context"
	"fmt"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/tools/cache"
)

func main() {
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		// Try to use in-cluster config if the kubeconfig is not available
		config, err = rest.InClusterConfig()
		if err != nil {
			panic(err.Error())
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Successfully connected to the Kubernetes API")

	watchlist := cache.NewListWatchFromClient(clientset.CoreV1().RESTClient(), "pods", metav1.NamespaceAll, fields.OneTermEqualSelector("spec.nodeName", ""))

	_, controller := cache.NewInformer(
		watchlist,
		&corev1.Pod{},
		0, // Duration is set to 0 for now
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				// Handle unscheduled pods (e.g. scheduling logic here)
				pod := obj.(*corev1.Pod)

				// Placeholder for your scheduling algorithm
				// This should determine the best node for the pod
				nodeName, err := findBestNodeForPod(pod, clientset, "example-label-key", "example-label-value")
				if err != nil {
					fmt.Printf("Error finding best node for pod %s: %v\n", pod.Name, err)
					return
				}

				// Assuming findBestNodeForPod is a function that returns a nodeName
				if nodeName != "" {
					err := bindPodToNode(pod, nodeName, clientset)
					if err != nil {
						fmt.Printf("Error scheduling pod %s to node %s: %v\n", pod.Name, nodeName, err)
					}
				}
			},
		},
	)

	stop := make(chan struct{})
	defer close(stop)
	go controller.Run(stop)

	// Keep the main thread running
	select {}
}

func bindPodToNode(pod *corev1.Pod, nodeName string, clientset *kubernetes.Clientset) error {
	binding := &corev1.Binding{
		ObjectMeta: metav1.ObjectMeta{
			Name:		pod.Name,
			Namespace:	pod.Namespace,
		},
		Target: corev1.ObjectReference{
			Kind: "Node",
			Name: nodeName,
		},
	}

	return clientset.CoreV1().Pods(pod.Namespace).Bind(context.Background(), binding, metav1.CreateOptions{})
}

/*func schedulePodToNode(pod *corev1.Pod, nodeName string, clientset *kubernetes.Clientset) error {
	// Update the pod's nodeName to schedule it
	pod.Spec.NodeName = nodeName
	_, err := clientset.CoreV1().Pods(pod.Namespace).Update(context.Background(), pod, metav1.UpdateOptions{})
	return err
}*/

func findBestNodeForPod(pod *corev1.Pod, clientset *kubernetes.Clientset, labelKey, labelValue string) (string, error) {
	// List all nodes
	nodes, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return "", err
	}

	// Filter nodes based on the specified label
	var filteredNodes []corev1.Node
	for _, node := range nodes.Items {
		if value, ok := node.Labels[labelKey]; ok && value == labelValue {
			filteredNodes = append(filteredNodes, node)
		}
	}

	// If no nodes match the label, return an error
	if len(filteredNodes) == 0 {
		return "", fmt.Errorf("no nodes found with label %s=%s", labelKey, labelValue)
	}

	// Placeholder for selecting the best node from the filtered list
	// For simplicity, we'll just return the first node in the filtered list
	// You can implement more sophisticated logic here
	bestNode := filteredNodes[0].Name

	return bestNode, nil
}
