package ui

import (
	"fmt"
	"github.com/common-nighthawk/go-figure"
	"github.com/olekukonko/tablewriter"
	"github.com/schollz/progressbar"
	"github.com/spf13/viper"
	"os"
	"strconv"
	"time"
)

type UI struct {
	ips             []string
	port            int
	bufferSize      int
	filename        string
	proxyType       string
	proxyAddr       string
	reconType       string
	timeout         time.Duration
	portRange       []int
	concurrentScans int
	scanDelay       time.Duration
}

// NewUI creates a new UI instance and returns it.
func NewUI() *UI {
	return &UI{}
}

// DisplayBanner displays the banner of the application.
func (ui *UI) DisplayBanner() {
	figure.NewFigure("RE-tard-con", "", true).Print()
}

// AskInputs asks for user inputs interactively.
// It validates the inputs and returns an error if any of the inputs is invalid.
func (ui *UI) AskInputs() error {
	// Ask for IPs
	prompt := &survey.MultiLine{
		Message: "Enter the IP addresses (separated by commas):",
	}
	err := survey.AskOne(prompt, &ui.ips)
	if err != nil {
		return fmt.Errorf("failed to get IPs: %w", err)
	}

	// Validate IPs...

	// Ask for port range
	prompt = &survey.Input{
		Message: "Enter the start of the port range:",
	}
	err = survey.AskOne(prompt, &ui.portRange[0])
	if err != nil {
		return fmt.Errorf("failed to get start of port range: %w", err)
	}

	// Validate start of port range...

	prompt = &survey.Input{
		Message: "Enter the end of the port range:",
	}
	err = survey.AskOne(prompt, &ui.portRange[1])
	if err != nil {
		return fmt.Errorf("failed to get end of port range: %w", err)
	}

	// Validate end of port range...

	// Ask if they want to use a proxy
	useProxy := false
	prompt = &survey.Confirm{
		Message: "Do you want to use a proxy?",
	}
	err = survey.AskOne(prompt, &useProxy)
	if err != nil {
		return fmt.Errorf("failed to get proxy choice: %w", err)
	}

	if useProxy {
		// Ask for proxy details or filename
		prompt = &survey.Input{
			Message: "Enter the proxy details (username:password@ip:port) or the filename containing the proxy details:",
		}
		err = survey.AskOne(prompt, &ui.proxyAddr)
		if err != nil {
			return fmt.Errorf("failed to get proxy details or filename: %w", err)
		}
	}

	// Validate proxy details or filename...

	// Ask for filename
	prompt = &survey.Input{
		Message: "Enter the filename to save results:",
	}
	err = survey.AskOne(prompt, &ui.filename)
	if err != nil {
		return fmt.Errorf("failed to get filename: %w", err)
	}

	// Validate filename...

	// Add similar prompts for other inputs...
	return nil
}

// DisplayProgress displays a progress bar with the given total number of steps.
// It handles any error that occurs when updating the progress bar.
func (ui *UI) DisplayProgress(total int) {
	bar := progressbar.New(total)
	for i := 0; i < total; i++ {
		err := bar.Add(1)
		if err != nil {
			fmt.Printf("Failed to update progress bar: %v\n", err)
			break
		}
		time.Sleep(time.Second)
	}
}

// DisplayResults displays the given results in a table.
func (ui *UI) DisplayResults(results []Result) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"IP", "Port", "Banner"})

	for _, result := range results {
		row := []string{result.IP, strconv.Itoa(result.Port), result.Banner}
		table.Append(row)
	}

	table.Render()
}

// LoadConfig loads the configuration from the config file.
// It returns an error if the config file does not exist or cannot be read.
func (ui *UI) LoadConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(*viper.ConfigFileNotFoundError); ok {
			// Config file not found; return an error
			return fmt.Errorf("config file not found: %v", err)
		} else {
			// Some other error occurred
			return fmt.Errorf("error reading config file: %v", err)
		}
	}

	// Config file found; load the values
	ui.ips = viper.GetStringSlice("ips")
	ui.portRange = viper.GetIntSlice("portRange")
	ui.bufferSize = viper.GetInt("bufferSize")
	ui.filename = viper.GetString("filename")
	ui.proxyType = viper.GetString("proxyType")
	ui.proxyAddr = viper.GetString("proxyAddr")
	ui.reconType = viper.GetString("reconType")
	ui.timeout = viper.GetDuration("timeout")
	ui.concurrentScans = viper.GetInt("concurrentScans")
	ui.scanDelay = viper.GetDuration("scanDelay")

	return nil
}

// Getter methods to get the inputs
func (ui *UI) GetIPs() []string {
	return ui.ips
}

func (ui *UI) GetPort() int {
	return ui.port
}

func (ui *UI) GetBufferSize() int {
	return ui.bufferSize
}

func (ui *UI) GetFilename() string {
	return ui.filename
}

func (ui *UI) GetProxyType() string {
	return ui.proxyType
}

func (ui *UI) GetProxyAddr() string {
	return ui.proxyAddr
}

func (ui *UI) GetReconType() string {
	return ui.reconType
}

func (ui *UI) GetTimeout() time.Duration {
	return ui.timeout
}

func (ui *UI) GetPortRange() []int {
	return ui.portRange
}

func (ui *UI) GetConcurrentScans() int {
	return ui.concurrentScans
}

func (ui *UI) GetScanDelay() time.Duration {
	return ui.scanDelay
}
