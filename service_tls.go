package go_service_util

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
)

func BuildTlsConfig(etcdClientKeyFile string, etcdClientCertFile string, caCertFile string) (*tls.Config, error) {
	etcdClientKeyPair, err := tls.LoadX509KeyPair(etcdClientCertFile, etcdClientKeyFile)
	if err != nil {
		return nil, err
	}

	caCertData, err := ioutil.ReadFile(caCertFile)
	if err != nil {
		return nil, err
	}
	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(caCertData)

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{etcdClientKeyPair},
		RootCAs:      certPool,
	}

	return tlsConfig, nil
}
