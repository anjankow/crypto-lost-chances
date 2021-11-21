package server

import (
	"api/internal/config"
	"log"
	"net"
)

const localAddress = ":8081"
const defaultProdPort = ":80"

func getIpAddress() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

func getAddress(env config.RunEnvironment) string {
	if env == config.Production {
		return getIpAddress().String() + defaultProdPort
	}

	return localAddress
}
