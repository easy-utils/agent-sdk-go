// Package agentsdk provides a strong-typed client for the abc AgentService.
//
// The abc agent is a shared backend. ext servers (and the abc gateway)
// use this SDK to talk to it over easyrpc (agent.v1.AgentService). Web
// frontends and Flutter never connect to the agent directly — they go through
// abc, which forwards this same surface. The underlying RPC client is
// generated from agent/v1/agent.proto into agent.easyrpc.go; this package adds
// auth + a few ergonomic helpers on top of the easyrpc Transport.
package agentsdk

import (
	"context"
	"errors"
	"fmt"
	"strings"

	easyrpc "github.com/easy-utils/easy-rpc-go"

	agentv1 "github.com/easy-utils/agent-sdk-go/agent/v1"
)

// Client is a thin, status-aware agent session/file client.
type Client struct {
	base string
	svc  *agentv1.AgentServiceClient
}

// Option customizes the client.
type Option func(*clientOptions)

type clientOptions struct {
	realm easyrpc.Realm
	// extra metadata headers (e.g. auth) attached to every request.
	headers easyrpc.Headers
}

// WithRealm overrides the protocol realm (RealmStd = h1+h2c+h2, RealmAuto =
// adds h3). Default is RealmStd.
func WithRealm(r easyrpc.Realm) Option {
	return func(o *clientOptions) { o.realm = r }
}

// WithAuthHeader attaches a raw auth header (e.g. "Authorization: Bearer x")
// to every request. For cross-language consistency default to Bearer token.
func WithAuthHeader(name, value string) Option {
	return func(o *clientOptions) {
		if o.headers == nil {
			o.headers = easyrpc.Headers{}
		}
		o.headers.Set(name, value)
	}
}

// WithToken adds a Bearer token header to every request.
func WithToken(token string) Option {
	return func(o *clientOptions) {
		if o.headers == nil {
			o.headers = easyrpc.Headers{}
		}
		o.headers.Set("Authorization", "Bearer "+token)
	}
}

// h2cTransport returns a cleartext-H2 transport: the agent backend serves
// HTTP/2 only, so over plain http:// Go defaults to HTTP/1.1 which the agent
// rejects. Use RealmStd (h1+h2c+h2) which negotiates h2c.
func withBase(base string, o clientOptions) *Client {
	base = strings.TrimRight(base, "/")
	rt := easyrpc.WithMetadata(o.headers, &baseTransport{rt: easyrpc.NewTransport(o.realm), base: base})
	return &Client{
		base: base,
		svc:  agentv1.NewAgentServiceClient(rt),
	}
}

// baseTransport prepends the base URL to relative RPC paths so the generated
// client can use the canonical Connect path (/agent.v1.AgentService/<Method>).
type baseTransport struct {
	rt   easyrpc.Transport
	base string
}

func (b *baseTransport) prepend(u string) string {
	if strings.HasPrefix(u, "http://") || strings.HasPrefix(u, "https://") {
		return u
	}
	return b.base + u
}

func (b *baseTransport) Send(ctx context.Context, req easyrpc.Request) (easyrpc.Response, error) {
	req.URL = b.prepend(req.URL)
	return b.rt.Send(ctx, req)
}

func (b *baseTransport) OpenStream(ctx context.Context, req easyrpc.Request) (easyrpc.Stream, error) {
	req.URL = b.prepend(req.URL)
	return b.rt.OpenStream(ctx, req)
}

// New builds an agent client. token, when non-empty, is sent as a Bearer
// header on every request. baseURL is protocol+host (no trailing slash).
//
// Default transport speaks cleartext HTTP/2 (h2c prior knowledge) via
// RealmStd; pass WithRealm(RealmAuto) for h3, WithToken / WithAuthHeader for
// per-request auth.
func New(baseURL, token string, opts ...Option) *Client {
	if token == "" {
		token = "devtoken"
	}
	o := clientOptions{realm: easyrpc.RealmStd, headers: easyrpc.Headers{}}
	o.headers.Set("Authorization", "Bearer "+token)
	for _, opt := range opts {
		opt(&o)
	}
	return withBase(baseURL, o)
}

// ---- session ----

// Session mirrors the agent's session row.
type Session struct {
	Name         string
	Model        string
	Preset       string
	TipID        string
	Org, Repo    string
	Branch       string
	SystemPrompt string
	MaxTurns     int32
}

func sessionFromPb(s *agentv1.Session) Session {
	if s == nil {
		return Session{}
	}
	return Session{
		Name:         s.GetName(),
		Model:        s.GetModel(),
		Preset:       s.GetPreset(),
		TipID:        s.GetTipId(),
		Org:          s.GetOrg(),
		Repo:         s.GetRepo(),
		Branch:       s.GetBranch(),
		SystemPrompt: s.GetSystemPrompt(),
		MaxTurns:     s.GetMaxTurns(),
	}
}

// EnsureSession creates the session; an existing session ("already exists",
// modelled as AlreadyExists/Conflict) is treated as idempotent success.
func (c *Client) EnsureSession(ctx context.Context, name string) error {
	_, err := c.svc.CreateSession(ctx, &agentv1.CreateSessionRequest{Name: name})
	if err == nil {
		return nil
	}
	if IsExists(err) {
		return nil
	}
	return errDownstream("agent", err)
}

// ListSessions returns every session name.
func (c *Client) ListSessions(ctx context.Context) (map[string]bool, error) {
	res, err := c.svc.ListSessions(ctx, &agentv1.ListSessionsRequest{})
	if err != nil {
		return nil, errDownstream("agent", err)
	}
	out := map[string]bool{}
	for _, s := range res.GetSessions() {
		if s.GetName() != "" {
			out[s.GetName()] = true
		}
	}
	return out, nil
}

// GetSession returns the session row (nil when absent).
func (c *Client) GetSession(ctx context.Context, name string) (*Session, error) {
	res, err := c.svc.GetSession(ctx, &agentv1.GetSessionRequest{Id: name})
	if err != nil {
		if IsNotFound(err) {
			return nil, nil
		}
		return nil, errDownstream("agent", err)
	}
	s := sessionFromPb(res.GetSession())
	return &s, nil
}

// Fork pins a session fork. messageID empty forks from the tip; preset
// overrides the forked session's role. Existing name is idempotent success.
func (c *Client) Fork(ctx context.Context, parentSID, name, messageID, preset string) error {
	_, err := c.svc.Fork(ctx, &agentv1.ForkRequest{Id: parentSID, Name: name, MessageId: messageID, Preset: preset})
	if err == nil || IsExists(err) {
		return nil
	}
	return errDownstream("agent", err)
}

// DeleteSession removes a session; an absent session is idempotent success.
func (c *Client) DeleteSession(ctx context.Context, name string) error {
	_, err := c.svc.DeleteSession(ctx, &agentv1.DeleteSessionRequest{Id: name})
	if err == nil || IsNotFound(err) {
		return nil
	}
	return errDownstream("agent", err)
}

// ---- file ----

// GetFile returns a stored file's bytes + metadata by code.
func (c *Client) GetFile(ctx context.Context, code string) (data []byte, name, mime string, err error) {
	res, rerr := c.svc.GetFile(ctx, &agentv1.GetFileRequest{Code: code})
	if rerr != nil {
		if IsNotFound(rerr) {
			return nil, "", "", ErrNotFound
		}
		return nil, "", "", errDownstream("agent", rerr)
	}
	return res.GetData(), res.GetName(), res.GetMime(), nil
}

// GetFileMeta returns a stored file's metadata by code.
func (c *Client) GetFileMeta(ctx context.Context, code string) (*agentv1.GetFileMetaResponse, error) {
	res, err := c.svc.GetFileMeta(ctx, &agentv1.GetFileMetaRequest{Code: code})
	if err != nil {
		return nil, errDownstream("agent", err)
	}
	return res, nil
}


// ListPresets returns all presets, with the system prompt resolved per locale
// (best-match i18n map, falling back to the default prompt).
func (c *Client) ListPresets(ctx context.Context, locale string) ([]*agentv1.Preset, error) {
	res, err := c.svc.ListPresets(ctx, &agentv1.ListPresetsRequest{Locale: locale})
	if err != nil {
		return nil, errDownstream("agent", err)
	}
	return res.GetPresets(), nil
}

// ---- errors ----

var (
	ErrNotFound = errors.New("not found")
	ErrExists   = errors.New("already exists")
)

// IsNotFound reports a Connect NotFound (code 5).
func IsNotFound(err error) bool {
	return codeOf(err) == 5
}

// IsExists reports an AlreadyExists / InvalidArgument (codes 6 / 3).
func IsExists(err error) bool {
	c := codeOf(err)
	return c == 6 || c == 3
}

func codeOf(err error) int {
	if err == nil {
		return 0
	}
	if re, ok := err.(*easyrpc.RPCError); ok {
		return re.Code
	}
	return 13
}

func errDownstream(svc string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", svc, err)
}
