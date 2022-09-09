package main

import (
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	"k8schaos/pkg/log"
	"k8schaos/pkg/podchaos"
	"k8schaos/utils"
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

	informerFactory := informers.NewSharedInformerFactory(clientset, time.Second*30)

	podChaosSvc := podchaos.NewK8sChaos(clientset, informerFactory,
		podchaos.WithNamespace(config.Namespace),
		podchaos.WithTimePeriod(config.TimePeriod),
	)

	informerFactory.Start(wait.NeverStop)
	informerFactory.WaitForCacheSync(wait.NeverStop)

	stopCh := make(chan bool)
	podChaosSvc.Run(stopCh)

}

// TODO: Instrumentation
// TODO: High availability
