package goods_crud

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Set2105/hezzl_test_goods_crud/internal/api"
	"github.com/Set2105/hezzl_test_goods_crud/internal/nats"
	"github.com/Set2105/hezzl_test_goods_crud/internal/postgres"
	"github.com/Set2105/hezzl_test_goods_crud/internal/redis"
)

const goodsTableName = "test.goods"

type GoodsCRUD struct {
	db *postgres.PostgresDb
	c  *redis.Redis
	n  *nats.Nats

	l *log.Logger // handler console logger

	cTime time.Duration // chahe time
}

// COMMON

// Init

func InitGoodsCRUD(db *postgres.PostgresDb, redis *redis.Redis, nats *nats.Nats, l *log.Logger, cTime time.Duration) (*GoodsCRUD, error) {
	if db == nil {
		return nil, fmt.Errorf("InitGoodsCRUD: db is nil")
	}
	if redis == nil {
		return nil, fmt.Errorf("InitGoodsCRUD: redis is nil")
	}
	if nats == nil {
		return nil, fmt.Errorf("InitGoodsCRUD: nats is nil")
	}
	return &GoodsCRUD{l: l, db: db, c: redis, cTime: cTime, n: nats}, nil
}

func (gc *GoodsCRUD) InitMux() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /good/create", gc.CreateGoodsHandler)
	mux.HandleFunc("PATCH /good/update", gc.UpdateGoodsHandler)
	mux.HandleFunc("DELETE /good/remove", gc.DeleteGoodsHandler)
	mux.HandleFunc("GET /good/list", gc.ListGoodsHandler)
	mux.HandleFunc("PATCH /good/reprioritiize", gc.ReprioritiizeGoodsHandler)
	return mux
}

// Handlers

func (gc *GoodsCRUD) CreateGoodsHandler(w http.ResponseWriter, req *http.Request) {
	var err error
	defer gc.ConsoleErrLog(&err)

	r := InitGoodsCreateRequest()
	if err = api.ParseAndValidate(req, r); err != nil {
		api.WriteErrorResponse(w, http.StatusBadRequest)
		err = fmt.Errorf("CreateGoodHandler.%s", err.Error())
		return
	}

	g, err := gc.CreateGood(goodsTableName, r.ProjectId, r.Payload.Name)
	if err != nil {
		api.WriteErrorResponse(w, http.StatusInternalServerError)
		err = fmt.Errorf("CreateGoodHandler.%s", err.Error())
		return
	}

	jsonG, err := json.Marshal(g)
	if err != nil {
		api.WriteErrorResponse(w, http.StatusInternalServerError)
		err = fmt.Errorf("CreateGoodHandler.json.Marshal: %s", err.Error())
		return
	}

	go gc.LogGoods(g)
	go gc.c.ChacheByte(generateGoodKey(g), jsonG, gc.cTime)
	w.Write(jsonG)
}

func (gc *GoodsCRUD) UpdateGoodsHandler(w http.ResponseWriter, req *http.Request) {
	var err error
	defer gc.ConsoleErrLog(&err)

	r := InitGoodsUpdateRequest()
	if err = api.ParseAndValidate(req, r); err != nil {
		api.WriteErrorResponse(w, http.StatusBadRequest)
		err = fmt.Errorf("UpdateGoodsHandler.%s", err.Error())
		return
	}

	exists, err := gc.checkGood(goodsTableName, r.Id, r.ProjectId)
	if err != nil {
		api.WriteErrorResponse(w, http.StatusInternalServerError)
		err = fmt.Errorf("UpdateGoodsHandler.%s", err.Error())
		return
	}
	if !exists {
		api.WriteErrorResponsePayload(w, http.StatusNotFound, 3, "errors.good.notFound", nil)
		return
	}

	g, err := gc.UpdateGood(goodsTableName, r.Id, r.ProjectId, r.Payload.Name, r.Payload.Description)
	if err != nil {
		api.WriteErrorResponse(w, http.StatusInternalServerError)
		err = fmt.Errorf("UpdateGoodsHandler.%s", err.Error())
		return
	}

	jsonG, err := json.Marshal(g)
	if err != nil {
		api.WriteErrorResponse(w, http.StatusInternalServerError)
		err = fmt.Errorf("CreateGoodHandler.json.Marshal: %s", err.Error())
		return
	}

	go gc.LogGoods(g)
	go gc.c.ChacheByte(generateGoodKey(g), jsonG, gc.cTime)
	w.Write(jsonG)
}

func (gc *GoodsCRUD) DeleteGoodsHandler(w http.ResponseWriter, req *http.Request) {
	var err error
	defer gc.ConsoleErrLog(&err)

	r := InitGoodsDeleteRequest()
	if err = api.ParseAndValidate(req, r); err != nil {
		api.WriteErrorResponse(w, http.StatusBadRequest)
		err = fmt.Errorf("DeleteGoodsHandler.%s", err.Error())
		return
	}

	exists, err := gc.checkGood(goodsTableName, r.Id, r.ProjectId)
	if err != nil {
		api.WriteErrorResponse(w, http.StatusInternalServerError)
		err = fmt.Errorf("DeleteGoodsHandler.%s", err.Error())
		return
	}
	if !exists {
		api.WriteErrorResponsePayload(w, http.StatusNotFound, 3, "errors.good.notFound", nil)
		return
	}

	g, res, err := gc.DeleteGood(goodsTableName, r.Id, r.ProjectId)
	if err != nil {
		api.WriteErrorResponse(w, http.StatusInternalServerError)
		err = fmt.Errorf("DeleteGoodsHandler.%s", err.Error())
		return
	}

	go gc.LogGoods(g)
	api.WriteJson(w, res)
}

func (gc *GoodsCRUD) ListGoodsHandler(w http.ResponseWriter, req *http.Request) {
	var err error
	defer gc.ConsoleErrLog(&err)

	r := InitGoodsListRequest()
	if err = api.ParseAndValidate(req, r); err != nil {
		api.WriteErrorResponse(w, http.StatusBadRequest)
		err = fmt.Errorf("ListGoodsHandler.%s", err.Error())
		return
	}

	var jsonResp []byte
	if jsonResp, err = gc.c.GetByte(generateGoodsListKeyReq(r)); err == nil {
		w.Write(jsonResp)
		return
	} else {
		gc.ConsoleErrLog(&err)
	}

	var resp *GoodsListResponse
	resp, err = gc.ListGoods(goodsTableName, r.Offset, r.Limit)
	if err != nil {
		api.WriteErrorResponse(w, http.StatusInternalServerError)
		err = fmt.Errorf("ListGoodsHandler.%s", err.Error())
		return
	}

	jsonResp, err = json.Marshal(resp)
	if err != nil {
		api.WriteErrorResponse(w, http.StatusInternalServerError)
		err = fmt.Errorf("CreateGoodHandler.json.Marshal: %s", err.Error())
		return
	}

	go gc.c.ChacheByte(generateGoodsListKey(resp), jsonResp, gc.cTime)
	w.Write(jsonResp)
}

func (gc *GoodsCRUD) ReprioritiizeGoodsHandler(w http.ResponseWriter, req *http.Request) {
	var err error
	defer gc.ConsoleErrLog(&err)

	r := InitGoodsReprioritiizeRequest()
	if err = api.ParseAndValidate(req, r); err != nil {
		api.WriteErrorResponse(w, http.StatusBadRequest)
		err = fmt.Errorf("ReprioritiizeGoodsHandler.%s", err.Error())
		return
	}

	exists, err := gc.checkGood(goodsTableName, r.Id, r.ProjectId)
	if err != nil {
		api.WriteErrorResponse(w, http.StatusInternalServerError)
		err = fmt.Errorf("ReprioritiizeGoodsHandler.%s", err.Error())
		return
	}
	if !exists {
		api.WriteErrorResponsePayload(w, http.StatusNotFound, 3, "errors.good.notFound", nil)
		return
	}

	goods, resp, err := gc.ReprioritizeGoods(goodsTableName, r.Id, r.ProjectId, r.Payload.NewPriority)
	if err != nil {
		api.WriteErrorResponse(w, http.StatusInternalServerError)
		err = fmt.Errorf("ReprioritiizeGoodsHandler.%s", err.Error())
		return
	}

	gc.LogGoods(goods...)
	for _, g := range goods {
		go gc.ChacheGood(g)
	}
	err = api.WriteJson(w, resp)
}
