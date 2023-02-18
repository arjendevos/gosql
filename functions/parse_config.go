package functions

import (
	"io/ioutil"

	"github.com/BurntSushi/toml"
	"golang.org/x/mod/modfile"
)

type SqlBoilerConfig struct {
	Output       string
	PkgName      string `toml:"pkgname"`
	Wipe         bool
	NoTests      bool `toml:"no-tests"`
	AddEnumTypes bool `toml:"add-enum-types"`
	Psql         struct {
		DbName    string `toml:"dbname"`
		Host      string
		Port      int
		User      string
		Pass      string
		SslMode   string `toml:"sslmode"`
		Schema    string
		Blacklist []string
	}
}

func parseSqlBoilerConfig() (*SqlBoilerConfig, error) {
	config := SqlBoilerConfig{}

	// Read the TOML file
	data, err := ioutil.ReadFile("sqlboiler.toml")
	if err != nil {
		return nil, err
	}

	// Parse the TOML data into the Config struct
	if _, err := toml.Decode(string(data), &config); err != nil {
		return nil, err
	}
	return &config, nil
}

func parseModFile() (*modfile.File, error) {
	content, err := ioutil.ReadFile("go.mod")
	if err != nil {
		return nil, err
	}

	// Parse the go.mod file
	mod, err := modfile.Parse("go.mod", content, nil)
	if err != nil {
		return nil, err
	}

	return mod, nil
}
