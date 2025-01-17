module git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git

go 1.23.2

require (
	git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git v0.0.0-20241002165513-602012adc104
	git-codecommit.eu-central-1.amazonaws.com/v1/repos/prng.git v0.0.0-20230405095258-31d142aea225
	github.com/goccy/go-json v0.10.2
	github.com/stretchr/testify v1.8.4
	golang.org/x/text v0.19.0
)

require (
	github.com/aead/chacha20 v0.0.0-20180709150244-8b13a72661da // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/sys v0.18.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

//replace (
//	git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git => ../go-utils
//	git-codecommit.eu-central-1.amazonaws.com/v1/repos/prng.git => ../prng
//)
