package am

import (
	"github.com/arjendevos/gosql/models/dm"
)

type ScrapedInstagramAccountChained struct {
	*dm.ScrapedInstagramAccountChained

	Account                 *dm.Account                 `boil:"Account" json:"Account" toml:"Account" yaml:"Account"`
	ScrapedInstagramAccount *dm.ScrapedInstagramAccount `boil:"ScrapedInstagramAccount" json:"ScrapedInstagramAccount" toml:"ScrapedInstagramAccount" yaml:"ScrapedInstagramAccount"`
}

type ScrapedInstagramAccountChainedSlice []*ScrapedInstagramAccountChained

func SqlBoilerScrapedInstagramAccountChainedsToApiScrapedInstagramAccountChaineds(a dm.ScrapedInstagramAccountChainedSlice) ScrapedInstagramAccountChainedSlice {
	var s ScrapedInstagramAccountChainedSlice
	for _, d := range a {
		s = append(s, SqlBoilerScrapedInstagramAccountChainedToApiScrapedInstagramAccountChained(d))
	}
	return s
}

func SqlBoilerScrapedInstagramAccountChainedToApiScrapedInstagramAccountChained(a *dm.ScrapedInstagramAccountChained) *ScrapedInstagramAccountChained {
	return &ScrapedInstagramAccountChained{
		ScrapedInstagramAccountChained: a,

		Account:                 a.R.Account,
		ScrapedInstagramAccount: a.R.ScrapedInstagramAccount,
	}
}
