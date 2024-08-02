package server

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	smc "github.com/punkzberryz/smartcard-reader-go/pkg/smc"
)

func (s *Server) handleHome(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("I'm okay"))
}

func (s *Server) auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//check api key
		authHeader := r.Header.Get("authorization")
		authFields := strings.Fields(authHeader)
		if len(authFields) != 2 || authFields[0] != "Bearer" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		apiKey := authFields[1]

		if apiKey != s.ServerConfig.ApiKey {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		//Proceed to next middleware or handler
		next.ServeHTTP(w, r)
	})
}

type readResponse struct {
	Error error
	Data  []byte
}

func (s *Server) handleRead(w http.ResponseWriter, r *http.Request) {
	//create a context with a 15-second timeout
	ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
	defer cancel()
	//create a channel to signal completion
	showFaceImage := r.URL.Query().Get("show-face-image") == "true"
	showNhsoData := r.URL.Query().Get("show-nhso-data") == "true"
	showLaserData := r.URL.Query().Get("show-laser-data") == "true"

	//start reading the card in a goroutine
	done := make(chan readResponse)
	go func(showFaceImage, showNhsoData, showLaserData bool, results chan<- readResponse) {
		// time.Sleep(15 * time.Second) // Simulating work that takes longer than 15 seconds
		resp := readResponse{Data: nil, Error: nil}
		data, err := s.Card.Read(nil, &smc.Options{
			ShowFaceImage: showFaceImage,
			ShowNhsoData:  showNhsoData,
			ShowLaserData: showLaserData,
		})
		if err != nil {
			//send back error
			resp.Error = err
			results <- resp
			return
		}
		dataJSON, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			//send back error
			resp.Error = err
			results <- resp
			return
		}
		//send back data
		resp.Data = dataJSON
		results <- resp

	}(showFaceImage, showNhsoData, showLaserData, done)

	//wait for either the work to complete or a timeout
	select {
	case resp := <-done:
		{
			if resp.Error != nil {
				if resp.Error == smc.ErrNoSmartCardReader {
					http.Error(w, "card reader not found", http.StatusInternalServerError)
					return
				}
				http.Error(w, resp.Error.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(200)
			w.Write(resp.Data)
		}

	case <-ctx.Done():
		http.Error(w, "request timed out", http.StatusRequestTimeout)
	}
}
