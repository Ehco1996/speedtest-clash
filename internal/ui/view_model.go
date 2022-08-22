package ui

import (
	"context"
	"errors"

	"github.com/Ehco1996/speedtest-clash/pkg/clash"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/sync/errgroup"
)

// TODO: maybe move all io operations here
func (m model) Init() tea.Cmd {
	return nil
}

func (m *model) FetchProxy(path string) error {
	cfg, err := clash.LoadConfig(path)
	if err != nil {
		return err
	}
	for _, p := range cfg.Proxies {
		m.ps.proxyNodeList = append(m.ps.proxyNodeList, p)
	}
	if len(m.ps.proxyNodeList) == 0 {
		return errors.New("not have enough proxy nodes")
	}

	// set proxy item
	items := []list.Item{}
	for idx := range m.ps.proxyNodeList {
		items = append(items, proxyItem{m.ps.proxyNodeList[idx]})
	}

	m.ps.uiList = list.New(items, list.NewDefaultDelegate(), 0, 0)
	m.ps.uiList.Title = "Choose your proxy node..."
	m.ps.uiList.SetShowFilter(false)
	m.ps.uiList.SetShowTitle(true)
	m.ps.uiList.SetShowHelp(true)
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
	m.ts.testServerList = serverList

	// fetch ping
	eg, ctx := errgroup.WithContext(ctx)
	for idx := range m.ts.testServerList {
		s := m.ts.testServerList[idx]
		eg.Go(func() error {
			return s.GetPingLatency(ctx, m.c.GetInnerClient())
		})
	}
	return eg.Wait()
}
