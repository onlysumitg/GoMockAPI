package models

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/onlysumitg/GoMockAPI/internal/validator"
	bolt "go.etcd.io/bbolt"
)

// -----------------------------------------------------------------
//
// -----------------------------------------------------------------

type Condition struct {
	ID           string `json:"id" db:"id" form:"id"`
	EndpointID   string `json:"endpointid" db:"endpointid" form:"endpointid"`
	Name         string `json:"name" db:"name" form:"name"`
	Variable     string `json:"variable" db:"variable" form:"variable"`
	VariableName string `json:"variablename" db:"variablename" form:"variablename"`

	Operator          string `json:"operator" db:"operator" form:"operator"`
	Compareto         string `json:"compareto" db:"compareto" form:"compareto"`
	ComparetoDataType string `json:"comparetodatatype" db:"comparetodatatype" form:"comparetodatatype"`

	RequestParam *EndPointRequestParam `json:"-" db:"-" form:"-"`

	validator.Validator
}

// -----------------------------------------------------------------
//
// -----------------------------------------------------------------
func (m *Condition) HasPassed(apiCall *ApiCall) bool {

	requestValue, found := apiCall.RequestFlatMap[m.RequestParam.Key]

	// if not var found in request ==> failed
	if !found {
		apiCall.LogError(fmt.Sprintf("Condition Failed. Request param not found %s", m.RequestParam.Key))

		return false
	}

	operatorFunc, found := OperatorFuncMap[m.Operator]

	// if operator function not found
	// return true
	if !found {
		apiCall.LogInfo(fmt.Sprintf("Condition Passed. Operator NOT FOUND %s", m.Operator))

		return true
	}

	hasPassed := operatorFunc(requestValue.Value, m.Compareto, m.RequestParam.DefaultDatatype)
	apiCall.LogInfo(fmt.Sprintf("Condition Passed? %t", hasPassed))

	return hasPassed
}

// -----------------------------------------------------------------
//
// -----------------------------------------------------------------
// Define a new UserModel type which wraps a database connection pool.
type ConditionModel struct {
	DB *bolt.DB
}

func (m *ConditionModel) getTableName() []byte {
	return []byte("endpointconditions")
}

// -----------------------------------------------------------------
//
// -----------------------------------------------------------------
// We'll use the Insert method to add a new record to the "users" table.
func (m *ConditionModel) Save(u *Condition) (string, error) {
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
func (m *ConditionModel) Update(u *Condition, clearCache bool) error {
	err := m.DB.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(m.getTableName())
		if err != nil {
			return err
		}
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
func (m *ConditionModel) Delete(id string) error {

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
func (m *ConditionModel) Get(id string) (*Condition, error) {

	if id == "" {
		return nil, errors.New("param blank id not allowed")
	}
	var conditionJSON []byte // = make([]byte, 0)

	err := m.DB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(m.getTableName())
		if bucket == nil {
			return errors.New("table does not exits")
		}
		conditionJSON = bucket.Get([]byte(strings.ToUpper(id)))

		return nil

	})
	condition := Condition{}
	if err != nil {
		return &condition, err
	}

	// log.Println("serverJSON >2 >>", serverJSON)

	if conditionJSON != nil {
		err := json.Unmarshal(conditionJSON, &condition)

		if err == nil {
			m.AssignRequestParam(&condition)
		}

		return &condition, err

	}

	return &condition, ErrServerNotFound

}

// -----------------------------------------------------------------
//
// -----------------------------------------------------------------
// We'll use the Exists method to check if a user exists with a specific ID.
func (m *ConditionModel) ListById(endpointID string) []*Condition {
	conditions := make([]*Condition, 0)
	_ = m.DB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(m.getTableName())
		if bucket == nil {
			return errors.New("table does not exits")
		}
		c := bucket.Cursor()
		prefix := []byte(strings.ToUpper(endpointID) + "_")
		for k, v := c.First(); k != nil; k, v = c.Next() {
			if bytes.HasPrefix(k, prefix) {

				condition := Condition{}
				err := json.Unmarshal(v, &condition)
				if err == nil {
					m.AssignRequestParam(&condition)
					conditions = append(conditions, &condition)
				}
			}
		}

		return nil
	})
	return conditions

}

// -----------------------------------------------------------------
//
// -----------------------------------------------------------------
// We'll use the Exists method to check if a user exists with a specific ID.
func (m *ConditionModel) DuplicateName(condition *Condition, forEndPoint EndPoint) bool {
	exists := false
	for _, server := range m.ListById(forEndPoint.ID) {
		if strings.EqualFold(server.Name, condition.Name) && !strings.EqualFold(server.ID, condition.ID) {
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
func (m *ConditionModel) AssignRequestParam(u *Condition) {
	rpDB := &EndPointRequestParamModel{DB: m.DB}
	requestParam, errx := rpDB.Get(u.Variable)
	if errx == nil {
		u.RequestParam = requestParam
	}
}

// -----------------------------------------------------------------
//
// -----------------------------------------------------------------
func (m *ConditionModel) ClearEndPointData(endpointID string) {
	rlist := m.ListById(endpointID)
	for _, r := range rlist {
		m.Delete(r.ID)
	}
}
