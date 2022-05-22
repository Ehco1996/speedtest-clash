package ui

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/Dreamacro/clash/constant"
	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/sync/errgroup"

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

	testPrecent float64
	quitting    bool

	c *speedtest.Client
}

func InitialModel() model {
	return model{proxyNodeList: []constant.Proxy{}}
}

func (m *model) FetchProxy(path string) error {
	cfg, err := clash.LoadConfig(path)
	if err != nil {
		return err
	}
	for _, p := range cfg.Proxies {
		m.proxyNodeList = append(m.proxyNodeList, p)
	}
	if len(m.proxyNodeList) == 0 {
		return errors.New("not have enough proxy nodes")
	}
	return nil
}

func (m *model) FetchTestServers() error {
	ctx := context.TODO()
	_, err := m.c.FetchUserInfo(ctx)
	if err != nil {
		return err
	}

	serverList, err := m.c.FetchServerList(ctx)
	if err != nil {
		return err
	}
	m.testServerList = serverList

	// fetch ping
	eg, ctx := errgroup.WithContext(ctx)
	for idx := range m.testServerList {
		s := m.testServerList[idx]
		eg.Go(func() error {
			return s.GetPingLatency(ctx)
		})
	}
	return eg.Wait()
}

func (m model) Init() tea.Cmd {
	return nil
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

	if m.selectedProxyNode == "" {
		return m.updateForProxyNode(msg)
	} else if m.selectedServer == "" {
		return m.updateForTestServer(msg)
	}
	return m.updateTestPrecent(msg)
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

func (m model) updateForTestServer(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			m.selectedServer = m.testServerList[m.serverIdx].Name
			return m, tickOneSecond()
		}
	}
	return m, nil
}

func (m model) updateTestPrecent(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case tickMsg:
		m.testPrecent += float64(0.2)
		if m.testPrecent >= float64(1) {
			m.testPrecent = float64(1)
		}
		return m, tickOneSecond()
	}
	return m, nil
}

type tickMsg struct{}

func tickOneSecond() tea.Cmd {
	return tea.Tick(time.Second/10, func(time.Time) tea.Msg {
		return tickMsg{}
	})
}
