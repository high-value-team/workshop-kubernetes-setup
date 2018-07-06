package main

import (
	"fmt"

	"github.com/high-value-team/workshop-kubernetes-setup/workshop-ci/interior_interactions"
	"github.com/high-value-team/workshop-kubernetes-setup/workshop-ci/interior_models"
	"github.com/high-value-team/workshop-kubernetes-setup/workshop-ci/portal_http"
)

// wird durch build.sh gesetzt
var VersionNumber string = ""

func main() {
	// config
	cliParams := NewCLIParams(VersionNumber)

	// // consturct
	interactions := interior_interactions.NewInteractions(cliParams)
	httpPortal := portal_http.NewHTTPPortal(interactions, VersionNumber)

	// run
	httpPortal.Run(8080)
}

func handleException() {
	r := recover()
	if r != nil {
		switch ex := r.(type) {
		case interior_models.SadException:
			fmt.Println(ex.Message())
		case interior_models.SuprisingException:
			fmt.Println(ex.Message())
		default:
			if err, ok := r.(error); !ok {
				fmt.Printf("unknown error:%s\n", err)
			} else {
				fmt.Printf("unknown error:%s\n", r)
			}
		}
	}
}
