package chat

import (
	"context"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/coder/websocket"
)

type ChatService struct {
	Chats map[string]*Chat
}

type Chat struct {
	SubMu sync.Mutex
	Subs  map[*Subscriber]struct{}
}

type Subscriber struct {
	Msg chan []byte
}

func NewChatServer() *ChatService {
	return &ChatService{
		Chats: make(map[string]*Chat),
	}
}

func NewChat() *Chat {
	c := &Chat{
		Subs: make(map[*Subscriber]struct{}),
	}
	return c
}

func (ch *Chat) Subscribe(w http.ResponseWriter, r *http.Request) error {
	var mu sync.Mutex
	var conn *websocket.Conn
	s := &Subscriber{
		Msg: make(chan []byte),
	}

	ch.addSubscriber(s)
	defer ch.deleteSubscriber(s)

	c, err := websocket.Accept(w, r, nil)
	if err != nil {
		slog.Error("could not accept web socket connection", "error", err)
		return err
	}

	mu.Lock()
	conn = c
	mu.Unlock()

	ctx := conn.CloseRead(context.Background())

	go pingLoop(conn, ctx, 10*time.Second)

	for {
		select {
		case msg := <-s.Msg:
			err := writeTimeout(ctx, 5*time.Second, conn, msg)
			if err != nil {
				return err
			}
		case <-ctx.Done():
			slog.Info("subscriber disconnected, no longer handling messages")
			return ctx.Err()
		}
	}
}

func pingLoop(c *websocket.Conn, ctx context.Context, duration time.Duration) {
	ticker := time.NewTicker(duration)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if c == nil {
				continue
			}
		}

		if err := pingConnection(c); err != nil {
			slog.Error("unsuccessful ping, closing connection", "error", err)
			err := c.Close(websocket.StatusNormalClosure, "could not ping connection")
			if err != nil {
				slog.Error("failed to close now from connection in ping loop", "error", err)
				return
			}
			return
		}
	}
}

func pingConnection(c *websocket.Conn) error {
	pctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := c.Ping(pctx); err != nil {
		return err
	}
	return nil
}

func (ch *Chat) deleteSubscriber(s *Subscriber) {
	ch.SubMu.Lock()
	delete(ch.Subs, s)
	ch.SubMu.Unlock()
}

func (ch *Chat) addSubscriber(s *Subscriber) {
	ch.SubMu.Lock()
	ch.Subs[s] = struct{}{}
	ch.SubMu.Unlock()
}

func writeTimeout(ctx context.Context, timeout time.Duration, c *websocket.Conn, msg []byte) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return c.Write(ctx, websocket.MessageText, msg)
}

func (ch *Chat) PublishMessage(message []byte) {
	ch.SubMu.Lock()
	defer ch.SubMu.Unlock()

	for s := range ch.Subs {
		select {
		case s.Msg <- message:
		default:
			slog.Warn("no message sent, channel is ful or no reciever")
		}
	}
}
