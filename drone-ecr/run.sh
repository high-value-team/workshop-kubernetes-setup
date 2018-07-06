#!/bin/sh

ok=1

# environment variables
if [[ -z ${DRONE_REPO_OWNER} ]]; then echo DRONE_REPO_OWNER not set; ok=0; fi
if [[ -z ${DRONE_REPO_NAME} ]]; then echo DRONE_REPO_NAME not set; ok=0; fi
if [[ -z ${DRONE_COMMIT_SHA} ]]; then echo DRONE_COMMIT_SHA not set; ok=0; fi
if [[ -z ${AWS_DEFAULT_REGION} ]]; then echo AWS_DEFAULT_REGION not set; ok=0; fi
if [[ -z ${AWS_ACCESS_KEY_ID} ]]; then echo AWS_ACCESS_KEY_ID not set; ok=0; fi
if [[ -z ${AWS_SECRET_ACCESS_KEY} ]]; then echo AWS_SECRET_ACCESS_KEY not set; ok=0; fi

echo; echo checking Environment variables
if [[ $ok -eq 0 ]]; then
    printf 'Invalid parameters!\n'
    exit 1;
fi

echo; echo starting docker daemon
#dockerd > /dev/null 2>&1 &
dockerd &

echo; echo loging into ECR registry
$(aws ecr get-login --no-include-email)

echo; echo get ECR_REGISTRY
export ECR_REGISTRY=$(aws ecr get-login --no-include-email | cut -d ' ' -f7 | sed -e "s/^https:\/\///")
if [[ -z ${ECR_REGISTRY} ]]; then echo ECR_REGISTRY could not be retrieved; exit 1; fi
echo ECR_REGISTRY:$ECR_REGISTRY

echo; echo docker build
echo docker build . --tag $DRONE_REPO_OWNER/$DRONE_REPO_NAME:$DRONE_COMMIT_SHA
docker build . --tag $DRONE_REPO_OWNER/$DRONE_REPO_NAME:$DRONE_COMMIT_SHA

echo; echo create ECR repository
echo aws ecr create-repository --repository-name $DRONE_REPO_OWNER/$DRONE_REPO_NAME
aws ecr create-repository --repository-name $DRONE_REPO_OWNER/$DRONE_REPO_NAME

echo; echo tag docker image
echo docker tag $DRONE_REPO_OWNER/$DRONE_REPO_NAME:$DRONE_COMMIT_SHA $ECR_REGISTRY/$DRONE_REPO_OWNER/$DRONE_REPO_NAME:$DRONE_COMMIT_SHA
docker tag $DRONE_REPO_OWNER/$DRONE_REPO_NAME:$DRONE_COMMIT_SHA $ECR_REGISTRY/$DRONE_REPO_OWNER/$DRONE_REPO_NAME:$DRONE_COMMIT_SHA

echo; echo docker image to ECR
echo docker push $ECR_REGISTRY/$DRONE_REPO_OWNER/$DRONE_REPO_NAME:$DRONE_COMMIT_SHA
docker push $ECR_REGISTRY/$DRONE_REPO_OWNER/$DRONE_REPO_NAME:$DRONE_COMMIT_SHA
