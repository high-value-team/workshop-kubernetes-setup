# drone-kubernetes

deploy to kubernetes cluster

**usage with drone.io**
```
pipeline:
    deploy:
        image: hvt1/drone-kubernetes:latest
        secrets: [ kubernetes_server, kubernetes_certificate_authority_data, kubernetes_client_certificate_data, kubernetes_client_key_data ]
        deployment: .hvt.zone/k8s-deployment.yml
```

**build**
```
docker build . --tag hvt1/drone-kubernetes
docker push hvt1/drone-kubernetes
```
