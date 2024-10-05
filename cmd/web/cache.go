package main

import (
	"fmt"
	"strings"

	"github.com/onlysumitg/GoMockAPI/internal/models"
)

func (app *application) invalidateEndPointCache() {
	app.invalidEndPointCache = true
}

// // ------------------------------------------------------
// //
// // ------------------------------------------------------
// func (app *application) GetEndPoint(namespace, endpointName, httpmethod string) (*storedProc.StoredProc, error) {
// 	endPointKey := fmt.Sprintf("%s_%s_%s", strings.ToUpper(namespace), strings.ToUpper(endpointName), strings.ToUpper(httpmethod))

// 	endPoint, found := app.endPointCache[endPointKey]
// 	if !found || app.invalidEndPointCache {
// 		app.endPointCache = make(map[string]*storedProc.StoredProc)
// 		app.endPointMutex.Lock()
// 		for _, sp := range app.storedProcs.List() {
// 			app.endPointCache[fmt.Sprintf("%s_%s_%s", strings.ToUpper(sp.Namespace), strings.ToUpper(sp.EndPointName), strings.ToUpper(sp.HttpMethod))] = sp
// 		}
// 		endPoint, found = app.endPointCache[endPointKey]
// 		app.invalidEndPointCache = false
// 		app.endPointMutex.Unlock()

// 		if !found {

// 			return nil, fmt.Errorf("Not Found: %s", strings.ReplaceAll(endPointKey, "_", " "))
// 		}

// 	}

// 	return endPoint, nil

// }

func (app *application) GetEndPoint(collection, endpointname, httpmethod string) (*models.EndPoint, error) {

	endPointKey := fmt.Sprintf("%s_%s_%s", strings.ToLower(collection), strings.ToLower(endpointname), strings.ToLower(httpmethod))

	endPoint, found := app.endPointCache[endPointKey]
	if !found || app.invalidEndPointCache {

		app.endPointMutex.Lock()
		app.endPointCache = app.endpoints.BuildEndPointCache(app.maxAllowedEndPoints)
		endPoint, found = app.endPointCache[endPointKey]
		app.invalidEndPointCache = false
		app.endPointMutex.Unlock()

		if !found {

			return nil, fmt.Errorf("not found: %s %s %s", collection, endpointname, httpmethod)
		}

	}

	return endPoint, nil

}
