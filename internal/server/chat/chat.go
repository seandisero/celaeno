package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/coder/websocket"
	"github.com/seandisero/celaeno/internal/server"
	"github.com/seandisero/celaeno/internal/shared"
)

type ChatService struct {
	Chats map[string]*Chat
}

type Chat struct {
	SubMu sync.Mutex
	Subs  map[*Subscriber]struct{}
}

type Subscriber struct {
	Msg chan shared.Message
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

func (ch *Chat) HandleMessages(conn *websocket.Conn, sub *Subscriber) error {
	slog.Info("handling messages")
	defer conn.CloseNow()
	ctx := conn.CloseRead(context.Background())
	for {
		select {
		case msg := <-sub.Msg:
			fmt.Printf(" > server got message %s", msg.Message)
			err := writeTimeout(ctx, 5*time.Second, conn, msg)
			if err != nil {
				return err
			}
		case <-ctx.Done():
			fmt.Println("no longer handling messages")
			return ctx.Err()
		}
	}
}

func (ch *Chat) Subscribe(w http.ResponseWriter, r *http.Request) error {
	var mu sync.Mutex
	var conn *websocket.Conn
	s := &Subscriber{
		Msg: make(chan shared.Message),
	}
	ch.addSubscriber(s)
	defer ch.deleteSubscriber(s)

	c, err := websocket.Accept(w, r, nil)
	if err != nil {
		server.RespondWithError(w, http.StatusInternalServerError, "could not accept request: %w", err)
		return err
	}
	mu.Lock()
	conn = c
	mu.Unlock()

	ctx := conn.CloseRead(context.Background())
	for {
		select {
		case msg := <-s.Msg:
			fmt.Printf(" > server got message %s", msg.Message)
			err := writeTimeout(ctx, 5*time.Second, conn, msg)
			if err != nil {
				return err
			}
		case <-ctx.Done():
			fmt.Println("no longer handling messages")
			return ctx.Err()
		}
	}
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

func writeTimeout(ctx context.Context, timeout time.Duration, c *websocket.Conn, msg shared.Message) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return c.Write(ctx, websocket.MessageText, data)
}

func (ch *Chat) PublishMessage(message shared.Message) {
	ch.SubMu.Lock()
	defer ch.SubMu.Unlock()

	for s := range ch.Subs {
		s.Msg <- message
	}
}

func (ch *Chat) processMessage(ctx context.Context, msg shared.Message) {
	for sub := range ch.Subs {
		ch.SubMu.Lock()
		sub.Msg <- msg
		ch.SubMu.Unlock()
	}
}
