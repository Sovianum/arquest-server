package model

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"os/user"
	"testing"
)

func TestUser_Unmarshal_ParseError(t *testing.T) {
	u := &User{}
	data := []byte("{\"age\": 10, \"sex\": \"F\"")
	err := json.Unmarshal(data, u)
	assert.NotNil(t, err)
}

func TestUser_Unmarshal_InvalidUser(t *testing.T) {
	u := &User{}
	data := []byte("{\"login\": \"login\", \"password\": \"pass\", \"sex\": \"p\"}")
	err := json.Unmarshal(data, &u)

	assert.NotNil(t, err)
	assert.Equal(t, RegistrationInvalidSex, err.Error())
}

func TestUser_Unmarshal_Success(t *testing.T) {
	u := User{}
	data := []byte("{\"id\": 10, \"login\": \"login\", \"password\": \"password\"}")
	err := json.Unmarshal(data, &u)

	assert.Nil(t, err)
	assert.Equal(t, "login", u.Login)
	assert.Equal(t, 10, u.Id)
}
