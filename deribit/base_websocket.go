package deribit

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

const (
	wssSchema = "wss"
)

type baseWS struct {
	conn              *websocket.Conn
	schema            string
	host              string
	path              string
	onOpenCb          func() error
	onMessageCb       func(interface{}) error
	pingDuration      time.Duration
	pingMessage       []byte
	reconnectInterval time.Duration
}

func newWebsocket(schema, host, path string, onOpenCb func() error, onMessageCb func(interface{}) error, pingDuration time.Duration, pingMessage []byte, reconnectInterval time.Duration) *baseWS {
	return &baseWS{
		schema:            schema,
		host:              host,
		path:              path,
		onOpenCb:          onOpenCb,
		onMessageCb:       onMessageCb,
		pingDuration:      pingDuration,
		pingMessage:       pingMessage,
		reconnectInterval: reconnectInterval,
	}
}

func (ws *baseWS) startStreaming(ctx context.Context) {
	for {
		err := ws.connect(ctx)
		if err != nil {
			log.Printf("Error while listening to websocket, err: %s", err.Error())
			time.Sleep(ws.reconnectInterval)
			log.Printf("Reconnecting websocket")
		}
	}
}

func (ws *baseWS) onClose() {
	ws.conn.Close()
	ws.conn = nil
}

func (ws *baseWS) connect(ctx context.Context) error {
	u := url.URL{Scheme: wssSchema, Host: ws.host, Path: ws.path}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return fmt.Errorf("Dial error: %s", err.Error())
	}
	ws.conn = c
	defer ws.onClose()
	pingCtx, cancelFn := context.WithCancel(ctx)

	go ws.startPingService(pingCtx)

	err = ws.onOpenCb()
	if err != nil {
		log.Printf("Error while calling on open callback function, err: %s", err.Error())
		cancelFn()
		return err
	}

	for {
		select {
		case <-ctx.Done():
			log.Println("Context done stopping websocket connection")
			cancelFn()
			return nil
		default:
			_, msg, err := ws.conn.ReadMessage()
			if err != nil {
				log.Printf("Error while reading message from websocket, err: %s", err.Error())
				cancelFn()
				return err
			}
			err = ws.onMessageCb(msg)
			if err != nil {
				log.Printf("Error while calling on message callback function, err: %s", err.Error())
				cancelFn()
				return err
			}
		}
	}
}

func (ws *baseWS) startPingService(ctx context.Context) {
	ticker := time.NewTicker(ws.pingDuration)
	for {
		select {
		case <-ticker.C:
			err := ws.conn.WriteControl(websocket.PingMessage, ws.pingMessage, time.Now().Add(ws.pingDuration))
			if err != nil {
				log.Printf("Error while writing ping message, err: %s", err.Error())
			}
		case <-ctx.Done():
			log.Println("Context done, stopping ping service")
			return
		}
	}
}
