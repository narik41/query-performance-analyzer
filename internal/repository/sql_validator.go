package repository

import (
	"errors"
	"regexp"
	"strings"
)

var ErrInvalidQuery = errors.New("invalid query: only read-only SELECT queries are allowed")

type SQLValidator struct {
	forbiddenKeywords []*regexp.Regexp
	allowedPrefixes   []string
}

func NewSQLValidator() *SQLValidator {
	return &SQLValidator{
		forbiddenKeywords: []*regexp.Regexp{
			regexp.MustCompile(`\bINSERT\b`),
			regexp.MustCompile(`\bUPDATE\b`),
			regexp.MustCompile(`\bDELETE\b`),
			regexp.MustCompile(`\bREPLACE\b`),
			regexp.MustCompile(`\bMERGE\b`),
			regexp.MustCompile(`\bUPSERT\b`),
			regexp.MustCompile(`\bCREATE\b`),
			regexp.MustCompile(`\bALTER\b`),
			regexp.MustCompile(`\bDROP\b`),
			regexp.MustCompile(`\bTRUNCATE\b`),
			regexp.MustCompile(`\bRENAME\b`),
			regexp.MustCompile(`\bGRANT\b`),
			regexp.MustCompile(`\bREVOKE\b`),
			regexp.MustCompile(`\bCALL\b`),
			regexp.MustCompile(`\bDO\b`),
			regexp.MustCompile(`\bHANDLER\b`),
			regexp.MustCompile(`\bLOAD\b\s+DATA\b`),
			regexp.MustCompile(`\bLOAD\b\s+XML\b`),
		},

		allowedPrefixes: []string{"select", "with", "explain", "show", "describe", "desc"},
	}
}

func (v *SQLValidator) Validate(query string) error {

	normalized := v.normalizeQuery(query)

	if !v.hasAllowedPrefix(normalized) {
		return ErrInvalidQuery
	}

	for _, re := range v.forbiddenKeywords {
		if re.MatchString(normalized) {
			return ErrInvalidQuery
		}
	}

	return nil
}

func (v *SQLValidator) normalizeQuery(query string) string {

	lines := strings.Split(query, "\n")
	var cleanedLines []string
	for _, line := range lines {
		if idx := strings.Index(line, "--"); idx != -1 {
			line = line[:idx]
		}
		cleanedLines = append(cleanedLines, line)
	}
	query = strings.Join(cleanedLines, " ")

	for {
		start := strings.Index(query, "/*")
		if start == -1 {
			break
		}
		end := strings.Index(query[start+2:], "*/")
		if end == -1 {
			query = query[:start]
			break
		}
		query = query[:start] + query[start+2+end+2:]
	}

	query = regexp.MustCompile(`\s+`).ReplaceAllString(query, " ")
	return strings.TrimSpace(strings.ToUpper(query))
}

func (v *SQLValidator) hasAllowedPrefix(query string) bool {
	upper := strings.ToUpper(query)
	for _, prefix := range v.allowedPrefixes {
		if strings.HasPrefix(upper, strings.ToUpper(prefix)) {
			return true
		}
	}
	return false
}

func ValidateReadOnlyQuery(query string) error {
	return NewSQLValidator().Validate(query)
}
