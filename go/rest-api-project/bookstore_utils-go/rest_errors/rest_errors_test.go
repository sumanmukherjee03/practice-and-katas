package rest_errors

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewError(t *testing.T) {
	assert := assert.New(t)
	err := NewError("Testing new error")
	assert.NotNil(err)
	assert.EqualValues("Testing new error", err.Error())
}

func TestNewBadRequestError(t *testing.T) {
	assert := assert.New(t)
	err := NewBadRequestError(fmt.Errorf("Testing bad request"))
	assert.NotNil(err)
	assert.EqualValues(http.StatusBadRequest, err.Status)
	assert.EqualValues("bad_request", err.Error)
}

func TestNewNotFoundError(t *testing.T) {
	assert := assert.New(t)
	err := NewNotFoundError(fmt.Errorf("Testing not found"))
	assert.NotNil(err)
	assert.EqualValues(http.StatusNotFound, err.Status)
	assert.EqualValues("not_found", err.Error)
}

func TestNewInternalServerError(t *testing.T) {
	assert := assert.New(t)
	err := NewInternalServerError(fmt.Errorf("Testing internal server error"))
	assert.NotNil(err)
	assert.EqualValues(http.StatusInternalServerError, err.Status)
	assert.EqualValues("internal_server_error", err.Error)
}
