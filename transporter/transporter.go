package transporter

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/sannonthachai/poc-sse-go/sse"
)

type SSERoute struct {
	Broker *sse.Broker
}

func NewAuthHttpRoute(broker *sse.Broker) *SSERoute {
	return &SSERoute{Broker: broker}
}

func (s *SSERoute) SSERoutePrivate(e *echo.Group) {
	sse := e.Group("/api/v1")
	sse.GET("/sse", s.handleSSE)
}

func (s *SSERoute) handleSSE(c echo.Context) error {
	c.Response().Header().Set("Content-Type", "text/event-stream")
	c.Response().Header().Set("Cache-Control", "no-cache")
	c.Response().Header().Set("Connection", "keep-alive")
	c.Response().Header().Set("Transfer-Encoding", "chunked")
	// c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	// c.Response().WriteHeader(http.StatusOK)

	messageChan := make(chan string)
	// enc := json.NewEncoder(c.Response())

	s.Broker.NewClients <- messageChan

	for {

		select {

		case message := <-messageChan:
			fmt.Fprintf(c.Response().Writer, "event: notice\n")
			fmt.Fprintf(c.Response().Writer, "data: Message: %s\n\n", message)

			fmt.Fprintf(c.Response().Writer, "event: add\n")
			fmt.Fprintf(c.Response().Writer, "data: add: %s\n\n", "add")
			// if err := enc.Encode(message); err != nil {
			// 	return err
			// }
			c.Response().Flush()

		case <-c.Request().Context().Done():
			s.Broker.DefunctClients <- messageChan
			return nil
		}

	}
}
