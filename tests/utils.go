package test

import "password-manager/pkg"

func TestVault() *pkg.Vault {
	v := pkg.Vault{
		Name:           "test",
		MasterPassword: "test",
	}
	return &v
}
