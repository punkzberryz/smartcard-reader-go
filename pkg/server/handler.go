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
	return CORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//check api key
		authHeader := r.Header.Get("authorization")
		authFields := strings.Fields(authHeader)
		if len(authFields) != 2 || authFields[0] != "Bearer" {
			handleError(w, http.StatusUnauthorized, "unauthorized")
			return
		}
		apiKey := authFields[1]

		if apiKey != s.ServerConfig.ApiKey {
			handleError(w, http.StatusUnauthorized, "unauthorized")
			return
		}
		//Proceed to next middleware or handler
		next.ServeHTTP(w, r)
	}))

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
					handleError(w, http.StatusInternalServerError, "card reader not found")
					return

				}
				handleError(w, http.StatusInternalServerError, resp.Error.Error())
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(200)
			w.Write(resp.Data)
		}

	case <-ctx.Done():
		handleError(w, http.StatusRequestTimeout, "request timed out")
	}
}

func CORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Credentials", "true")
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

		if r.Method == "OPTIONS" {
			http.Error(w, "No Content", http.StatusNoContent)
			return
		}

		next(w, r)
	}
}

func handleError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write([]byte(`{"error": "` + message + `"}`))
}
