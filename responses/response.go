package responses

type Response struct {
	Message string `json:"message"`
	Status  string `json:"status"`
	Data    any    `json:"data"`
}

func ResponseSuccess(message string, data any) Response {
	return Response{
		Message: message,
		Status:  "success",
		Data:    data,
	}
}

func ResponseError(message string) Response {
	return Response{
		Message: message,
		Status:  "error",
		Data:    nil,
	}
}
