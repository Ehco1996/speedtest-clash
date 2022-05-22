package ui

import (
	"fmt"
	"math"
	"strings"

	"github.com/muesli/termenv"
)

func (m model) View() string {
	if m.quitting {
		return "\n  See you later!\n\n"
	}
	if m.selectedProxyNode == "" {
		return m.viewSelectNode()
	} else if m.selectedServer == "" {
		return m.viewSelectServer()
	}
	return m.viewSpeedTest()
}

func (m model) viewSelectNode() string {
	// header
	tpl := "\nChoose your proxy node...\n"

	// body
	tpl += "\n%s\n"

	// footer
	tpl += subtle("up/down: select") + dot + subtle("enter: choose") + dot + subtle("q, esc: quit")

	nodes := ""
	for i, node := range m.proxyNodeList {
		nodes += fmt.Sprintf("%s\n", checkbox(node.Name(), m.proxyIdx == i))
	}
	return fmt.Sprintf(tpl, nodes)
}

func (m model) viewSelectServer() string {
	// header
	tpl := "\nChoose your speed test server...\n"

	// set user info
	tpl += "\n" + m.c.CurrentUser().String() + "\n"

	// body
	tpl += "\n%s\n"

	// footer
	tpl += subtle("up/down: select") + dot + subtle("enter: choose") + dot + subtle("q, esc: quit")

	server := ""
	for i, s := range m.testServerList {
		info := s.String() + fmt.Sprintf(" latency=[%dms]", s.Latency.Milliseconds())
		server += fmt.Sprintf("%s\n", checkbox(info, m.serverIdx == i))
	}
	return fmt.Sprintf(tpl, server)
}

func (m model) viewSpeedTest() string {
	label := "SpeedTesting..."
	msg := fmt.Sprintf("Proxy Node is %s and test server is %s...",
		keyword(m.selectedProxyNode), keyword(m.selectedServer))
	return subtle(msg) + "\n\n" + label + "\n" + progressbar(80, m.testPrecent) + "%"
}

func checkbox(label string, checked bool) string {
	if checked {
		return colorFg("[x] "+label, "212")
	}
	return fmt.Sprintf("[ ] %s", label)
}

func progressbar(width int, percent float64) string {
	w := float64(progressBarWidth)

	fullSize := int(math.Round(w * percent))
	var fullCells string
	for i := 0; i < fullSize; i++ {
		fullCells += termenv.String(progressFullChar).Foreground(term.Color(ramp[i])).String()
	}

	emptySize := int(w) - fullSize
	emptyCells := strings.Repeat(progressEmpty, emptySize)

	return fmt.Sprintf("%s%s %3.0f", fullCells, emptyCells, math.Round(percent*100))
}
