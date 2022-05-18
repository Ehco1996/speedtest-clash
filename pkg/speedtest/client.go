package speedtest

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"net/http"
	"sort"
	"strconv"
)

type Client struct {
	inner *http.Client

	user *User
}

func NewClient(c *http.Client) *Client {
	return &Client{inner: c}
}

func (c *Client) CurrentUser() *User {
	return c.user
}

func (c *Client) FetchUserInfo(ctx context.Context) (*User, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, speedTestConfigUrl, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.inner.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	decoder := xml.NewDecoder(resp.Body)
	var users Users
	if err := decoder.Decode(&users); err != nil {
		return nil, err
	}
	if len(users.Users) == 0 {
		return nil, errors.New("failed to fetch user information")
	}
	c.user = &users.Users[0]
	return &users.Users[0], nil
}

func (c *Client) FetchServerList(ctx context.Context) (ServerList, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, speedTestServersUrl, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.inner.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var serverList ServerList

	decoder := json.NewDecoder(resp.Body)

	if err := decoder.Decode(&serverList); err != nil {
		return serverList, err
	}
	if len(serverList) == 0 {
		return nil, errors.New("failed to fetch servers")
	}

	// Calculate distance
	for _, server := range serverList {
		sLat, _ := strconv.ParseFloat(server.Lat, 64)
		sLon, _ := strconv.ParseFloat(server.Lon, 64)
		uLat, _ := strconv.ParseFloat(c.user.Lat, 64)
		uLon, _ := strconv.ParseFloat(c.user.Lon, 64)
		server.Distance = distance(sLat, sLon, uLat, uLon)
	}

	// Sort by distance
	sort.Sort(serverList)

	return serverList, nil
}
