package am

import (
	"github.com/arjendevos/gosql/models/dm"
)

type ScrapedInstagramAccount struct {
	*dm.ScrapedInstagramAccount

	Account                         *dm.Account                            `boil:"Account" json:"Account" toml:"Account" yaml:"Account"`
	ScrapedInstagramAccountChaineds dm.ScrapedInstagramAccountChainedSlice `boil:"ScrapedInstagramAccountChaineds" json:"ScrapedInstagramAccountChaineds" toml:"ScrapedInstagramAccountChaineds" yaml:"ScrapedInstagramAccountChaineds"`
}

type ScrapedInstagramAccountSlice []*ScrapedInstagramAccount

func SqlBoilerScrapedInstagramAccountsToApiScrapedInstagramAccounts(a dm.ScrapedInstagramAccountSlice) ScrapedInstagramAccountSlice {
	var s ScrapedInstagramAccountSlice
	for _, d := range a {
		s = append(s, SqlBoilerScrapedInstagramAccountToApiScrapedInstagramAccount(d))
	}
	return s
}

func SqlBoilerScrapedInstagramAccountToApiScrapedInstagramAccount(a *dm.ScrapedInstagramAccount) *ScrapedInstagramAccount {
	return &ScrapedInstagramAccount{
		ScrapedInstagramAccount: a,

		Account:                         a.R.Account,
		ScrapedInstagramAccountChaineds: a.R.ScrapedInstagramAccountChaineds,
	}
}
