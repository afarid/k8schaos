package k8schaos

import "k8s.io/apimachinery/pkg/runtime"

type K8sChaos interface {
	GetRandomObject() (runtime.Object, error)
	DeleteObject(runtime.Object) error
	Run(stopCh <-chan bool)
}
