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
	http.MethodGet:     ScopeActionRead,
	http.MethodOptions: ScopeActionRead,
	http.MethodHead:    ScopeActionRead,
	http.MethodPost:    ScopeActionWrite,
	http.MethodPut:     ScopeActionWrite,
	http.MethodPatch:   ScopeActionWrite,
	http.MethodDelete:  ScopeActionDelete,
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
