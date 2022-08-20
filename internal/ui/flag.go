package ui

var ()

//hardcode in speedtest.com kB
var dlSizes = [...]int{350, 500, 750, 1000, 1500, 2000, 2500, 3000, 3500, 4000}
var ulSizes = [...]int{100, 300, 500, 800, 1000, 1500, 2500, 3000, 3500, 4000}

var (
	DownLoadConcurrency = 1

	requestCount = 10 // one thread send 10 request

	// TODO adjust size by ping metrics
	downloadSize = dlSizes[2] // 750*750.jpg ~= 500k one request
	uploadSize   = ulSizes[2] // 500*500.jpg ~= 240k one request
)
