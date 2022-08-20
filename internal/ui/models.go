package ui

import (
	"context"
	"net/http"

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
	ps modelProxyServer
	ts modelTestServer
	tp modelTestProgress
}

var _ tea.Model = (*model)(nil)

func InitialModel() model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return model{
		ps: modelProxyServer{},
		ts: modelTestServer{},
		tp: modelTestProgress{spinner: s, progress: progress.New(progress.WithDefaultGradient())},
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	if msg, ok := msg.(tea.WindowSizeMsg); ok {
		h, v := proxyNodeListStyle.GetFrameSize()
		m.ps.uiList.SetSize(msg.Width-h, msg.Height-v)
	}

	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()
		if k == "q" || k == "esc" || k == "ctrl+c" {
			m.quitting = true
			return m, tea.Quit
		}
	}

	if m.ps.selectedProxyNode == "" {
		return m.updateForSelectProxyNode(msg)
	} else if m.ts.selectedServer == "" {
		return m.updateForSelectTestServer(msg)
	}
	return m.updateForTestProgress(msg)
}

func (m model) updateForSelectProxyNode(msg tea.Msg) (tea.Model, tea.Cmd) {
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
	var cmd tea.Cmd
	m.ps.uiList, cmd = m.ps.uiList.Update(msg)
	return m, cmd
}

func (m model) updateForSelectTestServer(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			if m.ts.serverIdx > 0 {
				m.ts.serverIdx--
			}
		case "down":
			if m.ts.serverIdx < len(m.ts.testServerList)-1 {
				m.ts.serverIdx++
			}
		case "enter":
			s := m.ts.testServerList[m.ts.serverIdx]
			m.ts.selectedServer = s.Name
			// TODO handle err
			m.tp.ch, _ = s.DownLoadTest(context.TODO(), m.c.GetInnerClient(), DownLoadConcurrency, requestCount, downloadSize)
			return m, receiveTestResOnce(m.tp.ch)
		}
	}
	return m, nil
}

func (m model) updateForTestProgress(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case speedTestMsg:

		if m.tp.progress.Percent() >= 1.0 {
			m.quitting = true
			return m, nil
		}
		res := *msg.res
		m.tp.currentRes = res
		cmd := m.tp.progress.SetPercent(res.Percent)
		return m, tea.Batch(cmd, m.tp.spinner.Tick, receiveTestResOnce(m.tp.ch))

	case progress.FrameMsg:
		// FrameMsg is sent when the progress bar wants to animate itself
		progressModel, cmd := m.tp.progress.Update(msg)
		m.tp.progress = progressModel.(progress.Model)
		return m, cmd
	case spinner.TickMsg:
		if !m.quitting {
			s, cmd := m.tp.spinner.Update(spinner.Tick())
			m.tp.spinner = s
			return m, cmd
		}
	}
	return m, nil
}

// used for user select proxy server
type modelProxyServer struct {
	proxyIdx          int
	selectedProxyNode string
	proxyNodeList     []constant.Proxy

	uiList list.Model
}

// for ui list
type proxyItem struct {
	constant.Proxy
}

func (i proxyItem) FilterValue() string { return i.Name() }
func (i proxyItem) Title() string       { return i.Name() }
func (i proxyItem) Description() string {
	return i.Type().String() + "\t" + i.Addr()
}

// used for user select test server
type modelTestServer struct {
	serverIdx      int
	selectedServer string
	testServerList speedtest.ServerList
}

// used for trigger download/upload test
type modelTestProgress struct {
	progress   progress.Model
	currentRes speedtest.Result
	spinner    spinner.Model

	ch chan speedtest.Result
}

type speedTestMsg struct {
	res *speedtest.Result
}

func receiveTestResOnce(ch chan speedtest.Result) tea.Cmd {
	return func() tea.Msg {
		res := <-ch
		return speedTestMsg{res: &res}
	}
}
