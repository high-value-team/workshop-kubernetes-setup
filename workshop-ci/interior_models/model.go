package interior_models

type CLIParams struct {
	AwsRegion                          string
	AwsAccessKey                       string
	AwsSecretAccessKey                 string
	EcrRepositoryID                    string
	KubernetesServer                   string
	KubernetesCertificateAuthorityData string
	KubernetesClientCertificateData    string
	KubernetesClientKeyData            string
}
