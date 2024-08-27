package lsp

type ClientInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type InitializeRequestParams struct {
	ClientInfo *ClientInfo `json:"clientInfo"`
}

type InitializeRequest struct {
	Params InitializeRequestParams `json:"params"`
	Request
}

type InitializeResponse struct {
	Response
	Result InitializeResult `json:"result"`
}

type InitializeResult struct {
	ServerInfo   ServerInfo         `json:"serverInfo"`
	Capabilities ServerCapabilities `json:"capabilities"`
}

type (
	ServerCapabilities struct {
		TextDocumentSync   int  `json:"textDocumentSync"`
		HoverProvider      bool `json:"hoverProvider"`
		DefinitionProvider bool `json:"definitionProvider"`
	}
	ServerInfo struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	}
)

func NewInitializeResponse(id int) InitializeResponse {
	return InitializeResponse{
		Response: Response{
			RPC: "2.0",
			ID:  &id,
		},
		Result: InitializeResult{
			Capabilities: ServerCapabilities{
				TextDocumentSync:   1,
				HoverProvider:      true,
				DefinitionProvider: true,
			},
			ServerInfo: ServerInfo{
				Name:    "go-local-lsp",
				Version: "0.0.0.0.0.0-beta1.final",
			},
		},
	}
}
