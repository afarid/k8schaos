package podchaos

import (
	"context"
	"errors"
	"go.uber.org/zap"
	v1 "k8s.io/api/core/v1"
	v1meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/informers"
	v1informer "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8schaos/pkg/k8schaos"
	"k8schaos/pkg/log"
	"k8schaos/utils"
	"time"
)

type option func(*service)

func WithNamespace(namespace string) func(*service) {
	return func(s *service) {
		s.namespace = namespace
	}
}

func WithTimePeriod(timePeriod time.Duration) func(*service) {
	return func(s *service) {
		s.timePeriod = timePeriod
	}
}

type service struct {
	namespace   string
	timePeriod  time.Duration
	clientset   *kubernetes.Clientset
	podInformer v1informer.PodInformer
}

func (s *service) GetRandomObject() (runtime.Object, error) {
	pods, err := s.podInformer.Lister().Pods(s.namespace).List(labels.NewSelector())
	if err != nil {
		return nil, err
	}

	if len(pods) == 0 {
		return nil, ErrorNoMatch
	}

	randInt := utils.GenerateRandomInt(0, len(pods))

	return pods[randInt], nil
}

func (s *service) DeleteObject(object runtime.Object) error {
	podName := object.(*v1.Pod).Name
	err := s.clientset.CoreV1().Pods(s.namespace).Delete(context.Background(), podName, v1meta.DeleteOptions{})
	return err
}

func (s *service) Run(stopCh <-chan bool) {
	log.Logger.Info("starting pod chaos monkey", zap.String("time_interval", s.timePeriod.String()),
		zap.String("namespace", s.namespace))
	ticker := time.NewTicker(s.timePeriod)

loop:
	for {
		select {
		case <-ticker.C:
			object, err := s.GetRandomObject()
			if err != nil {
				log.Logger.Error(err.Error())
				continue
			}
			log.Logger.Info("deleting pod", zap.String("pod_name", object.(*v1.Pod).Name))
			err = s.DeleteObject(object)
			if err != nil {
				log.Logger.Error(err.Error())
				continue
			}
		case <-stopCh:
			break loop
		}
	}

}

func NewK8sChaos(clientset *kubernetes.Clientset, informerFactory informers.SharedInformerFactory, options ...option) k8schaos.K8sChaos {
	podChaos := &service{}
	for _, opt := range options {
		opt(podChaos)
	}

	if podChaos.timePeriod == 0 {
		podChaos.timePeriod = time.Minute * 5
	}

	podChaos.podInformer = informerFactory.Core().V1().Pods()
	podChaos.podInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    func(new interface{}) {},
		UpdateFunc: func(old, new interface{}) {},
		DeleteFunc: func(obj interface{}) {},
	})

	podChaos.clientset = clientset

	return podChaos

}

var ErrorNoMatch = errors.New("getting empty list for listing namespace")
