package api

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"tinyIM/api/router"
	"tinyIM/api/rpc"
	"tinyIM/config"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

type Chat struct {
}

func New() *Chat {
	return &Chat{}
}

func (c *Chat) Run() {
	rpc.InitLogicRpcClient()

	r := router.Register()
	runMode := config.GetGinRunMode()
	gin.SetMode(runMode)
	apiCfg := config.Conf.Api
	port := apiCfg.ApiBase.ListenPort
	logrus.Infof("chat server start in %s mode at :%d", runMode, port)

	// 使用http.Server封装一层gin是为了方便后续实现优雅退出
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logrus.Errorf("start chat server failed: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
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
