package utils

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"sync"
)

const (
	InfoLevel     = "INFO"
	ErrorLevel    = "ERROR"
	maxRetries    = 3
	maxGoroutines = 1000
)

type Logger interface {
	Log(level string, message string)
}

// SaveResults saves the given data into a file with the given filename.
// It logs any errors that occur during the process.
func SaveResults(filename string, data []byte, logger Logger) error {
	file, err := os.Create(filename)
	if err != nil {
		logger.Log(ErrorLevel, fmt.Sprintf("Failed to create file %s: %v", filename, err))
		return fmt.Errorf("failed to create file %s: %w", filename, err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	_, err = writer.Write(data)
	if err != nil {
		logger.Log(ErrorLevel, fmt.Sprintf("Failed to write data to file %s: %v", filename, err))
		return fmt.Errorf("failed to write data to file %s: %w", filename, err)
	}

	err = writer.Flush()
	if err != nil {
		logger.Log(ErrorLevel, fmt.Sprintf("Failed to flush data to file %s: %v", filename, err))
		return fmt.Errorf("failed to flush data to file %s: %w", filename, err)
	}

	logger.Log(InfoLevel, fmt.Sprintf("Data written to file %s", filename))
	return nil
}

// GrabBanner tries to connect to the given IP and port and reads the banner.
// It retries up to maxRetries times if the first attempt fails.
func GrabBanner(ctx context.Context, ip string, port int, bufferSize int, dialer func(ctx context.Context, network, addr string) (net.Conn, error), logger Logger) (string, error) {
	addr := fmt.Sprintf("%s:%d", ip, port)
	var banner string
	var err error

	for i := 0; i < maxRetries; i++ {
		conn, err := dialer(ctx, "tcp", addr)
		if err != nil {
			logger.Log(ErrorLevel, fmt.Sprintf("Failed to connect to %s: %v", addr, err))
			continue
		}
		defer conn.Close()

		buffer := make([]byte, bufferSize)
		n, err := conn.Read(buffer)
		if err != nil {
			if errors.Is(err, io.EOF) {
				logger.Log(InfoLevel, fmt.Sprintf("Connection to %s closed by server", addr))
			} else {
				logger.Log(ErrorLevel, fmt.Sprintf("Failed to read from %s: %v", addr, err))
				continue
			}
		}

		banner = string(buffer[:n])
		logger.Log(InfoLevel, fmt.Sprintf("Grabbed banner from %s: %s", addr, banner))
		break
	}

	if err != nil {
		return "", fmt.Errorf("failed to grab banner from %s after %d attempts: %w", addr, maxRetries, err)
	}

	return banner, nil
}

// GrabBanners grabs banners from the given IPs and port concurrently.
// It limits the number of concurrent goroutines to maxGoroutines.
func GrabBanners(ctx context.Context, ips []string, port int, bufferSize int, dialer func(ctx context.Context, network, addr string) (net.Conn, error), logger Logger) {
	var wg sync.WaitGroup
	sem := make(chan struct{}, maxGoroutines)

	for _, ip := range ips {
		wg.Add(1)
		sem <- struct{}{} // acquire a token
		go func(ip string) {
			defer wg.Done()
			defer func() { <-sem }() // release the token

			banner, err := GrabBanner(ctx, ip, port, bufferSize, dialer, logger)
			if err != nil {
				logger.Log(ErrorLevel, fmt.Sprintf("Failed to grab banner from %s: %v", ip, err))
				return
			}
			logger.Log(InfoLevel, fmt.Sprintf("Grabbed banner from %s: %s", ip, banner))
		}(ip)
	}

	wg.Wait()
}
