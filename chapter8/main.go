package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
)

func createPod() error {
	pod := createPodObject()
	serializer := getJSONSerializer()
	postBoday, err := serializerPodObject(serializer,pod)
	if err != nil{
		return err
	}
	reqCreate, err := buildPostRequest(postBoday)
	if err != nil{
		return err
	}
	client := &http.client{}
	resp, err := client.Do(reqCreate)
	if err != nil{
		retunr err
	}
	def resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil{
		return err
	}
	if resp.StatusCode <300{
		createPod, err := deserializePodBody(serialize, body)
		if err != nil{
			return err
		}
		json, err = json.MarshalIndent(createPod,""," ")
		if err != nil{
			return err
		}
		fmt.Printf("%s\n",json)
	}else{
		status. err :=
	}
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
func serializePodObject(serializer runtime.Serializer, pod *corev1.Pod,)(io.Reader,error){
	var buf bytes.Buffer
	err := serializer.Encode(pod,&buf)
	if err != nil {
		return nil,err
	}
	return &buf, nil
}

func buildPostRequest(body io.Reader)(*http.Request,error,){
	reqCreate, err :=http.NewRequest(
		"POST",
		"http://172.0.0.1:8001/api/v1/namespaces/default/pods",
		body,
	)
	if err != nil {
		return nil, err
	}
	reqCreate.Header.Add(
		"Accept",
		"application/json",
	)
	reqCreate.Header.Add(
		"Content-Type",
		"application/json",
	)
	return reqCreate,nil
}
func deserializePodBody(
	serilizer runtime.Serializer,
	body []byte,
)(*)
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
