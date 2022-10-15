// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package aladino

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
)

type AladinoLex struct {
	input string
	ast   Expr
}

const EOF = 0

type tokenDef struct {
	regex *regexp.Regexp
	kind  string
	token int
}

var tokens = []tokenDef{
	{
		// Allowed formats:
		// YYYYMMDD - e.g. 20220405
		// YYYY-MM-DD - e.g. 2022-04-05
		// YYYYMMDDTHH:MM:SS - e.g. 20220405T22:01:50
		// YYYY-MM-DDTHH:MM:SS - e.g. 2022-04-05T22:01:50
		// Recommended format:
		// RFC3339 (https://pkg.go.dev/time#pkg-constants)
		regex: regexp.MustCompile(`^\d{4}-?\d{2}-?\d{2}(T\d{2}:\d{2}:\d{2})?`),
		kind:  "timestamp",
		token: TIMESTAMP,
	},
	{
		// Examples:
		// 15 days ago
		// 3 months ago
		// 8 hours ago
		regex: regexp.MustCompile(`^[0-9]+\s(year(s?)|month(s?)|day(s?)|week(s?)|hour(s?)|minute(s?))\sago`),
		kind:  "relativeTimestamp",
		token: RELATIVETIMESTAMP,
	},
	{
		regex: regexp.MustCompile(`^[0-9]*\.?[0-9]+([eE][-+]?[0-9]+)?`),
		kind:  "number",
		token: NUMBER,
	},
	{
		regex: regexp.MustCompile(`^true`),
		kind:  "bool",
		token: TRUE,
	},
	{
		regex: regexp.MustCompile(`^false`),
		kind:  "bool",
		token: FALSE,
	},
	{
		regex: regexp.MustCompile(`^:\s?String*`),
		kind:  "stringType",
		token: TK_STRING_TYPE,
	},
	{
		regex: regexp.MustCompile(`^:\s?Int*`),
		kind:  "intType",
		token: TK_INT_TYPE,
	},
	{
		regex: regexp.MustCompile(`^:\s?Bool*`),
		kind:  "boolType",
		token: TK_BOOL_TYPE,
	},
	{
		regex: regexp.MustCompile(`^:\s?\[String\]*`),
		kind:  "stringArrayType",
		token: TK_STRING_ARRAY_TYPE,
	},
	{
		regex: regexp.MustCompile(`^:\s?\[Int\]*`),
		kind:  "intArrayType",
		token: TK_INT_ARRAY_TYPE,
	},
	{
		regex: regexp.MustCompile(`^:\s?\[Int\]*`),
		kind:  "boolArrayType",
		token: TK_BOOL_ARRAY_TYPE,
	},
	{
		regex: regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9]*`),
		kind:  "identifier",
		token: IDENTIFIER,
	},
	{
		regex: regexp.MustCompile(`^"([^"\\]|\\[\s\S])*"`),
		kind:  "stringLiteral",
		token: STRINGLITERAL,
	},
	{
		regex: regexp.MustCompile(`^(>|<)=?`),
		kind:  "binop",
		token: TK_CMPOP,
	},
	{
		regex: regexp.MustCompile(`^==`),
		kind:  "binop",
		token: TK_EQ,
	},
	{
		regex: regexp.MustCompile(`^!=`),
		kind:  "binop",
		token: TK_NEQ,
	},
	{
		regex: regexp.MustCompile(`^!`),
		kind:  "binop",
		token: TK_NOT,
	},
	{
		regex: regexp.MustCompile(`^&&`),
		kind:  "binop",
		token: TK_AND,
	},
	{
		regex: regexp.MustCompile(`^\|\|`),
		kind:  "binop",
		token: TK_OR,
	},
	{
		regex: regexp.MustCompile(`^=>`),
		kind:  "lambda",
		token: TK_LAMBDA,
	},
}

func (l *AladinoLex) Lex(lval *AladinoSymType) int {
	// fmt.Printf("lex: input: %v\n", l.input)
	// Skip spaces.
	for ; len(l.input) > 0 && isSpace(l.input[0]); l.input = l.input[1:] {
	}

	// Check if the input has ended.
	if len(l.input) == 0 {
		return EOF
	}

	// Check if one of the regular expressions matches.
	for _, tokDef := range tokens {
		str := tokDef.regex.FindString(l.input)

		if str == "" {
			continue
		}

		switch tokDef.kind {
		case "number":
			num, err := strconv.Atoi(str)
			if err != nil {
				fmt.Printf("error in atoi %v\n", err)
				log.Fatal(err)
			}
			lval.int = num
		case "stringLiteral":
			// Pass string content to the parser.
			lval.str = str[1 : len(str)-1]
		case "stringType", "intType", "boolType", "stringArrayType", "intArrayType", "boolArrayType":
			lval.str = strings.ReplaceAll(str, " ", "")[1:]
		default:
			lval.str = str
		}

		l.input = l.input[len(str):]
		return tokDef.token
	}

	// Otherwise return the next letter.
	ret := int(l.input[0])
	l.input = l.input[1:]
	return ret
}

func (l *AladinoLex) Error(s string) {
	fmt.Printf("syntax error on %s\n", l.input)
}

func isSpace(c byte) bool {
	return c == ' '
}
