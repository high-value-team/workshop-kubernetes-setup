#!/bin/bash

#  prepare_deployment:
#	image: fnbk/drone-envsubst
#	source: .hvt.zone/template.k8s-deployment.yaml
#	destination: .hvt.zone/k8s-deployment.yaml
#	secrets: [ecr_repository_id, aws_default_region]

# Test
export AWS_DEFAULT_REGION=ecr_region
export ECR_REPOSITORY_ID=ecr_repository_id
export PLUGIN_SOURCE=template.deployment.yml
export PLUGIN_DESTINATION=deployment.yml
export DRONE_WORKSPACE=$(PWD)
export DRONE_REPO_OWNER=fnbk
export DRONE_REPO_NAME=hello

../run.sh

cat ${DRONE_WORKSPACE}/${PLUGIN_DESTINATION}
