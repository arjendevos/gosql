package am

import (
	"github.com/arjendevos/gosql/models/dm"
)

type Account struct {
	*dm.Account

	InstagramAccounts               dm.InstagramAccountSlice               `boil:"InstagramAccounts" json:"InstagramAccounts" toml:"InstagramAccounts" yaml:"InstagramAccounts"`
	Proxies                         dm.ProxySlice                          `boil:"Proxies" json:"Proxies" toml:"Proxies" yaml:"Proxies"`
	ScrapedInstagramAccounts        dm.ScrapedInstagramAccountSlice        `boil:"ScrapedInstagramAccounts" json:"ScrapedInstagramAccounts" toml:"ScrapedInstagramAccounts" yaml:"ScrapedInstagramAccounts"`
	ScrapedInstagramAccountChaineds dm.ScrapedInstagramAccountChainedSlice `boil:"ScrapedInstagramAccountChaineds" json:"ScrapedInstagramAccountChaineds" toml:"ScrapedInstagramAccountChaineds" yaml:"ScrapedInstagramAccountChaineds"`
}

type AccountSlice []*Account

func SqlBoilerAccountsToApiAccounts(a dm.AccountSlice) AccountSlice {
	var s AccountSlice
	for _, d := range a {
		s = append(s, SqlBoilerAccountToApiAccount(d))
	}
	return s
}

func SqlBoilerAccountToApiAccount(a *dm.Account) *Account {
	return &Account{
		Account: a,

		InstagramAccounts:               a.R.InstagramAccounts,
		Proxies:                         a.R.Proxies,
		ScrapedInstagramAccounts:        a.R.ScrapedInstagramAccounts,
		ScrapedInstagramAccountChaineds: a.R.ScrapedInstagramAccountChaineds,
	}
}
