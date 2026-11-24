package binding

import (
	"net/http"

	"github.com/go-playground/form/v4"
)

var formDecoder = form.NewDecoder()

type queryBinding struct{}

func (queryBinding) Name() string {
	return "URL-QUERY"
}

func (queryBinding) Bind(req *http.Request, obj any) error {
	if req == nil || req.URL == nil {
		return newBindingError("invalid request")
	}
	values := req.URL.Query()
	err := formDecoder.Decode(obj, values)
	if err != nil {
		return newBindingError(err.Error())
	}
	return validate(obj)
}
