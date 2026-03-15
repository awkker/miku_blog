package dto

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type PagedData struct {
	Items interface{} `json:"items"`
	Total int64       `json:"total"`
	Page  int         `json:"page"`
	Size  int         `json:"size"`
}

func OK(data interface{}) *Response {
	return &Response{Code: 0, Message: "ok", Data: data}
}

func OKPaged(items interface{}, total int64, page, size int) *Response {
	return &Response{
		Code:    0,
		Message: "ok",
		Data: PagedData{
			Items: items,
			Total: total,
			Page:  page,
			Size:  size,
		},
	}
}

func Err(code int, msg string) *Response {
	return &Response{Code: code, Message: msg}
}
