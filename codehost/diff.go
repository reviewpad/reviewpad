// Copyright (C) 2022 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file

package codehost

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/reviewpad/reviewpad/v4/utils"
)

type diff struct {
	blocks []*diffBlock
}

type diffSpan struct {
	Start int32
	End   int32
}

type diffBlock struct {
	IsContext bool
	Old       *diffSpan
	New       *diffSpan
	oldLine   string
	newLine   string
}

type chunkLinesInfo struct {
	oldLine, newLine int32
	numOld, numNew   int32
}

func parseFilePatch(patch string) ([]*diffBlock, error) {
	lines := strings.Split(patch, "\n")
	return parseFilePatchLines(lines)
}

func parseFilePatchLines(lines []string) ([]*diffBlock, error) {
	diff := &diff{
		blocks: []*diffBlock{},
	}

	currentLine := 0
	numAdded, numRemoved := int32(0), int32(0)

	for currentLine < len(lines) {
		line := lines[currentLine]
		currentLine++
		if !strings.HasPrefix(line, "@@") {
			continue
		}

		chunkLines, err := parseLines(line)
		if err != nil {
			return nil, fmt.Errorf("error in chunk lines parsing (%d): %v\npatch: %s", currentLine, err, strings.Join(lines, "\n"))
		}

		oldIndex := int32(0)
		newIndex := int32(0)

		for oldIndex < chunkLines.numOld || newIndex < chunkLines.numNew {
			if currentLine >= len(lines) {
				break
			}

			chunkLine := lines[currentLine]
			currentLine++
			runes := []rune(chunkLine)
			firstChar := string(runes[0])
			text := string(runes[1:])
			switch firstChar {
			case " ":
				appendUnmodifiedLine(diff, chunkLines, text)
				oldIndex++
				newIndex++
			case "+":
				numAdded++
				appendAddedLine(diff, chunkLines, text)
				newIndex++
			case "-":
				numRemoved++
				appendRemovedLine(diff, chunkLines, text)
				oldIndex++
			}
		}
	}

	return diff.blocks, nil
}

func parseLines(line string) (*chunkLinesInfo, error) {
	sections := strings.Split(line, " ")
	if len(sections) < 3 {
		return nil, fmt.Errorf("missing lines info: %s", line)
	}

	oldLine, numOld, err := parseVersionSection(sections[1])
	if err != nil {
		return nil, fmt.Errorf("error when parsing old section %s: %v", sections[1], err)
	}

	newLine, numNew, err := parseVersionSection(sections[2])
	if err != nil {
		return nil, fmt.Errorf("error when parsing new section %s: %v", sections[2], err)
	}

	return &chunkLinesInfo{
		oldLine: oldLine,
		newLine: newLine,
		numOld:  numOld,
		numNew:  numNew,
	}, nil
}

func parseVersionSection(section string) (int32, int32, error) {
	blocks := strings.Split(section, ",")
	line, err := strconv.ParseInt(blocks[0], 10, 32)
	if err != nil {
		return 0, 0, fmt.Errorf("wrong line format (%s): %v", blocks[0], err)
	}
	replyLine := utils.AbsInt32(int32(line))
	if len(blocks) < 2 {
		return replyLine, replyLine, nil
	}

	numLines, err := strconv.ParseInt(blocks[1], 10, 32)
	if err != nil {
		return 0, 0, fmt.Errorf("wrong num old lines format (%s): %v", blocks[1], err)
	}
	replyNumLines := int32(numLines)

	return replyLine, replyNumLines, nil
}

func appendUnmodifiedLine(diff *diff, lines *chunkLinesInfo, lineText string) {
	var block *diffBlock
	if len(diff.blocks) > 0 {
		lastBlock := diff.blocks[len(diff.blocks)-1]
		if lastBlock.IsContext && lastBlock.Old.End == lines.oldLine-1 {
			block = lastBlock
			block.newLine = block.newLine + "\n" + lineText
			block.oldLine = block.oldLine + "\n" + lineText
		}
	}

	if block == nil {
		block = &diffBlock{
			Old: &diffSpan{
				Start: lines.oldLine,
				End:   0,
			},
			New: &diffSpan{
				Start: lines.newLine,
				End:   0,
			},
			IsContext: true,
			oldLine:   lineText,
			newLine:   lineText,
		}
		diff.blocks = append(diff.blocks, block)
	}
	block.Old.End = lines.oldLine
	block.New.End = lines.newLine
	lines.newLine++
	lines.oldLine++
}

func appendAddedLine(diff *diff, lines *chunkLinesInfo, lineText string) {
	var block *diffBlock
	if len(diff.blocks) > 0 && !diff.blocks[len(diff.blocks)-1].IsContext {
		block = diff.blocks[len(diff.blocks)-1]
		if block.New == nil {
			block.New = &diffSpan{
				Start: lines.newLine,
				End:   0,
			}
			block.newLine = lineText
		} else {
			block.newLine = block.newLine + "\n" + lineText
		}
	} else {
		block = &diffBlock{
			New: &diffSpan{
				Start: lines.newLine,
				End:   0,
			},
			IsContext: false,
			newLine:   lineText,
		}
		diff.blocks = append(diff.blocks, block)
	}
	block.New.End = lines.newLine
	lines.newLine++
}

func appendRemovedLine(diff *diff, lines *chunkLinesInfo, lineText string) {
	var block *diffBlock
	if len(diff.blocks) > 0 && !diff.blocks[len(diff.blocks)-1].IsContext && diff.blocks[len(diff.blocks)-1].New == nil {
		block = diff.blocks[len(diff.blocks)-1]
		block.oldLine = block.oldLine + "\n" + lineText
	} else {
		block = &diffBlock{
			Old: &diffSpan{
				Start: lines.oldLine,
				End:   0,
			},
			IsContext: false,
			oldLine:   lineText,
		}
		diff.blocks = append(diff.blocks, block)
	}
	block.Old.End = lines.oldLine
	lines.oldLine++
}
