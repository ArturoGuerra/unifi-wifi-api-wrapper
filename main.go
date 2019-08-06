package main

import (
    "github.com/labstack/echo"
    "github.com/labstack/echo/middleware"
    "./handlers"
    "./middlwares"
    "./config"
)

func main () {
    e := echo.New()
    Config, _ := config.LoadConfig()
    handlers := handlers.New(Config)
    middlewares := middlewares.New(Config)
    e.Use(middleware.Logger())
    e.Use(middleware.Recover())
    e.POST("/api/login", handlers.Login)
    jwtConfig := middleware.JWTConfig{
        Claims: &config.JwtCustomClaims{},
        SigningKey: []byte("43q4cqerasdfcd"),
    }

    auth := e.Group("/api/wlanconf")
    auth.Use(middleware.JWTWithConfig(jwtConfig))
    auth.Use(middlewares.TokenAuth)
    auth.PUT("", handlers.Wlanconf)
//    auth.GET("", handlers.WlanconfGet)
    e.Logger.Fatal(e.Start(":6969"))
}
