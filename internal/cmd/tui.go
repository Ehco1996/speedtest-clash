package cmd

import (
	"github.com/Ehco1996/clash-speed/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
)

func runTUI() error {
	m := ui.InitialModel()
	if err := m.FetchProxy(cfgFile); err != nil {
		return err
	}
	if err := m.FetchTestServers(); err != nil {
		return err
	}
	p := tea.NewProgram(m)
	return p.Start()
}
