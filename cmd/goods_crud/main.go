package main

import (
	"fmt"
	"log"
	"os"
	"time"

	api "github.com/Set2105/hezzl_test_goods_crud/internal/api"
	g_crud "github.com/Set2105/hezzl_test_goods_crud/internal/api/goods_crud"
	"github.com/Set2105/hezzl_test_goods_crud/internal/nats"
	psql "github.com/Set2105/hezzl_test_goods_crud/internal/postgres"
	redis "github.com/Set2105/hezzl_test_goods_crud/internal/redis"
)

func start() {
	pSql, err := psql.InitPostgresDb(&psql.PostgresSettings{}, 10)
	if err != nil {
		fmt.Println(err)
		return
	}

	n, err := nats.InitNats(&nats.NatsSettings{})
	if err != nil {
		fmt.Println(err)
		return
	}

	r, err := redis.InitRedis(&redis.RedisSettings{})
	if err != nil {
		fmt.Println(err)
		return
	}

	mf, err := g_crud.InitGoodsCRUD(pSql, r, n, log.New(os.Stdout, "Error: ", log.Ldate|log.Ltime), time.Minute)
	if err != nil {
		fmt.Println(err)
		return
	}
	goodsServer, err := api.InitServer(":20001", mf)
	if err != nil {
		fmt.Println(err)
		return
	}
	if err := goodsServer.Start(); err != nil {
		panic(err)
	}
}

func main() {
	for {
		start()
		time.Sleep(time.Second)
	}
}
