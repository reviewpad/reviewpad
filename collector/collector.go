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
}

type collector struct {
	Client    mixpanel.Mixpanel
	Id        string
	Token     string
	RunnerEnv string
	// Allows to identify a unique source of events
	RunnerId string
	// Allows to identify the correct order of events
	Order int
	// Allows to identify the type of the event
	EventType string
	// Allows to identify the url
	Url string
}

func NewCollector(token, id, eventType, url, runner string) (Collector, error) {
	c := collector{
		Client:    mixpanel.New(token, ""),
		Id:        id,
		Token:     token,
		RunnerEnv: runner,
		RunnerId:  uuid.NewString(),
		Order:     0,
		EventType: eventType,
		Url:       url,
	}

	if token != "" {
		if err := c.Client.UpdateUser(c.Id, &mixpanel.Update{
			Operation: "$set",
			Properties: map[string]interface{}{
				"name": id,
			},
		}); err != nil {
			return nil, err
		}
	}

	return &c, nil
}

func (c *collector) Collect(eventName string, properties map[string]interface{}) error {
	if c.Token == "" {
		return nil
	}
	properties["runner"] = c.RunnerEnv
	properties["runnerId"] = c.RunnerId
	properties["order"] = c.Order
	properties["eventType"] = c.EventType
	properties["url"] = c.Url
	c.Order = c.Order + 1
	return c.Client.Track(c.Id, eventName, &mixpanel.Event{
		Properties: properties,
	})
}
