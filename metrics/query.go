// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package metrics

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type QueryProfilesResponse struct {
	Page      int           `json:"page"`
	PageSize  int           `json:"page_size"`
	SessionId string        `json:"session_id"`
	Status    string        `json:"status"`
	Total     int           `json:"total"`
	Results   []interface{} `json:"results"`
}

func GetMixpanelUser(distinctId, mixpanelAPISecret, mixpanelPassword string) (*QueryProfilesResponse, error) {
	url := "https://eu.mixpanel.com/api/2.0/engage?project_id=2729960"

	payload := strings.NewReader(fmt.Sprintf("distinct_id=%v", distinctId))

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return nil, err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	req.Header.Add("authorization", fmt.Sprintf("Basic %v", base64.StdEncoding.EncodeToString([]byte(mixpanelAPISecret+":"+mixpanelPassword))))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var mixpanelUser QueryProfilesResponse
	err = json.Unmarshal(body, &mixpanelUser)
	if err != nil {
		return nil, err
	}

	return &mixpanelUser, nil
}
