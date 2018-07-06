# drone-ecr

push image to ECR registry

**usage with drone.io**
```
pipeline:
  push:
	image: fnbk/drone-ecr
	secrets: [github_username, github_repository, aws_default_region, aws_access_key_id, aws_secret_access_key]
```

**build**
```
docker build . --tag hvt1/drone-ecr
docker push hvt1/drone-ecr
```

**run**
```
docker run --rm \
    -e DRONE_REPO_OWNER=fnbk \
    -e DRONE_REPO_NAME=hello \
    -e DRONE_COMMIT_SHA=48bc6acaaea144b068a307b14cbdd19768861a08 \
    -e AWS_DEFAULT_REGION=eu-central-1 \
    -e AWS_ACCESS_KEY_ID=xxx \
    -e AWS_SECRET_ACCESS_KEY=xxx \
    -v (pwd):(pwd) \
    -w (pwd) \
    --privileged \
    hvt1/drone-ecr
```