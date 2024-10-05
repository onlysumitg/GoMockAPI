package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/onlysumitg/GoMockAPI/internal/models"
)

// ------------------------------------------------------
//
// ------------------------------------------------------
func (app *application) ResponseParamHandlers(router *chi.Mux) {
	router.Route("/responseparam/{endpointid}/{ownerid}", func(r chi.Router) {
		//r.With(paginate).Get("/", listArticles)
		r.Use(app.RequireAuthentication)

		r.Use(app.EndPointOwnership)
		r.Use(noSurf)

		r.Get("/", app.ResponseParamList)
		r.Get("/{paramid}", app.ResponseParamView)
		r.Get("/update/{paramid}", app.ResponseParamUpdate)
		r.Post("/update", app.ResponseParamUpdatePost)

	})

}

// ------------------------------------------------------
//
// ------------------------------------------------------
func (app *application) ResponseParamList(w http.ResponseWriter, r *http.Request) {
	endpointID := chi.URLParam(r, "endpointid")
	endpoint, err := app.endpoints.Get(endpointID)
	if err != nil {
		app.Http404(w, r)
		return
	}

	ownerid := chi.URLParam(r, "ownerid")

	data := app.newTemplateData(r)
	data.ResponseParams = app.responseParams.ListByOwnerId(ownerid)
	data.EndPoint = endpoint
	data.ResponseParamsOwnerId = ownerid
	app.render(w, r, http.StatusOK, "responseparam_list.tmpl", data)

}

// ------------------------------------------------------
// EndPoint details
// ------------------------------------------------------
func (app *application) ResponseParamView(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	endpointID := chi.URLParam(r, "endpointid")
	endpoint, err := app.endpoints.Get(endpointID)
	if err != nil {
		app.Http404(w, r)
		return
	}
	data.EndPoint = endpoint
	app.render(w, r, http.StatusOK, "ResponseParam_view.tmpl", data)

}

// ------------------------------------------------------
// add new endpoint
// ------------------------------------------------------
func (app *application) ResponseParamUpdate(w http.ResponseWriter, r *http.Request) {
	endpointID := chi.URLParam(r, "endpointid")
	endpoint, err := app.endpoints.Get(endpointID)
	if err != nil {
		app.Http404(w, r)
		return
	}

	requestParams := app.requestParams.ListById(endpoint.ID)
	paramKeys := make([]string, 0)

	for _, requestParam := range requestParams {
		paramKeys = append(paramKeys, fmt.Sprintf("REQUEST[%s]: %s", requestParam.DefaultDatatype, requestParam.Key))
	}

	randomFuncKeys := models.GetRandonFunctiolist()
	paramKeys = append(paramKeys, randomFuncKeys...)

	paramid := chi.URLParam(r, "paramid")
	data := app.newTemplateData(r)

	data.RequestParamAutoComplate = strings.Join(paramKeys, ",")
	data.RequestParamAutoComplateList = paramKeys

	responseParam, err := app.responseParams.Get(paramid)
	if err != nil {
		app.sessionManager.Put(r.Context(), "error", fmt.Sprintf("Error updating endpoint: %s", err.Error()))
		app.goBack(w, r, http.StatusBadRequest)
		return
	}

	data.ConditionGroups = app.conditionGroup.ListById(endpoint.ID)

	data.Form = responseParam
	data.EndPoint = endpoint

	data.ResponseParamsOwnerId = chi.URLParam(r, "ownerid")

	app.render(w, r, http.StatusOK, "responseparam_update.tmpl", data)

}

// ----------------------------------------------
func (app *application) ResponseParamUpdatePost(w http.ResponseWriter, r *http.Request) {
	app.invalidateEndPointCache()

	endpointID := chi.URLParam(r, "endpointid")

	err := r.ParseForm()
	if err != nil {
		app.sessionManager.Put(r.Context(), "error", fmt.Sprintf("001 Error processing form %s", err.Error()))
		app.goBack(w, r, http.StatusBadRequest)
		return
	}

	var submitedParam models.EndPointResponseParam

	err = app.formDecoder.Decode(&submitedParam, r.PostForm)
	if err != nil {
		app.sessionManager.Put(r.Context(), "error", fmt.Sprintf("002 Error processing form %s", err.Error()))
		app.goBack(w, r, http.StatusBadRequest)
		return
	}

	responseParam, err := app.responseParams.Get(submitedParam.ID)
	if err != nil {
		app.sessionManager.Put(r.Context(), "error", fmt.Sprintf("Error updating endpoint: %s", err.Error()))
		app.goBack(w, r, http.StatusBadRequest)
		return
	}

	responseParam.OverrideValue = submitedParam.OverrideValue

	_, err = app.responseParams.Save(responseParam)
	if err != nil {
		app.serverError500(w, r, err)
		return
	}

	app.invalidateEndPointCache()

	app.sessionManager.Put(r.Context(), "flash", fmt.Sprintf("EndPoint %s added sucessfully", responseParam.Key))

	ownerid := chi.URLParam(r, "ownerid")
	http.Redirect(w, r, fmt.Sprintf("/responseparam/%s/%s", endpointID, ownerid), http.StatusSeeOther)
}
