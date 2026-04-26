package apiresponse

type Response struct {
	Meta    interface{} `json:"meta,omitempty"`
	Message interface{} `json:"message"`
	Status  uint16      `json:"status"`
	Data    interface{} `json:"data,omitempty"`
	Error   XError      `json:"error,omitempty"`
}

type XError struct {
	Code    uint16 `json:"code"`
	Message string `json:"message"`
	Status  bool   `json:"status"`
}
