package http

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sumanmukherjee03/practice-and-katas/go/rest-api-project/bookstore_oauth-api/src/domain/access_token"
	access_token_service "github.com/sumanmukherjee03/practice-and-katas/go/rest-api-project/bookstore_oauth-api/src/services/access_token"
	"github.com/sumanmukherjee03/practice-and-katas/go/rest-api-project/bookstore_utils-go/rest_errors"
)

type AccessTokenHandler interface {
	GetById(*gin.Context)
	Create(*gin.Context)
	UpdateExpirationTime(*gin.Context)
}

type accessTokenHandler struct {
	service access_token_service.Service
}

func NewHandler(service access_token_service.Service) AccessTokenHandler {
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
	var atr access_token.AccessTokenRequest
	if err := ctx.ShouldBindJSON(&atr); err != nil {
		restErr := rest_errors.NewBadRequestError(fmt.Errorf("Bad input passed to create access token : %v", err))
		ctx.JSON(restErr.Status, restErr)
		return
	}
	at, createErr := h.service.Create(atr)
	if createErr != nil {
		ctx.JSON(createErr.Status, createErr)
		return
	}
	ctx.JSON(http.StatusCreated, at)
}

func (h *accessTokenHandler) UpdateExpirationTime(ctx *gin.Context) {
	atr := access_token.AccessTokenRequest{}
	at, err := h.service.Create(atr)
	if err != nil {
		ctx.JSON(err.Status, err)
		return
	}
	ctx.JSON(http.StatusOK, at)
}
