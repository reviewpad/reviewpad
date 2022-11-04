// Copyright (C) 2022 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package rp

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

func String(v string) *string {
	return &v
}

func Int(v int) *int {
	return &v
}

func JsonRawMessage(v json.RawMessage) *json.RawMessage {
	return &v
}

func MustAtoi(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		log.WithError(err).Panicf("failed to convert %q to int", s)
	}
	return i
}

func MustAtoi64(s string) int64 {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		log.WithError(err).Panicf("failed to convert %q to int", s)
	}
	return i
}

func SafeMarshal(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(b)
}

func StringifyMap(m map[string]string) string {
	var s strings.Builder
	for key, value := range m {
		s.WriteString(fmt.Sprintf("\"%v\":\"%v\"", key, value))
	}
	return s.String()
}
