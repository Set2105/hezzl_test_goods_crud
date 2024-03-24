package api

type ErrorPayload struct {
	Code     int       `json:"code"`
	Message  string    `json:"message"`
	Detailes *Detailes `json:"detailes"`
}

type Detailes struct {
}
