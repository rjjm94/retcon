package main

import (
	"your-project/logger"
	"your-project/proxy"
	"your-project/recon"
	"your-project/ui"
	"your-project/utils"
)

func main() {
	// Create a new logger
	log, err := logger.NewLogger("app.log", logger.INFO)
	if err != nil {
		log.Fatal(err)
	}

	// Load the UI
	ui := ui.NewUI()
	err = ui.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	// Create a new proxy dialer
	dialer, err := proxy.NewHTTPProxy(ui.GetProxyAddr(), "", "")
	if err != nil {
		log.Fatal(err)
	}

	// Create a new recon
	r := recon.NewTCPRecon("tcp", ui.GetIPs()[0], dialer.Dial, ".", log, false)

	// Perform the scan
	err = r.Scan()
	if err != nil {
		log.Fatal(err)
	}

	// Grab banners
	utils.GrabBanners(ui.GetIPs(), ui.GetPort(), ui.GetBufferSize(), dialer.Dial, log)

	// Save results
	err = utils.SaveResults(ui.GetFilename(), []byte("results"), log)
	if err != nil {
		log.Fatal(err)
	}

	// Close the logger
	err = log.Close()
	if err != nil {
		log.Fatal(err)
	}
}
