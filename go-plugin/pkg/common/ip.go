package common

import "net"

var IpV4 = ""

func init() {
	IpV4 = IPV4()
}

func IPV4() string {
	ipv4s := AllIPV4()
	if len(ipv4s) > 0 {
		return ipv4s[0]
	}
	return "no-hostname"
}

func AllIPV4() (ipv4s []string) {
	addresses, err := net.InterfaceAddrs()
	if err != nil {
		return
	}

	for _, addr := range addresses {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				ipv4 := ipNet.IP.String()
				if ipv4 == "127.0.0.1" || ipv4 == "localhost" {
					continue
				}
				ipv4s = append(ipv4s, ipv4)
			}
		}
	}
	return
}
