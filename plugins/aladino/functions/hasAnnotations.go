// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/explore-dev/atlas-common/go/api/entities"
	protobuf "github.com/explore-dev/atlas-common/go/api/services"
	"github.com/google/go-github/v45/github"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	"github.com/reviewpad/reviewpad/v3/utils"
)

func HasAnnotation() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type: aladino.BuildFunctionType([]aladino.Type{aladino.BuildStringType()}, aladino.BuildBoolType()),
		Code: patchHasAnnotationCode,
	}
}

func patchHasAnnotationCode(e aladino.Env, args []aladino.Value) (aladino.Value, error) {
	annotation := args[0].(*aladino.StringValue).Val

	symbolsByFileName, err := getSymbolsFromPatch(e)
	if err != nil {
		return nil, err
	}

	for _, fileSymbols := range symbolsByFileName {
		for _, symbol := range fileSymbols.Symbols {
			for _, symbolComment := range symbol.CodeComments {
				comment := symbolComment.Code
				if commentHasAnnotation(comment, annotation) {
					return aladino.BuildTrueValue(), nil
				}
			}
		}
	}

	return aladino.BuildFalseValue(), nil
}

func commentHasAnnotation(comment, annotation string) bool {
	normalizedComment := strings.ToLower(strings.Trim(comment, " "))
	normalizedAnnotation := strings.ToLower(annotation)
	anPrefix := "reviewpad-an: "
	if strings.HasPrefix(normalizedComment, anPrefix) {
		rest := normalizedComment[len(anPrefix)-1:]
		annotations := strings.Split(rest, " ")
		for _, annot := range annotations {
			if normalizedAnnotation == annot {
				return true
			}
		}
	}

	return false
}

// TODO: get files of patch with one request
func getSymbolsFromPatch(e aladino.Env) (map[string]*entities.Symbols, error) {
	res := make(map[string]*entities.Symbols)

	headOwner := utils.GetPullRequestHeadOwnerName(e.GetPullRequest())
	headRepo := utils.GetPullRequestHeadRepoName(e.GetPullRequest())

	baseOwner := utils.GetPullRequestBaseOwnerName(e.GetPullRequest())
	baseRepo := utils.GetPullRequestBaseRepoName(e.GetPullRequest())

	url := e.GetPullRequest().GetHead().GetRepo().GetURL()
	patch := e.GetPatch()

	lastCommit := e.GetPullRequest().Head.SHA

	for fp, commitFile := range patch {
		ioReader, _, err := e.GetClient().Repositories.DownloadContents(e.GetCtx(), headOwner, headRepo, fp, &github.RepositoryContentGetOptions{
			Ref: *e.GetPullRequest().Head.Ref,
		})
		if err != nil {
			// If fails to download the file from head, then tries to download it from base.
			// This happens when a file has been removed.
			ioReader, _, err = e.GetClient().Repositories.DownloadContents(e.GetCtx(), baseOwner, baseRepo, fp, &github.RepositoryContentGetOptions{
				Ref: *e.GetPullRequest().Base.Ref,
			})

			if err != nil {
				return nil, err
			}
		}

		// ReadFrom doesn't convert escaped unicode characters
		// Therefore we convert the bytes array to string, remove the escape and convert again to bytes

		buf := new(bytes.Buffer)
		buf.ReadFrom(ioReader)
		blob := buf.Bytes()

		str, _ := strconv.Unquote(strings.Replace(strconv.Quote(string(blob)), `\\u`, `\u`, -1))

		blob = []byte(str)

		blocks, err := getBlocks(commitFile)
		if err != nil {
			return nil, err
		}

		reply, err := e.GetSemanticClient().GetSymbols(e.GetCtx(), &protobuf.GetSymbolsRequest{
			Uri:      url,
			CommitId: *lastCommit,
			Filepath: fp,
			Blob:     blob,
			BlobId:   *commitFile.Repr.SHA,
			Diff:     &entities.ResolveFileDiff{Blocks: blocks},
		})

		if err != nil {
			return nil, err
		}

		res[fp] = reply.Symbols
	}

	return res, nil
}

func getBlocks(commitFile *aladino.File) ([]*entities.ResolveBlockDiff, error) {
	res := make([]*entities.ResolveBlockDiff, len(commitFile.Diff))
	for i, block := range commitFile.Diff {
		if block.New != nil {
			res[i] = &entities.ResolveBlockDiff{
				Start: block.New.Start,
				End:   block.New.End,
			}
		} else {
			if block.Old != nil {
				res[i] = &entities.ResolveBlockDiff{
					Start: block.Old.Start,
					End:   block.Old.End,
				}
			} else {
				return nil, fmt.Errorf("getBlocks fatal: block %+v is not valid", block)
			}
		}
	}

	return res, nil
}
