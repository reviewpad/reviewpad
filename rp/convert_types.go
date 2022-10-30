// Copyright (C) 2022 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package rp

import (
	"strconv"

	log "github.com/sirupsen/logrus"
)

// String returns a pointer to the string value passed in.
func String(v string) *string {
	return &v
}

// Int64 returns a pointer to the int64 value passed in.
func Int64(v int64) *int64 {
	return &v
}

// MustAtoi converts a string to an int, panicking if the conversion fails.
func MustAtoi(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		log.WithError(err).Panicf("failed to convert %q to int", s)
	}
	return i
}
