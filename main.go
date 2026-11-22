// Trimlet Auth: the two-factor codes (TOTP) of my accounts, kept in a file encrypted
// with a master password that the system keychain remembers.
package main

import (
	"net/url"
	"time"

	"github.com/trimlet/gui"
)

func main() {
	a, w := gui.NewWindow("local.trimletauth", "Trimlet Auth")
	gui.UsePalette("copper") // Auth's color
	m := newModel("")
	path, err := vaultPath()
	if err != nil {
		m.fail("Could not find the folder for the vault: " + err.Error())
	} else {
		m.path = path
		m.start()
	}

	m.copy = func(text string) { a.Clipboard().SetContent(text) }
	m.clipboard = func() string { return a.Clipboard().Content() }
	m.quit = a.Quit
	m.openURL = func(s string) {
		if u, err := url.Parse(s); err == nil {
			a.OpenURL(u)
		}
	}
	// A copied code does not stay in the clipboard after Trimlet Auth is closed.
	w.SetOnClosed(func() {
		if m.copied != "" && m.clipboard() == m.copied {
			m.copy("")
		}
	})
	screen := gui.NewScreen(a, w, m)
	m.ui = screen.Do

	// The codes and the seconds left change by themselves.
	go func() {
		for range time.Tick(time.Second) {
			screen.Do(m.tick)
		}
	}()
	gui.Run(a, w)
}
