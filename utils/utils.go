package utils

import (
	"fmt"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func GenerateRandomInt(min int, max int) int {
	return rand.Intn(max-min) + min
}

func GenerateRandomString(size int) string {
	chars := "abcdefghijklmnopqrstuvwxyz"
	result := ""
	for i := 0; i < size; i++ {
		index := GenerateRandomInt(0, len(chars))
		result = fmt.Sprintf("%s%s", result, string(chars[index]))

	}
	return result
}

func GetK8sClient() (*kubernetes.Clientset, error) {
	var kubeconfig string
	config := &rest.Config{}
	kubeconfig, ok := os.LookupEnv("KUBECONFIG")
	if !ok {
		kubeconfig = filepath.Join(homedir.HomeDir(), ".kube", "config")
	}

	_, err := os.Stat(kubeconfig)
	if err == nil {
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, err
		}
	} else {
		config, err = rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return clientset, nil
}
