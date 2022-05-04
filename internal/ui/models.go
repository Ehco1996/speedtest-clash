package ui

import (
	"errors"
	"time"

	"github.com/Dreamacro/clash/constant"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/Ehco1996/clash-speed/pkg/clash"
)

type model struct {
	proxyIdx          int
	selectedProxyNode string
	proxyNodeList     []constant.Proxy

	serverIdx      int
	selectedServer string
	testServerList []string

	testPrecent float64
	quitting    bool
}

func InitialModel() model {
	// for now everything is fake
	return model{
		proxyNodeList:  []constant.Proxy{},
		testServerList: []string{"江苏联通", "上海电信"},
	}
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
			m.selectedProxyNode = m.proxyNodeList[m.proxyIdx].Name()
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
			m.selectedServer = m.testServerList[m.serverIdx]
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
