// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package utils

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

type PrefixFormatter struct {
	logrus.TextFormatter
}

func (p *PrefixFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	if prefix, ok := entry.Data["prefix"].(string); ok {
		entry.Message = fmt.Sprintf("%s %s", prefix, entry.Message)
		delete(entry.Data, "prefix")
	}

	return p.TextFormatter.Format(entry)
}

func NewLogger(logLevel string) (*logrus.Entry, error) {
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		return nil, err
	}

	logger := logrus.New()
	logger.Formatter = &PrefixFormatter{
		TextFormatter: logrus.TextFormatter{
			ForceColors: true,
		},
	}
	logger.Level = level

	return logrus.NewEntry(logger), nil
}
