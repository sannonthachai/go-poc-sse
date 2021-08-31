package main

import (
	"fmt"
	"log"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sannonthachai/poc-sse-go/sse"
	"github.com/sannonthachai/poc-sse-go/transporter"
)

func main() {
	e := echo.New()
	routePrivate := e.Group("/private")
	sse := sse.NewBroker()
	sse.Start()

	router := transporter.NewAuthHttpRoute(sse)
	router.SSERoutePrivate(routePrivate)

	go func() {
		for i := 0; ; i++ {

			// Create a little message to send to clients,
			// including the current time.
			sse.Messages <- fmt.Sprintf("%d - the time is %v", i, time.Now())

			// Print a nice log message and sleep for 5s.
			log.Printf("Sent message %d ", i)
			time.Sleep(5e9)

		}
	}()

	e.Logger.Fatal(e.Start("localhost:8080"))
}
