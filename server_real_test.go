package main

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/zalando/go-keyring"
)

// TestRealServer runs the link and a change against a real trimletauth-server, only when asked:
//
//	TRIMLETAUTH_SERVER=http://127.0.0.1:8689 TRIMLETAUTH_TEST_INVITE=$(trimletauth-server invite) \
//	TRIMLETAUTH_TEST_USER=test-$RANDOM go test -run TestRealServer .   (in apps/auth)
func TestRealServer(t *testing.T) {
	server, invite, user := os.Getenv("TRIMLETAUTH_SERVER"), os.Getenv("TRIMLETAUTH_TEST_INVITE"), os.Getenv("TRIMLETAUTH_TEST_USER")
	if server == "" || invite == "" || user == "" {
		t.Skip("set TRIMLETAUTH_SERVER, TRIMLETAUTH_TEST_INVITE and TRIMLETAUTH_TEST_USER to run it")
	}
	keyring.MockInit()
	sm := newSyncModel(t, t.TempDir(), newFakeGitHub(t))
	sm.gl.api = server
	sm.newVault("master-pw")
	sm.withAccount()

	sm.linkServer(t, user, invite)
	sm.until(t, func() bool { return sm.screen == screenSync && !sm.sync.busy() })
	if sm.sync.done == "" {
		t.Fatal(sm.text())
	}
	sm.until(t, func() bool { return sm.link.Computer != "" })

	before := sm.link.SHA
	sm.screen = screenList
	sm.Rune('m')
	sm.typ("Renamed")
	sm.typ("")
	sm.until(t, func() bool { return !sm.saving && sm.link.SHA != before })
	if sm.vault.Accounts[0].Name != "Renamed" {
		t.Fatal(sm.text())
	}

	// A second device, approved by the first one.
	pc := serverPC(t, newFakeServer(t))
	pc.gl.api = server
	pc.Rune('2')
	pc.Rune('2')
	pc.typ(user)
	pc.typ("master-pw")
	pc.until(t, func() bool { return pc.screen == screenSync && pc.sync.pair != "" })
	sm.lastCheck = time.Time{}
	sm.tick()
	sm.until(t, func() bool { return sm.screen == screenApprove })
	sm.Rune('y')
	sm.typ("master-pw")
	pc.until(t, func() bool { return pc.screen == screenList })
	if pc.vault == nil || pc.vault.Accounts[0].Name != "Renamed" {
		t.Fatal(pc.text())
	}
	sm.until(t, func() bool { return sm.screen == screenList && !strings.Contains(sm.text(), "Approving") })

	sm.Rune('s')
	sm.Rune('1')
	sm.until(t, func() bool { return sm.screen == screenForm })
	sm.typ("delete")
	sm.typ("master-pw")
	sm.until(t, func() bool { return sm.screen == screenStart })
	if !strings.Contains(sm.text(), "still on trimlet.dev") {
		t.Fatal(sm.text())
	}
}
