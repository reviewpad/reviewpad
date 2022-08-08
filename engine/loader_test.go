// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package engine_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/reviewpad/reviewpad/v3/engine"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

type httpMockResponder struct {
	url       string
	responder httpmock.Responder
}

func TestLoad(t *testing.T) {
	tests := map[string]struct {
		httpMockResponders     []httpMockResponder
		inputReviewpadFilePath string
		wantReviewpadFilePath  string
		wantErr                string
	}{
		"when the file has a parsing error": {
			inputReviewpadFilePath: "testdata/loader/reviewpad_with_parse_error.yml",
			wantErr:                "yaml: unmarshal errors:\n  line 5: cannot unmarshal !!str `parse-e...` into engine.ReviewpadFile",
		},
		"when the file imports a nonexistent file": {
			inputReviewpadFilePath: "testdata/loader/reviewpad_with_import_of_nonexistent_file.yml",
			httpMockResponders: []httpMockResponder{
				{
					url:       "https://foo.bar/nonexistent_file",
					responder: httpmock.NewErrorResponder(fmt.Errorf("file doesn't exist")),
				},
			},
			wantErr: "Get \"https://foo.bar/nonexistent_file\": file doesn't exist",
		},
		"when the file imports a file that has a parsing error": {
			inputReviewpadFilePath: "testdata/loader/reviewpad_with_import_file_with_parse_error.yml",
			httpMockResponders: []httpMockResponder{
				{
					url:       "https://foo.bar/reviewpad_with_parse_error.yml",
					responder: httpmock.NewBytesResponder(200, httpmock.File("testdata/loader/reviewpad_with_parse_error.yml").Bytes()),
				},
			},
			wantErr: "yaml: unmarshal errors:\n  line 5: cannot unmarshal !!str `parse-e...` into engine.ReviewpadFile",
		},
		"when the file has cyclic imports": {
			inputReviewpadFilePath: "testdata/loader/reviewpad_with_cyclic_dependency_a.yml",
			httpMockResponders: []httpMockResponder{
				{
					url:       "https://foo.bar/reviewpad_with_cyclic_dependency_b.yml",
					responder: httpmock.NewBytesResponder(200, httpmock.File("testdata/loader/reviewpad_with_cyclic_dependency_b.yml").Bytes()),
				},
				{
					url:       "https://foo.bar/reviewpad_with_cyclic_dependency_a.yml",
					responder: httpmock.NewBytesResponder(200, httpmock.File("testdata/loader/reviewpad_with_cyclic_dependency_a.yml").Bytes()),
				},
			},
			wantErr: "loader: cyclic dependency",
		},
		"when the file has import chains": {
			inputReviewpadFilePath: "testdata/loader/reviewpad_with_imports_chain.yml",
			httpMockResponders: []httpMockResponder{
				{
					url:       "https://foo.bar/reviewpad_with_no_imports.yml",
					responder: httpmock.NewBytesResponder(200, httpmock.File("testdata/loader/reviewpad_with_no_imports.yml").Bytes()),
				},
				{
					url:       "https://foo.bar/reviewpad_with_one_import.yml",
					responder: httpmock.NewBytesResponder(200, httpmock.File("testdata/loader/reviewpad_with_one_import.yml").Bytes()),
				},
			},
			wantReviewpadFilePath: "testdata/loader/reviewpad_appended.yml",
		},
		"when the file has no issues": {
			inputReviewpadFilePath: "testdata/loader/reviewpad_with_no_imports.yml",
			wantReviewpadFilePath:  "testdata/loader/reviewpad_with_no_imports.yml",
		},
		"when the file requires action transformation": {
			inputReviewpadFilePath: "testdata/loader/transform/reviewpad_before_action_transform.yml",
			wantReviewpadFilePath:  "testdata/loader/transform/reviewpad_after_action_transform.yml",
		},
		"when the file requires extra action transformation": {
			inputReviewpadFilePath: "testdata/loader/transform/reviewpad_before_extra_action_transform.yml",
			wantReviewpadFilePath:  "testdata/loader/transform/reviewpad_after_extra_action_transform.yml",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			if test.httpMockResponders != nil {
				httpmock.Activate()
				defer httpmock.DeactivateAndReset()

				registerHttpResponders(test.httpMockResponders)
			}

			var wantReviewpadFile *engine.ReviewpadFile
			if test.wantReviewpadFilePath != "" {
				wantReviewpadFileData, err := loadReviewpadFile(test.wantReviewpadFilePath)
				if err != nil {
					assert.FailNow(t, "Error reading reviewpad file: %v", err)
				}

				wantReviewpadFile, err = parseReviewpadFile(wantReviewpadFileData)
				if err != nil {
					assert.FailNow(t, "Error parsing reviewpad file: %v", err)
				}
			}

			reviewpadFileData, err := loadReviewpadFile(test.inputReviewpadFilePath)
			if err != nil {
				assert.FailNow(t, "Error reading reviewpad file: %v", err)
			}

			gotReviewpadFile, gotErr := engine.Load(reviewpadFileData)

			if gotErr != nil && gotErr.Error() != test.wantErr {
				assert.FailNow(t, "Load() error = %v, wantErr %v", gotErr, test.wantErr)
			}
			assert.Equal(t, wantReviewpadFile, gotReviewpadFile)
		})
	}
}

func registerHttpResponders(httpMockResponders []httpMockResponder) {
	for _, httpMockResponder := range httpMockResponders {
		httpmock.RegisterResponder("GET", httpMockResponder.url, httpMockResponder.responder)
	}
}

func loadReviewpadFile(filepath string) ([]byte, error) {
	reviewpadFileData, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	return reviewpadFileData, nil
}

func parseReviewpadFile(data []byte) (*engine.ReviewpadFile, error) {
	reviewpadFile := &engine.ReviewpadFile{}
	err := yaml.Unmarshal(data, &reviewpadFile)
	if err != nil {
		return nil, err
	}

	// At the end of loading all imports from the file, its imports are reset to []engine.PadImport{}.
	// However, the parsing of the wanted reviewpad file, sets the imports to []engine.PadImport(nil).
	reviewpadFile.Imports = []engine.PadImport{}

	return reviewpadFile, nil
}
