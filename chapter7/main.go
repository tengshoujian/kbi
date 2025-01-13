package main

import (
	"context"
	"fmt"
	"log"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		log.Fatal(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}
	mylabels := map[string]string{
		"app.kubernetes.io/component": "my-component",
		"app.kubernetes.io/name":      "a-name",
	}
	pod := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:   "example-pod1",
			Labels: mylabels,
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:  "nginx-container",
					Image: "nginx",
				},
			},
		},
	}

	result, err := clientset.CoreV1().Pods("default").Create(context.TODO(), pod, metav1.CreateOptions{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("create pod %q in namespace %q.\n", result.GetObjectMeta().GetName(), result.GetObjectMeta().GetNamespace())

}
