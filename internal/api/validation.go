package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type JsonRequest interface {
	PointerMap() map[string]any
	Validate() error
	GetPayload() any
}

func ParseAndValidate(req *http.Request, r JsonRequest) error {
	if err := ParseHttpRequest(req, r); err != nil {
		return fmt.Errorf("ParseAndValidate.%s", err.Error())
	}
	if err := r.Validate(); err != nil {
		return fmt.Errorf("ParseAndValidate.%s", err.Error())
	}
	return nil
}

func ParseHttpRequest(req *http.Request, r JsonRequest) error {
	if err := parseParamsHttpRequest(req, r.PointerMap()); err != nil {
		return fmt.Errorf("ParseHttpRequest.%s", err.Error())
	}
	if err := parseMsgFromJsonBody(req, r.GetPayload()); err != nil {
		return fmt.Errorf("ParseHttpRequest.%s", err.Error())
	}
	return nil
}

func parseMsgFromJsonBody(req *http.Request, payload any) error {
	dec := json.NewDecoder(req.Body)
	if err := dec.Decode(payload); err != nil {
		return fmt.Errorf("parseMsgFromJsonBody.json.Decode: %s", err.Error())
	} else {
		return nil
	}
}

func parseParamsHttpRequest(req *http.Request, urlParseParams map[string]any) error {
	if urlParseParams == nil {
		return nil
	}
	params, err := url.ParseQuery(req.URL.RawQuery)
	if err != nil {
		return fmt.Errorf("parseParamsHttpRequest.url.ParseQuery: %s", err.Error())
	}
	for key, pointer := range urlParseParams {
		val := params.Get(key)
		if val == "" {
			continue
		}
		switch t := pointer.(type) {
		case *string:
			*pointer.(*string) = val
		case *int:
			intVal, err := strconv.Atoi(val)
			if err != nil {
				return fmt.Errorf("parseParamsHttpRequest.strconv.Atoi: %s", err.Error())
			}
			*pointer.(*int) = intVal
		case *int64:
			intVal, err := strconv.ParseInt(val, 10, 0)
			if err != nil {
				return fmt.Errorf("parseParamsHttpRequest.strconv.Atoi: %s", err.Error())
			}
			*pointer.(*int64) = intVal

		default:
			return fmt.Errorf("parseParamsHttpRequest: urlParseParams[\"%s\"] has unsupported type [\"%s\"]", key, t)
		}
	}
	return nil
}
