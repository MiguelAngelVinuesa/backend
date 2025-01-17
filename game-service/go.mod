module git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service

go 1.23.2

toolchain go1.23.3

require (
	git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-config.git v0.0.0-20241221002003-83490f0896e4
	git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git v0.0.0-20241221001840-4d295adec4db
	git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-manager.git v0.0.0-20241221001934-5a777252bbbf
	git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git v0.0.0-20241002165513-602012adc104
	github.com/goccy/go-json v0.10.2
	github.com/gofiber/fiber/v2 v2.52.2
	github.com/stretchr/testify v1.9.0
	github.com/valyala/fasthttp v1.52.0
)

require (
	git-codecommit.eu-central-1.amazonaws.com/v1/repos/prng.git v0.0.0-20230405095258-31d142aea225 // indirect
	github.com/aead/chacha20 v0.0.0-20180709150244-8b13a72661da // indirect
	github.com/andybalholm/brotli v1.1.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/klauspost/compress v1.17.7 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.15 // indirect
	github.com/pierrec/lz4/v4 v4.1.21 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/twmb/franz-go v1.16.1 // indirect
	github.com/twmb/franz-go/pkg/kmsg v1.7.0 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/tcplisten v1.0.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	golang.org/x/sys v0.18.0 // indirect
	golang.org/x/text v0.19.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

//replace (
//	git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-config.git => ../game-config
//	git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git => ../game-engine
//	git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-manager.git => ../game-manager
//	git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git => ../go-utils
//)
