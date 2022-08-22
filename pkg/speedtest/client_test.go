package speedtest

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestClient_FetchUserInfo(t *testing.T) {
	c := NewClient(http.DefaultClient)
	user, err := c.FetchUserInfo(context.TODO())
	require.NoError(t, err)
	require.NotEmpty(t, user)
	println(user.String())
}

func TestClient_FetchServerList(t *testing.T) {
	ctx := context.TODO()
	c := NewClient(http.DefaultClient)
	user, err := c.FetchUserInfo(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	println(user.String())

	serverList, err := c.FetchServerList(ctx)
	require.NoError(t, err)
	for idx := range serverList {
		println(serverList[idx].String())
	}
}

func TestClient_Server_PingTest(t *testing.T) {
	ctx := context.TODO()
	c := NewClient(http.DefaultClient)
	user, err := c.FetchUserInfo(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	println(user.String())

	serverList, err := c.FetchServerList(ctx)
	require.NoError(t, err)
	for idx := range serverList {
		s := serverList[idx]
		println(s.String())
		require.NoError(t, s.GetPingLatency(ctx, c.inner))
	}
}

func TestClient_Server_DownLoad(t *testing.T) {
	ctx := context.TODO()
	c := NewClient(http.DefaultClient)
	user, err := c.FetchUserInfo(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	println(user.String())

	serverList, err := c.FetchServerList(ctx)
	require.NoError(t, err)
	require.Greater(t, len(serverList), 0)
	s := serverList[0]

	ch, err := s.DownLoadTest(ctx, c.GetInnerClient(), 1, 500, time.Second*2)
	require.NoError(t, err)

	for res := range ch {
		fmt.Printf("current download speed is %.2f mpbs total bytes %d mb \n", res.CurrentSpeed, res.TotalBytes)
	}

	require.Greater(t, s.DLSpeed, float64(0))
	fmt.Printf("download speed is %.2f mpbs", s.DLSpeed)
}
