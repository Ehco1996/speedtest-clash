package ui

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/Dreamacro/clash/constant"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/Ehco1996/clash-speed/pkg/clash"
	"github.com/Ehco1996/clash-speed/pkg/speedtest"
)

type model struct {
	quitting bool
	c        *speedtest.Client

	// sub models
	ps  modelProxyServer
	str modelSpeedTestRes
	sts modelSpeedTestServer
}

var _ tea.Model = (*model)(nil)

func InitialModel() model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return model{
		ps:  modelProxyServer{},
		sts: modelSpeedTestServer{},
		str: modelSpeedTestRes{spinner: s, progress: progress.New(progress.WithDefaultGradient())},
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Make sure these keys always quit
	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()
		if k == "q" || k == "esc" || k == "ctrl+c" {
			m.quitting = true
			return m, tea.Quit
		}
	}

	if m.ps.selectedProxyNode == "" {
		return m.updateForProxyNode(msg)
	} else if m.sts.selectedServer == "" {
		return m.updateForSelectTestServer(msg)
	}
	return m.updateForTestProgress(msg)
}

func (m model) updateForProxyNode(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			if m.ps.proxyIdx > 0 {
				m.ps.proxyIdx--
			}
		case "down":
			if m.ps.proxyIdx < len(m.ps.proxyNodeList)-1 {
				m.ps.proxyIdx++
			}
		case "enter":
			// user selected one proxy node, let init data with this proxy
			m.ps.selectedProxyNode = m.ps.proxyNodeList[m.ps.proxyIdx].Name()
			// init inner speed test proxy client
			hc := &http.Client{Transport: clash.NewClashTransport(m.ps.proxyNodeList[m.ps.proxyIdx])}
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
			if m.sts.serverIdx > 0 {
				m.sts.serverIdx--
			}
		case "down":
			if m.sts.serverIdx < len(m.sts.testServerList)-1 {
				m.sts.serverIdx++
			}
		case "enter":
			s := m.sts.testServerList[m.sts.serverIdx]
			m.sts.selectedServer = s.Name
			// TODO handle err
			m.str.resChan, _ = s.DownLoadTest(context.TODO(), m.c.GetInnerClient(), DownLoadConcurrency, requestCount, downloadSize)
			return m, tickOnceForProgress()
		}
	}
	return m, nil
}

func (m model) updateForTestProgress(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case speedTestTickMsg:
		if m.str.progress.Percent() >= 1.0 {
			m.quitting = true
			return m, nil
		}
		res := <-m.str.resChan
		log.Printf("tick once res=%s ", res.String())
		m.str.currentRes = res
		cmd := m.str.progress.SetPercent(res.Percent)
		return m, tea.Batch(cmd, tickOnceForProgress(), m.str.spinner.Tick)
		// FrameMsg is sent when the progress bar wants to animate itself
	case progress.FrameMsg:
		progressModel, cmd := m.str.progress.Update(msg)
		m.str.progress = progressModel.(progress.Model)
		return m, cmd
	case spinner.TickMsg:
		if !m.quitting {
			s, cmd := m.str.spinner.Update(spinner.Tick())
			m.str.spinner = s
			return m, cmd
		}
	}
	return m, nil
}

type speedTestTickMsg struct{}

func tickOnceForProgress() tea.Cmd {
	return tea.Tick(time.Microsecond, func(time.Time) tea.Msg {
		return speedTestTickMsg{}
	})
}

// used for trigger download/upload test
type modelSpeedTestRes struct {
	progress   progress.Model
	spinner    spinner.Model
	resChan    chan speedtest.Result
	currentRes speedtest.Result
}

// used for user select test server
type modelSpeedTestServer struct {
	serverIdx      int
	selectedServer string
	testServerList speedtest.ServerList
}

// used for user select proxy server
type modelProxyServer struct {
	proxyIdx          int
	selectedProxyNode string
	uiList            list.Model
	proxyNodeList     []constant.Proxy
}
