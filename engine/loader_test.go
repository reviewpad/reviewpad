// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package engine_test

import (
	"fmt"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/reviewpad/reviewpad/v3/engine"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestLoad(t *testing.T) {
	type httpMockResponder struct {
		url       string
		responder httpmock.Responder
	}

	tests := map[string]struct {
		reviewpadFilePath     string
		httpMockResponders    []httpMockResponder
		wantReviewpadFilePath string
		wantErr               string
	}{
		"when the file has a parsing error": {
			reviewpadFilePath: "../testdata/engine/loader/reviewpad_with_parse_error.yml",
			wantErr:           "yaml: unmarshal errors:\n  line 5: cannot unmarshal !!str `parse-e...` into engine.ReviewpadFile",
		},
		"when the file has an import of a non existing file": {
			reviewpadFilePath: "../testdata/engine/loader/reviewpad_with_import_of_nonexisting_file.yml",
			httpMockResponders: []httpMockResponder{
				{
					url:       "https://foo.bar/invalid-url",
					responder: httpmock.NewErrorResponder(fmt.Errorf("invalid import url")),
				},
			},
			wantErr: "Get \"https://foo.bar/invalid-url\": invalid import url",
		},
		"when the file has an import of a file with a parsing error": {
			reviewpadFilePath: "../testdata/engine/loader/reviewpad_with_import_file_with_parse_error.yml",
			httpMockResponders: []httpMockResponder{
				{
					url:       "https://foo.bar/reviewpad_with_parse_error.yml",
					responder: httpmock.NewBytesResponder(200, httpmock.File("../testdata/engine/loader/reviewpad_with_parse_error.yml").Bytes()),
				},
			},
			wantErr: "yaml: unmarshal errors:\n  line 5: cannot unmarshal !!str `parse-e...` into engine.ReviewpadFile",
		},
		"when the file imports other files": {
			reviewpadFilePath: "../testdata/engine/loader/reviewpad_with_chain_of_imports.yml",
			httpMockResponders: []httpMockResponder{
				{
					url:       "https://foo.bar/reviewpad_with_no_imports.yml",
					responder: httpmock.NewBytesResponder(200, httpmock.File("../testdata/engine/loader/reviewpad_with_no_imports.yml").Bytes()),
				},
				{
					url:       "https://foo.bar/reviewpad_with_one_import.yml",
					responder: httpmock.NewBytesResponder(200, httpmock.File("../testdata/engine/loader/reviewpad_with_one_import.yml").Bytes()),
				},
			},
			wantReviewpadFilePath: "../testdata/engine/loader/reviewpad_appended.yml",
		},
		"when the file has cyclic imports": {
			reviewpadFilePath: "../testdata/engine/loader/reviewpad_with_cyclic_dependency_a.yml",
			httpMockResponders: []httpMockResponder{
				{
					url:       "https://foo.bar/reviewpad_with_cyclic_dependency_b.yml",
					responder: httpmock.NewBytesResponder(200, httpmock.File("../testdata/engine/loader/reviewpad_with_cyclic_dependency_b.yml").Bytes()),
				},
				{
					url:       "https://foo.bar/reviewpad_with_cyclic_dependency_a.yml",
					responder: httpmock.NewBytesResponder(200, httpmock.File("../testdata/engine/loader/reviewpad_with_cyclic_dependency_a.yml").Bytes()),
				},
			},
			wantErr: "loader: cyclic dependency",
		},
		"when the file has no issues": {
			reviewpadFilePath:     "../testdata/engine/loader/reviewpad_with_no_imports.yml",
			wantReviewpadFilePath: "../testdata/engine/loader/reviewpad_with_no_imports.yml",
		},
		"when the file requires action transformation": {
			reviewpadFilePath:     "../testdata/engine/loader/reviewpad_before_action_transform.yml",
			wantReviewpadFilePath: "../testdata/engine/loader/reviewpad_after_action_transform.yml",
		},
		"when the file requires extra action transformation": {
			reviewpadFilePath:     "../testdata/engine/loader/reviewpad_before_extra_action_transform.yml",
			wantReviewpadFilePath: "../testdata/engine/loader/reviewpad_after_extra_action_transform.yml",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			reviewpadFileData := httpmock.File(test.reviewpadFilePath).Bytes()
			if test.httpMockResponders != nil {
				httpmock.Activate()
				defer httpmock.DeactivateAndReset()

				for _, httpMockResponder := range test.httpMockResponders {
					httpmock.RegisterResponder("GET", httpMockResponder.url, httpMockResponder.responder)
				}
			}

			var wantReviewpadFile *engine.ReviewpadFile
			isExpectedReviewpadFile := test.wantReviewpadFilePath != ""
			if isExpectedReviewpadFile {
				wantReviewpadFileData := httpmock.File(test.wantReviewpadFilePath).Bytes()
				wantReviewpadFile = &engine.ReviewpadFile{}
				yaml.Unmarshal(wantReviewpadFileData, &wantReviewpadFile)
				// At the end of loading all imports from the file, its imports are reset to []engine.PadImport{}.
				// However, the parsing of the wanted reviewpad file, sets the imports to []engine.PadImport(nil).
				wantReviewpadFile.Imports = []engine.PadImport{}
			}

			gotReviewpadFile, err := engine.Load(reviewpadFileData)

			if err != nil && err.Error() != test.wantErr {
				assert.FailNow(t, "Load() error = %v, wantErr %v", err, test.wantErr)
			}
			assert.Equal(t, wantReviewpadFile, gotReviewpadFile)
		})
	}
}
