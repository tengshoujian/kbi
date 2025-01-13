package main

import (
	"encoding/json"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
)

func createPod() error {
	pod := createPodObject()
	serializer := getJSONSerializer()

}

func createPodObject() *corev1.Pod {
	pod := corev1.Pod{
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "runtime",
					Image: "nginx",
				},
			},
		},
	}
	pod.SetName("my-pod")
	pod.SetLabels(map[string]string{
		"app.kubernetes.io/componet": "my-component",
		"app.kubernetes.io/name":     "a-name",
	})
	return &pod
}

func getJSONSerializer() runtime.Serializer {
	scheme := runtime.NewScheme()
	scheme.AddKnownTypes(
		schema.GroupVersion{
			Group:   "",
			Version: "v1",
		},
		&corev1.Pod{},
		&metav1.Status{},
	)
	return json.NewSerializeWithOptions(
		json.SimpleMetaFactory{},
		nil,
		scheme,
		json.SerializerOptions{},
	)

}
