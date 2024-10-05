package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/justinas/nosurf"
	"github.com/onlysumitg/GoMockAPI/internal/models"
	"github.com/onlysumitg/GoMockAPI/ui"
	"github.com/onlysumitg/GoMockAPI/utils/httputils"
)

type templateData struct {
	CurrentYear int

	HostUrl string

	Form any //use this Form field to pass the validation errors and previously submitted data back to the template when we re-display the form.

	// differnt notifications
	Flash   string
	Warning string
	Error   string

	IsAuthenticated bool
	IsSuperuser     bool

	CSRFToken string // Add a CSRFToken field.   <input type='hidden' name='csrf_token' value='{{.CSRFToken}}'>

	EndPoint  *models.EndPoint
	EndPoints []*models.EndPoint

	RequestParam  *models.EndPointRequestParam
	RequestParams []*models.EndPointRequestParam

	ResponseParam         *models.EndPointResponseParam
	ResponseParams        []*models.EndPointResponseParam
	ResponseParamsOwnerId string

	RequestParamAutoComplate     string
	RequestParamAutoComplateList []string

	Condition  *models.Condition
	Conditions []*models.Condition

	EndPointResponse  *models.EndPointResponse
	EndPointResponses []*models.EndPointResponse

	ConditionGroup  *models.ConditionGroup
	ConditionGroups []*models.ConditionGroup

	Collection  *models.Collection
	Collections []*models.Collection

	ComparisonOperators []string

	LogEntries []string
	Next       string

	Messages []string

	RbacRoles                   []string
	RbacRole                    string
	RbacPermissions             []string
	RbacRolePermissionsIncluded []string
	RbacRolePermissionsExcluded []string

	Users       []*models.User
	User        *models.User
	CurrentUser *models.User

	TestMode bool

	HttpCodes map[int]string
}

func ListComparisonOperators() []string {
	return_List := []string{
		"EQUALS_TO",
		"NOT_EQUALS_TO",

		"LESS_THAN",
		"LESS_THAN_OR_EQUALS_TO",
		"GREATER_THAN",
		"GREATER_THAN_OR_EQUALS_TO",
		"CONTAINS",
		"STARTS_WITH",
		"ENDS_WITH",
	}
	sort.Strings(return_List)
	return return_List
}

// -----------------------------------------------------------------
//
// -----------------------------------------------------------------
func (app *application) newTemplateData(r *http.Request) *templateData {

	td := &templateData{
		CurrentYear:         time.Now().Year(),
		CSRFToken:           nosurf.Token(r), // Add the CSRF token.
		HostUrl:             app.hostURL,
		ComparisonOperators: ListComparisonOperators(),
		IsAuthenticated:     app.isAuthenticated(r), // use {{if .IsAuthenticated}} in template
		TestMode:            app.testMode,
		HttpCodes:           httputils.Codes,
	}
	user, err := app.GetUser(r)
	if err == nil {
		td.CurrentUser = user

	}

	return td
}

// -----------------------------------------------------------------
//
// -----------------------------------------------------------------
func (app *application) setTemplateDefaults(r *http.Request, templateData *templateData) {

	// Add the flash message to the template data, if one exists.
	templateData.Flash = app.sessionManager.PopString(r.Context(), "flash")
	templateData.Warning = app.sessionManager.PopString(r.Context(), "warning")
	templateData.Error = app.sessionManager.PopString(r.Context(), "error")

}

// -----------------------------------------------------------------
//
// -----------------------------------------------------------------
func (app *application) render(w http.ResponseWriter, r *http.Request, status int, page string, data *templateData) {
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError500(w, r, err)
		return
	}
	// Initialize a new buffer.
	buf := new(bytes.Buffer)

	app.setTemplateDefaults(r, data)

	baseTemplateName := "base"

	if strings.HasPrefix(page, "account_") {
		baseTemplateName = "account_base"
	}

	if strings.HasPrefix(page, "public_") {
		baseTemplateName = "public_base"
	}
	// Write the template to the buffer, instead of straight to the
	// http.ResponseWriter. If there's an error, call our serverError() helper
	// and then return.
	err := ts.ExecuteTemplate(buf, baseTemplateName, data)
	if err != nil {
		app.serverError500(w, r, err)
		return
	}
	// If the template is written to the buffer without any errors, we are safe
	// to go ahead and write the HTTP status code to http.ResponseWriter.
	w.WriteHeader(status)
	// Write the contents of the buffer to the http.ResponseWriter. Note: this
	// is another time where we pass our http.ResponseWriter to a function that
	// takes an io.Writer.
	buf.WriteTo(w)
}

// -----------------------------------------------------------------
//
// -----------------------------------------------------------------
func (app *application) templateToString(page string, data any) (string, error) {
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)

		return "", err
	}
	// Initialize a new buffer.
	buf := new(bytes.Buffer)

	baseTemplateName := "base"

	if strings.HasPrefix(page, "account_") {
		baseTemplateName = "account_base"
	}

	if strings.HasPrefix(page, "email_") {
		baseTemplateName = "email_base"
	}

	// Write the template to the buffer, instead of straight to the
	// http.ResponseWriter. If there's an error, call our serverError() helper
	// and then return.
	err := ts.ExecuteTemplate(buf, baseTemplateName, data)
	if err != nil {

		return "", err
	}

	return buf.String(), nil
}

// -----------------------------------------------------------------
//
// -----------------------------------------------------------------
func (app *application) newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}
	pages, err := fs.Glob(ui.Files, "html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}
	app.loadPages(cache, pages)

	pages, err = fs.Glob(ui.Files, "html/accounts/*.tmpl")
	if err != nil {
		return nil, err
	}
	app.loadPages(cache, pages)

	pages, err = fs.Glob(ui.Files, "html/emails/*.tmpl")
	if err != nil {
		return nil, err
	}
	app.loadPages(cache, pages)

	pages, err = fs.Glob(ui.Files, "html/public/*.tmpl")
	if err != nil {
		return nil, err
	}
	app.loadPages(cache, pages)

	return cache, nil
}

// -----------------------------------------------------------------
//
// -----------------------------------------------------------------
func (app *application) loadPages(cache map[string]*template.Template, pages []string) {
	// Use fs.Glob() to get a slice of all filepaths in the ui.Files embedded
	// filesystem which match the pattern 'html/pages/*.tmpl'. This essentially
	// gives us a slice of all the 'page' templates for the application, just
	// like before.

	for _, page := range pages {
		name := filepath.Base(page)
		// Create a slice containing the filepath patterns for the templates we
		// want to parse.
		patterns := []string{
			"html/base.tmpl",
			"html/account_base.tmpl",
			"html/email_base.tmpl",
			"html/public_base.tmpl",

			"html/partials/*.tmpl",

			page,
		}
		// Use ParseFS() instead of ParseFiles() to parse the template files
		// from the ui.Files embedded filesystem.
		ts, err := template.New(name).Funcs(app.getFunctionMap()).ParseFS(ui.Files, patterns...)
		if err != nil {
			log.Fatalln("Error loading template", err.Error())
		}
		cache[name] = ts

	}
}
