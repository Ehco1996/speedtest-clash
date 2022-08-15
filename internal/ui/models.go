package ui

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/Dreamacro/clash/constant"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/Ehco1996/clash-speed/pkg/clash"
	"github.com/Ehco1996/clash-speed/pkg/speedtest"
)

type model struct {
	proxyIdx          int
	selectedProxyNode string
	proxyNodeList     []constant.Proxy

	serverIdx      int
	selectedServer string
	testServerList speedtest.ServerList

	quitting bool

	c *speedtest.Client

	// sub models
	progress progress.Model
	sp       modelSpeedTest
}

var _ tea.Model = (*model)(nil)

func InitialModel() model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return model{
		proxyNodeList: []constant.Proxy{},
		sp:            modelSpeedTest{spinner: s},
		progress:      progress.New(progress.WithDefaultGradient()),
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	log.Printf("updated %v", msg)
	// Make sure these keys always quit
	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()
		if k == "q" || k == "esc" || k == "ctrl+c" {
			m.quitting = true
			return m, tea.Quit
		}
	}

	if m.selectedProxyNode == "" {
		return m.updateForProxyNode(msg)
	} else if m.selectedServer == "" {
		return m.updateForSelectTestServer(msg)
	}
	return m.updateForTestProgress(msg)
}

func (m model) updateForProxyNode(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			if m.proxyIdx > 0 {
				m.proxyIdx--
			}
		case "down":
			if m.proxyIdx < len(m.proxyNodeList)-1 {
				m.proxyIdx++
			}
		case "enter":
			// user selected one proxy node, let init data with this proxy
			m.selectedProxyNode = m.proxyNodeList[m.proxyIdx].Name()
			// init inner speed test proxy client
			hc := &http.Client{Transport: clash.NewClashTransport(m.proxyNodeList[m.proxyIdx])}
			m.c = speedtest.NewClient(hc)
			// TODO: this is a slow io, maybe add some hints in ui
			if err := m.FetchTestServers(); err != nil {
				panic(err) // TODO: new a ui to show error
			}
		}
	}
	return m, nil
}

func (m model) updateForSelectTestServer(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			if m.serverIdx > 0 {
				m.serverIdx--
			}
		case "down":
			if m.serverIdx < len(m.testServerList)-1 {
				m.serverIdx++
			}
		case "enter":
			s := m.testServerList[m.serverIdx]
			m.selectedServer = s.Name
			// after select, we start draw progress
			// TODO handle err
			m.sp.resChan, _ = s.DownLoadTest(context.TODO(), m.c.GetInnerClient())
			return m, tickOnceForProgress()
		}
	}
	return m, nil
}

func (m model) updateForTestProgress(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case tickMsg:
		if m.progress.Percent() >= 1.0 {
			m.quitting = true
			return m, nil
		}
		// TODO this is a fake progress
		m.sp.currentRes = <-m.sp.resChan
		cmd := m.progress.IncrPercent(0.1)
		return m, tea.Batch(tickOnceForProgress(), cmd, m.sp.spinner.Tick)
		// FrameMsg is sent when the progress bar wants to animate itself
	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd
	case spinner.TickMsg:
		if !m.quitting {
			s, cmd := m.sp.spinner.Update(spinner.Tick())
			m.sp.spinner = s
			return m, cmd
		}
	}
	return m, nil
}

type tickMsg struct{}

func tickOnceForProgress() tea.Cmd {
	return tea.Tick(time.Second, func(time.Time) tea.Msg {
		return tickMsg{}
	})
}

type modelSpeedTest struct {
	spinner    spinner.Model
	resChan    chan speedtest.Result
	currentRes speedtest.Result
}
