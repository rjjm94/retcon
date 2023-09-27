package recon

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

type UDPRecon struct {
	network         string
	ip              string
	dialer          func(network, addr string) (net.Conn, error)
	scanDir         string
	logger          *log.Logger
	timeout         time.Duration
	portRange       [2]int
	concurrentScans int
	scanDelay       time.Duration
	debug           bool
}

func NewUDPRecon(network, ip string, dialer func(network, addr string) (net.Conn, error), scanDir string, logger *log.Logger, debug bool) *UDPRecon {
	return &UDPRecon{
		network:         network,
		ip:              ip,
		dialer:          dialer,
		scanDir:         scanDir,
		logger:          logger,
		timeout:         30 * time.Second,
		portRange:       [2]int{1, 65535}, // default port range
		concurrentScans: 100,              // default number of concurrent scans
		scanDelay:       0,                // default scan delay
		debug:           debug,
	}
}

func (r *UDPRecon) SetTimeout(timeout time.Duration) {
	r.timeout = timeout
}

func (r *UDPRecon) SetPortRange(start, end int) {
	r.portRange = [2]int{start, end}
}

func (r *UDPRecon) SetConcurrentScans(concurrentScans int) {
	r.concurrentScans = concurrentScans
}

func (r *UDPRecon) SetScanDelay(scanDelay time.Duration) {
	r.scanDelay = scanDelay
}

func (r *UDPRecon) Scan() error {
	var wg sync.WaitGroup
	sem := make(chan struct{}, r.concurrentScans)

	for i := r.portRange[0]; i <= r.portRange[1]; i++ {
		wg.Add(1)
		sem <- struct{}{} // acquire a token
		go func(port int) {
			defer wg.Done()
			defer func() { <-sem }() // release the token

			addr := fmt.Sprintf("%s:%d", r.ip, port)
			conn, err := r.dialer(r.network, addr)
			if err != nil {
				r.DebugPrint(fmt.Sprintf("Failed to connect to %s", addr))
				return
			}
			r.DebugPrint(fmt.Sprintf("Connected to %s", addr))
			conn.Close()
		}(i)
	}

	wg.Wait()
	return nil
}

func (r *UDPRecon) DebugPrint(message string) {
	if r.debug && r.logger != nil {
		r.logger.Println(message)
	}
}
