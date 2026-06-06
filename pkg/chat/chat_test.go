package chat

import (
	"context"
	"strings"
	"testing"

	"github.com/agennext/agent-chat/pkg/capability"
	"github.com/agennext/agent-chat/pkg/event"
	"github.com/agennext/agent-chat/pkg/guard"
	"github.com/agennext/agent-chat/pkg/loop"
	"github.com/agennext/agent-chat/pkg/store"
)

type answerReasoner struct{ text string }

func (a answerReasoner) Reason(_ context.Context, _ loop.State) (loop.Action, error) {
	return loop.Action{Kind: loop.ActionAnswer, Answer: a.text}, nil
}

func newEngine(r loop.Reasoner) *loop.Engine {
	return &loop.Engine{
		Reasoner: r,
		Registry: capability.NewRegistry(),
		Screener: guard.NewHeuristicScreener(),
		Decider:  guard.NewStaticDecider(),
		Ctx:      store.NewMemContextStore(),
		Mem:      store.NewMemMemoryStore(),
		Dedupe:   loop.NewMemDeduper(),
		Budget:   loop.DefaultBudget(),
	}
}

func TestChatCoreRunsAdmittedTurn(t *testing.T) {
	t.Parallel()
	core := New(newEngine(answerReasoner{"hello there"}))
	ev := event.New("e1", "web", "chat.message.v1", "s1", "acme", "u1", "", []byte("hi"))
	admitted := loop.AdmittedEvent{Event: ev, Principal: "u1", Scope: capability.Scope{Tenants: []string{"acme"}}}
	res, err := core.Run(context.Background(), admitted)
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if !strings.Contains(res.Answer, "hello there") {
		t.Fatalf("unexpected answer: %q", res.Answer)
	}
	if res.Escalated {
		t.Fatal("a clean answer must not escalate")
	}
}
