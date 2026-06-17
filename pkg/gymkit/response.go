package gymkit

import (
	"encoding/json"

	"github.com/kataras/iris/v12"
)

type Detail struct {
	Error interface{} `json:"errors,omitempty"`
	Data  interface{} `json:"data,omitempty"`
}

type Messages struct {
	Status  int  `json:"status"`
	Success bool `json:"success"`
	*Detail
}

// NewResponse sends a standard JSON response.
func NewResponse(ctx iris.Context, code int, content interface{}) {
	status := IsSuccessCode(code)
	var res Messages
	res.Status = code
	res.Success = status

	body := make(map[string]interface{})
	var returns interface{}

	if s, ok := content.(string); ok {
		if IsJSON(s) {
			json.Unmarshal([]byte(s), &body)
		} else {
			body["message"] = s
		}
	} else {
		returns = content
	}

	if status {
		if returns != nil {
			res.Detail = &Detail{Data: returns}
		} else {
			res.Detail = &Detail{Data: body}
		}
	} else {
		if returns != nil {
			res.Detail = &Detail{Error: returns}
		} else {
			res.Detail = &Detail{Error: body}
		}
	}

	if code == iris.StatusNoContent {
		code = iris.StatusOK
	}
	ctx.StatusCode(code)
	ctx.JSON(res)
}
