
package database

type BadgerStore struct{}

func NewBadgerStore() *BadgerStore {
	return &BadgerStore{}
}