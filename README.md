# agent-sdk-go

Typed Go client for the standalone agent backend (`agent.v1.AgentService` +
`AdminService`) over **easy-rpc** (Connect wire, true server-streaming).

The generated protobuf types + easy-rpc client are vendored in `agent/v1/`
(regenerate from `agent/v1/agent.proto`, easy-rpc 幂等生成的 `agent.pb.go` /
`agent.easyrpc.go`); this package adds auth plus a few ergonomic helpers.

```go
c := agentsdk.New("https://agent.example.com", "your-tenant-token")
presets, _ := c.ListPresets(ctx, "en")
```
