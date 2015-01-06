package sprintly

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestItems_List(t *testing.T) {
	client, server, mux := setup()
	defer server.Close()

	mux.HandleFunc("/products/1/items.json", func(w http.ResponseWriter, r *http.Request) {
		ensureMethod(t, r, "GET")
		fmt.Fprintf(w, `
		[
			{
				"status": "backlog",
				"product": {
					"archived": false,
					"id": 1,
					"name": "sprint.ly"
				},
				"progress": {
					"accepted_at": "2013-06-14T22:52:07+00:00",
					"closed_at": "2013-06-14T21:53:43+00:00",
					"started_at": "2013-06-14T21:50:36+00:00"
				},
				"description": "Require people to estimate the score of an item before they can start working on it.",
				"tags": [
					"scoring",
					"backlog"
				],
				"number": 188,
				"archived": false,
				"title": "Don't let un-scored items out of the backlog.",
				"created_by": {
					"first_name": "Joe",
					"last_name": "Stump",
					"id": 1,
					"email": "joe@joestump.net"
				},
				"score": "M",
				"assigned_to": {
					"first_name": "Joe",
					"last_name": "Stump",
					"id": 1,
					"email": "joe@joestump.net"
				},
				"type": "task"
			}
		]`)
	})

	items, _, err := client.Items.List(1, nil)
	if err != nil {
		t.Errorf("Items.List failed: %v", err)
		return
	}

	joe := &User{
		Id:        1,
		Email:     "joe@joestump.net",
		FirstName: "Joe",
		LastName:  "Stump",
	}

	layout := "2006-01-02T15:04:05-07:00"
	acceptedAt, err := time.Parse(layout, "2013-06-14T22:52:07+00:00")
	if err != nil {
		t.Error(err)
		return
	}
	closedAt, err := time.Parse(layout, "2013-06-14T21:53:43+00:00")
	if err != nil {
		t.Error(err)
		return
	}
	startedAt, err := time.Parse(layout, "2013-06-14T21:50:36+00:00")
	if err != nil {
		t.Error(err)
		return
	}

	want := []Item{
		{
			Number:      188,
			Title:       "Don't let un-scored items out of the backlog.",
			Description: "Require people to estimate the score of an item before they can start working on it.",
			Score:       "M",
			Status:      ItemStatusBacklog,
			Tags: []string{
				"scoring",
				"backlog",
			},
			Product: &Product{
				Archived: false,
				Id:       1,
				Name:     "sprint.ly",
			},
			Progress: &Progress{
				StartedAt:  &startedAt,
				AcceptedAt: &acceptedAt,
				ClosedAt:   &closedAt,
			},
			CreatedBy:  joe,
			AssignedTo: joe,
			Archived:   false,
			Type:       "task",
		},
	}

	ensureEqual(t, items, want)
}
