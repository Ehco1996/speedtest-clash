package ui

import (
	"fmt"
)

func (m model) View() string {
	if m.quitting {
		return "\n  See you later!\n\n Press ctrl + c to quit"
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
	return subtle(msg) + "\n\n" + label + "\n" + m.progress.View()
}

func checkbox(label string, checked bool) string {
	if checked {
		return colorFg("[x] "+label, "212")
	}
	return fmt.Sprintf("[ ] %s", label)
}
