package models

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/google/uuid"
	"github.com/onlysumitg/GoMockAPI/internal/validator"
	bolt "go.etcd.io/bbolt"
)

// -----------------------------------------------------------------
//
// -----------------------------------------------------------------

type ConditionGroupParameter struct {
	ResponseVariable     string `json:"responsevariable" db:"responsevariable" form:"responsevariable"`
	ResponseVariableName string `json:"responsevariablename" db:"responsevariablename" form:"responsevariablename"`

	ResponseVariableDatatype string `json:"responsevariabledatatype" db:"responsevariabledatatype" form:"responsevariabledatatype"`

	AssgineValue string `json:"assignvalue" db:"assignvalue" form:"assignvalue"`

	ResponseParam *EndPointResponseParam `json:"-" db:"-" form:"-"`

	validator.Validator `json:"-" db:"-" form:"-"`
}

// -----------------------------------------------------------------
//
// -----------------------------------------------------------------

type ConditionGroup struct {
	ID         string `json:"id" db:"id" form:"id"`
	EndpointID string `json:"endpointid" db:"endpointid" form:"endpointid"`
	Name       string `json:"name" db:"name" form:"name"`

	ConditionIDs []string `json:"conditionids" db:"conditionids" form:"conditionids"`

	Conditions []*Condition `json:"-" db:"-" form:"-"`

	ConditionGroupParameters []*ConditionGroupParameter `json:"responseandconditiongroupmap" db:"responseandconditiongroupmap" form:"responseandconditiongroupmap"`

	validator.Validator
	CallActualEndPoint bool `json:"callactualendpoint" db:"callactualendpoint" form:"callactualendpoint"`

	//HttpStatusCode int `json:"httpstatuscode" db:"httpstatuscode" form:"httpstatuscode"`

	ResponseID string `json:"httpstatuscode" db:"httpstatuscode" form:"httpstatuscode"`
}

// -----------------------------------------------------------------
//
// -----------------------------------------------------------------

func BuildConditionGroupParameter(e *EndPoint, id string) []*ConditionGroupParameter {

	responseParams := make([]*EndPointResponseParam, 0)

	cgParams := make([]*ConditionGroupParameter, 0)

	for _, r := range e.ResponseMap {
		if strings.EqualFold(r.ID, id) {
			responseParams = r.ResponseParams
		}
	}

	for _, responseParam := range responseParams {

		cgParams = append(cgParams, &ConditionGroupParameter{
			ResponseVariable:         responseParam.ID,
			ResponseVariableName:     responseParam.Key,
			ResponseParam:            responseParam,
			ResponseVariableDatatype: responseParam.DefaultDatatype,
		})

	}

	return cgParams
}

// -----------------------------------------------------------------
//
// -----------------------------------------------------------------

func (c *ConditionGroup) Initialize(e *EndPoint, responseid string) *ConditionGroup {

	// httpCodeToUse := c.HttpStatusCode

	// if httpcode > 0 {
	// 	httpCodeToUse = httpcode
	// }

	responseParams := make([]*EndPointResponseParam, 0)

	for _, r := range e.ResponseMap {
		if strings.EqualFold(r.ID, responseid) {
			responseParams = r.ResponseParams
		}
	}

	for _, responseParam := range responseParams {
		alreadyAssigned := false
		for _, mapeed := range c.ConditionGroupParameters {
			if responseParam.ID == mapeed.ResponseVariable {
				alreadyAssigned = true
				mapeed.ResponseParam = responseParam
			}
		}

		if !alreadyAssigned {
			c.ConditionGroupParameters = append(c.ConditionGroupParameters, &ConditionGroupParameter{
				ResponseVariable:         responseParam.ID,
				ResponseVariableName:     responseParam.Key,
				ResponseParam:            responseParam,
				ResponseVariableDatatype: responseParam.DefaultDatatype,
			})

		}

	}

	c.RemoveInvalidResponseParam(responseParams)
	return c
}

// -----------------------------------------------------------------
//
// -----------------------------------------------------------------

func (c *ConditionGroup) RemoveInvalidResponseParam(params []*EndPointResponseParam) {
	tempMap := make([]*ConditionGroupParameter, 0)

	for _, xx := range c.ConditionGroupParameters {

		if xx == nil || xx.ResponseParam == nil {
			continue
		}

		// check if this parameter in p
		found := false
		for _, p := range params {
			if strings.EqualFold(p.ID, xx.ResponseParam.ID) {
				found = true
			}
		}

		if xx.ResponseParam != nil && found {
			tempMap = append(tempMap, xx)
		}
	}

	c.ConditionGroupParameters = tempMap
	sort.Slice(c.ConditionGroupParameters, func(i, j int) bool {
		return c.ConditionGroupParameters[i].ResponseVariableName < c.ConditionGroupParameters[j].ResponseVariableName
	})
}

// -----------------------------------------------------------------
//
// -----------------------------------------------------------------
func (cg *ConditionGroup) Execute(apiCall *ApiCall) bool {

	if len(cg.Conditions) <= 0 {
		apiCall.LogInfo(fmt.Sprintf("SKIPPED Condition Group: %s. No condition to process.", cg.Name))
		return false
	}

	conditionFailed := false

	apiCall.LogInfo(fmt.Sprintf("Processing Condition Group: %s", cg.Name))

	for _, c := range cg.Conditions {
		apiCall.LogInfo(fmt.Sprintf("====> Started Processing Condition: %s", c.Name))
		conditionPassed := c.HasPassed(apiCall)
		if !conditionPassed {
			apiCall.LogError(fmt.Sprintf("Condition Condition: %s. Passed? %t. Condition Group Failed", c.Name, conditionPassed))
			conditionFailed = true
			break
		} else {
			apiCall.LogInfo(fmt.Sprintf("Condition Condition: %s. Passed? %t", c.Name, conditionPassed))
		}

		apiCall.LogInfo(fmt.Sprintf("====> Finised Processing Condition: %s", c.Name))
		apiCall.LogInfo("--")

	}

	if !conditionFailed {

		if cg.CallActualEndPoint {
			apiCall.LogInfo(fmt.Sprintf("Set to call ActualEndPoint %s", apiCall.CurrentEndPoint.ActualURL))
			apiCall.UseActualURL = true

		}

		apiCall.LogInfo("Condition Group Passes. Starting assignement")

		// set status code bases on condition group
		if cg.ResponseID != "" {
			if !apiCall.HasSet("*HTTP_STATUS_CODE") {
				apiCall.ResponseID = cg.ResponseID
				apiCall.SetKey("*HTTP_STATUS_CODE")
				//apiCall.LogInfo(fmt.Sprintf("Setting Http code: %d", cg.HttpStatusCode))
			} else {
				apiCall.LogInfo("Http code already assigned. No overrides")

			}
		}

		// set parameter value
		for _, p := range cg.ConditionGroupParameters {

			// update response param
			if strings.TrimSpace(p.AssgineValue) != "" {
				p.ResponseParam.process(apiCall, cg.ResponseID, p.AssgineValue)
			} else {
				apiCall.LogInfo(fmt.Sprintf("Param %s assignement skipped. Blank value", p.ResponseParam.Key))

			}
		}
		return true
	}

	return false

}

// -----------------------------------------------------------------
//
// -----------------------------------------------------------------
// Define a new UserModel type which wraps a database connection pool.
type ConditionGroupModel struct {
	DB *bolt.DB
}

func (m *ConditionGroupModel) getTableName() []byte {
	return []byte("endpointconditionsgroup")
}

// -----------------------------------------------------------------
//
// -----------------------------------------------------------------
// We'll use the Insert method to add a new record to the "users" table.
func (m *ConditionGroupModel) Save(u *ConditionGroup) (string, error) {
	if u.ID == "" {
		var id string = uuid.NewString()
		u.ID = fmt.Sprintf("%s_%s", u.EndpointID, id)
	}
	err := m.Update(u, false)

	return u.ID, err
}

// -----------------------------------------------------------------
//
// -----------------------------------------------------------------
// We'll use the Insert method to add a new record to the "users" table.
func (m *ConditionGroupModel) Update(u *ConditionGroup, clearCache bool) error {
	err := m.DB.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(m.getTableName())
		if err != nil {
			return err
		}
		u.Name = strings.ToUpper(u.Name)
		buf, err := json.Marshal(u)
		if err != nil {
			return err
		}

		// key = > user.name+ user.id
		key := strings.ToUpper(u.ID) // + string(itob(u.ID))

		return bucket.Put([]byte(key), buf)
	})

	return err

}

// -----------------------------------------------------------------
//
// -----------------------------------------------------------------
// We'll use the Insert method to add a new record to the "users" table.
func (m *ConditionGroupModel) Delete(id string) error {

	err := m.DB.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(m.getTableName())
		if err != nil {
			return err
		}
		key := strings.ToUpper(id)
		dbDeleteError := bucket.Delete([]byte(key))
		return dbDeleteError
	})

	return err
}

// -----------------------------------------------------------------
//
// -----------------------------------------------------------------
// We'll use the Exists method to check if a user exists with a specific ID.
func (m *ConditionGroupModel) Get(id string) (*ConditionGroup, error) {

	if id == "" {
		return nil, errors.New("param blank id not allowed")
	}
	var paramJSON []byte // = make([]byte, 0)

	err := m.DB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(m.getTableName())
		if bucket == nil {
			return errors.New("table does not exits")
		}
		paramJSON = bucket.Get([]byte(strings.ToUpper(id)))

		return nil

	})
	param := ConditionGroup{}
	if err != nil {
		return &param, err
	}

	if paramJSON != nil {
		err := json.Unmarshal(paramJSON, &param)
		if err == nil {
			m.AssignConditions(&param)
			m.AssignResponseParam(&param)
			m.RemoveBlankConditionIds(&param)

		}
		return &param, err

	}

	return &param, ErrServerNotFound

}

// -----------------------------------------------------------------
//
// -----------------------------------------------------------------
// We'll use the Exists method to check if a user exists with a specific ID.
func (m *ConditionGroupModel) ListById(endpointid string) []*ConditionGroup {
	params := make([]*ConditionGroup, 0)
	_ = m.DB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(m.getTableName())
		if bucket == nil {
			return errors.New("table does not exits")
		}
		c := bucket.Cursor()
		prefix := []byte(strings.ToUpper(endpointid) + "_")

		for k, v := c.First(); k != nil; k, v = c.Next() {
			if bytes.HasPrefix(k, prefix) {

				param := ConditionGroup{}
				err := json.Unmarshal(v, &param)
				if err == nil {
					m.AssignConditions(&param)
					m.AssignResponseParam(&param)
					m.RemoveBlankConditionIds(&param)
					params = append(params, &param)
				}
			}
		}

		return nil
	})

	sort.Slice(params, func(i, j int) bool {
		return params[i].Name < params[j].Name
	})

	return params

}

// -----------------------------------------------------------------
//
// -----------------------------------------------------------------
// We'll use the Exists method to check if a user exists with a specific ID.
func (m *ConditionGroupModel) DuplicateName(serverToCheck *ConditionGroup, forEndPoint EndPoint) bool {
	exists := false
	for _, server := range m.ListById(forEndPoint.ID) {
		if strings.EqualFold(server.Name, serverToCheck.Name) && !strings.EqualFold(server.ID, serverToCheck.ID) {
			exists = true
			break
		}
	}

	return exists
}

// -----------------------------------------------------------------
//
// -----------------------------------------------------------------
// We'll use the Insert method to add a new record to the "users" table.
func (m *ConditionGroupModel) AssignConditions(u *ConditionGroup) {

	u.Conditions = make([]*Condition, 0)
	conditionDB := &ConditionModel{DB: m.DB}

	for _, conditionID := range u.ConditionIDs {
		condition, err := conditionDB.Get(conditionID)
		if err == nil && condition.RequestParam != nil {
			u.Conditions = append(u.Conditions, condition)
		}
	}
}

// -----------------------------------------------------------------
//
// -----------------------------------------------------------------
// We'll use the Insert method to add a new record to the "users" table.
func (m *ConditionGroupModel) AssignResponseParam(u *ConditionGroup) {
	responseParamDB := &EndPointResponseParamModel{DB: m.DB}

	for _, mappedParam := range u.ConditionGroupParameters {

		if mappedParam == nil { //TODO > why nill
			continue
		}
		responseId := mappedParam.ResponseVariable

		r, err := responseParamDB.Get(responseId)
		if err == nil {
			mappedParam.ResponseParam = r

		}
	}

	// TODO >>>  u.RemoveInvalidResponseParam()

}

// -----------------------------------------------------------------
//
// -----------------------------------------------------------------
// We'll use the Insert method to add a new record to the "users" table.
func (m *ConditionGroupModel) RemoveBlankConditionIds(u *ConditionGroup) {

	newList := make([]string, 0)
	for _, conditionId := range u.ConditionIDs {
		if conditionId != "" {
			newList = append(newList, conditionId)
		}
	}

	u.ConditionIDs = newList

}

// -----------------------------------------------------------------
//
// -----------------------------------------------------------------
func (m *ConditionGroupModel) ClearEndPointData(endpointID string) {
	rlist := m.ListById(endpointID)
	for _, r := range rlist {
		m.Delete(r.ID)
	}
}
