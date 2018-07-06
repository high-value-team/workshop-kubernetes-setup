apiVersion: v1
kind: ConfigMap
metadata:
  name: workshop-ci-config
  namespace: default
data:
  AWS_REGION: "eu-central-1"
  AWS_ACCESS_KEY_ID: "TODO"
  AWS_SECRET_ACCESS_KEY: "TODO"
  ECR_REPOSITORY_ID: "690729310209"
  KUBERNETES_SERVER: "see kubeconfig_ip file"
  KUBERNETES_CERTIFICATE_AUTHORITY_DATA: "see kubeconfig_ip file"
  KUBERNETES_CLIENT_CERTIFICATE_DATA: "see kubeconfig_ip file"
  KUBERNETES_CLIENT_KEY_DATA: "see kubeconfig_ip file"
