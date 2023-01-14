// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package collector

import (
	"github.com/dukex/mixpanel"
	"github.com/google/uuid"
)

type Collector interface {
	Collect(eventName string, properties map[string]interface{}) error
	CollectError(err error) error
}

type collector struct {
	Client mixpanel.Mixpanel
	// Unique identifier that is connected to every event.
	// It is used to identify the user.
	DistinctId string
	// Token to authenticate the requests and identify the project.
	Token string
	// Runner environment where the event is being collected.
	// (i.e. github app, github action, playground)
	RunnerEnv string
	// Unique identifier for the events' source.
	RunnerId string
	// Order of the event in the events' source.
	Order int
	// Type of the event.
	// For more details see https://docs.github.com/en/actions/using-workflows/events-that-trigger-workflows
	EventType string
	// Optional properties that are added to every event.
	Optional *map[string]string
}

// NewCollector creates a collector instance.
// If the mixpanelToken is empty, the collector will not send any events.
// The distinctId identifies the user.
// The eventType identifies the type of the event.
// The runnerName identifies the runner environment where the event is being collected (e.g. github app, github action, playground).
// The options are optional properties that are added to every event.
func NewCollector(mixpanelToken, distinctId, eventType, runnerName string, options *map[string]string) (Collector, error) {
	c := collector{
		Client:     mixpanel.New(mixpanelToken, ""),
		DistinctId: distinctId,
		Token:      mixpanelToken,
		RunnerEnv:  runnerName,
		RunnerId:   uuid.NewString(),
		Order:      0,
		EventType:  eventType,
		Optional:   options,
	}

	if mixpanelToken != "" {
		// Define the user within Mixpanel
		err := c.Client.UpdateUser(c.DistinctId, &mixpanel.Update{
			Operation: "$set",
			Properties: map[string]interface{}{
				"name": distinctId,
			},
		})

		if err != nil {
			return nil, err
		}
	}

	return &c, nil
}

// Collect sends an event to mixpanel.
func (c *collector) Collect(eventName string, properties map[string]interface{}) error {
	if c.Token == "" {
		return nil
	}

	properties["runner"] = c.RunnerEnv
	properties["runnerId"] = c.RunnerId
	properties["order"] = c.Order
	properties["eventType"] = c.EventType

	if c.Optional != nil {
		for key, value := range *c.Optional {
			properties[key] = value
		}
	}

	c.Order = c.Order + 1

	return c.Client.Track(c.DistinctId, eventName, &mixpanel.Event{
		Properties: properties,
	})
}

// CollectError sends an error event to mixpanel.
func (c *collector) CollectError(err error) error {
	return c.Collect("Error", map[string]interface{}{
		"details": err.Error(),
	})
}
