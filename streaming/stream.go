package streaming

import (
    "net/http"
    "time"

    "github.com/gorilla/websocket"
)

const (
    writeWait   = 10 * time.Second
    pongWait    = 60 * time.Second
    pingPeriod  = (pongWait * 9) / 10
    readLimit   = 512
)

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    // TODO: restrict origins in production
    CheckOrigin: func(r *http.Request) bool { return true },
}

type Client struct {
    hub        *Hub
    conn       *websocket.Conn
    send       chan []byte
    instrument string
}

func ServeWS(hub *Hub, w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        return
    }
    inst := r.URL.Query().Get("instrument")
    c := &Client{
        hub:        hub,
        conn:       conn,
        send:       make(chan []byte, 256),
        instrument: inst,
    }

    c.hub.register <- c
    go c.writePump()
    go c.readPump()
}

func (c *Client) readPump() {
    defer func() {
        c.hub.unregister <- c
        c.conn.Close()
    }()
    c.conn.SetReadLimit(readLimit)
    c.conn.SetReadDeadline(time.Now().Add(pongWait))
    c.conn.SetPongHandler(func(string) error {
        c.conn.SetReadDeadline(time.Now().Add(pongWait))
        return nil
    })
    for {
        // We ignore client messages; this keeps the connection alive
        if _, _, err := c.conn.ReadMessage(); err != nil {
            break
        }
    }
}

func (c *Client) writePump() {
    ticker := time.NewTicker(pingPeriod)
    defer func() {
        ticker.Stop()
        c.conn.Close()
    }()
    for {
        select {
        case msg, ok := <-c.send:
            c.conn.SetWriteDeadline(time.Now().Add(writeWait))
            if !ok {
                _ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
                return
            }
            if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
                return
            }
        case <-ticker.C:
            c.conn.SetWriteDeadline(time.Now().Add(writeWait))
            if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
                return
            }
        }
    }
}