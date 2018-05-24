package httphandler

import (
	"net/http"

	"github.com/Brickchain/go-crypto.v2"
	"github.com/pkg/errors"
	controller "gitlab.brickchain.com/libs/go-controller.v2"
)

// AuthenticatedRequestWithBinding extends the AuthenticatedRequest with information about the active Binding
type AuthenticatedRequestWithBinding interface {
	AuthenticatedRequest
	Binding() controller.Binding
}

// RequestWithBinding extends the Request with information about the active Binding
type RequestWithBinding interface {
	Request
	Binding() controller.Binding
}

// standardRequestWithBinding is a basic implementation of RequestWithBinding
type standardRequestWithBinding struct {
	Request
	binding controller.Binding
}

// Binding returns the Binding object
func (s *standardRequestWithBinding) Binding() controller.Binding {
	return s.binding
}

// standardAuthenticatedRequestWithBinding is a basic implementation of AuthenticatedRequestWithBinding
type standardAuthenticatedRequestWithBinding struct {
	AuthenticatedRequest
	binding controller.Binding
}

// Binding returns the Binding object
func (s *standardAuthenticatedRequestWithBinding) Binding() controller.Binding {
	return s.binding
}

// ActionRequestWithBinding extends the ActionRequest with information about the active Binding
type ActionRequestWithBinding interface {
	ActionRequest
	Binding() controller.Binding
}

// standardActionRequestWithBinding is a basic implementation of ActionRequestWithBinding
type standardActionRequestWithBinding struct {
	ActionRequest
	binding controller.Binding
}

// Binding returns the Binding object
func (s *standardActionRequestWithBinding) Binding() controller.Binding {
	return s.binding
}

// addBinding is the wrapper to lookup the active binding and add to the request object
func addBinding(bm controller.BindingService, h func(RequestWithBinding) Response) func(Request) Response {
	return func(req Request) Response {
		bindID := req.URL().Query().Get("binding")
		if bindID == "" {
			return NewErrorResponse(http.StatusBadRequest, errors.New("No binding in request"))
		}

		binding, err := bm.Get(bindID)
		if err != nil {
			return NewErrorResponse(http.StatusInternalServerError, errors.Wrap(err, "could not lookup binding"))
		}

		return h(&standardRequestWithBinding{
			Request: req,
			binding: binding,
		})
	}
}

// addAuthenticatedBinding is the wrapper to lookup the active binding and add to the request object
func addAuthenticatedBinding(bm controller.BindingService, h func(AuthenticatedRequestWithBinding) Response) func(AuthenticatedRequest) Response {
	return func(req AuthenticatedRequest) Response {
		bindID := req.URL().Query().Get("binding")
		if bindID == "" {
			return NewErrorResponse(http.StatusBadRequest, errors.New("No binding in request"))
		}

		binding, err := bm.Get(bindID)
		if err != nil {
			return NewErrorResponse(http.StatusInternalServerError, errors.Wrap(err, "could not lookup binding"))
		}

		realmTP := crypto.Thumbprint(binding.Realm().PublicKey)

		success := false
		for _, mandate := range req.Mandates() {
			signerTP := crypto.Thumbprint(mandate.Signer)
			if realmTP == signerTP && hasRole(binding.AdminRoles(), mandate.Mandate.Role) {
				success = true
			}
		}

		if !success {
			return NewErrorResponse(http.StatusForbidden, errors.New("Mandate not signed by binding realm"))
		}

		return h(&standardAuthenticatedRequestWithBinding{
			AuthenticatedRequest: req,
			binding:              binding,
		})
	}
}

// addActionBinding is the wrapper to lookup the active binding and add to the request object
func addActionBinding(bm controller.BindingService, h func(ActionRequestWithBinding) Response) func(ActionRequest) Response {
	return func(req ActionRequest) Response {
		bindID := req.URL().Query().Get("binding")
		if bindID == "" {
			if req.Action().Params == nil {
				return NewErrorResponse(http.StatusBadRequest, errors.New("No binding in request"))
			}

			var ok bool
			bindID, ok = req.Action().Params["binding"]
			if !ok {
				return NewErrorResponse(http.StatusBadRequest, errors.New("No binding in request"))
			}
		}

		binding, err := bm.Get(bindID)
		if err != nil {
			return NewErrorResponse(http.StatusInternalServerError, errors.Wrap(err, "could not lookup binding"))
		}

		realmTP := crypto.Thumbprint(binding.Realm().PublicKey)

		success := false
		for _, mandate := range req.Mandates() {
			signerTP := crypto.Thumbprint(mandate.Signer)
			if realmTP == signerTP && hasRole(binding.AdminRoles(), mandate.Mandate.Role) {
				success = true
			}
		}

		if !success {
			return NewErrorResponse(http.StatusForbidden, errors.New("Mandate not signed by binding realm"))
		}

		return h(&standardActionRequestWithBinding{
			ActionRequest: req,
			binding:       binding,
		})
	}
}

func hasRole(roleList []string, role string) bool {
	for _, r := range roleList {
		if r == role {
			return true
		}
	}

	return false
}
