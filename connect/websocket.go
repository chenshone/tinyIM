package connect

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"tinyIM/config"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

func (c *Connect) InitWebSocket() {
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		c.serveWS(nil, w, r)
	})

	srv := &http.Server{
		Addr:    config.Conf.Connect.ConnectWebsocket.Bind,
		Handler: mux,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logrus.Errorf("Connect layer InitWebsocket() error: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit
	logrus.Infof("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logrus.Errorf("Server Shutdown:", err)
	}
	logrus.Infof("Server exiting")
	os.Exit(0)
}

func (c *Connect) serveWS(server *Server, w http.ResponseWriter, r *http.Request) {
	upGrader := websocket.Upgrader{
		ReadBufferSize:  server.Option.ReadBufferSize,
		WriteBufferSize: server.Option.WriteBufferSize,
	}
	upGrader.CheckOrigin = func(r *http.Request) bool { return true }

	conn, err := upGrader.Upgrade(w, r, nil)

	if err != nil {
		logrus.Errorf("serverWs err:%s", err.Error())
		return
	}

	//default broadcast size eq 512
	ch := NewChannel(server.Option.BroadcastSize)
	ch.conn = conn

	//send data to websocket conn
	go server.writePump(ch, c)
	//get data from websocket conn
	go server.readPump(ch, c)
}
