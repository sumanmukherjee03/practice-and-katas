package controllers

import (
	"fmt"
	"net/http"

	"github.com/sumanmukherjee03/practice-and-katas/go/rest-api-project/bookstore_items-api/domain/items"
	"github.com/sumanmukherjee03/practice-and-katas/go/rest-api-project/bookstore_items-api/services"
	"github.com/sumanmukherjee03/practice-and-katas/go/rest-api-project/bookstore_oauth-go/oauth"
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
		return
	}
	item := items.Item{
		Seller: oauth.GetCallerId(r),
	}

	result, err := services.ItemsService.Create(item)
	if err != nil {
		return
	}

	fmt.Println(result)
}

func (c *itemsController) Get(w http.ResponseWriter, r *http.Request) {
}
