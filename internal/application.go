package internal

import (
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/spf13/viper"

	"pets/internal/config"
	"pets/internal/repository"
	"pets/internal/server"
	"pets/internal/service"
	"pets/pkg/logger"
)

// App is a main app struct
type App struct {
	config     *config.Scheme
	repository repository.IRepository
	server     *server.HttpServer
}

// NewApp is used to get new App instance
func NewApp() *App {
	a := &App{}

	a.initConfig()

	a.repository = repository.NewRepository(a.config.DB)

	srv := service.NewService(a.repository)
	a.server = server.NewServer(a.config.Http, srv)

	return a
}

// initConfig is used to init new config using viper. If no config in env vars will set default values
func (a *App) initConfig() {
	a.config = &config.Scheme{}

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			logger.Log().WithField("layer", "App").Fatalf("viper read config error: %v", err.Error())
		}
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	viper.AllowEmptyEnv(true)

	if err := viper.Unmarshal(a.config); err != nil {
		logger.Log().WithField("layer", "App").Fatalf("viper unmarshal config error: %v", err.Error())
	}

	logger.Log().WithField("layer", "App").Infof("config initialaized")
}

// Run is used to run app
func (a *App) Run() {
	go a.server.ListenAndServ()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	<-quit
	return
}

// Stop is used to stop app
func (a *App) Stop() {
	if a.repository != nil {
		a.repository.Stop()
	}
}
