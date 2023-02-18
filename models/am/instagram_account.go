package am

import (
	"github.com/arjendevos/gosql/models/dm"
)

type InstagramAccount struct {
	*dm.InstagramAccount

	Account    *dm.Account    `boil:"Account" json:"Account" toml:"Account" yaml:"Account"`
	Credential *dm.Credential `boil:"Credential" json:"Credential" toml:"Credential" yaml:"Credential"`
	Job        *dm.Job        `boil:"Job" json:"Job" toml:"Job" yaml:"Job"`
	Proxy      *dm.Proxy      `boil:"Proxy" json:"Proxy" toml:"Proxy" yaml:"Proxy"`
}

type InstagramAccountSlice []*InstagramAccount

func SqlBoilerInstagramAccountsToApiInstagramAccounts(a dm.InstagramAccountSlice) InstagramAccountSlice {
	var s InstagramAccountSlice
	for _, d := range a {
		s = append(s, SqlBoilerInstagramAccountToApiInstagramAccount(d))
	}
	return s
}

func SqlBoilerInstagramAccountToApiInstagramAccount(a *dm.InstagramAccount) *InstagramAccount {
	return &InstagramAccount{
		InstagramAccount: a,

		Account:    a.R.Account,
		Credential: a.R.Credential,
		Job:        a.R.Job,
		Proxy:      a.R.Proxy,
	}
}
