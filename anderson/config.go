package anderson

type Config struct {
	Whitelist  []string `yaml:"whitelist"`
	Blacklist  []string `yaml:"blacklist"`
	Exceptions []string `yaml:"exceptions"`
}
