##Kubernetes Chaos monkey

This is deployment to delete random pods from random target namespace or all namespaces


####Prerequisites
* Local kind cluster
* Docker 
* Helm

####Configuration
By default, k9sChaos removes random pods from default namespace, If you want to change the target name space or time interval,
Please edit `app.env` file before deployment
```env
NAMESPACE=<different-namespace>
TIME_PERIOD=<different-time-internal>
```

####Deployment
By default, The Makefile deploys k8sChaos in kube-system namespace in context `kind-kind`, If you want to change deployment 
namespace or cluster, please edit `deploy.env` before deployment. 
After tweaking the config, You just need below command to deploy the helm chart for k8sChaos
```shell
make deploy
```