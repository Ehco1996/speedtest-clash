package speedtest

const (
	speedTestConfigUrl  = "https://www.speedtest.net/speedtest-config.php"
	speedTestServersUrl = "https://www.speedtest.net/api/js/servers?engine=js&limit=3"
)

//hardcode in speedtest.com kB
var dlSizes = [...]int{350, 500, 750, 1000, 1500, 2000, 2500, 3000, 3500, 4000}
var ulSizes = [...]int{100, 300, 500, 800, 1000, 1500, 2500, 3000, 3500, 4000}
