package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"

	"github.com/manishmeganathan/essensio/jsonrpc"
)

// TODO:
// 1. RPC instead CLI -  Done
// 2. Tx Model
// 3. Tx Pool
// 4. Update the RPC

const SERVER_PORT = 8080

func main() {
	// Create a new RPC Server and register the JSON Codec
	server := rpc.NewServer()
	server.RegisterCodec(json.NewCodec(), "application/json")
	server.RegisterCodec(json.NewCodec(), "application/json;charset=UTF-8")

	// Create a new JSON-RPC API for Essensio
	api := jsonrpc.NewAPI()
	defer api.Stop()

	// Register the Essensio API with the Server
	if err := server.RegisterService(api, ""); err != nil {
		log.Fatalln("Failed to Register Essensio API:", err)
	}

	// Set up a new Multiplexed Router
	router := mux.NewRouter()
	router.Handle("/rpc", server)

	// HTTP Listen & Serve
	fmt.Println("Server Starting...")
	if err := http.ListenAndServe(fmt.Sprintf(":%v", SERVER_PORT), router); err != nil {
		log.Fatalln(err)
	}
}
