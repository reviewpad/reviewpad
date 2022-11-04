// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file

package utils

import (
	"fmt"
	"strings"
	"time"

	"github.com/icza/gox/timex"
)

func pluralize(part int, label string) string {
	if part > 1 {
		return label + "s"
	}

	return label
}

func FormatTimeDiff(x time.Time, y time.Time) string {
	years, months, days, hours, minutes, seconds := timex.Diff(x, y)
	parts := make([]string, 0, 6)

	if years > 0 {
		parts = append(parts, fmt.Sprintf("%d %s", years, pluralize(years, "year")))
	}

	if months > 0 {
		parts = append(parts, fmt.Sprintf("%d %s", months, pluralize(months, "month")))
	}

	if days > 0 {
		parts = append(parts, fmt.Sprintf("%d %s", days, pluralize(days, "day")))
	}

	if hours > 0 {
		parts = append(parts, fmt.Sprintf("%d %s", hours, pluralize(hours, "hour")))
	}

	if minutes > 0 {
		parts = append(parts, fmt.Sprintf("%d %s", minutes, pluralize(minutes, "minute")))
	}

	if seconds > 0 {
		parts = append(parts, fmt.Sprintf("%d %s", seconds, pluralize(seconds, "second")))
	}

	l := len(parts)

	if l == 1 {
		return parts[0]
	} else if l == 2 {
		return parts[0] + " and " + parts[1]
	}

	return strings.Join(parts[:l-1], ", ") + " and " + parts[l-1]
}
