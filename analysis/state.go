package analysis

import (
	"fmt"
	"log"
	"regexp"
	"simple-lsp/lsp"
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

var worstEditos = []string{"VS Code", "JetBrains", "InteliJ"}

func getWorstEditorsIndex(line string, logger *log.Logger) (int, string) {
	idx := -1
	editor := ""

	for _, editorName := range worstEditos {
		idx = strings.Index(line, editorName)
		logger.Printf("Trying to find %s, on line: %s, idx: %d", editorName, line, idx)
		if idx >= 0 {
			editor = editorName
			break
		}
	}

	return idx, editor
}

func (s *State) TextDocumentCodeAction(id int, uri string, logger *log.Logger) lsp.TextDocumentCodeActionResponse {
	text := s.Documents[uri]

	actions := []lsp.CodeAction{}

	for row, line := range strings.Split(text, "\n") {
		idx, editorName := getWorstEditorsIndex(line, logger)
		if idx >= 0 && len(editorName) > 0 {
			logger.Printf("Editor %s found on line %d", editorName, line)
			replaceChange := map[string][]lsp.TextEdit{}
			replaceChange[uri] = []lsp.TextEdit{
				{
					Range:   LineRange(row, idx, idx+len(editorName)),
					NewText: "NeoVim",
				},
			}

			actions = append(actions, lsp.CodeAction{
				Title: fmt.Sprintf("Replace %s with the Best editor out there", editorName),
				Edit:  &lsp.WorkspaceEdit{Changes: replaceChange},
			})

			removeShitIDE := map[string][]lsp.TextEdit{}
			removeShitIDE[uri] = []lsp.TextEdit{
				{
					Range:   LineRange(row, idx, idx+len(editorName)),
					NewText: "Shit IDE Removed",
				},
			}

			actions = append(actions, lsp.CodeAction{
				Title: fmt.Sprintf("Remove %s shit IDE from this file", editorName),
				Edit:  &lsp.WorkspaceEdit{Changes: removeShitIDE},
			})
		}
	}

	response := lsp.TextDocumentCodeActionResponse{
		Response: lsp.Response{
			RPC: "2.0",
			ID:  &id,
		},
		Result: actions,
	}
	return response
}

func LineRange(line, start, end int) lsp.Range {
	return lsp.Range{
		Start: lsp.Position{
			Line:      line,
			Character: start,
		},
		End: lsp.Position{
			Line:      line,
			Character: end,
		},
	}
}
