package server

import (
	"api/internal/config"
	"log"
	"net"
)

const localListenPort = ":8081"
const prodListenPort = ":80"

func getIpAddress() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

func getExternalAddress(env config.RunEnvironment) string {
	if env == config.Production {
		return getIpAddress().String() + prodListenPort
	}

	return "localhost" + localListenPort
}

func getListenAddr(env config.RunEnvironment) string {
	if env == config.Production {
		return prodListenPort
	}

	return localListenPort
}
