package todo

import (
	"fmt"
	"log"
	"net/http"

	"github.com/benschw/opin-go/rest"
	"github.com/benschw/opin-go/vault"
)

type TodoResource struct {
	Db vault.DbProvider
}

func (r *TodoResource) Health(res http.ResponseWriter, req *http.Request) {
	//2xx => pass, 429 => warn, anything else => critical

	// if we can't get a db connection, set health to Critical
	db, err := r.Db.Get()
	if err != nil {
		rest.SetInternalServerErrorResponse(res, err)
		return
	}

	var todos []Todo
	assoc := db.Find(&todos)
	if assoc.Error != nil {
		rest.SetInternalServerErrorResponse(res, err)
		return
	}
	// set health to OK
	if err := rest.SetOKResponse(res, nil); err != nil {
		rest.SetInternalServerErrorResponse(res, err)
		return
	}
}

func (r *TodoResource) Add(res http.ResponseWriter, req *http.Request) {
	var todo Todo

	db, err := r.Db.Get()
	if err != nil {
		rest.SetInternalServerErrorResponse(res, err)
		return
	}

	if err := rest.Bind(req, &todo); err != nil {
		log.Print(err)
		rest.SetBadRequestResponse(res)
		return
	}

	db.Create(&todo)

	if err := rest.SetCreatedResponse(res, todo, fmt.Sprintf("todo/%d", todo.Id)); err != nil {
		rest.SetInternalServerErrorResponse(res, err)
		return
	}
}

func (r *TodoResource) Get(res http.ResponseWriter, req *http.Request) {
	id, err := rest.PathInt(req, "id")
	if err != nil {
		rest.SetBadRequestResponse(res)
		return
	}
	var todo Todo

	db, err := r.Db.Get()
	if err != nil {
		rest.SetInternalServerErrorResponse(res, err)
		return
	}
	if db.First(&todo, id).RecordNotFound() {
		rest.SetNotFoundResponse(res)
		return
	}

	if err := rest.SetOKResponse(res, todo); err != nil {
		rest.SetInternalServerErrorResponse(res, err)
		return
	}
}

func (r *TodoResource) GetAll(res http.ResponseWriter, req *http.Request) {
	var todos []Todo

	db, err := r.Db.Get()
	if err != nil {
		rest.SetInternalServerErrorResponse(res, err)
		return
	}
	db.Find(&todos)

	if err := rest.SetOKResponse(res, todos); err != nil {
		rest.SetInternalServerErrorResponse(res, err)
		return
	}
}
func (r *TodoResource) Update(res http.ResponseWriter, req *http.Request) {
	var todo Todo

	id, err := rest.PathInt(req, "id")
	if err != nil {
		rest.SetBadRequestResponse(res)
		return
	}
	if err := rest.Bind(req, &todo); err != nil {
		rest.SetBadRequestResponse(res)
		return
	}
	todo.Id = id

	var found Todo
	db, err := r.Db.Get()
	if err != nil {
		rest.SetInternalServerErrorResponse(res, err)
		return
	}
	if db.First(&found, id).RecordNotFound() {
		rest.SetNotFoundResponse(res)
		return
	}

	db.Save(&todo)

	if err := rest.SetOKResponse(res, todo); err != nil {
		rest.SetInternalServerErrorResponse(res, err)
		return
	}

}
func (r *TodoResource) Delete(res http.ResponseWriter, req *http.Request) {
	id, err := rest.PathInt(req, "id")
	if err != nil {
		rest.SetBadRequestResponse(res)
		return
	}
	var todo Todo

	db, err := r.Db.Get()
	if err != nil {
		rest.SetInternalServerErrorResponse(res, err)
		return
	}
	if db.First(&todo, id).RecordNotFound() {
		rest.SetNotFoundResponse(res)
		return
	}

	db.Delete(&todo)

	if err := rest.SetNoContentResponse(res); err != nil {
		rest.SetInternalServerErrorResponse(res, err)
		return
	}
}
