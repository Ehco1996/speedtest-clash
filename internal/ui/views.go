package ui

import (
	"fmt"
	"log"
)

func (m model) View() string {
	if m.ps.selectedProxyNode == "" {
		return m.viewSelectNode()
	} else if m.sts.selectedServer == "" {
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
	for i, node := range m.ps.proxyNodeList {
		nodes += fmt.Sprintf("%s\n", checkbox(node.Name(), m.ps.proxyIdx == i))
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
	for i, s := range m.sts.testServerList {
		info := s.String() + fmt.Sprintf(" latency=[%dms]", s.Latency.Milliseconds())
		server += fmt.Sprintf("%s\n", checkbox(info, m.sts.serverIdx == i))
	}
	return fmt.Sprintf(tpl, server)
}

func (m model) viewSpeedTest() string {
	log.Printf("refresh once res=%s ", m.str.currentRes.String())

	label := "SpeedTesting..."

	title := fmt.Sprintf("Proxy Node is %s and the test server is %s",
		keyword(m.ps.selectedProxyNode), keyword(m.sts.selectedServer))

	speed := m.str.currentRes.CurrentSpeed
	if m.quitting {
		speed = m.sts.testServerList[m.sts.serverIdx].DLSpeed
	}
	speedDownload := fmt.Sprintf("\nDownloading %s ....  %.2f mbps", m.str.spinner.View(), speed)
	// speedUpload := fmt.Sprintf("\nUploading %s ....  %d Mbps", m.sp.spinner.View(), m.sp.upload)

	content := subtle(title) + "\n\n" + label + "\n" + m.str.progress.View() + "\n" + speedDownload

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
