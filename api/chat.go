package api

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"tinyIM/api/router"
	"tinyIM/config"
)

type Chat struct {
}

func New() *Chat {
	return &Chat{}
}

func (c *Chat) Run() {
	r := router.Register()
	runMode := config.GetGinRunMode()
	gin.SetMode(runMode)
	apiCfg := config.Conf.Api
	port := apiCfg.ApiBase.ListenPort
	logrus.Infof("chat server start in %s mode at :%d", runMode, port)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logrus.Errorf("start chat server failed: %s\n", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit
	logrus.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logrus.Errorf("chat server shutdown failed: %s\n", err)
	}

	logrus.Info("chat server exiting")
	os.Exit(0)
}
