package tunio

import (
	"bufio"
	"errors"
	"net"
	"net/http"
	"time"
)

var (
	timeout   = time.Second * 60
	keepAlive = time.Second * 60
)

// NewLanternDialer creates the dialer function used by tunio to connect to
// external sites.
func NewLanternDialer(proxyAddr string, dial dialer) dialer {
	if dial == nil {
		// A simple transparent dialer.
		d := net.Dialer{
			Timeout:   timeout,
			KeepAlive: keepAlive,
		}
		dial = d.Dial
	}
	// A CONNECT proxy dialer that works with a Lantern client.
	return func(proto, addr string) (net.Conn, error) {
		conn, err := dial("tcp", proxyAddr)
		if err != nil {
			return nil, err
		}

		req, err := http.NewRequest("CONNECT", addr, nil)
		if err != nil {
			return nil, err
		}

		req.Host = addr
		if err := req.Write(conn); err != nil {
			return nil, err
		}

		r := bufio.NewReader(conn)
		resp, err := http.ReadResponse(r, req)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
			return conn, nil
		}

		return nil, errors.New("Could not CONNECT to Lantern proxy.")
	}
}
