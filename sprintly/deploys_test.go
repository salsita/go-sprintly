package sprintly

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gorilla/schema"
)

var testingDeploy = Deploy{
	Environment: "staging",
	Items: []Item{
		{
			Number: 188,
			Title:  "Who knows ...",
		},
	},
}

var testingDeployJson = `
{
	"environment": "staging",
	"items": [
		{
			"number": 188,
			"title": "Who knows ..."
		}
	]
}
`

var (
	testingDeploySlice     = []Deploy{testingDeploy}
	testingDeploySliceJson = fmt.Sprintf("[%v]", testingDeployString)
)

func TestDeploys_List(t *testing.T) {
	client, server, mux := setup()
	defer server.Close()

	mux.HandleFunc("/products/1/deploys.json", func(w http.ResponseWriter, r *http.Request) {
		ensureMethod(t, r, "GET")
		fmt.Fprint(w, testingDeploySliceJson)
	})

	deploys, _, err := client.Deploys.List(1, nil)
	if err != nil {
		t.Errorf("Deploys.List failed: %v", err)
		return
	}

	ensureEqual(t, deploys, testingDeploySlice)
}

func TestDeploys_Create(t *testing.T) {
	client, server, mux := setup()
	defer server.Close()

	args := DeployCreateArgs{
		Environment: "staging",
		ItemNumbers: []int{1, 2, 3, 4, 5},
	}

	mux.HandleFunc("/products/1/deploys.json", func(w http.ResponseWriter, r *http.Request) {
		ensureMethod(t, r, "POST")

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Error(err)
			return
		}

		values, err := url.ParseQuery(body.String())
		if err != nil {
			t.Error(err)
			return
		}

		var receivedArgs DeployCreateArgs
		if err := schema.NewDecoder().Decode(&receivedArgs, values); err != nil {
			t.Error(err)
			return
		}

		ensureEqual(t, &receivedArgs, &args)
		fmt.Fprint(w, testingDeployJson)
	})

	deploy, _, err := client.Deploys.Create(1, &args)
	if err != nil {
		t.Errorf("Deploys.Create failed: %v", err)
		return
	}

	ensureEqual(t, &deploy, &testingDeploy)
}
