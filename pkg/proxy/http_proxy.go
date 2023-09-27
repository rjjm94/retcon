package proxy

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"
)

const (
	networkType = "tcp"
	maxRetries  = 3
)

// HTTPProxyDialer is a dialer for HTTP proxies.
type HTTPProxyDialer struct {
	proxyAddr   string
	username    string
	password    string
	isConnected bool
	timeout     time.Duration
	conn        net.Conn
}

// NewHTTPProxy creates a new HTTPProxyDialer.
func NewHTTPProxy(proxyAddr, username, password string) (*HTTPProxyDialer, error) {
	if proxyAddr == "" {
		return nil, errors.New("proxy address cannot be empty")
	}

	return &HTTPProxyDialer{
		proxyAddr:   proxyAddr,
		username:    username,
		password:    password,
		isConnected: false,
		timeout:     30 * time.Second,
	}, nil
}

// Dial connects to the address on the named network.
func (d *HTTPProxyDialer) Dial(network, addr string) (net.Conn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), d.timeout)
	defer cancel()

	var err error
	retries := 0

	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("connection attempt timed out")
		default:
			d.conn, err = net.DialTimeout(network, d.proxyAddr, d.timeout)
			if err == nil {
				d.isConnected = true
				return d.conn, nil
			}
			retries++
			if retries >= maxRetries {
				return nil, fmt.Errorf("failed to connect after %d attempts", maxRetries)
			}
			fmt.Println("Failed to connect, retrying...")
			time.Sleep(1 * time.Second)
		}
	}
}

// SetCredentials sets the username and password for the proxy.
func (d *HTTPProxyDialer) SetCredentials(username, password string) {
	d.username = username
	d.password = password
}

// GetProxyAddress returns the address of the proxy.
func (d *HTTPProxyDialer) GetProxyAddress() string {
	return d.proxyAddr
}

// IsConnected checks if the dialer is connected to the proxy.
func (d *HTTPProxyDialer) IsConnected() bool {
	return d.isConnected
}

// Reconnect reconnects to the proxy.
func (d *HTTPProxyDialer) Reconnect() error {
	if d.conn != nil {
		d.conn.Close()
	}
	conn, err := d.Dial(networkType, d.proxyAddr)
	if err != nil {
		return err
	}
	d.conn = conn
	return nil
}

// SetTimeout sets the timeout for the dialer.
func (d *HTTPProxyDialer) SetTimeout(timeout time.Duration) {
	d.timeout = timeout
}
