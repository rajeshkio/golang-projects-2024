package main

import (
	"crypto/tls"

	"github.com/rk280392/customCICDTool/runServer"
)

func main() {
	port := "30480"
	certFile := "certs/fullchain.pem"
	privKey := "certs/privkey.pem"
	cert, _ := tls.LoadX509KeyPair(certFile, privKey)
	runServer.RunServer(port, cert)
}
