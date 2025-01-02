package main

import (
	"code-generator-demo/pkg/generated/clientset/versioned/typed/appcontroller/v1alpha1"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		panic(err.Error())
	}

	// 在 pkg/generated/clientset/versioned/typed/appcontroller/v1alpha1/appcontroller_client.go 中
	// 自动生成的 AppcontrollerV1alpha1Client，用于操作这种GVR
	appClient, err := v1alpha1.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// 获取default命名空间下的所有Application
	appList, err := appClient.Applications("default").List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	for _, app := range appList.Items {
		fmt.Println(app.Name)
	}
}
