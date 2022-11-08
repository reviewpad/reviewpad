// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file

package utils_test

import (
	"testing"
	"time"

	"github.com/reviewpad/reviewpad/v3/utils"
	"github.com/stretchr/testify/assert"
)

func TestFormatTimeDiff(t *testing.T) {
	tests := map[string]struct {
		name       string
		firstTime  time.Time
		secondTime time.Time
		result     string
	}{
		"1 second": {
			firstTime:  time.Date(2020, 1, 1, 1, 1, 1, 1, time.UTC),
			secondTime: time.Date(2020, 1, 1, 1, 1, 2, 1, time.UTC),
			result:     "1 second",
		},
		"1 minute": {
			firstTime:  time.Date(2020, 1, 1, 1, 1, 1, 1, time.UTC),
			secondTime: time.Date(2020, 1, 1, 1, 2, 1, 1, time.UTC),
			result:     "1 minute",
		},
		"1 hour": {
			firstTime:  time.Date(2020, 1, 1, 1, 1, 1, 1, time.UTC),
			secondTime: time.Date(2020, 1, 1, 2, 1, 1, 1, time.UTC),
			result:     "1 hour",
		},
		"1 day": {
			firstTime:  time.Date(2020, 1, 1, 1, 1, 1, 1, time.UTC),
			secondTime: time.Date(2020, 1, 2, 1, 1, 1, 1, time.UTC),
			result:     "1 day",
		},
		"1 month": {
			firstTime:  time.Date(2020, 1, 1, 1, 1, 1, 1, time.UTC),
			secondTime: time.Date(2020, 2, 1, 1, 1, 1, 1, time.UTC),
			result:     "1 month",
		},
		"1 year": {
			firstTime:  time.Date(2020, 1, 1, 1, 1, 1, 1, time.UTC),
			secondTime: time.Date(2021, 1, 1, 1, 1, 1, 1, time.UTC),
			result:     "1 year",
		},
		"2 years": {
			firstTime:  time.Date(2020, 1, 1, 1, 1, 1, 1, time.UTC),
			secondTime: time.Date(2022, 1, 1, 1, 1, 1, 1, time.UTC),
			result:     "2 years",
		},
		"1 year, and 1 month": {
			firstTime:  time.Date(2020, 1, 1, 1, 1, 1, 1, time.UTC),
			secondTime: time.Date(2021, 2, 1, 1, 1, 1, 1, time.UTC),
			result:     "1 year and 1 month",
		},
		"1 year, and 2 months": {
			firstTime:  time.Date(2020, 1, 1, 1, 1, 1, 1, time.UTC),
			secondTime: time.Date(2021, 3, 1, 1, 1, 1, 1, time.UTC),
			result:     "1 year and 2 months",
		},
		"2 years, 4 months, and 10 days": {
			firstTime:  time.Date(2022, 5, 11, 1, 1, 1, 1, time.UTC),
			secondTime: time.Date(2020, 1, 1, 1, 1, 1, 1, time.UTC),
			result:     "2 years, 4 months and 10 days",
		},
		"1 hour, 10 minutes, and 5 seconds": {
			firstTime:  time.Date(2020, 1, 1, 2, 11, 6, 1, time.UTC),
			secondTime: time.Date(2020, 1, 1, 1, 1, 1, 1, time.UTC),
			result:     "1 hour, 10 minutes and 5 seconds",
		},
		"1 minute, and 5 seconds": {
			firstTime:  time.Date(2020, 1, 1, 1, 1, 1, 1, time.UTC),
			secondTime: time.Date(2020, 1, 1, 1, 2, 6, 1, time.UTC),
			result:     "1 minute and 5 seconds",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := utils.ReadableTimeDiff(test.firstTime, test.secondTime)
			assert.Equal(t, test.result, result)
		})
	}
}
