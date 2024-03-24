package goods_crud

import (
	"fmt"
)

const goods_fields = "id,project_id,\"name\",description,priority,removed,created_at"

// Exists

func ExistsGoodStatement(table string, id int64, projectId int64) string {
	s := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE id=%d AND project_id=%d AND removed=false)",
		table, id, projectId)
	return s + ";"
}

// Transactions

func LockRowStatement(table string) string {
	s := fmt.Sprintf("LOCK TABLE %s IN ROW EXCLUSIVE MODE", table)
	return s + ";"
}

func SelectForStatement(tableName string, id int64, projectId int64, forWhat string) string {
	s := fmt.Sprintf("SELECT 1 FROM %s WHERE id=%d AND project_id=%d FOR %s", tableName, id, projectId, forWhat)
	return s + ";"
}

// Create

func CreateGoodStatement(table string, projectId int64, name string) string {
	s := fmt.Sprintf(`INSERT INTO %s (project_id,"name") VALUES (%d,'%s') RETURNING %s`,
		table, projectId, name, goods_fields)
	return s + ";"
}

// Read

func ListGoodsStatement(table string, Offset int64, Limit int64) string {
	s := fmt.Sprintf(`SELECT %s FROM %s ORDER BY id, project_id LIMIT %d OFFSET %d`,
		goods_fields, table, Limit, Offset)
	return s + ";"
}

// Delete

func DeleteGoodStatement(table string, id int64, projectId int64) string {
	s := fmt.Sprintf("UPDATE %s SET removed = TRUE WHERE id=%d AND project_id=%d RETURNING %s",
		table, id, projectId, goods_fields)
	return s + ";"
}

// Update

func UpdateGoodStatement(table string, name string, description string, id int64, projectId int64) string {
	descLine := ""
	if description != "" {
		descLine = " AND description = '" + description + "'"
	}
	s := fmt.Sprintf("UPDATE %s SET \"name\"='%s' %s WHERE id=%d AND project_id=%d RETURNING %s",
		table, name, descLine, id, projectId, goods_fields)
	return s + ";"
}

// Reprioritize

func ReprioritiizeGoodsStatement(tableName string, priority int64) string {
	s := fmt.Sprintf("UPDATE %s SET priority=priority+1 WHERE priority>=%d RETURNING %s",
		tableName, priority, goods_fields)
	return s + ";"
}

func SetPriorityGoodStatement(tableName string, id int64, projectId int64, newPriority int64) string {
	s := fmt.Sprintf("UPDATE %s SET priority=%d WHERE id=%d AND project_id=%d RETURNING %s",
		tableName, newPriority, id, projectId, goods_fields)
	return s + ";"
}

func SelectForReprioritiizeStatement(tableName string, priority int64, forWhat string) string {
	s := fmt.Sprintf("SELECT * FROM %s WHERE priority>=%d FOR %s", tableName, priority, forWhat)
	return s + ";"
}
