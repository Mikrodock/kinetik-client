package tools

import (
	"errors"
	"net"
)

func GetIPAddress() (string, error) {
	netIf, err := net.InterfaceByName("eth0")
	if err != nil {
		return "", err
	}

	addrs, err := netIf.Addrs()
	if err != nil {
		return "", err
	}

	var ip net.IP

	for _, addr := range addrs {
		switch v := addr.(type) {
		case *net.IPNet:
			ip = v.IP
		case *net.IPAddr:
			ip = v.IP
		}
		if ip != nil && ip.To4() != nil {
			break
		}
	}

	if ip != nil {
		return ip.String(), nil
	} else {
		return "", errors.New("Cannot resolve IP of eth0 interface")
	}

}
