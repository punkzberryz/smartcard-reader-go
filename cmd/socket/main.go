package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/punkzberryz/smartcard-reader-go/pkg/model"
	"github.com/punkzberryz/smartcard-reader-go/pkg/server"
	"github.com/punkzberryz/smartcard-reader-go/pkg/smc"
	"github.com/punkzberryz/smartcard-reader-go/pkg/util"
)

func main() {

	//load env
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}

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

	broadcast := make(chan model.Message)
	socketSrv := server.SocketServer{
		Server:    srv,
		Broadcast: broadcast,
	}

	fmt.Printf("Server is running on  http://localhost:%s/.\nPress Ctrl+C to exit.\n", config.Port)
	if err := socketSrv.RunServerWithWebSocket(); err != nil {
		log.Fatal("Server Error: ", err)
	}

	go func() {
		smc := smc.NewSmartCard()
		for {
			err := smc.StartDaemon(broadcast, srv.CardConfig)
			if err != nil {
				log.Printf("Error occurred in daemon process (%v), wait 2 seconds to retry or press Ctrl+C to exit.", err.Error())

				message := model.Message{
					Event: "smc-error",
					Payload: map[string]string{
						"message": fmt.Sprintf("Error occurred in daemon process, %v.", err.Error()),
					},
				}
				broadcast <- message

				time.Sleep(2 * time.Second)
			}
		}
	}()

	// Listen for syscall signals for process to interrupt/quit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	s := <-sig
	log.Printf("Received %v signal to shutdown.", s)

}
