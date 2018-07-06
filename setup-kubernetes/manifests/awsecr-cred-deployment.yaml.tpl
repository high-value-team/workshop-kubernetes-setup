apiVersion: v1
kind: Secret
metadata:
  name: registry-creds-ecr
  namespace: kube-system
  labels:
    app: registry-creds
    cloud: ecr
data:
  AWS_ACCESS_KEY_ID: BASE64_ENCODED(AWS_ACCESS_KEY_ID)
  AWS_SECRET_ACCESS_KEY: BASE64_ENCODED(AWS_SECRET_ACCESS_KEY)
  aws-account: BASE64_ENCODED(aws-account-id)
  aws-region: BASE64_ENCODED(aws-region)
type: Opaque
---
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRole
metadata:
    name: registry-creds
rules:
- apiGroups:
  - ""
  resources:
  - namespaces
  verbs:
  - list
  - watch
- apiGroups:
  - ""
  resources:
  - secrets
  verbs:
  - create
  - get
  - update
- apiGroups:
  - ""
  resources:
  - serviceaccounts
  verbs:
  - get
  - update
---
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRoleBinding
metadata:
  name: registry-creds
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: registry-creds
subjects:
  - kind: ServiceAccount
    name: registry-creds
    namespace: kube-system
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: registry-creds
  namespace: kube-system
---
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  labels:
    version: v1.7
  name: registry-creds
  namespace: kube-system
spec:
  replicas: 1
  template:
    metadata:
      labels:
        name: registry-creds
        version: v1.7
    spec:
      serviceAccountName: registry-creds
      containers:
      - image: upmcenterprises/registry-creds:1.7
        name: registry-creds
        imagePullPolicy: Always
        env:
          - name: AWS_ACCESS_KEY_ID
            valueFrom:
              secretKeyRef:
                name: registry-creds-ecr
                key: AWS_ACCESS_KEY_ID
          - name: AWS_SECRET_ACCESS_KEY
            valueFrom:
              secretKeyRef:
                name: registry-creds-ecr
                key: AWS_SECRET_ACCESS_KEY
          - name: awsaccount
            valueFrom:
              secretKeyRef:
                name: registry-creds-ecr
                key: aws-account
          - name: awsregion
            valueFrom:
              secretKeyRef:
                name: registry-creds-ecr
                key: aws-region