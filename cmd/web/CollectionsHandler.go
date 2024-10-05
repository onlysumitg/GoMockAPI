package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/onlysumitg/GoMockAPI/internal/models"
	"github.com/onlysumitg/GoMockAPI/internal/validator"
	"github.com/onlysumitg/GoMockAPI/utils/stringutils"
)

// ------------------------------------------------------
//
// ------------------------------------------------------
func (app *application) CollectionsHandlers(router *chi.Mux) {
	router.Route("/collections", func(r chi.Router) {
		r.Use(app.sessionManager.LoadAndSave)

		r.Use(app.RequireAuthentication)

		// CSRF
		r.Use(noSurf)
		r.Get("/", app.collectionsList)
		r.Get("/add", app.collectionsAdd)
		r.Post("/add", app.collectionsAdd)

		g1 := r.Group(nil)
		g1.Use(app.RequireSuperAdmin)
		g1.Get("/edit/{id}", app.collectionsAdd)
		g1.Post("/edit/{id}", app.collectionsAdd)

		g1.Get("/delete/{id}", app.collectionsDelete)
		g1.Post("/delete", app.collectionsDeleteConfirm)

	})

}

// ------------------------------------------------------
//
// ------------------------------------------------------
func (app *application) collectionsList(w http.ResponseWriter, r *http.Request) {

	data := app.newTemplateData(r)
	data.Collections = app.collectionsModel.List()
	app.render(w, r, http.StatusOK, "collection_list.tmpl", data)
}

// ------------------------------------------------------
//
// ------------------------------------------------------
func (app *application) collectionsAdd(w http.ResponseWriter, r *http.Request) {

	collection := &models.Collection{}

	id := chi.URLParam(r, "id")
	if id != "" {
		u, err := app.collectionsModel.Get(id)
		if err == nil {
			collection = u

		}
	}

	if r.Method == http.MethodPost {

		err := app.decodePostForm(r, &collection)
		if err != nil {
			app.clientError(w, http.StatusBadRequest, err)
			return
		}

		collection.Name = stringutils.RemoveSpecialChars(stringutils.RemoveMultipleSpaces(collection.Name))

		collection.CheckField(validator.NotBlank(collection.Name), "name", "This field cannot be blank")
		collection.CheckField(validator.CanNotBe(collection.Name, "V1"), "name", "Can not use reserved name V1.")

		collection.CheckField(validator.NotBlank(collection.Desc), "desc", "This field cannot be blank")
		collection.CheckField(!app.collectionsModel.DuplicateName(collection), "name", "Duplicate Name")

		if collection.Valid() {
			app.collectionsModel.Save(collection)

			endpoints := app.endpoints.ListByCollectionID(collection.ID)
			for _, ep := range endpoints {
				ep.CollectionName = collection.Name
				app.endpoints.ReBuildURL(ep)
			}

			app.invalidateEndPointCache()
			app.sessionManager.Put(r.Context(), "flash", "Saved sucessfully")

			http.Redirect(w, r, "/collections", http.StatusSeeOther)
			return
		}

	}

	data := app.newTemplateData(r)
	data.Form = collection

	app.render(w, r, http.StatusOK, "collection_add.tmpl", data)
}

// ------------------------------------------------------
//
// ------------------------------------------------------
func (app *application) collectionsDelete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	pr, err := app.collectionsModel.Get(id)
	if err != nil {
		app.clientError(w, http.StatusNotFound, err)
		return
	}

	data := app.newTemplateData(r)
	data.Collection = pr

	app.render(w, r, http.StatusOK, "collection_delete.tmpl", data)

}

// ------------------------------------------------------
//
// ------------------------------------------------------
func (app *application) collectionsDeleteConfirm(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		app.sessionManager.Put(r.Context(), "error", fmt.Sprintf("001 Error processing form %s", err.Error()))
		app.goBack(w, r, http.StatusSeeOther)
		return
	}

	id := r.PostForm.Get("id")

	err = app.collectionsModel.Delete(id)
	if err != nil {

		app.sessionManager.Put(r.Context(), "error", fmt.Sprintf("delete failed:: %s", err.Error()))
		app.goBack(w, r, http.StatusSeeOther)
		return
	} else {
		endpoints := app.endpoints.ListByCollectionID(id)
		for _, ep := range endpoints {
			app.endpoints.Delete(ep.ID)
		}
	}
	app.invalidateEndPointCache()
	app.sessionManager.Put(r.Context(), "flash", "Deleted sucessfully")

	http.Redirect(w, r, "/collections", http.StatusSeeOther)

}
