package main

import (
	"context"
	log "github.com/sirupsen/logrus"
 	"os"
	"os/signal"
	"wb_l0/internals/app"
	"wb_l0/internals/cfg"
)

func main()  {


	config := cfg.LoadAndStoreConfig()

	ctx, cancel := context.WithCancel(context.Background()) // создаем контекст для работы контекстно зависимых частей системы

	//grays full shut down
	c:= make(chan os.Signal, 1) //coздаем канал для сигналов системы
	signal.Notify(c, os.Interrupt) //ждем сигнал ОС что наше приложение завершилось

	server := app.NewServer(config, ctx) // создаем сервер

	go func() { //горутина отслеживания сигнала сообщений системы
		oscall := <-c // ловим сигнал
		log.Printf("system call:%+v", oscall)
		server.Shutdown() //выключаем сервер
		cancel() // отменяем все контексты которые использовали для бд и т.д.
	}()

	server.Serve() //запускаем сервер



}
