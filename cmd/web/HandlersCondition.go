package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/onlysumitg/GoMockAPI/internal/models"
	"github.com/onlysumitg/GoMockAPI/internal/validator"
)

// ------------------------------------------------------
//
// ------------------------------------------------------
func (app *application) ConditionHandlers(router *chi.Mux) {
	router.Route("/conditions/{endpointid}", func(r chi.Router) {
		//r.With(paginate).Get("/", listArticles)
		// CSRF
		r.Use(app.RequireAuthentication)

		r.Use(app.EndPointOwnership)
		r.Use(noSurf)
		r.Get("/", app.ConditionList)
		r.Get("/{paramid}", app.ConditionView)

		r.Get("/add", app.ConditionUpdate)
		r.Post("/add", app.ConditionUpdatePost)

		r.Get("/update/{paramid}", app.ConditionUpdate)
		r.Post("/update/{paramid}", app.ConditionUpdatePost)

		r.Get("/delete/{objectid}", app.ConditionDelete)
		r.Post("/delete", app.ConditionDeleteConfirm)

	})

}

// ------------------------------------------------------
//
// ------------------------------------------------------
func (app *application) ConditionList(w http.ResponseWriter, r *http.Request) {
	endpointID := chi.URLParam(r, "endpointid")

	endpoint, err := app.endpoints.Get(endpointID)
	if err != nil {
		app.Http404(w, r)
		return
	}

	data := app.newTemplateData(r)
	data.Conditions = app.condition.ListById(endpoint.ID)
	data.EndPoint = endpoint
	app.render(w, r, http.StatusOK, "condition_list.tmpl", data)

}

// ------------------------------------------------------
// Delete servet
// ------------------------------------------------------
func (app *application) ConditionDelete(w http.ResponseWriter, r *http.Request) {
	endpointID := chi.URLParam(r, "endpointid")

	endpoint, err := app.endpoints.Get(endpointID)
	if err != nil {
		app.Http404(w, r)
		return
	}

	cgID := chi.URLParam(r, "objectid")

	cg, err := app.condition.Get(cgID)
	if err != nil {
		app.Http404(w, r)
		return
	}

	data := app.newTemplateData(r)
	data.EndPoint = endpoint
	data.Condition = cg

	app.render(w, r, http.StatusOK, "condition_delete.tmpl", data)

}

// ------------------------------------------------------
// Delete servet
// ------------------------------------------------------
func (app *application) ConditionDeleteConfirm(w http.ResponseWriter, r *http.Request) {
	app.invalidateEndPointCache()

	endpointID := chi.URLParam(r, "endpointid")

	err := r.ParseForm()
	if err != nil {
		app.sessionManager.Put(r.Context(), "error", fmt.Sprintf("001 Error processing form %s", err.Error()))
		app.goBack(w, r, http.StatusBadRequest)
		return
	}

	objectid := r.PostForm.Get("objectid")

	err = app.condition.Delete(objectid)
	if err != nil {

		app.sessionManager.Put(r.Context(), "error", fmt.Sprintf("Error deleting Condition: %s", err.Error()))
		app.goBack(w, r, http.StatusBadRequest)
		return
	}
	app.sessionManager.Put(r.Context(), "flash", "Deleted sucessfully")

	http.Redirect(w, r, fmt.Sprintf("/conditions/%s", endpointID), http.StatusSeeOther)

}

// ------------------------------------------------------
// EndPoint details
// ------------------------------------------------------
func (app *application) ConditionView(w http.ResponseWriter, r *http.Request) {

	endpointID := chi.URLParam(r, "endpointid")

	endpoint, err := app.endpoints.Get(endpointID)
	if err != nil {
		app.Http404(w, r)
		return
	}
	data := app.newTemplateData(r)
	data.EndPoint = endpoint
	app.render(w, r, http.StatusOK, "condition_view.tmpl", data)

}

// ------------------------------------------------------
// add new endpoint
// ------------------------------------------------------
func (app *application) ConditionUpdate(w http.ResponseWriter, r *http.Request) {
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

	paramKeys = append(paramKeys, "Client: IP")

	data := app.newTemplateData(r)

	data.RequestParamAutoComplate = strings.Join(paramKeys, ",")

	paramid := chi.URLParam(r, "paramid")

	if paramid == "" {
		data.Form = models.Condition{}
	} else {
		condition, err := app.condition.Get(paramid)
		if err != nil {
			app.sessionManager.Put(r.Context(), "error", fmt.Sprintf("Error updating endpoint: %s", err.Error()))
			app.goBack(w, r, http.StatusBadRequest)
			return
		}
		data.Form = condition
	}
	data.EndPoint = endpoint

	data.RequestParams = app.requestParams.ListById(endpoint.ID)

	app.render(w, r, http.StatusOK, "condition_add.tmpl", data)

}

// ----------------------------------------------
func (app *application) ConditionUpdatePost(w http.ResponseWriter, r *http.Request) {
	app.invalidateEndPointCache()
	endpointID := chi.URLParam(r, "endpointid")

	endpoint, err := app.endpoints.Get(endpointID)
	if err != nil {
		app.Http404(w, r)
		return
	}

	err = r.ParseForm()
	if err != nil {
		app.sessionManager.Put(r.Context(), "error", fmt.Sprintf("001 Error processing form %s", err.Error()))
		app.goBack(w, r, http.StatusBadRequest)
		return
	}

	var condition models.Condition

	err = app.formDecoder.Decode(&condition, r.PostForm)
	if err != nil {
		app.sessionManager.Put(r.Context(), "error", fmt.Sprintf("002 Error processing form %s", err.Error()))
		app.goBack(w, r, http.StatusBadRequest)
		return
	}

	requestParam, err := app.requestParams.Get(condition.Variable)
	if err != nil {
		condition.CheckField(false, "variable", "Please select a value")

	}

	condition.CheckField(validator.NotBlank(condition.Compareto), "compareto", "This field cannot be blank")

	condition.CheckField(validator.MustBeOfType(condition.Compareto, requestParam.DefaultDatatype), "compareto", fmt.Sprintf("Make sure value is compatible with %s", requestParam.DefaultDatatype))

	condition.Name = fmt.Sprintf("%s %s %s", condition.VariableName, condition.Operator, condition.Compareto)
	//condition.CheckField(validator.NotBlank(condition.Name), "name", "This field cannot be blank")

	//TODO duplicate name check not working ===>
	condition.CheckField(!app.condition.DuplicateName(&condition, *endpoint), "compareto", "Duplicate condition")

	if !condition.Valid() {
		data := app.newTemplateData(r)
		data.Form = condition
		data.EndPoint = endpoint
		data.RequestParams = app.requestParams.ListById(endpoint.ID)

		app.sessionManager.Put(r.Context(), "error", "Please fix error(s) and resubmit")

		app.render(w, r, http.StatusUnprocessableEntity, "condition_add.tmpl", data)

		return
	}

	condition.EndpointID = endpointID
	condition.VariableName = requestParam.Key

	condition.Name = fmt.Sprintf("%s %s %s", condition.VariableName, condition.Operator, condition.Compareto)

	_, err = app.condition.Save(&condition)
	if err != nil {
		app.serverError500(w, r, err)
		return
	}
	app.invalidateEndPointCache()

	app.sessionManager.Put(r.Context(), "flash", fmt.Sprintf("Condition %s added sucessfully", condition.Name))

	http.Redirect(w, r, fmt.Sprintf("/conditions/%s", endpointID), http.StatusSeeOther)
}
