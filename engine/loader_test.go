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
		httpMockResponders                                     []httpMockResponder
		inputReviewpadFilePath, wantReviewpadFilePath, wantErr string
	}{
		// Fail
		"when the file has a parsing error": {
			inputReviewpadFilePath: "../testdata/engine/loader/reviewpad_with_parse_error.yml",
			wantErr:                "yaml: unmarshal errors:\n  line 5: cannot unmarshal !!str `parse-e...` into engine.ReviewpadFile",
		},
		"when the file has an import of a nonexisting file": {
			inputReviewpadFilePath: "../testdata/engine/loader/reviewpad_with_import_of_nonexisting_file.yml",
			httpMockResponders: []httpMockResponder{
				{
					url:       "https://foo.bar/nonexistent_file",
					responder: httpmock.NewErrorResponder(fmt.Errorf("file doesn't exist")),
				},
			},
			wantErr: "Get \"https://foo.bar/nonexistent_file\": file doesn't exist",
		},
		"when the file imports a file that has a parsing error": {
			inputReviewpadFilePath: "../testdata/engine/loader/reviewpad_with_import_file_with_parse_error.yml",
			httpMockResponders: []httpMockResponder{
				{
					url:       "https://foo.bar/reviewpad_with_parse_error.yml",
					responder: httpmock.NewBytesResponder(200, httpmock.File("../testdata/engine/loader/reviewpad_with_parse_error.yml").Bytes()),
				},
			},
			wantErr: "yaml: unmarshal errors:\n  line 5: cannot unmarshal !!str `parse-e...` into engine.ReviewpadFile",
		},
		"when the file has cyclic imports": {
			inputReviewpadFilePath: "../testdata/engine/loader/reviewpad_with_cyclic_dependency_a.yml",
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
		// Pass
		"when the file imports other files": {
			inputReviewpadFilePath: "../testdata/engine/loader/reviewpad_with_imports_chain.yml",
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
		"when the file has no issues": {
			inputReviewpadFilePath: "../testdata/engine/loader/reviewpad_with_no_imports.yml",
			wantReviewpadFilePath:  "../testdata/engine/loader/reviewpad_with_no_imports.yml",
		},
		"when the file requires action transformation": {
			inputReviewpadFilePath: "../testdata/engine/loader/reviewpad_before_action_transform.yml",
			wantReviewpadFilePath:  "../testdata/engine/loader/reviewpad_after_action_transform.yml",
		},
		"when the file requires extra action transformation": {
			inputReviewpadFilePath: "../testdata/engine/loader/reviewpad_before_extra_action_transform.yml",
			wantReviewpadFilePath:  "../testdata/engine/loader/reviewpad_after_extra_action_transform.yml",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			reviewpadFileData := httpmock.File(test.inputReviewpadFilePath).Bytes()
			if test.httpMockResponders != nil {
				httpmock.Activate()
				defer httpmock.DeactivateAndReset()

				for _, httpMockResponder := range test.httpMockResponders {
					httpmock.RegisterResponder("GET", httpMockResponder.url, httpMockResponder.responder)
				}
			}

			var wantReviewpadFile *engine.ReviewpadFile

			isReviewpadFileExpected := test.wantReviewpadFilePath != ""

			if isReviewpadFileExpected {
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
