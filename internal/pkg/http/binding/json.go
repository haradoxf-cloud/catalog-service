package binding

import (
	"net/http"

	"github.com/haradoxf-cloud/catalog-service/internal/pkg/http/httph"
)

type jsonBinding struct{}

func (jsonBinding) Name() string {
	return "JSON"
}

func (jsonBinding) Bind(req *http.Request, obj any) error {
	if req == nil || req.Body == nil {
		return newBindingError("invalid request")
	}

	err := httph.DecodeJSON(req, obj)
	if err != nil {
		return newBindingError(err.Error())
	}
	return validate(obj)
}
