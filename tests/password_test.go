package test

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestAddingPassword(t *testing.T) {
	vault := TestVault()
	vault.NewPassword("Facebook", "https://facebook.com", "username", "password", "basic")

	assert.Assert(t, len(vault.Passwords) == 1, "Password was not added")
	assert.Assert(t, vault.Passwords[0].Name == "Facebook", "Password name was not added")
	assert.Assert(t, vault.Passwords[0].Url == "https://facebook.com", "Password url was not added")
	assert.Assert(t, vault.Passwords[0].Username == "username", "Password username was not added")
	assert.Assert(t, vault.Passwords[0].Hint == "basic", "Password hint was not added")

	assert.Assert(t, vault.Passwords[0].Hash != "password", "Password was not encrypted")
}

func TestReadingPassword(t *testing.T) {
	vault := TestVault()
	vault.NewPassword("Facebook", "https://facebook.com", "username", "password", "basic")

	password, err := vault.GetPassword("Facebook")
	assert.NilError(t, err, "Password should be read")

	passwordString := vault.ReadPassword(&password)
	assert.Assert(t, passwordString == "password", "Password was read incorrectly")
}

func TestUpdatingPassword(t *testing.T) {
	vault := TestVault()
	vault.NewPassword("Facebook", "https://facebook.com", "username", "password", "basic")

	err := vault.UpdatePassword("Facebook", "newpassword", "newpassword")
	assert.NilError(t, err, "Password should be updated")

	passwordMatch, err := vault.VerifyPassword("Facebook", "newpassword")
	assert.Assert(t, passwordMatch, "Password was not updated")
}
