// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package utils

import (
	"github.com/sirupsen/logrus"
)

// NewLogger creates a new logger with the given log level
// The log levels cab be: trace, debug, info, warning, error, fatal and panic.
// The logger is configured to output JSON logs
func NewLogger(logLevel logrus.Level) *logrus.Entry {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetReportCaller(true)
	logger.SetLevel(logLevel)

	return logrus.NewEntry(logger)
}
