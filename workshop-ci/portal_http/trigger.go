package portal_http

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/high-value-team/workshop-kubernetes-setup/workshop-ci/interior_interactions"
	"github.com/high-value-team/workshop-kubernetes-setup/workshop-ci/interior_models"

	"github.com/google/go-github/github"
)

type TriggerPipeline struct {
	Interactions *interior_interactions.Interactions
}

func NewTriggerPipelineHandler(interactions *interior_interactions.Interactions) http.HandlerFunc {
	triggerPipeline := TriggerPipeline{Interactions: interactions}
	return triggerPipeline.Handle
}

func (handler *TriggerPipeline) Handle(writer http.ResponseWriter, reader *http.Request) {
	onPing := handlePingEvent
	onPush := handler.handlePushEvent
	handleWebHookEvent(reader, onPing, onPush)
}

func handleWebHookEvent(reader *http.Request, onPing func(*http.Request), onPush func(*http.Request)) {
	switch github.WebHookType(reader) {
	case "ping":
		onPing(reader)
	case "push":
		onPush(reader)
	default:
		panic(interior_models.SadException{Err: fmt.Errorf("Webhook Type not supported")})
	}
}

func (handler *TriggerPipeline) handlePushEvent(reader *http.Request) {
	githubName, githubRepository, commitSHA, createdAt := parsePushEventBody(reader)
	pipeline := getPipeline(githubName, githubRepository, commitSHA)
	go handler.Interactions.TriggerPipeline(githubName, githubRepository, commitSHA, pipeline, createdAt)
}

func handlePingEvent(reader *http.Request) {}

func parsePushEventBody(reader *http.Request) (string, string, string, time.Time) {

	payload, err := getPayload(reader)
	if err != nil {
		log.Printf("ValidatePayload:%s", err)
		panic(interior_models.SuprisingException{Err: err})
	}

	event, err := github.ParseWebHook(github.WebHookType(reader), payload)
	if err != nil {
		log.Printf("ParseWebHook:%s", err)
		panic(interior_models.SuprisingException{Err: err})
	}

	pushEvent := event.(*github.PushEvent)

	// fmt.Printf("\n\n#\n# New PushEvent\n#\n\nUsername:%s\nRepository:%s\nSHA:%s\nTimestamp:%s\n\n", *pushEvent.Repo.Owner.Name, *pushEvent.Repo.Name, *pushEvent.HeadCommit.ID, pushEvent.Repo.CreatedAt.Time)
	return *pushEvent.Repo.Owner.Name, *pushEvent.Repo.Name, *pushEvent.HeadCommit.ID, pushEvent.Repo.CreatedAt.Time
}

func getPipeline(githubName, githubRepository, commitSHA string) string {
	url := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s/.hvt.zone/drone-pipeline.yml", githubName, githubRepository, commitSHA)
	response, err := http.Get(url)
	if err != nil {
		panic(interior_models.SuprisingException{Err: err})
	}
	if response.StatusCode != http.StatusOK {
		panic(interior_models.SadException{Err: fmt.Errorf("Response Code:%d Status:%", response.StatusCode, response.Status)})
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(interior_models.SuprisingException{Err: err})
	}
	return string(contents)
}

func getPayload(r *http.Request) (payload []byte, err error) {
	var body []byte // Raw body that GitHub uses to calculate the signature.

	switch ct := r.Header.Get("Content-Type"); ct {
	case "application/json":
		var err error
		if body, err = ioutil.ReadAll(r.Body); err != nil {
			return nil, err
		}

		// If the content type is application/json,
		// the JSON payload is just the original body.
		payload = body

	case "application/x-www-form-urlencoded":
		// payloadFormParam is the name of the form parameter that the JSON payload
		// will be in if a webhook has its content type set to application/x-www-form-urlencoded.
		const payloadFormParam = "payload"

		var err error
		if body, err = ioutil.ReadAll(r.Body); err != nil {
			return nil, err
		}

		// If the content type is application/x-www-form-urlencoded,
		// the JSON payload will be under the "payload" form param.
		form, err := url.ParseQuery(string(body))
		if err != nil {
			return nil, err
		}
		payload = []byte(form.Get(payloadFormParam))

	default:
		return nil, fmt.Errorf("Webhook request has unsupported Content-Type %q", ct)
	}

	// sig := r.Header.Get(signatureHeader)
	// if err := validateSignature(sig, body, secretKey); err != nil {
	// 	return nil, err
	// }
	return payload, nil
}
