package agentsdk

import (
	"context"
	"testing"

	agentv1 "github.com/easy-utils/agent-sdk-go/agent/v1"
	"github.com/easy-utils/easy-rpc-go"
)

// stubAgent implements the AgentService surface the SDK exercises.
type stubAgent struct {
	sessions []string
}

func (s *stubAgent) Health(ctx context.Context, _ *agentv1.HealthRequest) (*agentv1.HealthResponse, error) {
	return &agentv1.HealthResponse{Ok: true, Name: "stub-agent"}, nil
}

func (s *stubAgent) ListSessions(ctx context.Context, _ *agentv1.ListSessionsRequest) (*agentv1.ListSessionsResponse, error) {
	res := &agentv1.ListSessionsResponse{}
	for _, n := range s.sessions {
		res.Sessions = append(res.Sessions, &agentv1.Session{Name: n})
	}
	return res, nil
}

func (s *stubAgent) CreateSession(ctx context.Context, in *agentv1.CreateSessionRequest) (*agentv1.CreateSessionResponse, error) {
	for _, n := range s.sessions {
		if n == in.Name {
			return nil, &easyrpc.RPCError{Code: 6, Message: "already exists"}
		}
	}
	s.sessions = append(s.sessions, in.Name)
	return &agentv1.CreateSessionResponse{Ok: true, SessionName: in.Name}, nil
}

func (s *stubAgent) GetSession(ctx context.Context, in *agentv1.GetSessionRequest) (*agentv1.GetSessionResponse, error) {
	for _, n := range s.sessions {
		if n == in.Id {
			return &agentv1.GetSessionResponse{Session: &agentv1.Session{Name: n}}, nil
		}
	}
	return nil, &easyrpc.RPCError{Code: 5, Message: "not found"}
}

func (s *stubAgent) DeleteSession(ctx context.Context, in *agentv1.DeleteSessionRequest) (*agentv1.DeleteSessionResponse, error) {
	out := s.sessions[:0]
	for _, n := range s.sessions {
		if n != in.Id {
			out = append(out, n)
		}
	}
	s.sessions = out
	return &agentv1.DeleteSessionResponse{Ok: true}, nil
}

func (s *stubAgent) Fork(ctx context.Context, in *agentv1.ForkRequest) (*agentv1.ForkResponse, error) {
	return &agentv1.ForkResponse{Session: &agentv1.Session{Name: in.Name}}, nil
}

// remaining methods are not exercised by the SDK unit tests; implement via
// a helper that maps them to nil (compiles fine without connectrpc).
func (s *stubAgent) ListMessages(ctx context.Context, _ *agentv1.ListMessagesRequest) (*agentv1.ListMessagesResponse, error) { return nil, nil }
func (s *stubAgent) Prompt(ctx context.Context, _ *agentv1.PromptRequest, _ func(*agentv1.PromptResponse) error) error { return nil }
func (s *stubAgent) WatchSession(ctx context.Context, _ *agentv1.WatchSessionRequest, _ func(*agentv1.WatchSessionResponse) error) error { return nil }
func (s *stubAgent) WatchSessions(ctx context.Context, _ *agentv1.WatchSessionsRequest, _ func(*agentv1.WatchSessionsResponse) error) error { return nil }
func (s *stubAgent) Rename(ctx context.Context, _ *agentv1.RenameRequest) (*agentv1.RenameResponse, error) { return nil, nil }
func (s *stubAgent) SetModel(ctx context.Context, _ *agentv1.SetModelRequest) (*agentv1.SetModelResponse, error) { return nil, nil }
func (s *stubAgent) Undo(ctx context.Context, _ *agentv1.UndoRequest) (*agentv1.UndoResponse, error) { return nil, nil }
func (s *stubAgent) State(ctx context.Context, _ *agentv1.StateRequest) (*agentv1.StateResponse, error) { return nil, nil }
func (s *stubAgent) Mailbox(ctx context.Context, _ *agentv1.MailboxRequest) (*agentv1.MailboxResponse, error) { return nil, nil }
func (s *stubAgent) UpdateSettings(ctx context.Context, _ *agentv1.UpdateSettingsRequest) (*agentv1.UpdateSettingsResponse, error) { return nil, nil }
func (s *stubAgent) Interrupt(ctx context.Context, _ *agentv1.InterruptRequest) (*agentv1.InterruptResponse, error) { return nil, nil }
func (s *stubAgent) Compact(ctx context.Context, _ *agentv1.CompactRequest) (*agentv1.CompactResponse, error) { return nil, nil }
func (s *stubAgent) ListProviders(ctx context.Context, _ *agentv1.ListProvidersRequest) (*agentv1.ListProvidersResponse, error) { return nil, nil }
func (s *stubAgent) ListProvidersCatalog(ctx context.Context, _ *agentv1.ListProvidersCatalogRequest) (*agentv1.ListProvidersCatalogResponse, error) { return nil, nil }
func (s *stubAgent) RegisterProvider(ctx context.Context, _ *agentv1.RegisterProviderRequest) (*agentv1.RegisterProviderResponse, error) { return nil, nil }
func (s *stubAgent) DeleteProvider(ctx context.Context, _ *agentv1.DeleteProviderRequest) (*agentv1.DeleteProviderResponse, error) { return nil, nil }
func (s *stubAgent) TestProvider(ctx context.Context, _ *agentv1.TestProviderRequest) (*agentv1.TestProviderResponse, error) { return nil, nil }
func (s *stubAgent) ListModels(ctx context.Context, _ *agentv1.ListModelsRequest) (*agentv1.ListModelsResponse, error) { return nil, nil }
func (s *stubAgent) ListPresets(ctx context.Context, _ *agentv1.ListPresetsRequest) (*agentv1.ListPresetsResponse, error) { return nil, nil }
func (s *stubAgent) UpsertPreset(ctx context.Context, _ *agentv1.UpsertPresetRequest) (*agentv1.UpsertPresetResponse, error) { return nil, nil }
func (s *stubAgent) DeletePreset(ctx context.Context, _ *agentv1.DeletePresetRequest) (*agentv1.DeletePresetResponse, error) { return nil, nil }
func (s *stubAgent) PreviewPreset(ctx context.Context, _ *agentv1.PreviewPresetRequest) (*agentv1.PreviewPresetResponse, error) { return nil, nil }
func (s *stubAgent) GetConfig(ctx context.Context, _ *agentv1.GetConfigRequest) (*agentv1.GetConfigResponse, error) { return nil, nil }
func (s *stubAgent) SetConfig(ctx context.Context, _ *agentv1.SetConfigRequest) (*agentv1.SetConfigResponse, error) { return nil, nil }
func (s *stubAgent) ListTools(ctx context.Context, _ *agentv1.ListToolsRequest) (*agentv1.ListToolsResponse, error) { return nil, nil }
func (s *stubAgent) GetToolConfig(ctx context.Context, _ *agentv1.GetToolConfigRequest) (*agentv1.GetToolConfigResponse, error) { return nil, nil }
func (s *stubAgent) SetToolConfig(ctx context.Context, _ *agentv1.SetToolConfigRequest) (*agentv1.SetToolConfigResponse, error) { return nil, nil }
func (s *stubAgent) SetExtensionConfig(ctx context.Context, _ *agentv1.SetExtensionConfigRequest) (*agentv1.SetExtensionConfigResponse, error) { return nil, nil }
func (s *stubAgent) UploadFile(ctx context.Context, _ *agentv1.UploadFileRequest) (*agentv1.UploadFileResponse, error) { return nil, nil }
func (s *stubAgent) IngestFile(ctx context.Context, _ *agentv1.IngestFileRequest) (*agentv1.IngestFileResponse, error) { return nil, nil }
func (s *stubAgent) GetFile(ctx context.Context, _ *agentv1.GetFileRequest) (*agentv1.GetFileResponse, error) { return nil, nil }
func (s *stubAgent) GetFileMeta(ctx context.Context, _ *agentv1.GetFileMetaRequest) (*agentv1.GetFileMetaResponse, error) { return nil, nil }
func (s *stubAgent) GetFileStream(ctx context.Context, _ *agentv1.GetFileRequest, _ func(*agentv1.FileChunk) error) error {
	return nil
}
func (s *stubAgent) GetAgentConfig(ctx context.Context, _ *agentv1.GetAgentConfigRequest) (*agentv1.GetAgentConfigResponse, error) { return nil, nil }

// inMemoryTransport dispatches to a registered AgentService in-process.
type inMemoryTransport struct {
	svc agentv1.AgentServiceService
	reg *easyrpc.ServiceRegistry
}

func newInMemory(impl agentv1.AgentServiceService) *inMemoryTransport {
	return &inMemoryTransport{svc: impl, reg: agentv1.RegisterAgentServiceService(impl)}
}

// memWriter captures a push-dispatch response into an easyrpc.Response.
type memWriter struct {
	status  int
	headers easyrpc.Headers
	body    []byte
}

func (w *memWriter) Status(code int)                    { w.status = code }
func (w *memWriter) Header(h easyrpc.Headers)           { w.headers = h }
func (w *memWriter) WriteFrame(payload []byte) error    { w.body = append(w.body, payload...); return nil }

func (t *inMemoryTransport) Send(ctx context.Context, req easyrpc.Request) (easyrpc.Response, error) {
	w := &memWriter{status: 200}
	// The real net/http bridge sets the unary content-type when the caller
	// omitted it; Dispatch (easy-rpc-go v1.4) requires it for codec/shape
	// negotiation, so mirror that here.
	if req.Headers.Get("Content-Type") == "" {
		if req.Headers == nil {
			req.Headers = easyrpc.Headers{}
		}
		req.Headers.Set("Content-Type", easyrpc.ContentTypeUnary)
	}
	_ = easyrpc.Dispatch(ctx, req, agentv1.AgentService_Methods(), t.reg, w)
	resp := easyrpc.Response{Status: w.status, Headers: w.headers, Body: w.body}
	// Mirror the real bridge: reconstruct the connect error from the
	// connect-code header when present, else from the JSON error body (which
	// carries the exact code — several Connect codes share an HTTP status).
	if w.status >= 300 {
		if e := easyrpc.StatusFromHeader(w.headers); e != nil {
			resp.Error = e
		} else if c, m, ds := easyrpc.DecodeErrorJSON(w.body); c != 0 {
			resp.Error = &easyrpc.RPCError{Code: c, Message: m, Details: ds}
		} else {
			resp.Error = &easyrpc.RPCError{Code: easyrpc.ConnectFromStatus(w.status), Message: string(w.body)}
		}
	}
	return resp, nil
}
func (t *inMemoryTransport) OpenStream(ctx context.Context, req easyrpc.Request) (easyrpc.Stream, error) {
	panic("not used")
}

func TestEnsureSessionIdempotent(t *testing.T) {
	stub := &stubAgent{sessions: []string{"existing"}}
	rt := &inMemoryTransport{svc: stub, reg: agentv1.RegisterAgentServiceService(stub)}
	c := &Client{base: "mem://", svc: agentv1.NewAgentServiceClient(rt)}
	if err := c.EnsureSession(context.Background(), "existing"); err != nil {
		t.Fatalf("existing session should be idempotent: %v", err)
	}
	if err := c.EnsureSession(context.Background(), "new"); err != nil {
		t.Fatalf("create new session: %v", err)
	}
}

func TestListSessions(t *testing.T) {
	stub := &stubAgent{sessions: []string{"a", "b"}}
	rt := &inMemoryTransport{svc: stub, reg: agentv1.RegisterAgentServiceService(stub)}
	c := &Client{base: "mem://", svc: agentv1.NewAgentServiceClient(rt)}
	got, err := c.ListSessions(context.Background())
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if !got["a"] || !got["b"] || len(got) != 2 {
		t.Fatalf("unexpected sessions: %v", got)
	}
}

func TestGetSessionNotFound(t *testing.T) {
	stub := &stubAgent{sessions: []string{"a"}}
	rt := &inMemoryTransport{svc: stub, reg: agentv1.RegisterAgentServiceService(stub)}
	c := &Client{base: "mem://", svc: agentv1.NewAgentServiceClient(rt)}
	s, err := c.GetSession(context.Background(), "missing")
	if err != nil {
		t.Fatalf("get missing should be nil,nil: %v", err)
	}
	if s != nil {
		t.Fatalf("expected nil session")
	}
}
