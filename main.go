package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"golang-lsp/analysis"
	"golang-lsp/lsp"
	"golang-lsp/rpc"
	"io"
	"log"
	"os"
	"strconv"
)

func main() {
	logger := getLogger("/Users/douglasmendes/Git/personal/golang-lsp/log.txt")
	logger.Println("Hey, I started")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(Split)

	state := analysis.NewState()
	writer := os.Stdout

	for scanner.Scan() {
		msg := scanner.Bytes()

		method, contents, err := rpc.DecodeMessage(msg)
		if err != nil {
			logger.Println("GOt an error: %s:", err)
		}
		handleMessage(logger, writer, state, method, contents)
	}
}

func handleMessage(logger *log.Logger, writer io.Writer, state analysis.State, method string, contents []byte) {
	logger.Printf("Received msg with method: %s", method)

	switch method {
	case "initialize":
		var request lsp.InitializeRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("Problem to parse msg %s", err)
		}

		logger.Printf("Connected to: %s %s",
			request.Params.ClientInfo.Name,
			request.Params.ClientInfo.Version)
		msg := lsp.NewInitializeResponse(request.ID)
		writeResponse(writer, msg)
		logger.Println("Sent the reply")
	case "textDocument/didOpen":
		var request lsp.DidOpenTextDocumentNotification
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("textDocument/didOpen: %s", err)
			return
		}
		logger.Printf("textDocument/didOpen: %s", request.Params.TextDocument.Uri)
		state.OpenDocument(
			request.Params.TextDocument.Uri,
			request.Params.TextDocument.Text)
	case "textDocument/didChange":
		var request lsp.DidChangeTextDocumentNotification
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("textDocument/didChange: %s", err)
			return
		}

		logger.Printf("textDocument/didChange: %s", request.Params.TextDocument.Uri)
		for _, change := range request.Params.ContentChanges {
			state.UpdateDocument(request.Params.TextDocument.Uri, change.Text)
			logger.Printf("%#v\n", change)

		}
	case "textDocument/hover":
		var request lsp.HoverRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("textDocument/hover: %s", err)
			return
		}

		logger.Printf("textDocument/hover: %s", request.Params.TextDocumentPositionParams.Textdocument.Uri)

		response := state.Hover(request.ID, request.Params.Textdocument.Uri, request.Params.Position, logger)
		writeResponse(writer, response)
	case "textDocument/definition":
		var request lsp.DefinitionRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("textDocument/definition: %s", err)
			return
		}

		logger.Printf("textDocument/definition: uri: %s, pos: %s", request.Params.Textdocument.Uri, request.Params.Position)
		response := state.Definition(request.ID, request.Params.Textdocument.Uri, request.Params.Position)

		writeResponse(writer, response)

	}
}

func writeResponse(writer io.Writer, msg any) {
	reply := rpc.EncodeMessage(msg)
	writer.Write([]byte(reply))
}

// type SplitFunc func(data []byte, atEOF bool) (advance int, token []byte, err error)
func Split(data []byte, _ bool) (advance int, token []byte, err error) {
	header, content, found := bytes.Cut(data, []byte{'\r', '\n', '\r', '\n'})

	if !found {
		return 0, nil, nil
	}

	contentLengthBytes := header[len("Content-Length: "):]
	contentLength, err := strconv.Atoi(string(contentLengthBytes))
	if err != nil {
		return 0, nil, err
	}

	if len(content) < contentLength {
		return 0, nil, nil
	}

	totalLength := len(header) + 4 + contentLength
	return totalLength, data[:totalLength], nil
}

func getLogger(filename string) *log.Logger {
	logfile, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		panic("Hey dude give the right file!")
	}

	return log.New(logfile, "[go-lsp]", log.Ldate|log.Ltime|log.Lshortfile)
}
