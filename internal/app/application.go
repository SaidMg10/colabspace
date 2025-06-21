package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/SaidMg10/colabspace/internal/config"
	"github.com/SaidMg10/colabspace/internal/storage"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Application struct {
	Config  *config.Config
	Logger  *zap.SugaredLogger
	Router  *gin.Engine
	Store   storage.Storage
	Service Services
}

func (app *Application) Run() error {
	addr := fmt.Sprintf("%s:%d", app.Config.Server.Host, app.Config.Server.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      app.Router,
		WriteTimeout: app.Config.Server.WriteTimeout,
		ReadTimeout:  app.Config.Server.ReadTimeout,
		IdleTimeout:  app.Config.Server.IdleTimeout,
	}

	// Canal para shutdown, donde enviaremos errores si ocurre uno
	shutdown := make(chan error)

	// Creamos una rutina que esperará señales como Ctrl+C
	go func() {
		quit := make(chan os.Signal, 1)                      // recibe señales del OS
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // escúchalas
		s := <-quit                                          // esperamos una señal

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel() // liberamos el contexto cuando termine

		app.Logger.Infow("signal caught", "signal", s.String())

		shutdown <- srv.Shutdown(ctx) // apagamos el servidor y enviamos resultado
	}()

	app.Logger.Infow("server has started", "addr", addr, "env", app.Config.App.Env)

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdown
	if err != nil {
		return err
	}

	app.Logger.Infow("server has stopped", "addr", addr, "env", app.Config.App.Env)

	return nil
}
