# Setup Kubernetes Cluster on single EC2 instance

**What needs to be prepared?**

My ssh keys are located in the folder `~/.ssh/`. For this setup I generated a key pair in the AWS console `kubernetes-user`, so it is located under `~/.ssh/kubernetes-user.pem`.

* activate centos image on AWS Console
* setup ssh key pair for EC2 instance
* edit terraform.tfvars
* edit manifests/ci-config.yaml


**How to provision AWS?**
```
# provision EC2 instance
terraform init
echo yes | terraform apply

# wait for EC2 instance to finish initializing (about 5 minutes)
# copy /home/centos/kubeconfig_ip from EC2 instance to local machine
scp ... # see terraform output
export KUBECONFIG=(PWD)/kubeconfig_ip
kubectl get nodes
```

**How to configure the kubernetes cluster?**
```
# install ingress
kubectl create -f manifests/ingress-mandatory.yaml
kubectl create -f manifests/ingress-service-nodeport.yaml

# install hello world app
kubectl create -f manifests/hello.yaml
kubectl create -f manifests/hello-ingress.yaml
# => visit hello.hvt.zone

# aws credentials service (awsecr-cred)
kubectl create -f manifests/awsecr-cred-deployment.yaml

# install hvt ci
kubectl create -f manifests/workshop-ci-config.yaml
kubectl create -f manifests/workshop-ci-deployment.yaml

```


**How to cleanup all EC2 resources?**
```
echo yes | terraform destroy
```