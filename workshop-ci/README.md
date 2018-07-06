# Workshop continuous-integration tool

DESCRIPTION TODO


**What do I need to setup?**
* S3 bucket


dependency management: dep
* `dep init`
* `dep ensure`

**run on local machine**
* edit .env file
* `go build`
* ./workshop-ci
* `ngrok http 8080`

**build**
```
docker build . --tag hvt1/workshop-ci
docker push hvt1/workshop-ci
```

**improvements**

hard coded values: 
* s3 bucket
* privileged pipeline images (DRONE_ESCALATE)