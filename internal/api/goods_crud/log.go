package goods_crud

import (
	"encoding/json"
	"fmt"
	"time"

	ch_m "github.com/Set2105/hezzl_test_goods_crud/models/clickhouse"
	psql_m "github.com/Set2105/hezzl_test_goods_crud/models/postgres"
)

var chSubj = "goods_log"

func (gc *GoodsCRUD) ConsoleErrLog(errPointer *error) {
	if gc.l != nil {
		var err error
		if errPointer != nil {
			err = *errPointer
		} else {
			err = fmt.Errorf("ConsoleErrLog: error nil pointer")
		}
		if err != nil {
			gc.errLog(err)
		}
	}
}

func (gc *GoodsCRUD) errLog(err error) {
	gc.l.Println(err)
}

func (gc *GoodsCRUD) LogGoods(goods ...*psql_m.Good) error {
	goodsLog := chLogFromGoods(goods...)
	for _, gl := range goodsLog {
		go gc.LogGood(gl)
	}
	return nil
}

func (gc *GoodsCRUD) LogGood(g *ch_m.GoodsLog) error {
	data, err := json.Marshal(g)
	if err != nil {
		gc.errLog(fmt.Errorf("LogGood.Marshal: %s", err))
	}
	err = gc.n.Publish(chSubj, data)
	if err != nil {
		gc.errLog(fmt.Errorf("LogGood.Publish: %s", err))
	}
	return nil
}

func chLogFromGoods(goods ...*psql_m.Good) []*ch_m.GoodsLog {
	t := time.Now()
	goodLogs := make([]*ch_m.GoodsLog, len(goods))
	for i, g := range goods {
		goodLogs[i] = chLogFromGood(g, t)
	}
	return goodLogs
}

func chLogFromGood(g *psql_m.Good, t time.Time) *ch_m.GoodsLog {
	if g == nil {
		return nil
	}
	return &ch_m.GoodsLog{
		Id:          g.Id,
		ProjectId:   g.ProjectId,
		Name:        g.Name,
		Description: g.Description,
		Priority:    g.Priority,
		Removed:     g.Removed,
		EventTime:   t,
	}
}
