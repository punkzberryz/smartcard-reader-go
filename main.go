package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "embed"

	"github.com/punkzberryz/smartcard-reader-go/pkg/server"
	"github.com/punkzberryz/smartcard-reader-go/pkg/smc"
	"github.com/punkzberryz/smartcard-reader-go/pkg/util"
)

//go:embed env.txt
var env []byte

func main() {

	configStr := string(env)
	config, err := util.LoadConfigFromString(configStr)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	log.Printf("Config: %+v", config)

	serverCfg := server.ServerConfig{
		Port:   config.Port,
		ApiKey: config.ApiKey,
	}

	srv := server.Server{
		ServerConfig: serverCfg,
		Card:         smc.NewSmartCard(),
		CardConfig: &smc.Options{
			ShowFaceImage: config.ShowImage,
			ShowNhsoData:  config.ShowNhso,
			ShowLaserData: config.ShowLaser,
		},
	}

	go func() {
		if err := srv.RunServer(); err != nil {
			fmt.Println("Server Error: ", err)
		}
	}()

	fmt.Printf("Server is running on  http://localhost:%s/. Press Ctrl+C to exit.\n", config.Port)

	// Set up channel to listen for interrupt signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	//Block until a signal is received
	<-sigChan
	fmt.Println("\nShutting down gracefully...")
	//wait for 2 seconds before shutting down
	time.Sleep(2 * time.Second)
	// Perform any cleanup here if necessary
	fmt.Println("Server stopped")
}

/*
	Endpoint: http://localhost:8080/read with optional query parameters of
		show-face-image = true / false
		show-nhso-data = true / false
		show-laser-data = true / false

		And require header of
		Authorization : Bearer <API_KEY>

	Example:
	curl -X GET "http://localhost:8080/read?show-face-image=true&show-nhso-data=true&show-laser-data=true" \
     -H "Authorization: Bearer <API_KEY>"


*/
