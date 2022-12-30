package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/Osedisc/assessment/handler"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	handler.InitDatabase(os.Getenv("DATABASE_URL"))
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			auth := c.Request().Header.Get("Authorization")
			if strings.Contains(auth, "wrong") {
				return echo.ErrUnauthorized
			}
			return next(c)
		}
	})

	e.POST("/expenses", handler.PostExpenses)
	e.GET("/expenses/:id", handler.GetExpensebyid)
	e.PUT("/expenses/:id", handler.UpdateExpense)
	e.GET("/expenses", handler.GetAllExpenses)

	fmt.Println("Please use server.go for main file")
	fmt.Println("start at port:", e.Start(os.Getenv("PORT")))

	//Graceful shutdown
	go func() {
		if err := e.Start(os.Getenv("PORT")); err != nil && err != http.ErrServerClosed { // Start server
			e.Logger.Fatal("shutting down the server")
		}
	}()
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	<-shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
