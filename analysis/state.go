package analysis

import (
	"fmt"
	"golang-lsp/lsp"
	"log"
	"regexp"
	"strings"
	"unicode"
)

type State struct {
	Documents map[string]string
}

func NewState() State {
	return State{Documents: map[string]string{}}
}

func (s *State) OpenDocument(uri, text string) {
	s.Documents[uri] = text
}

func (s *State) UpdateDocument(uri, text string) {
	s.Documents[uri] = text
}

func (s *State) Hover(id int, uri string, position lsp.Position, logger *log.Logger) lsp.HoverResponse {
	document := s.Documents[uri]

	documentation := ""

	switch getTagName(document, position, logger) {
	case "douglas":
		documentation = "Engeneering Manager at Xpertsea"
	case "meamed":
		documentation = "CTO at Xpertsea"
	case "onildo":
		documentation = "Tiozao bruxo da programação"
	default:
		documentation = fmt.Sprintf("File %s, Characters: %d", uri, len(document))

	}

	return lsp.HoverResponse{
		Response: lsp.Response{
			ID:  &id,
			RPC: "2.0",
		},
		Result: lsp.HoverResult{
			Contents: documentation,
		},
	}
}

func getTagName(content string, position lsp.Position, logger *log.Logger) string {
	// Split the content into lines
	lines := strings.Split(content, "\n")
	if position.Line < 0 || position.Line >= len(lines) {
		logger.Fatal("Position is out of range %s", content)
		return ""
	}

	line := lines[position.Line]
	if position.Character < 0 || position.Character >= len(line) {
		logger.Fatal("Position is out of range %s", content)

		return ""
	}

	start := position.Character
	end := position.Character

	for start > 0 && isWordChar(line[start-1]) {
		start--
	}

	for end < len(line) && isWordChar(line[end]) {
		end++
	}

	// Extract the word at the cursor position
	word := line[start:end]
	logger.Printf("Word to validate %s", word)

	// Check if the word matches the pattern "@<username>"
	pattern := regexp.MustCompile(`^@(\w+)$`)
	matches := pattern.FindStringSubmatch(word)

	if len(matches) == 2 {
		logger.Println("Word: %s", matches)
		return matches[1] // Return the captured username without the "@" symbol
	}

	logger.Println("Word not found: %s")
	return "" // Return an empty string if the cursor is not on a tag name
}

func isWordChar(char byte) bool {
	return unicode.IsLetter(rune(char)) || unicode.IsDigit(rune(char)) || char == '_' || char == '@'
}

func (s *State) Definition(id int, uri string, position lsp.Position) lsp.DefinitionResponse {
	// document := s.Documents[uri]
	//
	newLinePosition := position.Line - 1

	if newLinePosition < 0 {
		newLinePosition = 0
	}

	return lsp.DefinitionResponse{
		Response: lsp.Response{
			ID:  &id,
			RPC: "2.0",
		},
		Result: lsp.Location{
			Uri: uri,
			Range: lsp.Range{
				Start: lsp.Position{
					Line:      newLinePosition,
					Character: 0,
				},
				End: lsp.Position{
					Line:      newLinePosition,
					Character: 0,
				},
			},
		},
	}
}
