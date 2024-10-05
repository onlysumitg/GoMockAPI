package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// ------------------------------------------------------
//
// ------------------------------------------------------
func (app *application) RequestParamHandlers(router *chi.Mux) {
	router.Route("/requestparam/{endpointid}", func(r chi.Router) {

		r.Use(app.RequireAuthentication)

		r.Use(app.EndPointOwnership)
		r.Use(noSurf)

		r.Get("/", app.RequestParamList)
		r.Get("/{paramid}", app.EndPointView)

		r.Get("/update/{paramid}", app.EndPointAdd)

	})

}

// ------------------------------------------------------
//
// ------------------------------------------------------
func (app *application) RequestParamList(w http.ResponseWriter, r *http.Request) {
	endpointID := chi.URLParam(r, "endpointid")
	endpoint, err := app.endpoints.Get(endpointID)
	if err != nil {
		app.Http404(w, r)
		return
	}

	data := app.newTemplateData(r)
	data.RequestParams = app.requestParams.ListById(endpoint.ID)
	data.EndPoint = endpoint
	app.render(w, r, http.StatusOK, "requestparam_list.tmpl", data)

}

// ------------------------------------------------------
// EndPoint details
// ------------------------------------------------------
func (app *application) RequestParamView(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	endpointID := chi.URLParam(r, "endpointid")
	endpoint, err := app.endpoints.Get(endpointID)
	if err != nil {
		app.Http404(w, r)
		return
	}
	data.EndPoint = endpoint
	app.render(w, r, http.StatusOK, "RequestParam_view.tmpl", data)

}
