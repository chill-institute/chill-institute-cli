package main

const (
	rootPackageName = "@chill-institute/cli"
	binaryName      = "chilly"
	repositoryURL   = "git+https://github.com/chill-institute/chill-cli.git"
	homepageURL     = "https://github.com/chill-institute/chill-cli#readme"
)

type options struct {
	distDir string
	outDir  string
	version string
}
