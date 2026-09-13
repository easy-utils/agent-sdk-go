package agentsdk

import (
	"context"
	"testing"
	"time"

	agentv1 "github.com/easy-utils/agent-proto/agent/v1"
)

// TestLiveLocal exercises the regenerated SDK against the locally running agent
// (supervisor agent-server on :18080).
func TestLiveLocal(t *testing.T) {
	c := New("http://127.0.0.1:18080", "devtoken")
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()

	presets, err := c.ListPresets(ctx, "zh")
	if err != nil { t.Fatal("listPresets:", err) }
	if len(presets) == 0 { t.Fatal("no presets") }
	t.Log("presets:", len(presets), presets[0].Id, presets[0].SystemPrompt)

	ss, err := c.ListSessions(ctx)
	if err != nil { t.Fatal("listSessions:", err) }
	t.Log("sessions:", len(ss))

	// New RPC: WatchSessions (stream) — first frame must be quick.
	start := time.Now()
	st, err := c.svc.WatchSessions(ctx, &agentv1.WatchSessionsRequest{})
	if err != nil { t.Fatal("watchSessions:", err) }
	if _, err := st.Recv(); err != nil { t.Fatal("recv:", err) }
	t.Log("watchSessions first frame after", time.Since(start))
}
