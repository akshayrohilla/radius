// ------------------------------------------------------------
// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.
// ------------------------------------------------------------

package handlerv3

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Azure/radius/pkg/azure/azresources"
	"github.com/Azure/radius/pkg/radlogger"
	"github.com/Azure/radius/pkg/radrp/armerrors"
	"github.com/Azure/radius/pkg/radrp/frontend/resourceproviderv3"
	"github.com/Azure/radius/pkg/radrp/rest"
	"github.com/Azure/radius/pkg/radrp/schemav3"
)

// A brief note on error handling... The handler is responsible for all of the direct actions
// with HTTP request/reponse.
//
// The RP returns the rest.Response type for "known" or "expected" error conditions:
// - validation error
// - missing data
//
// The RP returns an error for "unexpected" error conditions:
// - DB failure
// - I/O failure
//
// This code will assume that any Golang error returned from the RP represents a reliability error
// within the RP or a bug.

type handler struct {
	rp               resourceproviderv3.ResourceProvider
	validatorFactory ValidatorFactory
}

func (h *handler) ListApplications(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	response, err := h.rp.ListApplications(ctx, resourceID(req))
	if err != nil {
		internalServerError(ctx, w, req, err)
		return
	}

	err = response.Apply(ctx, w, req)
	if err != nil {
		internalServerError(ctx, w, req, err)
		return
	}
}

func (h *handler) GetApplication(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	response, err := h.rp.GetApplication(ctx, resourceID(req))
	if err != nil {
		internalServerError(ctx, w, req, err)
		return
	}

	err = response.Apply(ctx, w, req)
	if err != nil {
		internalServerError(ctx, w, req, err)
		return
	}
}

func (h *handler) UpdateApplication(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	body, err := readJSONBody(req)
	if err != nil {
		badRequest(ctx, w, req, err)
		return
	}

	id := resourceID(req)
	validator, err := h.findValidator(id)
	if err != nil {
		badRequest(ctx, w, req, err)
		return
	}

	validationErrs := validator.ValidateJSON(body)
	if len(validationErrs) > 0 {
		validationError(ctx, w, req, validationErrs)
		return
	}

	response, err := h.rp.UpdateApplication(ctx, id, body)
	if err != nil {
		internalServerError(ctx, w, req, err)
		return
	}

	err = response.Apply(ctx, w, req)
	if err != nil {
		internalServerError(ctx, w, req, err)
		return
	}
}

func (h *handler) DeleteApplication(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	response, err := h.rp.DeleteApplication(ctx, resourceID(req))
	if err != nil {
		internalServerError(ctx, w, req, err)
		return
	}

	err = response.Apply(ctx, w, req)
	if err != nil {
		internalServerError(ctx, w, req, err)
		return
	}
}

func (h *handler) ListAllV3ResourcesByApplication(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	response, err := h.rp.ListAllV3ResourcesByApplication(ctx, resourceID(req))
	if err != nil {
		internalServerError(ctx, w, req, err)
		return
	}

	err = response.Apply(ctx, w, req)
	if err != nil {
		internalServerError(ctx, w, req, err)
		return
	}
}

func (h *handler) ListResources(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	response, err := h.rp.ListResources(ctx, resourceID(req))
	if err != nil {
		internalServerError(ctx, w, req, err)
		return
	}

	err = response.Apply(ctx, w, req)
	if err != nil {
		internalServerError(ctx, w, req, err)
		return
	}
}

func (h *handler) GetResource(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	response, err := h.rp.GetResource(ctx, resourceID(req))
	if err != nil {
		internalServerError(ctx, w, req, err)
		return
	}

	err = response.Apply(ctx, w, req)
	if err != nil {
		internalServerError(ctx, w, req, err)
		return
	}
}

func (h *handler) UpdateResource(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	body, err := readJSONBody(req)
	if err != nil {
		badRequest(ctx, w, req, err)
		return
	}

	id := resourceID(req)
	validator, err := h.findValidator(id)
	if err != nil {
		badRequest(ctx, w, req, err)
		return
	}

	validationErrs := validator.ValidateJSON(body)
	if len(validationErrs) > 0 {
		validationError(ctx, w, req, validationErrs)
		return
	}

	response, err := h.rp.UpdateResource(ctx, id, body)
	if err != nil {
		internalServerError(ctx, w, req, err)
		return
	}

	err = response.Apply(ctx, w, req)
	if err != nil {
		internalServerError(ctx, w, req, err)
		return
	}
}

func (h *handler) DeleteResource(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	response, err := h.rp.DeleteResource(ctx, resourceID(req))
	if err != nil {
		internalServerError(ctx, w, req, err)
		return
	}

	err = response.Apply(ctx, w, req)
	if err != nil {
		internalServerError(ctx, w, req, err)
		return
	}
}

func (h *handler) GetOperation(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	response, err := h.rp.GetOperation(ctx, resourceID(req))
	if err != nil {
		internalServerError(ctx, w, req, err)
		return
	}

	err = response.Apply(ctx, w, req)
	if err != nil {
		internalServerError(ctx, w, req, err)
		return
	}
}

func (h *handler) ListSecrets(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	body, err := readJSONBody(req)
	if err != nil {
		badRequest(ctx, w, req, err)
		return
	}

	input := resourceproviderv3.ListSecretsInput{}
	err = json.Unmarshal(body, &input)
	if err != nil {
		badRequest(ctx, w, req, err)
		return
	}

	response, err := h.rp.ListSecrets(ctx, input)
	if err != nil {
		internalServerError(ctx, w, req, err)
		return
	}

	err = response.Apply(ctx, w, req)
	if err != nil {
		internalServerError(ctx, w, req, err)
		return
	}
}

func (h *handler) findValidator(id azresources.ResourceID) (schemav3.Validator, error) {
	resourceType := id.Types[len(id.Types)-1].Type
	return h.validatorFactory(resourceType)
}

func resourceID(req *http.Request) azresources.ResourceID {
	logger := radlogger.GetLogger(req.Context())
	id, err := azresources.Parse(req.URL.Path)
	if err != nil {
		logger.Info("URL was not a valid resource id: %v", req.URL.Path)
		// just log the error - it will be handled in the RP layer.
	}
	return id
}

func validationError(ctx context.Context, w http.ResponseWriter, req *http.Request, validationErrs []schemav3.ValidationError) {
	logger := radlogger.GetLogger(ctx)

	body := armerrors.ErrorResponse{
		Error: armerrors.ErrorDetails{
			Code:    armerrors.Invalid,
			Message: "Validation error",
			Details: make([]armerrors.ErrorDetails, len(validationErrs)),
		},
	}

	for i, err := range validationErrs {
		if err.JSONError != nil {
			// The given document isn't even JSON.
			body.Error.Details[i].Message = fmt.Sprintf("%s: %v", err.Message, err.JSONError)
		} else {
			body.Error.Details[i].Message = fmt.Sprintf("%s: %s", err.Position, err.Message)
		}
	}

	response := rest.NewBadRequestARMResponse(body)
	err := response.Apply(ctx, w, req)
	if err != nil {
		// There's no way to recover if we fail writing here, we likly partially wrote to the response stream.
		w.WriteHeader(http.StatusInternalServerError)
		logger.Error(err, fmt.Sprintf("error writing marshaled %T bytes to output", body))
	}
}

func badRequest(ctx context.Context, w http.ResponseWriter, req *http.Request, err error) {
	logger := radlogger.GetLogger(ctx)
	// Try to use the ARM format to send back the error info
	body := armerrors.ErrorResponse{
		Error: armerrors.ErrorDetails{
			Code:    armerrors.Invalid,
			Message: err.Error(),
		},
	}

	response := rest.NewBadRequestARMResponse(body)
	err = response.Apply(ctx, w, req)
	if err != nil {
		// There's no way to recover if we fail writing here, we likly partially wrote to the response stream.
		w.WriteHeader(http.StatusInternalServerError)
		logger.Error(err, fmt.Sprintf("error writing marshaled %T bytes to output", body))
	}
}

// Responds with an HTTP 500
func internalServerError(ctx context.Context, w http.ResponseWriter, req *http.Request, err error) {
	logger := radlogger.GetLogger(ctx)
	logger.Error(err, "unhandled error")

	// Try to use the ARM format to send back the error info
	body := armerrors.ErrorResponse{
		Error: armerrors.ErrorDetails{
			Message: err.Error(),
		},
	}

	response := rest.NewInternalServerErrorARMResponse(body)
	err = response.Apply(ctx, w, req)
	if err != nil {
		// There's no way to recover if we fail writing here, we likly partially wrote to the response stream.
		w.WriteHeader(http.StatusInternalServerError)
		logger.Error(err, fmt.Sprintf("error writing marshaled %T bytes to output", body))
	}
}

func readJSONBody(req *http.Request) ([]byte, error) {
	defer req.Body.Close()
	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading request body: %w", err)
	}

	return data, nil
}
