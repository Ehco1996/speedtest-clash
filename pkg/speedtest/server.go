package speedtest

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go.uber.org/atomic"
	"golang.org/x/sync/errgroup"
)

type Server struct {
	URL  string `xml:"url,attr" json:"url"`
	URL2 string `xml:"url2,attr" json:"url_2"`

	Lat     string `xml:"lat,attr" json:"lat"`
	Lon     string `xml:"lon,attr" json:"lon"`
	Name    string `xml:"name,attr" json:"name"`
	Country string `xml:"country,attr" json:"country"`
	Sponsor string `xml:"sponsor,attr" json:"sponsor"`
	ID      string `xml:"id,attr" json:"id"`
	Host    string `xml:"host,attr" json:"host"`

	Distance float64       `json:"distance"`
	Latency  time.Duration `json:"latency"`
	DLSpeed  float64       `json:"dl_speed"` // Mbps
	ULSpeed  float64       `json:"ul_speed"` // Mbps

	//one test metrics
	downLoadTestReceivedBytes *atomic.Int64
}

func (s *Server) String() string {
	return fmt.Sprintf("<Server> id=[%4s] name=[%s] country=[%s] dis=[%fkm]", s.ID, s.Name, s.Country, s.Distance)
}

func (s *Server) GetPingLatency(ctx context.Context, c *http.Client) error {
	pingURL := strings.Split(s.URL, "/upload.php")[0] + "/latency.txt"

	sTime := time.Now()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, pingURL, nil)
	if err != nil {
		return err
	}

	resp, err := c.Do(req)
	if err != nil {
		return err
	}

	fTime := time.Now()

	resp.Body.Close()

	// trans to  ms
	s.Latency = fTime.Sub(sTime)
	return nil
}

func (s *Server) DownLoadTest(ctx context.Context, c *http.Client, concurrency, downloadSize int, duration time.Duration) (chan Result, error) {
	sTime := time.Now()
	endTime := time.Now().Add(duration)
	ctx2, cancel := context.WithDeadline(ctx, endTime)
	// base download url
	dlURL := strings.Split(s.URL, "/upload.php")[0] + "/random" + strconv.Itoa(downloadSize) + "x" + strconv.Itoa(downloadSize) + ".jpg"
	log.Printf("start download test url: %s", dlURL)
	resChan := make(chan Result, 10)

	// init one test metrics
	s.downLoadTestReceivedBytes.Store(0)

	eg, ctx := errgroup.WithContext(ctx2)
	for idx := 0; idx < concurrency; idx++ {
		eg.Go(func() error {
			for {
				size, err := downloadRequest(ctx, c, dlURL)
				if err == nil {
					s.downLoadTestReceivedBytes.Add(size)
					totalBytes := s.downLoadTestReceivedBytes.Load()
					res := Result{
						CurrentSpeed: calcMbpsSpeed(totalBytes, sTime),
						TotalBytes:   totalBytes,
						Percent:      calcPercent(sTime, endTime),
					}
					resChan <- res
				} else {
					return err
				}
			}
		})
	}

	// start speed test thread
	go func() {
		if err := eg.Wait(); err != nil {
			log.Printf("DownLoadTest meet error=%v", err)
			// send final res
			resChan <- Result{Percent: 1.0}
		}
		close(resChan)
		cancel()
		s.DLSpeed = calcMbpsSpeed(s.downLoadTestReceivedBytes.Load(), sTime)
		// clear one test metrics
		s.downLoadTestReceivedBytes.Store(0)
	}()

	return resChan, nil
}

func downloadRequest(ctx context.Context, c *http.Client, dlURL string) (int64, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, dlURL, nil)
	if err != nil {
		return 0, err
	}
	resp, err := c.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	_, err = io.Copy(ioutil.Discard, resp.Body)
	return resp.ContentLength, err
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

type Result struct {
	CurrentSpeed float64
	TotalBytes   int64
	Percent      float64
}

func (r *Result) String() string {
	return fmt.Sprintf("cur speed=%f cur bytes=%d cur percent=%f", r.CurrentSpeed, r.TotalBytes, r.Percent)
}

func calcMbpsSpeed(bytes int64, startTime time.Time) float64 {
	fTime := time.Now()
	MB := float64(1000 * 1000)
	// MBps(MB per second)
	MBps := float64(bytes) / MB / fTime.Sub(startTime).Seconds()
	return math.Round(MBps * 8)
}

func calcPercent(startTime, endTime time.Time) float64 {
	return float64(time.Now().Sub(startTime)) / float64(endTime.Sub(startTime))
}
