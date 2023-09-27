package proxy

import (
	"errors"
	"fmt"
	"golang.org/x/net/proxy"
	"net"
	"time"
)

const (
	defaultNetworkType = "tcp"
	maxRetries         = 3
)

// SOCKS5ProxyDialer is a struct that holds the necessary details for a SOCKS5 proxy connection.
type SOCKS5ProxyDialer struct {
	proxyAddr   string
	username    string
	password    string
	timeout     time.Duration
	isConnected bool
	conn        net.Conn
}

// NewSOCKS5Proxy creates a new SOCKS5 proxy dialer with the given proxy address, username, and password.
func NewSOCKS5Proxy(proxyAddr, username, password string) (*SOCKS5ProxyDialer, error) {
	if proxyAddr == "" {
		return nil, errors.New("proxy address cannot be empty")
	}

	return &SOCKS5ProxyDialer{
		proxyAddr:   proxyAddr,
		username:    username,
		password:    password,
		timeout:     30 * time.Second,
		isConnected: false,
	}, nil
}

// SetCredentials sets the username and password for the proxy connection.
func (d *SOCKS5ProxyDialer) SetCredentials(username, password string) {
	d.username = username
	d.password = password
}

// GetProxyAddress returns the proxy address.
func (d *SOCKS5ProxyDialer) GetProxyAddress() string {
	return d.proxyAddr
}

// IsConnected returns whether the proxy connection is established.
func (d *SOCKS5ProxyDialer) IsConnected() bool {
	return d.isConnected
}

// Dial establishes a connection to the given network address via the proxy.
func (d *SOCKS5ProxyDialer) Dial(network, addr string) (net.Conn, error) {
	retries := 0
	var conn net.Conn

	for retries < maxRetries {
		dialer, err := proxy.SOCKS5(network, d.proxyAddr, &proxy.Auth{User: d.username, Password: d.password}, proxy.Direct)
		if err != nil {
			retries++
			time.Sleep(time.Second)
			continue
		}
		conn, err = dialer.Dial(network, addr)
		if err != nil {
			retries++
			time.Sleep(time.Second)
			continue
		}
		d.isConnected = true
		d.conn = conn
		return conn, nil
	}
	return nil, fmt.Errorf("failed to connect after %d attempts", maxRetries)
}

// Reconnect re-establishes the proxy connection.
func (d *SOCKS5ProxyDialer) Reconnect() error {
	if d.conn != nil {
		d.conn.Close()
	}
	conn, err := d.Dial(defaultNetworkType, d.proxyAddr)
	if err != nil {
		return err
	}
	d.conn = conn
	d.isConnected = false
	return nil
}

// SetTimeout sets the timeout duration for the proxy connection.
func (d *SOCKS5ProxyDialer) SetTimeout(timeout time.Duration) {
	d.timeout = timeout
}
