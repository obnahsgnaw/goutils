package netutil

import (
	"net"
	"sort"
)

// GetIp Returns the IP address of the external access
func GetIp() string {
	var conn net.Conn
	var err error

	if conn, err = net.Dial("udp", "8.8.8.8:80"); err != nil {
		return "127.0.0.1"
	}
	defer func() { _ = conn.Close() }()

	return conn.LocalAddr().(*net.UDPAddr).IP.String()
}

// GetIps Returns the IP addresses of all NICs
func GetIps() []string {
	var ips []string
	var addresses []net.Addr
	var err error

	interfaces, _ := net.Interfaces()
	for _, i := range interfaces {
		// 过滤虚拟网卡
		if i.Flags&net.FlagUp == 0 || i.Flags&net.FlagLoopback != 0 || i.Flags&net.FlagPointToPoint != 0 {
			continue
		}
		if addresses, err = i.Addrs(); err != nil {
			continue
		}
		for _, addr := range addresses {
			ip := addr.(*net.IPNet).IP
			if ip.To4() != nil {
				ips = append(ips, ip.String())
			}
		}
	}

	if len(ips) == 0 {
		ips = []string{"127.0.0.1"}
	}

	sort.Strings(ips)

	return ips
}
