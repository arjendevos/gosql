package am

import (
	"github.com/arjendevos/gosql/models/dm"
)

type Credential struct {
	*dm.Credential

	InstagramAccounts dm.InstagramAccountSlice `boil:"InstagramAccounts" json:"InstagramAccounts" toml:"InstagramAccounts" yaml:"InstagramAccounts"`
}

type CredentialSlice []*Credential

func SqlBoilerCredentialsToApiCredentials(a dm.CredentialSlice) CredentialSlice {
	var s CredentialSlice
	for _, d := range a {
		s = append(s, SqlBoilerCredentialToApiCredential(d))
	}
	return s
}

func SqlBoilerCredentialToApiCredential(a *dm.Credential) *Credential {
	return &Credential{
		Credential: a,

		InstagramAccounts: a.R.InstagramAccounts,
	}
}
