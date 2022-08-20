package ui

import (
	"fmt"
	"log"
)

func (m model) View() string {
	if m.ps.selectedProxyNode == "" {
		return m.viewSelectProxyNode()
	} else if m.ts.selectedServer == "" {
		return m.viewSelectServer()
	}
	return m.viewSpeedTest()
}

func (m model) viewSelectProxyNode() string {
	return proxyNodeListStyle.Render(m.ps.uiList.View())
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
	for i, s := range m.ts.testServerList {
		info := s.String() + fmt.Sprintf(" latency=[%dms]", s.Latency.Milliseconds())
		server += fmt.Sprintf("%s\n", checkbox(info, m.ts.serverIdx == i))
	}
	return fmt.Sprintf(tpl, server)
}

func (m model) viewSpeedTest() string {
	log.Printf("refresh once res=%s", m.tp.currentRes.String())

	label := fmt.Sprintf("Running SpeedTesting Type: %s ...", m.tp.currentRes.Type)

	title := fmt.Sprintf("Proxy Node is %s and the test server is %s",
		keyword(m.ps.selectedProxyNode), keyword(m.ts.selectedServer))

	downLoadSpeed := m.tp.currentRes.CurrentSpeed
	if m.tp.finishDownloadTest {
		downLoadSpeed = m.ts.testServerList[m.ts.serverIdx].DLSpeed
	}
	speedDownloadContent := fmt.Sprintf("\nDownloading %s ....  %.2f mbps", m.tp.spinner.View(), downLoadSpeed)

	upLoadSpeed := m.tp.currentRes.CurrentSpeed
	if m.tp.finishUploadTest {
		upLoadSpeed = m.ts.testServerList[m.ts.serverIdx].ULSpeed
	}
	speedUploadContent := fmt.Sprintf("\nUploading %s ....  %.2f mbps", m.tp.spinner.View(), upLoadSpeed)

	content := subtle(title) + "\n\n" + label + "\n" + m.tp.progress.View() + "\n" + speedDownloadContent + "\n" + speedUploadContent

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
