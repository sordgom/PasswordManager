package test

import (
	"testing"

	"github.com/sordgom/PasswordManager/cli/pkg"

	"gotest.tools/v3/assert"
)

func TestPasswordEncryption(t *testing.T) {
	v := pkg.Vault{
		Name:           "test",
		MasterPassword: "test",
	}
	encodedPassword, err := v.EncryptPassword("test")
	assert.NilError(t, err, "Password should be encrypted")
	assert.Assert(t, encodedPassword != "test", "Password should be encrypted")
}

func TestPasswordDecryption(t *testing.T) {
	v := pkg.Vault{
		Name:           "test",
		MasterPassword: "test",
	}
	v.NewPassword("test", "test", "test", "test", "test")

	decryptedPassword, err := v.DecryptPassword(v.Passwords[0].Hash)
	assert.NilError(t, err, "Password should be decrypted")
	assert.Assert(t, decryptedPassword == "test", "Password should be decrypted")
}
