package cmd

import (
	"fmt"
	"os"

	"github.com/Ehco1996/clash-speed/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
)

func runTUI() error {
	f, err := tea.LogToFile("debug.log", "")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()

	m := ui.InitialModel()
	if err := m.FetchProxy(cfgFile); err != nil {
		return err
	}
	p := tea.NewProgram(
		m,
		tea.WithAltScreen(),       // use the full size of the terminal in its "alternate screen buffer"
		tea.WithMouseCellMotion(), // turn on mouse support so we can track the mouse wheel)
	)
	return p.Start()
}
