package http

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sumanmukherjee03/practice-and-katas/go/rest-api-project/bookstore_oauth-api/src/domain/access_token"
	"github.com/sumanmukherjee03/practice-and-katas/go/rest-api-project/bookstore_oauth-api/src/utils/errors"
)

type AccessTokenHandler interface {
	GetById(*gin.Context)
	Create(*gin.Context)
	UpdateExpirationTime(*gin.Context)
}

type accessTokenHandler struct {
	service access_token.Service
}

func NewHandler(service access_token.Service) AccessTokenHandler {
	return &accessTokenHandler{
		service: service,
	}
}

func (h *accessTokenHandler) GetById(ctx *gin.Context) {
	atId := ctx.Param("access_token_id")
	at, err := h.service.GetById(atId)
	if err != nil {
		ctx.JSON(err.Status, err)
		return
	}
	ctx.JSON(http.StatusOK, at)
}

func (h *accessTokenHandler) Create(ctx *gin.Context) {
	var at access_token.AccessToken
	if err := ctx.ShouldBindJSON(&at); err != nil {
		restErr := errors.NewBadRequestError(fmt.Errorf("Bad input passed to create access token : %v", err))
		ctx.JSON(restErr.Status, restErr)
		return
	}
	if createErr := h.service.Create(at); createErr != nil {
		ctx.JSON(createErr.Status, createErr)
		return
	}
	ctx.JSON(http.StatusCreated, at)
}

func (h *accessTokenHandler) UpdateExpirationTime(ctx *gin.Context) {
	at := access_token.AccessToken{}
	if err := h.service.Create(at); err != nil {
		ctx.JSON(err.Status, err)
		return
	}
	ctx.JSON(http.StatusOK, at)
}