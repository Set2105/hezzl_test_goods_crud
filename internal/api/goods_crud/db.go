package goods_crud

import (
	"fmt"

	psql_m "github.com/Set2105/hezzl_test_goods_crud/models/postgres"
)

// common

func (gc *GoodsCRUD) checkGood(tableName string, id int64, projectId int64) (b bool, err error) {
	rows, err := gc.db.Db.Query(ExistsGoodStatement(tableName, id, projectId))
	if err != nil {
		return b, fmt.Errorf("CheckGood.Query.ExistsGoodStatement: %s", err)
	}
	rows.Next()
	err = rows.Scan(&b)
	if err != nil {
		return b, fmt.Errorf("CheckGood.Scan: %s", err)
	}
	return b, nil
}

// CRUD

func (gc *GoodsCRUD) CreateGood(tableName string, projectId int64, name string) (*psql_m.Good, error) {
	db := gc.db.Db
	rows, err := db.Query(CreateGoodStatement(tableName, projectId, name))
	if err != nil {
		return nil, fmt.Errorf("CreateGood.Query: %s", err)
	}
	var g psql_m.Good
	rows.Next()
	err = rows.Scan(&g.Id, &g.ProjectId, &g.Name, &g.Description, &g.Priority, &g.Removed, &g.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("CreateGood.Scan: %s", err)
	}
	return &g, nil
}

func (gc *GoodsCRUD) UpdateGood(tableName string, id int64, projectId int64, name string, description string) (GoodsUpdateResponse, error) {
	db := gc.db.Db
	tx, err := db.Begin()
	defer tx.Rollback()
	if err != nil {
		return nil, fmt.Errorf("UpdateGood.Db.Begin: %s", err)
	}
	_, err = tx.Exec(LockRowStatement(tableName))
	if err != nil {
		return nil, fmt.Errorf("UpdateGood.Tx.Exec.LockRowStatement: %s", err)
	}
	_, err = tx.Exec(SelectForStatement(tableName, id, projectId, "UPDATE"))
	if err != nil {
		return nil, fmt.Errorf("UpdateGood.Tx.Exec.SelectForStatement: %s", err)
	}
	rows, err := tx.Query(UpdateGoodStatement(tableName, name, description, id, projectId))
	if err != nil {
		return nil, fmt.Errorf("UpdateGood.Tx.Query.UpdateGoodStatement: %s", err)
	}
	var g psql_m.Good
	rows.Next()
	err = rows.Scan(&g.Id, &g.ProjectId, &g.Name, &g.Description, &g.Priority, &g.Removed, &g.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("UpdateGood.Scan: %s", err)
	}
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("UpdateGood.Tx.Commit: %s", err)
	}
	return &g, nil
}

func (gc *GoodsCRUD) DeleteGood(tableName string, id int64, projectId int64) (*psql_m.Good, *GoodsDeleteResponse, error) {
	db := gc.db.Db
	tx, err := db.Begin()
	defer tx.Rollback()
	if err != nil {
		return nil, nil, fmt.Errorf("DeleteGood.Db.Begin: %s", err)
	}
	_, err = tx.Exec(LockRowStatement(tableName))
	if err != nil {
		return nil, nil, fmt.Errorf("DeleteGood.Tx.Exec.LockRowStatement: %s", err)
	}
	_, err = tx.Exec(SelectForStatement(tableName, id, projectId, "UPDATE"))
	if err != nil {
		return nil, nil, fmt.Errorf("DeleteGood.Tx.Exec.SelectForStatement: %s", err)
	}
	rows, err := tx.Query(DeleteGoodStatement(tableName, id, projectId))
	if err != nil {
		return nil, nil, fmt.Errorf("DeleteGood.Tx.Query.DeleteGoodStatement: %s", err)
	}
	var g psql_m.Good
	rows.Next()
	err = rows.Scan(&g.Id, &g.ProjectId, &g.Name, &g.Description, &g.Priority, &g.Removed, &g.CreatedAt)
	if err != nil {
		return nil, nil, fmt.Errorf("DeleteGood.Scan: %s", err)
	}
	if err := tx.Commit(); err != nil {
		return nil, nil, fmt.Errorf("DeleteGood.Tx.Commit: %s", err)
	}
	resp := GoodsDeleteResponse{
		Id:        g.Id,
		ProjectId: g.ProjectId,
		Removed:   g.Removed,
	}
	return &g, &resp, nil
}

func (gc *GoodsCRUD) ListGoods(tableName string, offset, limit int64) (resp *GoodsListResponse, err error) {
	db := gc.db.Db
	rows, err := db.Query(ListGoodsStatement(tableName, offset, limit))
	if err != nil {
		return nil, fmt.Errorf("ListGoods.Query: %s", err)
	}

	resp = &GoodsListResponse{Meta: &GoodsListMeta{Offset: offset, Limit: limit}, Goods: []*psql_m.Good{}}
	for rows.Next() {
		var g psql_m.Good
		err = rows.Scan(&g.Id, &g.ProjectId, &g.Name, &g.Description, &g.Priority, &g.Removed, &g.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("ListGoods.Scan: %s", err)
		}
		resp.Meta.Total++
		if g.Removed {
			resp.Meta.Removed++
		}
		resp.Goods = append(resp.Goods, &g)
	}
	return resp, nil
}

func (gc *GoodsCRUD) ReprioritizeGoods(tableName string, id int64, projectId int64, newPriority int64) ([]*psql_m.Good, *GoodsReprioritiizeResponse, error) {
	db := gc.db.Db
	tx, err := db.Begin()
	defer tx.Rollback()
	if err != nil {
		return nil, nil, fmt.Errorf("ReprioritizeGoods.Db.Begin: %s", err)
	}

	_, err = tx.Exec(LockRowStatement(tableName))
	if err != nil {
		return nil, nil, fmt.Errorf("ReprioritizeGoods.Tx.Exec.LockRowStatement: %s", err)
	}

	_, err = tx.Exec(SelectForStatement(tableName, id, projectId, "UPDATE"))
	if err != nil {
		return nil, nil, fmt.Errorf("ReprioritizeGoods.Tx.Exec.SelectForStatement: %s", err)
	}
	_, err = tx.Exec(SelectForReprioritiizeStatement(tableName, newPriority, "UPDATE"))
	if err != nil {
		return nil, nil, fmt.Errorf("ReprioritizeGoods.Tx.Exec.SelectForStatement: %s", err)
	}

	resp := GoodsReprioritiizeResponse{}
	goods := []*psql_m.Good{}

	// Set other obj priority
	rows, err := tx.Query(ReprioritiizeGoodsStatement(tableName, newPriority))
	if err != nil {
		return nil, nil, fmt.Errorf("ReprioritizeGoods.Tx.Query.ReprioritiizeGoodsStatement: %s", err)
	}
	for rows.Next() {
		var g psql_m.Good
		err = rows.Scan(&g.Id, &g.ProjectId, &g.Name, &g.Description, &g.Priority, &g.Removed, &g.CreatedAt)
		if err != nil {
			return nil, nil, fmt.Errorf("ReprioritizeGoods.rows.Scan: %s", err)
		}
		resp.Priorities = append(resp.Priorities, &Priority{Id: g.Id, Priority: g.Priority})
		goods = append(goods, &g)
	}

	// Set this obj priority
	var g psql_m.Good
	rows, err = tx.Query(SetPriorityGoodStatement(tableName, id, projectId, newPriority))
	if err != nil {
		return nil, nil, fmt.Errorf("ReprioritizeGoods.Tx.Query.UpdateGoodStatement: %s", err)
	}
	rows.Next()
	err = rows.Scan(&g.Id, &g.ProjectId, &g.Name, &g.Description, &g.Priority, &g.Removed, &g.CreatedAt)
	if err != nil {
		return nil, nil, fmt.Errorf("ReprioritizeGoods.rows.Scan: %s", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, nil, fmt.Errorf("ReprioritizeGoods.Tx.Commit: %s", err)
	}
	return goods, &resp, nil
}
