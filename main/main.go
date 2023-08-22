package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	//_ "github.com/jackc/pgx"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	_ "github.com/lib/pq"

	"github.com/almaz91/todo-app"
	"github.com/almaz91/todo-app/pkg/handler"
	"github.com/almaz91/todo-app/pkg/repository"
	"github.com/almaz91/todo-app/pkg/service"
	"github.com/spf13/viper"
)

// @title           Todo App API
// @version         1.0
// @description     API Server for TodoList Application

// @host      almazneft123.ru:80
// @BasePath  /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))
	if err := initConfig(); err != nil {
		logrus.Fatalf("error initializing configs: %s", err.Error())
	}

	if err := godotenv.Load(".env"); err != nil {
		logrus.Fatalf("error loading env variables: %s", err.Error())
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
	})
	if err != nil {
		logrus.Fatalf("failed to initialize db: %s", err.Error())
	}
	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)
	srv := new(todo.Server)
	go func() {
		if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
			logrus.Fatalf("error occured while running hhtp server: %s", err.Error())
		}
	}()

	logrus.Printf("TodoApp server start listen om port: %s", viper.GetString("port"))

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	logrus.Print("TodoApp Shutting Down")
	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Print("error occurred while shutting down http server: %s\n", err.Error())
	}

	if err := db.Close(); err != nil {
		logrus.Print("error occurred while closing db connection: %s\n", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
