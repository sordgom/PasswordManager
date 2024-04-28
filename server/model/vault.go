package model

type Vault struct {
	Name           string
	MasterPassword string

	Passwords []Password
}

func New(name, masterPassword string) Vault {
	return Vault{
		Name:           name,
		MasterPassword: masterPassword,
	}
}

func (v *Vault) VerifyMasterPassword(newPassword string) bool {
	//logger
	return Compare(newPassword, v.MasterPassword)
}
