package main

import(
    "fmt"
    "io"
    "time"
    "net/http"
    "golang.org/x/net/websocket"
)

type Server struct {
    conns map[*websocket.Conn]bool
}

func NewServer() *Server {
    return &Server{
        conns: make(map[*websocket.Conn]bool),
    }
}

func (s *Server) handleWSOrderbook(ws *websocket.Conn) {
    fmt.Println("new connsection from client to orderbook feed:", ws.RemoteAddr())

    for {
        payload := fmt.Sprintf("orderbook data -> %d\n", time.Now().UnixNano())
        ws.Write([]byte(payload))
        time.Sleep(time.Second * 2)
    }
}

func (s *Server) handleWS(ws *websocket.Conn) {
    fmt.Println("new connsection from client:", ws.RemoteAddr())

    s.conns[ws] = true

    s.readLoop(ws)
}

func (s *Server) readLoop(ws *websocket.Conn) {
    buf := make([]byte, 1024)
    for {
        n, err := ws.Read(buf)
        if err != nil {
            if err == io.EOF {
                break
            }
            fmt.Println("read error:", err)
            continue
        }
        msg := buf[:n]

        s.broadcast(msg)
    }
}

func (s *Server) broadcast(b []byte) {
    for ws := range s.conns {
        go func(ws *websocket.Conn) {
            if _, err := ws.Write(b); err != nil {
                fmt.Println("write error:", err)
            }
        }(ws)
    }
}

func main() {
    server := NewServer()
    http.Handle("/ws", websocket.Handler(server.handleWS))
    http.Handle("/orderbookfeed", websocket.Handler(server.handleWSOrderbook))
    http.ListenAndServe(":8080", nil)
}

