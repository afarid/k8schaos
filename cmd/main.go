package main

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	v1meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	"k8schaos/pkg/log"
	"k8schaos/pkg/podchaos"
	"k8schaos/utils"
	"net/http"
	"os"
	"time"
)

func main() {
	// Load configuration from config file and environment variables.
	configPath := os.Getenv("CONFIG_PATH")
	config, err := utils.LoadConfig(configPath)
	if err != nil {
		log.Logger.Fatal("Unable to read config file", zap.Error(err))

	}

	// Get Kubernetes rest client
	clientset, err := utils.GetK8sClient()
	if err != nil {
		log.Logger.Fatal("Unable to get k8s client", zap.Error(err))
	}

	informerFactory := informers.NewSharedInformerFactoryWithOptions(clientset, time.Second*30, informers.WithTweakListOptions(func(options *v1meta.ListOptions) {
		options.LabelSelector = "chaos-enabled=true"
	}))

	podChaosSvc := podchaos.NewK8sChaos(clientset, informerFactory,
		podchaos.WithNamespace(config.Namespace),
		podchaos.WithTimePeriod(config.TimePeriod),
	)

	informerFactory.Start(wait.NeverStop)
	informerFactory.WaitForCacheSync(wait.NeverStop)

	stopCh := make(chan bool)
	go podChaosSvc.Run(stopCh)

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":2112", nil)
	}()

	<-stopCh

}

// TODO: Instrumentation
// TODO: High availability
