package models

type Model interface {
	GetId() string
	GetTableName() string
}
