package ui

import (
	"context"
	"errors"

	"github.com/Ehco1996/clash-speed/pkg/clash"
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
			return s.GetPingLatency(ctx, m.c.GetInnerClient())
		})
	}
	return eg.Wait()
}
