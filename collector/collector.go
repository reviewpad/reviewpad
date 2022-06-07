// Copyright (C) 2019-2022 Explore.dev, Unipessoal Lda - All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package collector

import (
	"github.com/dukex/mixpanel"
)

type Collector interface {
	Collect(eventName string, properties *map[string]interface{}) error
}

type collector struct {
	Client mixpanel.Mixpanel
	Id     string
	Token  string
}

func NewCollector(token string, id string) Collector {
	return &collector{
		Client: mixpanel.New(token, ""),
		Id:     id,
		Token:  token,
	}
}

func (c *collector) Collect(eventName string, properties *map[string]interface{}) error {
	if c.Token == "" {
		return nil
	}
	return c.Client.Track(c.Id, eventName, &mixpanel.Event{
		Properties: *properties,
	})
}
