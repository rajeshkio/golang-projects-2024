package main

import (
	"crypto/tls"

	"github.com/rk280392/customCICDTool/runServer"
)

func main() {
	port := "30480"
	cert, _ := tls.LoadX509KeyPair("certs/fullchain.pem", "certs/privkey.pem")
	runServer.RunServer(port, cert)
}
