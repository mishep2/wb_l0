package main

import (
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/mishep2/wb_l0"
	"github.com/mishep2/wb_l0/controllers"
	"github.com/mishep2/wb_l0/models"
	"github.com/mishep2/wb_l0/pkg/handler"
	"github.com/mishep2/wb_l0/pkg/repository"
	"github.com/mishep2/wb_l0/pkg/repository/cache"
	"github.com/mishep2/wb_l0/pkg/repository/postgres"
	"github.com/mishep2/wb_l0/pkg/service"
	"github.com/nats-io/stan.go"
	"log"
	"os"
	"time"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading env variables: %s", err.Error())
	}

	orderCache := cache.NewCache() // создается кеш с помощью вызова фунции NewCache() из пакета cache и присваиваеется переменная OrderCache, зачет можно использовать для хранения данных временно, для ускорения доступа к ним

	db, err := postgres.NewPostgresDB() //Создает новый экземпляр базы данных Postgres с помощью postgres.NewPostgresDB(). Если произошла ошибка при инициализации базы данных, программа завершится с фатальным логированием сообщения об ошибке.
	if err != nil {
		log.Fatalf("failed to initialize db: %s", err.Error())
	}

	var orders []models.Order //Объявляет переменную orders как срез моделей заказов (models.Order).

	orders, err = postgres.GetOrders(db) // Вызывает функцию GetOrders() из пакета postgres, передавая в нее созданный экземпляр базы данных db. Результаты запроса к базе данных (список заказов) и ошибка сохраняются в переменные orders и err.
	if err != nil {
		log.Fatalf("error by getting orders from postgres :%s", err)
	}
	// set orders into cache
	for _, order := range orders {
		orderCache.Set(order.Order_uid, order)
	} //фукнция проходит по каждому элементу массиву с помощью цикла range  Для каждого элемента она вызывает метод "Set" объекта "orderCache", используя ключ "order.Order_uid" и значение "order". Таким образом, эта функция добавляет или обновляет элементы в кэше "orderCache" с помощью данных из массива "orders".

	// nats

	sc, err := stan.Connect("test-cluster", "test", stan.NatsURL("nats://nats_streaming:4222"))
	if err != nil {
		log.Fatalf("failed to connect to nats-streaming : %s", err)
	}

	// Simple Async Subscriber
	preTime, _ := time.ParseDuration("1m")
	sub, _ := sc.Subscribe("orders", controllers.MsgHandler(orderCache, db), stan.StartAtTimeDelta(preTime))

	repos := repository.NewRepository(orderCache)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(wb_l0.Server)
	if err := srv.Run(os.Getenv("SERVER_PORT"), handlers.InitRoutes()); err != nil {
		log.Fatalf("error occured while running http server: %s", err.Error())
	}

	// Unsubscribe sudo docker-compose -f nats&nats-streaming.yaml up -d
	sub.Unsubscribe()

	// Close connection
	sc.Close()
}
