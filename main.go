package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"

	tls "github.com/refraction-networking/utls"
)

func main() {
	client := CreateHTTPClient()

	req, _ := http.NewRequest("GET", "https://http2.pro/api/v1", nil)
	req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.111 Safari/537.36")
	req.Header.Add("accept", "*/*")

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(bodyBytes))
}

func CreateHTTPClient() http.Client {
	// initialize the http client
	client := http.Client{
		Transport: &http.Transport{
			DialTLSContext: func(ctx context.Context, network, addr string) (net.Conn, error) {

				//initialize the tcp connection
				tcpConn, err := (&net.Dialer{}).DialContext(ctx, network, addr)
				if err != nil {
					return nil, err
				}

				//initialize the conifg for tls
				config := tls.Config{
					ServerName: strings.Split(addr, ":")[0], //set the server name with the provided addr
				}

				//initialize a tls connection with the underlying tcop connection and config
				tlsConn := tls.UClient(tcpConn, &config, tls.HelloChrome_83)

				//start the tls handshake between servers
				err = tlsConn.Handshake()
				if err != nil {
					return nil, fmt.Errorf("uTlsConn.Handshake() error: %w", err)
				}

				return tlsConn, nil
			},
			ForceAttemptHTTP2: true,
		},
	}

	return client
}
