package main

import (
	"fmt"
	"github.com/obnahsgnaw/goutils/netutil"
	"strings"
)

func main() {
	ip := netutil.GetIp()
	fmt.Println(ip)

	ips := netutil.GetIps()
	fmt.Println(strings.Join(ips, ","))
}
