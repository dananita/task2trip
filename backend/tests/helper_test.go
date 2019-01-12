package tests

import (
	"fmt"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
)

type TokenAuth struct {
	AuthToken string
}

func (a *TokenAuth) AuthenticateRequest(r runtime.ClientRequest, _ strfmt.Registry) error {
	return r.SetHeaderParam("Authorization", fmt.Sprintf("Bearer %s", a.AuthToken))
}
