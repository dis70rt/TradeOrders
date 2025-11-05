package streaming

type Broadcast struct {
    Instrument string
    Data       []byte
}

type Hub struct {
    clients    map[*Client]bool
    register   chan *Client
    unregister chan *Client
    broadcast  chan Broadcast
}

func NewHub() *Hub {
    return &Hub{
        clients:    make(map[*Client]bool),
        register:   make(chan *Client),
        unregister: make(chan *Client),
        broadcast:  make(chan Broadcast, 1024),
    }
}

func (h *Hub) Run() {
    for {
        select {
        case c := <-h.register:
            h.clients[c] = true
        case c := <-h.unregister:
            if _, ok := h.clients[c]; ok {
                delete(h.clients, c)
                close(c.send)
            }
        case msg := <-h.broadcast:
            for c := range h.clients {
                // If client subscribed to a specific instrument, filter
                if c.instrument != "" && c.instrument != msg.Instrument {
                    continue
                }
                select {
                case c.send <- msg.Data:
                default:
                    // Drop slow client
                    delete(h.clients, c)
                    close(c.send)
                }
            }
        }
    }
}

func (h *Hub) Broadcast(b Broadcast) {
    h.broadcast <- b
}