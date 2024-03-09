package main

import (
	"fmt"
	"log"
	"rinha-de-backend-2024-q1/src/config"
	"rinha-de-backend-2024-q1/src/controller"
	"rinha-de-backend-2024-q1/src/database"
	"time"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

func init() {
	config.LoadEnvironment()

}

func main() {

	db, err := database.NewDBConn(config.Strdbconn, config.MaxConn)
	if err != nil {
		log.Fatalf("error in database.NewDBConn: %v\n", err)
	}
	ctrl := controller.InjectDB(db)
	defer db.Close()

	router := fasthttprouter.New()
	router.GET("/clientes/:customer_id/extrato", ctrl.Statement)
	router.POST("/clientes/:customer_id/transacoes", ctrl.Transactions)

	srv := fasthttp.Server{
		Handler:            router.Handler,
		ReadTimeout:        time.Second + 29000*time.Millisecond,
		WriteTimeout:       time.Second + 29000*time.Millisecond,
		MaxRequestBodySize: 2 * 1024,
	}

	fmt.Printf("Running on %d\n", config.Apiport)
	go func() {
		if err := srv.ListenAndServe(fmt.Sprintf(":%d", config.Apiport)); err != nil {
			log.Fatalf("error in fasthttp.ListenAndServe: %v\n", err)
		}
	}()
	select {}

}
