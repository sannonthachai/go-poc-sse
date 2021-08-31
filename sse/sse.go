package sse

type Broker struct {
	Clients        map[chan string]bool
	NewClients     chan chan string
	DefunctClients chan chan string
	Messages       chan string
}

func NewBroker() *Broker {
	return &Broker{
		Clients:        make(map[chan string]bool),
		NewClients:     make(chan (chan string)),
		DefunctClients: make(chan (chan string)),
		Messages:       make(chan string),
	}
}

func (b *Broker) Start() {
	go func() {
		for {
			select {

			case s := <-b.NewClients:
				b.Clients[s] = true

			case s := <-b.DefunctClients:
				delete(b.Clients, s)
				close(s)

			case msg := <-b.Messages:
				for s := range b.Clients {
					s <- msg
				}

			}
		}
	}()
}
