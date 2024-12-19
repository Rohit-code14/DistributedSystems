package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type logger struct {
	Log *Log
}

func StartServer() {

	Logger := logger{Log: NewLog()}

	http.HandleFunc("POST /", func(w http.ResponseWriter, r *http.Request) {
		var requestJson ProduceRequest
		err := json.NewDecoder(r.Body).Decode(&requestJson)
		if err != nil {
			fmt.Print(err)
			errorJson := Error{Message: "Body is either empty or unparsable", ErrorCode: http.StatusBadRequest}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(errorJson)
		}

		offset, err := Logger.Log.Append(requestJson.Record)

		if err != nil {
			errorJson := Error{Message: "Unable to commit log", ErrorCode: http.StatusInternalServerError}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(errorJson)
		}

		var responseJson ProduceResponse
		responseJson.Offset = offset
		err = json.NewEncoder(w).Encode(responseJson)

		if err != nil {
			errorJson := Error{Message: "Unable to generate response", ErrorCode: http.StatusInternalServerError}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(errorJson)
		}
	})

	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		var requestJson ConsumeRequest
		err := json.NewDecoder(r.Body).Decode(&requestJson)
		if err != nil {
			fmt.Print(err)
			errorJson := Error{Message: "Body is either empty or unparsable", ErrorCode: http.StatusBadRequest}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(errorJson)
		}

		record, err := Logger.Log.Read(requestJson.Offset)

		if err != nil {
			errorJson := Error{Message: "Offset not found", ErrorCode: http.StatusNotFound}
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(errorJson)
		}

		var responseJson ConsumeResponse
		responseJson.Record = record
		err = json.NewEncoder(w).Encode(responseJson)

		if err != nil {
			errorJson := Error{Message: "Unable to generate response", ErrorCode: http.StatusInternalServerError}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(errorJson)
		}
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

type ProduceRequest struct {
	Record Record `json:"record"`
}
type ProduceResponse struct {
	Offset uint64 `json:"offset"`
}
type ConsumeRequest struct {
	Offset uint64 `json:"offset"`
}
type ConsumeResponse struct {
	Record Record `json:"record"`
}

type Error struct {
	Message   string `json:"message"`
	ErrorCode int    `json:"errorCode"`
}
