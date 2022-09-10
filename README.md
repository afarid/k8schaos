# Kubernetes Chaos monkey

This is deployment to delete random pods from random target namespace or all namespaces


## Prerequisites
* Kind cluster
* Docker  
* Helm

### Configuration
By default, k9sChaos removes random pods from `default` namespace, If you want to change the target namespace or time interval,
Please edit `app.env` file before deployment
```env
NAMESPACE=<different-namespace>
TIME_PERIOD=<different-time-internal>
```

### Deployment
By default, The Makefile deploys k8sChaos in kube-system namespace in context `kind-kind`, If you want to change deployment 
namespace or cluster, please edit `deploy.env` before deployment. 
After tweaking the config, You just need below command to deploy the helm chart for k8sChaos. This will perform 3 actions
- Create docker image
- Load the image to kind cluster master node
- Deploy helm chart for k8schaos

```shell
make deploy
```

### TODO: 
 - [ ] Implement distributed locking so we can have more k8schaos pod
 - [ ] Adding prometheus metrics for the killed pods
 - [ ] Adding unit/integrations testing
 - [ ] Create ci/cd pipeline
 - [ ] Adding chaos money for more k8s objects