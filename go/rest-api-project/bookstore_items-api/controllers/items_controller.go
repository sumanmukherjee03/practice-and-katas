package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/sumanmukherjee03/practice-and-katas/go/rest-api-project/bookstore_items-api/domain/items"
	"github.com/sumanmukherjee03/practice-and-katas/go/rest-api-project/bookstore_items-api/services"
	"github.com/sumanmukherjee03/practice-and-katas/go/rest-api-project/bookstore_items-api/utils/http_utils"
	"github.com/sumanmukherjee03/practice-and-katas/go/rest-api-project/bookstore_oauth-go/oauth"
	"github.com/sumanmukherjee03/practice-and-katas/go/rest-api-project/bookstore_utils-go/rest_errors"
)

var (
	ItemsController itemsControllerInterface = &itemsController{}
)

// Instead of using a package to separate the controllers like we did in the users api
// we are using the same controller here but controlling the grouping of methods on a controller with an interface
// and a concrete struct implementing that interface. An object of this struct now represents a specific controller.
type itemsControllerInterface interface {
	Create(http.ResponseWriter, *http.Request)
	Get(http.ResponseWriter, *http.Request)
}

type itemsController struct {
}

func (c *itemsController) Create(w http.ResponseWriter, r *http.Request) {
	if err := oauth.Authenticate(r); err != nil {
		http_utils.RespondError(w, err)
		return
	}

	var itemRequest items.Item
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http_utils.RespondError(w, rest_errors.NewBadRequestError(fmt.Errorf("invalid request body : %v", err)))
		return
	}
	defer r.Body.Close()

	if err := json.Unmarshal(requestBody, &itemRequest); err != nil {
		http_utils.RespondError(w, rest_errors.NewBadRequestError(fmt.Errorf("invalid item json : %v", err)))
		return
	}

	itemRequest.Seller = oauth.GetCallerId(r)

	item, createErr := services.ItemsService.Create(itemRequest)
	if createErr != nil {
		http_utils.RespondError(w, createErr)
		return
	}

	http_utils.RespondJson(w, http.StatusCreated, item)
}

func (c *itemsController) Get(w http.ResponseWriter, r *http.Request) {
}
