package speedtest

import (
	"fmt"
)

type User struct {
	IP  string `xml:"ip,attr"`
	Lat string `xml:"lat,attr"`
	Lon string `xml:"lon,attr"`
	Isp string `xml:"isp,attr"`
}

type Users struct {
	Users []User `xml:"client"`
}

func (u *User) String() string {
	return fmt.Sprintf("<User> ip=[%s] isp=[%s] lat=[%s] lon=[%s]", u.IP, u.Isp, u.Lat, u.Lon)
}
