module github.com/bengarrett/zipcmt

// NOTE: .github/workflows/release.yml also requires the Go version
go 1.26.8

require (
	github.com/bengarrett/retrotxtgo v1.2.2
	github.com/bengarrett/sauce v1.2.9
	github.com/dustin/go-humanize v1.0.1
	github.com/gookit/color v1.6.1
	github.com/muesli/go-app-paths v0.2.2
	golang.org/x/text v0.42.0
)

require (
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/xo/terminfo v1.2.0 // indirect
	go.uber.org/nilaway v0.0.0-20251021214447-34f56b8c16b9 // indirect
	golang.org/x/mod v0.41.0 // indirect
	golang.org/x/sync v0.23.0 // indirect
	golang.org/x/sys v0.48.0 // indirect
	golang.org/x/tools v0.49.0 // indirect
)

tool go.uber.org/nilaway/cmd/nilaway
