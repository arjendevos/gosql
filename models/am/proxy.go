package am

import (
	"github.com/arjendevos/gosql/models/dm"
)

type Proxy struct {
	*dm.Proxy

	Account           *dm.Account              `boil:"Account" json:"Account" toml:"Account" yaml:"Account"`
	InstagramAccounts dm.InstagramAccountSlice `boil:"InstagramAccounts" json:"InstagramAccounts" toml:"InstagramAccounts" yaml:"InstagramAccounts"`
}

type ProxySlice []*Proxy

func SqlBoilerProxiesToApiProxies(a dm.ProxySlice) ProxySlice {
	var s ProxySlice
	for _, d := range a {
		s = append(s, SqlBoilerProxyToApiProxy(d))
	}
	return s
}

func SqlBoilerProxyToApiProxy(a *dm.Proxy) *Proxy {
	return &Proxy{
		Proxy: a,

		Account:           a.R.Account,
		InstagramAccounts: a.R.InstagramAccounts,
	}
}
