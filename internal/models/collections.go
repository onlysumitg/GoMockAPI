package models

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/onlysumitg/GoMockAPI/internal/validator"
	"github.com/onlysumitg/GoMockAPI/utils/stringutils"
	bolt "go.etcd.io/bbolt"
)

type Collection struct {
	ID string `json:"id" db:"id" form:"id"`

	Name string `json:"name" db:"name" form:"name"`
	Desc string `json:"desc" db:"desc" form:"desc"`

	validator.Validator `json:"-" db:"-" form:"-"`
}

// -----------------------------------------------------------------
//
// -----------------------------------------------------------------
// Define a new UserModel type which wraps a database connection pool.
type CollectionModel struct {
	DB *bolt.DB
}

func (m *CollectionModel) getTableName() []byte {
	return []byte("collections")
}

// -----------------------------------------------------------------
//
// -----------------------------------------------------------------
// We'll use the Insert method to add a new record to the "users" table.
func (m *CollectionModel) Save(u *Collection) error {
	if u.ID == "" {
		u.ID = uuid.NewString()
	}

	err := m.DB.Update(func(tx *bolt.Tx) error {
		u.Name = strings.ToUpper(u.Name)
		u.Name = stringutils.RemoveSpecialChars(stringutils.RemoveMultipleSpaces(u.Name))

		bucket, err := tx.CreateBucketIfNotExists(m.getTableName())
		if err != nil {
			return err
		}

		buf, err := json.Marshal(u)
		if err != nil {
			return err
		}

		key := strings.ToUpper(u.ID)

		return bucket.Put([]byte(key), buf)
	})

	return err

}

// -----------------------------------------------------------------
//
// -----------------------------------------------------------------
// We'll use the Insert method to add a new record to the "users" table.
func (m *CollectionModel) Delete(id string) error {

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
func (m *CollectionModel) Get(id string) (*Collection, error) {

	if id == "" {
		return nil, errors.New("blank id not allowed")
	}
	var calllogJSON []byte // = make([]byte, 0)

	err := m.DB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(m.getTableName())
		if bucket == nil {
			return errors.New("table does not exits")
		}
		calllogJSON = bucket.Get([]byte(strings.ToUpper(id)))

		return nil

	})
	pr := Collection{}
	if err != nil {
		return &pr, err
	}

	// log.Println("calllogJSON >2 >>", calllogJSON)

	if calllogJSON != nil {
		err := json.Unmarshal(calllogJSON, &pr)
		return &pr, err

	}

	return &pr, errors.New("not found")

}

// -----------------------------------------------------------------
//
// -----------------------------------------------------------------
// We'll use the Exists method to check if a user exists with a specific ID.
func (m *CollectionModel) List() []*Collection {
	prs := make([]*Collection, 0)
	_ = m.DB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(m.getTableName())
		if bucket == nil {
			return errors.New("table does not exits")
		}
		c := bucket.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {

			pr := Collection{}
			err := json.Unmarshal(v, &pr)
			if err == nil {
				prs = append(prs, &pr)
			}
		}

		return nil
	})
	return prs

}

// -----------------------------------------------------------------
//
// -----------------------------------------------------------------
// We'll use the Exists method to check if a user exists with a specific ID.
func (m *CollectionModel) DuplicateName(collectionToCheck *Collection) bool {
	exists := false
	for _, collection := range m.List() {
		if strings.EqualFold(collection.Name, collectionToCheck.Name) && !strings.EqualFold(collection.ID, collectionToCheck.ID) {
			exists = true
			break
		}
	}

	return exists
}
