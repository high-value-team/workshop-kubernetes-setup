package main

import (
	"os"

	"github.com/high-value-team/workshop-kubernetes-setup/workshop-ci/interior_models"
	"github.com/urfave/cli"

	_ "github.com/joho/godotenv/autoload"
)

func NewCLIParams(version string) *interior_models.CLIParams {
	cliParams := &interior_models.CLIParams{}

	app := cli.NewApp()
	app.Version = version
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "AwsRegion",
			EnvVar: "AWS_REGION",
			Value:  "eu-central-1",
		},
		cli.StringFlag{
			Name:   "AwsAccessKey",
			EnvVar: "AWS_ACCESS_KEY_ID",
			Value:  "",
		},
		cli.StringFlag{
			Name:   "AwsSecretAccessKey",
			EnvVar: "AWS_SECRET_ACCESS_KEY",
			Value:  "",
		},
		cli.StringFlag{
			Name:   "EcrRepositoryID",
			EnvVar: "ECR_REPOSITORY_ID",
			Value:  "",
		},
		cli.StringFlag{
			Name:   "KubernetesServer",
			EnvVar: "KUBERNETES_SERVER",
			Value:  "",
		},
		cli.StringFlag{
			Name:   "KubernetesCertificateAuthorityData",
			EnvVar: "KUBERNETES_CERTIFICATE_AUTHORITY_DATA",
			Value:  "",
		},
		cli.StringFlag{
			Name:   "KubernetesClientCertificateData",
			EnvVar: "KUBERNETES_CLIENT_CERTIFICATE_DATA",
			Value:  "",
		},
		cli.StringFlag{
			Name:   "KubernetesClientKeyData",
			EnvVar: "KUBERNETES_CLIENT_KEY_DATA",
			Value:  "",
		},
	}

	app.Action = func(c *cli.Context) error {
		cliParams.AwsRegion = c.String("AwsRegion")
		cliParams.AwsAccessKey = c.String("AwsAccessKey")
		cliParams.AwsSecretAccessKey = c.String("AwsSecretAccessKey")
		cliParams.EcrRepositoryID = c.String("EcrRepositoryID")
		cliParams.KubernetesServer = c.String("KubernetesServer")
		cliParams.KubernetesCertificateAuthorityData = c.String("KubernetesCertificateAuthorityData")
		cliParams.KubernetesClientCertificateData = c.String("KubernetesClientCertificateData")
		cliParams.KubernetesClientKeyData = c.String("KubernetesClientKeyData")
		return nil
	}
	app.Run(os.Args)

	return cliParams
}
