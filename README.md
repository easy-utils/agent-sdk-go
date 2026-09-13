# agent-sdk-go

Typed Go client for the standalone agent backend (`agent.v1.AgentService` +
`AdminService`) over **easy-rpc** (Connect wire, true server-streaming).

```go
c := agentsdk.New("https://agent.example.com", "your-tenant-token")
presets, _ := c.ListPresets(ctx, "en")
```
