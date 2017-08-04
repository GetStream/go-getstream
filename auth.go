package getstream

import (
	"fmt"
	"net/http"
)

type AuthenticationMethod uint32

const (
	ApplicationAuthentication AuthenticationMethod = iota
	FeedAuthentication
)

var authenticators = map[AuthenticationMethod]Authenticator{
	ApplicationAuthentication: applicationAuthenticator{},
	FeedAuthentication:        feedAuthenticator{},
}

var scopeActions = map[string]ScopeAction{
	"GET":     ScopeActionRead,
	"OPTIONS": ScopeActionRead,
	"HEAD":    ScopeActionRead,
	"POST":    ScopeActionWrite,
	"PUT":     ScopeActionWrite,
	"PATCH":   ScopeActionWrite,
	"DELETE":  ScopeActionDelete,
}

type Authenticator interface {
	Authenticate(signer *Signer, request *http.Request, context ScopeContext, feed Feed) error
}

type applicationAuthenticator struct {
	key    string
	secret string
}

func (a applicationAuthenticator) Authenticate(signer *Signer, request *http.Request, context ScopeContext, feed Feed) error {
	return signer.SignHTTP(request)
}

type feedAuthenticator struct {
	signer *Signer
}

func (a feedAuthenticator) Authenticate(signer *Signer, request *http.Request, context ScopeContext, feed Feed) error {
	action, ok := scopeActions[request.Method]
	if !ok {
		return fmt.Errorf("missing action")
	}
	return signer.SignJWT(request, context, action, feed.FeedIDWithoutColon())
}
