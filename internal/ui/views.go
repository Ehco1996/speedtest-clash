package ui

import (
	"fmt"
)

func (m model) View() string {
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

	title := fmt.Sprintf("Proxy Node is %s and the test server is %s",
		keyword(m.selectedProxyNode), keyword(m.selectedServer))

	speedDownload := fmt.Sprintf("\nDownloading %s ....  %.2f mbps", m.sp.spinner.View(), m.sp.currentRes.CurrentSpeed)
	// speedUpload := fmt.Sprintf("\nUploading %s ....  %d Mbps", m.sp.spinner.View(), m.sp.upload)

	content := subtle(title) + "\n\n" + label + "\n" + m.progress.View() + "\n" + speedDownload

	if m.quitting {
		content += m.viewQuit()
	}
	return content
}

func (m model) viewQuit() string {
	if m.quitting {
		return "\n\nSee you later!\n\nPress ctrl + c to quit！"
	}
	return ""
}
