package test

import "github.com/sordgom/PasswordManager/cli/pkg"

func TestVault() *pkg.Vault {
	v := pkg.Vault{
		Name:           "test",
		MasterPassword: "test",
	}
	return &v
}
