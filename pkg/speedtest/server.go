package speedtest

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"strings"
	"time"
)

type Server struct {
	URL     string `xml:"url,attr" json:"url"`
	Lat     string `xml:"lat,attr" json:"lat"`
	Lon     string `xml:"lon,attr" json:"lon"`
	Name    string `xml:"name,attr" json:"name"`
	Country string `xml:"country,attr" json:"country"`
	Sponsor string `xml:"sponsor,attr" json:"sponsor"`
	ID      string `xml:"id,attr" json:"id"`
	URL2    string `xml:"url2,attr" json:"url_2"`
	Host    string `xml:"host,attr" json:"host"`

	Distance float64       `json:"distance"`
	Latency  time.Duration `json:"latency"`
	DLSpeed  float64       `json:"dl_speed"`
	ULSpeed  float64       `json:"ul_speed"`

	client *http.Client
}

func (s *Server) String() string {
	return fmt.Sprintf("<Server> id=[%4s] name=[%s] country=[%s] dis=[%fkm]", s.ID, s.Name, s.Country, s.Distance)
}

func (s *Server) GetPingLatency(ctx context.Context) error {
	pingURL := strings.Split(s.URL, "/upload.php")[0] + "/latency.txt"

	sTime := time.Now()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, pingURL, nil)
	if err != nil {
		return err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}

	fTime := time.Now()

	resp.Body.Close()

	// trans to  ms
	s.Latency = fTime.Sub(sTime)
	return nil
}

type ServerList []*Server

// Len finds length of servers. For sorting servers.
func (svrs ServerList) Len() int {
	return len(svrs)
}

// Swap swaps i-th and j-th. For sorting servers.
func (svrs ServerList) Swap(i, j int) {
	svrs[i], svrs[j] = svrs[j], svrs[i]
}

// Swap swaps i-th and j-th. For sorting servers.
func (svrs ServerList) Less(i, j int) bool {
	return svrs[i].Distance < svrs[j].Distance
}

func Distance(lat1 float64, lon1 float64, lat2 float64, lon2 float64) float64 {
	radius := 6378.137

	a1 := lat1 * math.Pi / 180.0
	b1 := lon1 * math.Pi / 180.0
	a2 := lat2 * math.Pi / 180.0
	b2 := lon2 * math.Pi / 180.0

	x := math.Sin(a1)*math.Sin(a2) + math.Cos(a1)*math.Cos(a2)*math.Cos(b2-b1)
	return radius * math.Acos(x)
}
