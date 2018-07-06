# drone-envsubst

substitute environment variables in a file

**usage with drone.io**
```
pipeline:
  substitute:
	image: fnbk/drone-envsubst
	source: .hvt.zone/template.k8s-deployment.yaml
	destination: .hvt.zone/k8s-deployment.yaml
	secrets: [github_username, github_repository, ecr_repository_id, ecr_region]
```

**build**
```
docker build . --tag hvt1/drone-envsubst
docker push hvt1/drone-envsubst
```