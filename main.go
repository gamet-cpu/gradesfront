package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"

	handler "grades/handlers"
	"grades/middleware"
	"grades/repository"
)

func main() {
	viper.SetConfigFile("config.yaml")
	_ = viper.ReadInConfig()
	viper.BindEnv("db.url", "DATABASE_URL")
	viper.BindEnv("jwt.secret", "JWT_SECRET")
	viper.BindEnv("port", "PORT")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}

	dbpool, err := pgxpool.New(context.Background(),
		viper.GetString("db.url"))
	if err != nil {
		log.Fatal(err)
	}

	userRepo := repository.NewUserRepository(dbpool)
	userHandler := handler.NewUserHandler(userRepo)
	r := gin.Default()

	r.StaticFile("/", "./static/index.html")
	r.Static("/static", "./static")
	r.POST("/register", userHandler.Register)
	r.POST("/login", userHandler.Login)

	auth := r.Group("/")
	auth.Use(middleware.AuthMiddleware())
	{
		auth.GET("/user/me", userHandler.GetMe)          //Получение данных о пользователе
		auth.PUT("/user", userHandler.Update)            //Обновление данных пользователя
		auth.DELETE("/user", userHandler.Delete)         //Удаление пользователя
		auth.POST("/subjects", userHandler.CreateSubj)   //Создание данных о предметах
		auth.GET("/subjects", userHandler.GetSubj)       //Получение данных о предметах
		auth.PUT("/subjects", userHandler.UpdateSubj)    //Обновление данных о предметах
		auth.DELETE("/subjects", userHandler.DeleteSubj) //Удаление предметов
	}

	r.Run(":" + viper.GetString("port"))
}
