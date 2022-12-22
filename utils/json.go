// Copyright (C) 2022 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file

package utils

import (
	"bytes"
	"encoding/json"
	"io"
)

func MustRead(r io.Reader) string {
	b, err := io.ReadAll(r)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func MustWrite(w io.Writer, s string) {
	_, err := io.WriteString(w, s)
	if err != nil {
		panic(err)
	}
}

func MustWriteBytes(w io.Writer, data []byte) {
	_, err := w.Write(data)
	if err != nil {
		panic(err)
	}
}

func MustUnmarshal(data []byte, v interface{}) {
	err := json.Unmarshal(data, v)
	if err != nil {
		panic(err)
	}
}

func CompactJSONString(data string) (string, error) {
	dst := &bytes.Buffer{}
	err := json.Compact(dst, []byte(data))
	if err != nil {
		return "", err
	}

	return dst.String(), nil
}
