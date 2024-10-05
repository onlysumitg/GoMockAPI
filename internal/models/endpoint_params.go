package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/onlysumitg/GoMockAPI/utils/jsonutils"
	"github.com/onlysumitg/GoMockAPI/utils/xmlutils"
)

// ------------------------------------------------------------
//
// ------------------------------------------------------------
func (s EndPointModel) UpdateRequestParamFromApiCall(wg *sync.WaitGroup, endPoint *EndPoint, requestJson map[string]any) error {
	jsonString, err := json.Marshal(requestJson)
	if err != nil {
		return err
	}

	endPoint.SampleRequest = string(jsonString)
	endPoint.AppendReqParam = true
	return s.RebuildRequestParams(wg, endPoint)
}

// ------------------------------------------------------------
//
// ------------------------------------------------------------
func (s EndPointModel) RebuildRequestParams(wg *sync.WaitGroup, endPoint *EndPoint) error {

	var err error
	var flatmap map[string]xmlutils.ValueDatatype
	//xmlPlaceholder := ""

	switch endPoint.SampleRequestType {
	case "JSON":
		flatmap, err = jsonutils.JsonToFlatMap(endPoint.SampleRequest)
	case "XML":
		flatmap, _, err = xmlutils.XmlToFlatMapAndPlaceholder(endPoint.SampleRequest)

	default:
		err = errors.New("Unknow Request Type")

	}

	if err != nil {

	} else {

		endPointRequestParamModel := &EndPointRequestParamModel{DB: s.DB}

		paramMap := make(map[string]*EndPointRequestParam)

		// add param for client ip
		endPointRequestParam := &EndPointRequestParam{
			EndpointID:      endPoint.ID,
			Key:             "*CLIENT_IP",
			DefaultValue:    "",
			DefaultDatatype: "string",
		}
		paramMap["*CLIENT_IP"] = endPointRequestParam

		// create param based on current json
		for key, jsonVal := range flatmap {
			endPointRequestParam := &EndPointRequestParam{
				EndpointID:      endPoint.ID,
				Key:             key,
				DefaultValue:    jsonVal.Value,
				DefaultDatatype: jsonVal.DataType,
			}
			paramMap[key] = endPointRequestParam

		}

		// process already saved params
		savedParameters := endPointRequestParamModel.ListById(endPoint.ID)

		for _, savedParam := range savedParameters {

			if strings.HasPrefix(savedParam.Key, "*HEADER_") || strings.HasPrefix(savedParam.Key, "*PATH_") {
				continue
			}
			orgParam, found := paramMap[savedParam.Key]
			if !found {
				if !endPoint.AppendReqParam {
					endPointRequestParamModel.Delete(savedParam.ID)
				}
				continue
			}

			orgParam.OverrideValue = savedParam.OverrideValue
			orgParam.OverrrideDatatype = savedParam.OverrrideDatatype
			orgParam.ID = savedParam.ID // KEEP the ID

		}

		for _, v := range paramMap {
			endPointRequestParamModel.Save(v)
		}
	}
	if wg != nil {
		wg.Done()
	}

	return err
}

// ------------------------------------------------------------
//
// ------------------------------------------------------------
func (s EndPointModel) RebuildRequestHeaderParams(wg *sync.WaitGroup, endPoint *EndPoint) error {

	var err error
	var flatmap map[string]xmlutils.ValueDatatype
	//xmlPlaceholder := ""

	switch endPoint.SampleRequestHeaderType {
	case "JSON":
		flatmap, err = jsonutils.JsonToFlatMap(endPoint.SampleRequestHeader)
	case "XML":
		flatmap, _, err = xmlutils.XmlToFlatMapAndPlaceholder(endPoint.SampleRequestHeader)

	default:
		err = errors.New("Unknow Request Type")

	}

	//flatmap, err := jsonutils.JsonToFlatMap(endPoint.SampleRequestHeader)
	if err != nil {

	} else {

		endPointRequestParamModel := &EndPointRequestParamModel{DB: s.DB}

		paramMap := make(map[string]*EndPointRequestParam)

		// create param based on current json
		for key, jsonVal := range flatmap {
			keyToUse := fmt.Sprintf("*HEADER_%s", key)
			endPointRequestParam := &EndPointRequestParam{
				EndpointID:      endPoint.ID,
				Key:             keyToUse,
				DefaultValue:    jsonVal.Value,
				DefaultDatatype: jsonVal.DataType,
			}
			paramMap[keyToUse] = endPointRequestParam

		}

		// process already saved params
		savedParameters := endPointRequestParamModel.ListById(endPoint.ID)

		for _, savedParam := range savedParameters {

			if strings.HasPrefix(savedParam.Key, "*HEADER_") {

				orgParam, found := paramMap[savedParam.Key]
				if !found {

					if !endPoint.AppendReqHeader {
						endPointRequestParamModel.Delete(savedParam.ID)
					}
					continue
				}

				orgParam.OverrideValue = savedParam.OverrideValue
				orgParam.OverrrideDatatype = savedParam.OverrrideDatatype
				orgParam.ID = savedParam.ID // KEEP the ID
			}
		}

		for _, v := range paramMap {
			endPointRequestParamModel.Save(v)
		}
	}
	if wg != nil {
		wg.Done()
	}
	return err
}

// ------------------------------------------------------------
//
// ------------------------------------------------------------
func (s EndPointModel) RebuildRequestPathParams(wg *sync.WaitGroup, endPoint *EndPoint) error {

	endPointRequestParamModel := &EndPointRequestParamModel{DB: s.DB}

	paramMap := make(map[string]*EndPointRequestParam)

	// create param based on current json
	for _, pathPram := range endPoint.PathParams {
		keyToUse := pathPram.Name
		endPointRequestParam := &EndPointRequestParam{
			EndpointID:      endPoint.ID,
			Key:             keyToUse,
			DefaultValue:    pathPram.Value,
			DefaultDatatype: pathPram.DataType,
		}
		paramMap[keyToUse] = endPointRequestParam

	}

	// process already saved params
	savedParameters := endPointRequestParamModel.ListById(endPoint.ID)

	for _, savedParam := range savedParameters {

		if strings.HasPrefix(savedParam.Key, "*PATH_") {

			orgParam, found := paramMap[savedParam.Key]
			if !found {

				if !endPoint.AppendReqParam {
					endPointRequestParamModel.Delete(savedParam.ID)
				}
				continue
			}

			orgParam.OverrideValue = savedParam.OverrideValue
			orgParam.OverrrideDatatype = savedParam.OverrrideDatatype
			orgParam.ID = savedParam.ID // KEEP the ID
		}
	}

	for _, v := range paramMap {
		endPointRequestParamModel.Save(v)
	}

	if wg != nil {
		wg.Done()
	}
	return nil
}

// ------------------------------------------------------------
//
// ------------------------------------------------------------
func (s EndPointModel) RebuildResponseParams(wg *sync.WaitGroup, endPoint *EndPoint) error {
	//flatmap, err := jsonutils.JsonToFlatMap(endPoint.SampleResponse)
	defer func() {
		if wg != nil {
			wg.Done()
		}
	}()

	var err error

	var wglocal sync.WaitGroup
	for _, r := range endPoint.ResponseMap {
		wglocal.Add(1)
		go r.RebuildResponseParams(&wglocal, s.DB)

	}
	wglocal.Wait()
	return err
}

// ------------------------------------------------------------
//
// ------------------------------------------------------------
func (s EndPointModel) RebuildResponseHeaderParams(wg *sync.WaitGroup, endPoint *EndPoint) error {
	//flatmap, err := jsonutils.JsonToFlatMap(endPoint.SampleResponseHeader)

	defer func() {
		if wg != nil {
			wg.Done()
		}
	}()

	var err error
	var wglocal sync.WaitGroup

	for _, r := range endPoint.ResponseMap {
		wglocal.Add(1)
		go r.RebuildResponseHeaderParams(&wglocal, s.DB)

	}
	wglocal.Wait()

	return err
}
