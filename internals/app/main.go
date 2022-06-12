package app

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/nats-io/stan.go"
	"log"
	"net/http"
	"time"
	"wb_l0/api"
	dbOrder "wb_l0/internals/app/db"
	"wb_l0/internals/app/handlers"
	"wb_l0/internals/app/processors"
	"wb_l0/internals/cfg"
)

//сборка всех зависимостей для сервера
type Server struct {
	config cfg.Cfg
	ctx    context.Context
	srv    *http.Server
	db     *pgxpool.Pool
}

//задаем поля нашего сервера, для его старта нам нужен контекст и конфигурация
func NewServer(config cfg.Cfg, ctx context.Context) *Server {
	server := new(Server)
	server.ctx = ctx
	server.config = config
	return server
}

func (server *Server) Serve() {
	log.Println("Starting server")
	var err error
	log.Println(server.config.GetDBString())

	//несколько соединений к бд и сохраним его для закрытия при остановке приложения
	server.db, err = pgxpool.Connect(server.ctx, server.config.GetDBString())
	if err != nil {
		log.Fatal(err)
	}

	//создаем экземпляр storage для работы с БД и всем что связано с orders
	ordersStorage := dbOrder.NewOrdersStorage(server.db)
	//процессоры используют наш storage (логику работы с бд)
	ordersProcessor := processors.NewUserProcessor(ordersStorage)
	//хэндлеры используют процессоры
	ordersHandler := handlers.NewOrdersHandler(ordersProcessor)

	//создаем роуты
	routes := api.CreateRoutes(ordersHandler)

	//routes.Use() //todo дописать мидлвер

	//загрузим кеш из БД
	if err := ordersStorage.UploadCache(context.Background()); err != nil {
		log.Printf("Upload Cache: ", err)
	}

	//стартуем сервер, передаем порт и наш mux
	server.srv = &http.Server{ //в отличие от примеров http, здесь мы передаем наш server в поле структуры, для работы в Shutdown
		Addr:    ":" + server.config.Port,
		Handler: routes,
	}

	//nats streaming
	sNat, err := stan.Connect(
		"test-cluster",
		"order-consumer", stan.NatsURL(cfg.NatsURI),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer sNat.Close()

	//пишем в postgresql
	// Subscribe starting with most recently published value
	if _, err = sNat.Subscribe("orders", ordersProcessor.CreateOrder, stan.StartWithLastReceived()); err != nil {
		return
	}

	log.Println("server started")
	err = server.srv.ListenAndServe() //запускаем сервер
	if err != nil {
		log.Fatalln(err)
	}

	return
}

//закрывает пулл соединения к базе данных, вызывает cancel, останавливает сервер
func (server *Server) Shutdown() {
	log.Printf("server stopped")

	ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	server.db.Close()
	defer func() {
		cancel()
	}()

	var err error
	if err = server.srv.Shutdown(ctxShutdown); err != nil { //выключаем сервер, с ограниченным по времени контекстом
		log.Fatalf("server Shutdown Failed %q", err)
	}

	log.Printf("server exited properly")

	if err == http.ErrServerClosed {
		err = nil
	}
}
