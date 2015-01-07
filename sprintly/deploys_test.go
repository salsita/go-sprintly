package sprintly

import (
	"fmt"
	"net/http"
	"testing"
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

var testingDeployString = `
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
	testingDeploySlice       = []Deploy{testingDeploy}
	testingDeploySliceString = fmt.Sprintf("[%v]", testingDeployString)
)

func TestDeploys_List(t *testing.T) {
	client, server, mux := setup()
	defer server.Close()

	mux.HandleFunc("/products/1/deploys.json", func(w http.ResponseWriter, r *http.Request) {
		ensureMethod(t, r, "GET")
		fmt.Fprint(w, testingDeploySliceString)
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

		var got DeployCreateArgs
		if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
			t.Error(err)
			return
		}

		ensureEqual(t, &got, &args)
		fmt.Fprint(w, testingTaskString)
	})

	item, _, err := client.Deploys.Create(testingProduct.Id, &args)
	if err != nil {
		t.Errorf("Deploys.Create failed: %v", err)
		return
	}

	ensureEqual(t, item, &testingTask)
}
