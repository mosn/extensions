package trace

import "net"

func GetOutboundIP() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	var tmpIp string
	var tmpErr error
	for _, interfaceName := range interfaces {
		if interfaceName.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if interfaceName.Flags&net.FlagLoopback != 0 {
			continue // loop back interface
		}
		addresses, err := interfaceName.Addrs()
		if err != nil {
			tmpIp = ""
			tmpErr = err
			continue
		}
		for _, addr := range addresses {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			tmpIp = ip.String()
			tmpErr = nil
			return tmpIp, tmpErr
		}
	}
	return tmpIp, tmpErr
}
