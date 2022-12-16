package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/cloudfoundry-community/worlds-simplest-service-broker/pkg/broker"

	"code.cloudfoundry.org/lager"
	"github.com/gorilla/mux"
	"github.com/pivotal-cf/brokerapi"
)

func statusAPI(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK")
}

func main() {
	logger := lager.NewLogger("worlds-simplest-service-broker")
	logger.RegisterSink(lager.NewWriterSink(os.Stdout, lager.DEBUG))
	logger.RegisterSink(lager.NewWriterSink(os.Stderr, lager.ERROR))

	servicebroker := broker.NewBrokerImpl(logger)

	brokerCredentials := brokerapi.BrokerCredentials{
		Username: os.Getenv("AUTH_USER"),
		Password: os.Getenv("AUTH_PASSWORD"),
	}
	brokerAPI := brokerapi.New(servicebroker, logger, brokerCredentials)

	r := mux.NewRouter()
	r.HandleFunc("/health", statusAPI).Methods(http.MethodGet)
	r.PathPrefix("/v2").Handler(brokerAPI)

	// http.HandleFunc("/health", statusAPI)
	// http.Handle("/", brokerAPI)

	port := os.Getenv("PORT")
	if port == "" {
		port = ":3000"
	}
	// fmt.Println("\n\nStarting World's Simplest Service Broker on 0.0.0.0:" + port)

	logger.Info("Serving", lager.Data{"port": port})
	http.ListenAndServe(port, r)
	// logger.Fatal("http-listen", http.ListenAndServe("0.0.0.0:"+port, nil))
}
