package interior_interactions

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/cncd/pipeline/pipeline"
	"github.com/cncd/pipeline/pipeline/backend"
	"github.com/cncd/pipeline/pipeline/backend/docker"
	"github.com/cncd/pipeline/pipeline/frontend"
	"github.com/cncd/pipeline/pipeline/frontend/yaml"
	"github.com/cncd/pipeline/pipeline/frontend/yaml/compiler"
	"github.com/cncd/pipeline/pipeline/interrupt"
	"github.com/cncd/pipeline/pipeline/multipart"
	"github.com/high-value-team/workshop-kubernetes-setup/workshop-ci/interior_models"
)

type Interactions struct {
	env *interior_models.CLIParams
}

func NewInteractions(env *interior_models.CLIParams) *Interactions {
	return &Interactions{env}
}

type LogAppender struct {
	Log string
}

func NewLogAppender() *LogAppender {
	return &LogAppender{Log: ""}
}

func (a *LogAppender) Append(line string) {
	a.Log = fmt.Sprintf("%s%s", a.Log, line)
}

func handleException(statusPath string) {
	r := recover()
	if r != nil {
		switch r.(type) {
		case interior_models.SadException:
			WriteToS3Private("dwx2018", statusPath, "failure")
		case interior_models.SuprisingException:
			WriteToS3Private("dwx2018", statusPath, "failure")
		default:
			if err, ok := r.(error); !ok {
				log.Println(err)
			} else {
				log.Println(err)
			}
		}
	}
}

func (i *Interactions) TriggerPipeline(githubUsername, githubRepository, commitSHA, pipeline string, createdAt time.Time) {

	logPath := fmt.Sprintf("%s/%s/logs/%s", githubUsername, githubRepository, commitSHA)
	pipelinePath := fmt.Sprintf("%s/%s/pipelines/%s.json", githubUsername, githubRepository, commitSHA)
	statusPath := fmt.Sprintf("%s/%s/status/%s", githubUsername, githubRepository, commitSHA)
	timestampPath := fmt.Sprintf("%s/%s/timestamps/%s", githubUsername, githubRepository, commitSHA)

	defer handleException(statusPath)

	logAppender := NewLogAppender()
	writeToLog := func(part string) {
		writeToStdOut(part)
		logAppender.Append(part)
		WriteToS3Public("dwx2018", logPath, logAppender.Log)
	}
	writeStatusSuccess := func() { WriteToS3Private("dwx2018", statusPath, "success") }
	writeStatusFailure := func() { WriteToS3Private("dwx2018", statusPath, "failure") }
	writeStatusPending := func() { WriteToS3Private("dwx2018", statusPath, "pending") }
	writeTimestamp := func() { WriteToS3Private("dwx2018", timestampPath, createdAt.String()) }
	writePipeline := func(outputCompile string) { WriteToS3Private("dwx2018", pipelinePath, outputCompile) }
	writeStartLogHeader := func() {
		writeToLog(fmt.Sprintf("\n\n#\n# New PushEvent\n#\n\nUsername: %s\nRepository: %s\nSHA: %s\nTimestamp: %s\n\n", githubUsername, githubRepository, commitSHA, createdAt))
	}

	writeTimestamp()
	writeStatusPending()
	writeStartLogHeader()
	environmentVariables := i.buildEnvironmentVariables(githubUsername, githubRepository)
	compiledPipelineString := compile(pipeline, githubUsername, githubRepository, commitSHA, environmentVariables)
	writePipeline(compiledPipelineString)
	execute(compiledPipelineString, writeToLog, writeStatusSuccess, writeStatusFailure)
}

func (i *Interactions) buildEnvironmentVariables(githubUsername, githubRepository string) map[string]string {
	return map[string]string{
		"KUBERNETES_SERVER":                     i.env.KubernetesServer,
		"KUBERNETES_CERTIFICATE_AUTHORITY_DATA": i.env.KubernetesCertificateAuthorityData,
		"KUBERNETES_CLIENT_CERTIFICATE_DATA":    i.env.KubernetesClientCertificateData,
		"KUBERNETES_CLIENT_KEY_DATA":            i.env.KubernetesClientKeyData,
		"AWS_ACCESS_KEY_ID":                     i.env.AwsAccessKey,
		"AWS_SECRET_ACCESS_KEY":                 i.env.AwsSecretAccessKey,
		"AWS_DEFAULT_REGION":                    i.env.AwsRegion,
		"ECR_REPOSITORY_ID":                     i.env.EcrRepositoryID,
	}
}

func compile(inputPipeline, githubUsername, githubRepository, commitSHA string, environmentVariables map[string]string) string {
	pipelinePrefix := fmt.Sprintf("%d_%s", time.Now().Unix(), commitSHA)

	conf, err := yaml.ParseString(inputPipeline)
	if err != nil {
		panic(interior_models.SuprisingException{Err: err})
	}

	cmp := compiler.New(
		// DRONE_ESCALATE
		compiler.WithEscalated(
			"plugins/docker",
			"plugins/gcr",
			"plugins/ecr",
			"hvt1/drone-ecr", // TODO clean up hard coded whitelist for privileged pipeline steps
		),
		compiler.WithSecret([]compiler.Secret{}...),
		compiler.WithEnviron(environmentVariables),
		compiler.WithVolumes([]string{}...),
		compiler.WithWorkspace(
			"/pipeline",
			"src",
		),
		compiler.WithPrefix(pipelinePrefix),
		compiler.WithLocal(false),
		compiler.WithMetadata(frontend.Metadata{
			ID: "",
			Repo: frontend.Repo{
				Name:    fmt.Sprintf("%s/%s", githubUsername, githubRepository),                        /// "fnbk/hello",                        // CI_REPO_NAME
				Link:    fmt.Sprintf("https://github.com/%s/%s", githubUsername, githubRepository),     // "https://github.com/fnbk/hello",     // CI_REPO_LINK
				Remote:  fmt.Sprintf("https://github.com/%s/%s.git", githubUsername, githubRepository), // "https://github.com/fnbk/hello.git", // CI_REPO_REMOTE
				Private: false,                                                                         // CI_REPO_PRIVATE
				Secrets: []frontend.Secret{},
				Branch:  "",
			},
			Curr: frontend.Build{
				Number:   0, // CI_BUILD_NUMBER
				Created:  0, // CI_BUILD_CREATED
				Started:  0, // CI_BUILD_STARTED
				Finished: 0,
				Timeout:  0,
				Status:   "",
				Event:    "", // false CI_BUILD_EVENT
				Link:     "", // CI_BUILD_LINK
				Target:   "",
				Trusted:  false,
				Commit: frontend.Commit{
					Sha:     commitSHA,           // CI_COMMIT_SHA
					Ref:     "refs/heads/master", // CI_COMMIT_REF
					Refspec: "",                  // CI_COMMIT_REFSPEC
					Branch:  "master",            // CI_COMMIT_BRANCH
					Message: "",                  // CI_COMMIT_MESSAGE
					Author: frontend.Author{
						Name:   "", // CI_COMMIT_AUTHOR_NAME
						Email:  "",
						Avatar: "",
					},
				},
				Parent: 0,
			},
			Prev: frontend.Build{
				Number:   0,
				Created:  0,
				Started:  0,
				Finished: 0,
				Timeout:  0,
				Status:   "",
				Event:    "",
				Link:     "",
				Target:   "",
				Trusted:  false,
				Commit: frontend.Commit{
					Sha:     "",
					Ref:     "",
					Refspec: "",
					Branch:  "",
					Message: "",
					Author: frontend.Author{
						Name:   "",
						Email:  "",
						Avatar: "",
					},
				},
				Parent: 0,
			},
			Job: frontend.Job{
				Number: 0,
				Matrix: map[string]string{},
			},
			Sys: frontend.System{
				Name:    "",
				Host:    "",
				Link:    "",
				Arch:    "linux/amd64", // CI_SYSTEM_ARCH
				Version: "",
			},
		}),
	)

	compiled := cmp.Compile(conf)

	// marshal the compiled spec to formatted yaml
	out, err := json.MarshalIndent(compiled, "", "  ")
	if err != nil {
		panic(interior_models.SuprisingException{Err: err})
	}

	// fmt.Fprintf(os.Stdout, "Successfully compiled\n")

	return string(out)
}

func execute(compileOutput string, onWriteToLog func(string), onSuccess, onFailure func()) {
	timeout := time.Hour // docker deamon timeout

	config, err := pipeline.ParseString(compileOutput)
	if err != nil {
		panic(interior_models.SuprisingException{Err: err})
	}

	engine, err := docker.NewEnv()
	if err != nil {
		panic(interior_models.SuprisingException{Err: err})
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	ctx = interrupt.WithContext(ctx)

	buildStatusOK := true
	var logger = pipeline.LogFunc(func(proc *backend.Step, rc multipart.Reader) error {
		part, err := rc.NextPart()
		if err != nil {
			return err
		}

		var buf [512]byte
		for {
			_, err := part.Read(buf[:])

			onWriteToLog(string(buf[:]))

			buf = [512]byte{}

			if err == io.EOF {
				break
			}
			if err == io.ErrClosedPipe {
				break
			}
			if err != nil {
				return err
			}
		}

		return nil
	})

	var tracer = pipeline.TraceFunc(func(state *pipeline.State) error {
		if state.Process.Exited {
			line := fmt.Sprintf("proc %q exited with status %d\n", state.Pipeline.Step.Name, state.Process.ExitCode)
			onWriteToLog(line)
		} else {
			line := fmt.Sprintf("\n\n#\n# %q\n#\n\n", state.Pipeline.Step.Name)
			onWriteToLog(line)
			state.Pipeline.Step.Environment["CI_BUILD_STATUS"] = "success"
			state.Pipeline.Step.Environment["CI_BUILD_FINISHED"] = strconv.FormatInt(time.Now().Unix(), 10)
			if state.Pipeline.Error != nil {
				buildStatusOK = false
				state.Pipeline.Step.Environment["CI_BUILD_STATUS"] = "failure"
			}
		}
		return nil
	})

	runtime := pipeline.New(config,
		pipeline.WithContext(ctx),
		pipeline.WithLogger(logger),
		pipeline.WithTracer(tracer),
		pipeline.WithEngine(engine),
	)
	err = runtime.Run()
	if err != nil {
		onFailure()
		if _, ok := err.(*pipeline.ExitError); ok {
			fmt.Printf("ExitError:%s\n", err)
		} else {
			fmt.Printf("Panic:%s\n", err)
			// panic(err)
		}
	}

	onSuccess()
}

func writeToFile(filename, outputExecute string) {
	var writer = os.Stdout
	writer, err := os.Create(filename)
	if err != nil {
		panic(interior_models.SuprisingException{Err: err})
	}
	defer writer.Close()

	_, err = writer.WriteString(outputExecute)
	if err != nil {
		panic(interior_models.SuprisingException{Err: err})
	}
}

func writeToStdOut(outputExecute string) {
	_, err := fmt.Printf(outputExecute)
	if err != nil {
		panic(interior_models.SuprisingException{Err: err})
	}
}

func WriteToS3Private(bucket, path, content string) {
	writeToS3(false, bucket, path, content)
}

func WriteToS3Public(bucket, path, content string) {
	writeToS3(true, bucket, path, content)
}

func writeToS3(public bool, bucket, path, content string) {
	sess := session.Must(session.NewSession())
	svc := s3.New(sess, &aws.Config{Region: aws.String("eu-central-1")})

	acl := aws.String("private")
	if public {
		acl = aws.String("public-read")
	}

	params := &s3.PutObjectInput{
		Bucket:             aws.String(bucket), // Required
		Key:                aws.String(path),   // Required
		ACL:                acl,
		Body:               bytes.NewReader([]byte(content)),
		ContentType:        aws.String("text/plain"),
		ContentDisposition: aws.String("inline"),
	}
	_, err := svc.PutObject(params)
	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}
}
