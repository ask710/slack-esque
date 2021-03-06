package users

import (
	"errors"

	"github.com/info344-s18/challenges-ask710/servers/gateway/indexes"
)

//MockStore is a struct for a mock user store
type MockStore struct {
	TriggerError bool
	Result       *User
}

//NewMockStore creates a new MockStore struct
func NewMockStore(triggerError bool, result *User) *MockStore {
	return &MockStore{
		TriggerError: triggerError,
		Result:       result,
	}
}

//GetByID returns the User with the given ID
func (m *MockStore) GetByID(id int64) (*User, error) {
	if m.TriggerError {
		return nil, errors.New("Error with GetByID")
	}
	return m.Result, nil
}

//GetByEmail returns the User with the given email
func (m *MockStore) GetByEmail(email string) (*User, error) {
	if m.TriggerError {
		return nil, errors.New("Error with GetByEmail")
	}
	return m.Result, nil
}

//GetByUserName returns the User with the given Username
func (m *MockStore) GetByUserName(username string) (*User, error) {
	if m.TriggerError {
		return nil, errors.New("Error with GetByUserName")
	}
	return m.Result, nil
}

//Insert inserts the user into the database, and returns
//the newly-inserted User, complete with the DBMS-assigned ID
func (m *MockStore) Insert(user *User) (*User, error) {
	if m.TriggerError {
		return nil, errors.New("Error with Insert")
	}
	return m.Result, nil
}

//Update applies UserUpdates to the given user ID
//and returns the newly-updated user
func (m *MockStore) Update(id int64, updates *Updates) (*User, error) {
	if m.TriggerError {
		return nil, errors.New("Error with Update")
	}
	return m.Result, nil
}

//UpdatePhoto updates the photourl for a user
func (m *MockStore) UpdatePhoto(id int64, photourl string) (*User, error) {
	if m.TriggerError {
		return nil, errors.New("Error with UpdatePhoto")
	}
	return m.Result, nil
}

//Delete deletes the user with the given ID
func (m *MockStore) Delete(id int64) error {
	if m.TriggerError {
		return errors.New("Error with Delete")
	}
	return nil
}

//InsertLogin inserts login activity
func (m *MockStore) InsertLogin(login *Login) (*Login, error) {
	return nil, nil
}

//UpdatePassword updates password after resetting it.
func (m *MockStore) UpdatePassword(id int64, passHash []byte) (*User, error) {
	return nil, nil
}

//LoadUsers gets all users to add to the trie
func (m *MockStore) LoadUsers() (*indexes.Trie, error) {
	return nil, nil
}

//GetSearchUsers gets all users based on the found Ids
func (m *MockStore) GetSearchUsers(found []int64) (*[]User, error) {
	return nil, nil
}
