package models

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	bolt "go.etcd.io/bbolt"
)

// -----------------------------------------------------------------
//
// -----------------------------------------------------------------
// Define a new User type. Notice how the field names and types align
// with the columns in the database "users" table?
type EndPointRequestParam struct {
	ID                string `json:"id" db:"id" form:"id"`
	EndpointID        string `json:"endpointid" db:"endpointid" form:"endpointid"`
	Key               string `json:"key" db:"key" form:"key"`
	DefaultValue      any    `json:"defaultvalue" db:"defaultvalue" form:"defaultvalue"`
	OverrideValue     any    `json:"overridevalue" db:"overridevalue" form:"overridevalue"`
	DefaultDatatype   string `json:"defaultdatatype" db:"defaultdatatype" form:"defaultdatatype"`
	OverrrideDatatype string `json:"overridedatatype" db:"overridedatatype" form:"overridedatatype"`
}

// -----------------------------------------------------------------
//
// -----------------------------------------------------------------
// Define a new UserModel type which wraps a database connection pool.
type EndPointRequestParamModel struct {
	DB *bolt.DB
}

func (m *EndPointRequestParamModel) getTableName() []byte {
	return []byte("requestparam")
}

// -----------------------------------------------------------------
//
// -----------------------------------------------------------------
// We'll use the Insert method to add a new record to the "users" table.
func (m *EndPointRequestParamModel) Save(u *EndPointRequestParam) (string, error) {
	if u.ID == "" {
		var id string = u.Key // uuid.NewString()
		u.ID = fmt.Sprintf("%s_%s", u.EndpointID, id)
	}
	err := m.Update(u, false)

	return u.ID, err
}

// -----------------------------------------------------------------
//
// -----------------------------------------------------------------
// We'll use the Insert method to add a new record to the "users" table.
func (m *EndPointRequestParamModel) Update(u *EndPointRequestParam, clearCache bool) error {
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
func (m *EndPointRequestParamModel) Delete(id string) error {

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
func (m *EndPointRequestParamModel) Get(id string) (*EndPointRequestParam, error) {

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
	param := EndPointRequestParam{}
	if err != nil {
		return &param, err
	}

 

	if paramJSON != nil {
		err := json.Unmarshal(paramJSON, &param)
		return &param, err

	}

	return &param, ErrServerNotFound

}

// -----------------------------------------------------------------
//
// -----------------------------------------------------------------
// We'll use the Exists method to check if a user exists with a specific ID.
func (m *EndPointRequestParamModel) ListById(endpointID string) []*EndPointRequestParam {
	params := make([]*EndPointRequestParam, 0)
	_ = m.DB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(m.getTableName())
		if bucket == nil {
			return errors.New("table does not exits")
		}
		c := bucket.Cursor()
		prefix := []byte(strings.ToUpper(endpointID) + "_")

		for k, v := c.First(); k != nil; k, v = c.Next() {
			if bytes.HasPrefix(k, prefix) {
				param := EndPointRequestParam{}
				err := json.Unmarshal(v, &param)
				if err == nil {
					params = append(params, &param)
				}
			}
		}

		return nil
	})
	return params

}

// -----------------------------------------------------------------
//
// -----------------------------------------------------------------
func (m *EndPointRequestParamModel) ClearEndPointData(endpointID string)  { 
	rlist:= m.ListById(endpointID)
	for _, r:= range rlist{
		m.Delete(r.ID)
	}
}