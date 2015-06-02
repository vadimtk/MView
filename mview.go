package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/vadimtk/mview/controllers"
	"log"
	"net/http"
)

func main() {
	metricsC := controllers.Metrics

	router := httprouter.New()

	router.GET("/metrics", metricsC.Perform(metricsC.Index))
	router.GET("/browse", metricsC.Perform(metricsC.Browse))
	router.GET("/metric/hostid/:hostid/mid/:mid/sid/:sid", metricsC.Perform(metricsC.Metric))

	log.Println("Starting server on :3000")
	log.Fatal(http.ListenAndServe(":3000", router))
}
