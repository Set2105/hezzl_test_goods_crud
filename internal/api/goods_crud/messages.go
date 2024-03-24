package goods_crud

import (
	"fmt"

	psql_m "github.com/Set2105/hezzl_test_goods_crud/models/postgres"
)

// -----------------------------------------------
// POST URL:projectID=int /good/create/
// -----------------------------------------------

type GoodsCreateRequest struct {
	ProjectId int64
	Payload   *GoodsCreatePayload
}

func InitGoodsCreateRequest() *GoodsCreateRequest {
	var r GoodsCreateRequest
	var p GoodsCreatePayload
	r.Payload = &p
	return &r
}

func (r *GoodsCreateRequest) Validate() error {
	if r.ProjectId == 0 {
		return fmt.Errorf("GoodsCreateRequest.Validate: projectId is 0")
	}
	if err := r.Payload.Validate(); err != nil {
		return fmt.Errorf("GoodsCreateRequest.Validate.%s", err.Error())
	}
	return nil
}

func (r *GoodsCreateRequest) GetPayload() any {
	return r.Payload
}

func (r *GoodsCreateRequest) PointerMap() map[string]any {
	return map[string]any{
		"projectId": &r.ProjectId,
	}
}

type GoodsCreatePayload struct {
	Name string `json:"name"`
}

func (p *GoodsCreatePayload) Validate() error {
	if p.Name == "" {
		return fmt.Errorf("GoodsCreatePayload.Validate: name field is empty")
	}
	return nil
}

type GoodsCreateResponse *psql_m.Good

// -----------------------------------------------
// PATCH URL:id=int&projectId=int /good/update
// -----------------------------------------------

type GoodsUpdateRequest struct {
	Id        int64
	ProjectId int64
	Payload   *GoodsUpdatePayload
}

func InitGoodsUpdateRequest() *GoodsUpdateRequest {
	var r GoodsUpdateRequest
	var p GoodsUpdatePayload
	r.Payload = &p

	return &r
}

func (r *GoodsUpdateRequest) Validate() error {
	if r.Id == 0 {
		return fmt.Errorf("GoodsUpdateRequest.Validate.Payload: Id is 0")
	}
	if r.ProjectId == 0 {
		return fmt.Errorf("GoodsUpdateRequest.Validate.Payload: ProjectId is 0")
	}
	if r.Payload.Name == "" {
		return fmt.Errorf("GoodsUpdateRequest.Validate.Payload: name field is empty")
	}
	return nil
}

func (r *GoodsUpdateRequest) PointerMap() map[string]any {
	return map[string]any{
		"projectId": &r.ProjectId,
		"id":        &r.Id,
	}
}

func (r *GoodsUpdateRequest) GetPayload() any {
	return &r.Payload
}

type GoodsUpdatePayload struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type GoodsUpdateResponse *psql_m.Good

// -----------------------------------------------
// DELETE URL:id=int&projectId=int /good/remove
// -----------------------------------------------

type GoodsDeleteRequest struct {
	Id        int64
	ProjectId int64
	Payload   *GoodsDeletePayload
}

func InitGoodsDeleteRequest() *GoodsDeleteRequest {
	var r GoodsDeleteRequest
	var p GoodsDeletePayload
	r.Payload = &p

	return &r
}

func (r *GoodsDeleteRequest) Validate() error {
	return nil
}

func (r *GoodsDeleteRequest) PointerMap() map[string]any {
	return map[string]any{
		"projectId": &r.ProjectId,
		"id":        &r.Id,
	}
}

func (r *GoodsDeleteRequest) GetPayload() any {
	return &r.Payload
}

type GoodsDeletePayload struct {
}

type GoodsDeleteResponse struct {
	Id        int64 `json:"id"`
	ProjectId int64 `json:"campaignId"`
	Removed   bool  `json:"removed"`
}

// -----------------------------------------------
// GET URL:offset=int&limit=int /goods/list
// -----------------------------------------------

type GoodsListRequest struct {
	Limit   int64
	Offset  int64
	Payload *GoodsListPayload
}

func InitGoodsListRequest() *GoodsListRequest {
	var r GoodsListRequest
	var p GoodsListPayload
	r.Payload = &p

	return &r
}

func (r *GoodsListRequest) PointerMap() map[string]any {
	return map[string]any{
		"limit":  &r.Limit,
		"offset": &r.Offset,
	}
}

func (r *GoodsListRequest) Validate() error {
	if r.Limit == 0 {
		r.Limit = 1
	}
	if r.Offset == 0 {
		r.Offset = 10
	}
	return nil
}

func (r *GoodsListRequest) GetPayload() any {
	return &r.Payload
}

type GoodsListPayload struct {
}

type GoodsListResponse struct {
	Meta  *GoodsListMeta `json:"meta"`
	Goods []*psql_m.Good `json:"goods"`
}

type GoodsListMeta struct {
	Total   int64 `json:"total"`
	Removed int64 `json:"removed"`
	Limit   int64 `json:"limit"`
	Offset  int64 `json:"offset"`
}

// -----------------------------------------------
// PATCH URL:id=int&projectId=int /good/reprioritiize
// -----------------------------------------------

type GoodsReprioritiizeRequest struct {
	Id        int64
	ProjectId int64
	Payload   *GoodsReprioritiizePayload
}

func InitGoodsReprioritiizeRequest() *GoodsReprioritiizeRequest {
	var r GoodsReprioritiizeRequest
	var p GoodsReprioritiizePayload
	r.Payload = &p

	return &r
}

func (r *GoodsReprioritiizeRequest) Validate() error {
	if r.Id == 0 {
		return fmt.Errorf("GoodsReprioritiizeRequest.Validate: incorrect id")
	}
	if r.ProjectId == 0 {
		return fmt.Errorf("GoodsReprioritiizeRequest.Validate: incorrect projectId")
	}
	if r.Payload.NewPriority == 0 {
		return fmt.Errorf("GoodsReprioritiizeRequest.Validate: incorrect newPriority")
	}
	return nil
}

func (r *GoodsReprioritiizeRequest) PointerMap() map[string]any {
	return map[string]any{
		"projectId": &r.ProjectId,
		"id":        &r.Id,
	}
}

func (r *GoodsReprioritiizeRequest) GetPayload() any {
	return &r.Payload
}

type GoodsReprioritiizePayload struct {
	NewPriority int64 `json:"newPriority"`
}

type GoodsReprioritiizeResponse struct {
	Priorities []*Priority `json:"priorities"`
}

type Priority struct {
	Id       int64 `json:"id"`
	Priority int64 `json:"priority"`
}
