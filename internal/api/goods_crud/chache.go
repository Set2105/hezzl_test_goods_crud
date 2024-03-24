package goods_crud

import (
	"encoding/json"
	"fmt"

	psql_m "github.com/Set2105/hezzl_test_goods_crud/models/postgres"
)

func generateGoodKey(g *psql_m.Good) string {
	return fmt.Sprintf("good_%d_%d", g.Id, g.ProjectId)
}

func generateGoodsListKey(r *GoodsListResponse) string {
	return fmt.Sprintf("goods_list_%d_%d", r.Meta.Limit, r.Meta.Offset)
}

func generateGoodsListKeyReq(r *GoodsListRequest) string {
	return fmt.Sprintf("goods_list_%d_%d", r.Limit, r.Offset)
}

func (gc *GoodsCRUD) ChacheGood(g *psql_m.Good) {
	data, err := json.Marshal(g)
	if err != nil {
		gc.errLog(fmt.Errorf("LogGood.Marshal: %s", err))
	}
	err = gc.c.ChacheByte(generateGoodKey(g), data, gc.cTime)
	if err != nil {
		gc.errLog(fmt.Errorf("LogGood.Publish: %s", err))
	}
}
