package main

type target struct {
	goOS       string
	goArch     string
	npmOS      string
	npmArch    string
	suffix     string
	binaryFile string
}

var targets = []target{
	{goOS: "darwin", goArch: "amd64", npmOS: "darwin", npmArch: "x64", suffix: "darwin-x64", binaryFile: binaryName},
	{goOS: "darwin", goArch: "arm64", npmOS: "darwin", npmArch: "arm64", suffix: "darwin-arm64", binaryFile: binaryName},
	{goOS: "linux", goArch: "amd64", npmOS: "linux", npmArch: "x64", suffix: "linux-x64", binaryFile: binaryName},
	{goOS: "linux", goArch: "arm64", npmOS: "linux", npmArch: "arm64", suffix: "linux-arm64", binaryFile: binaryName},
	{goOS: "windows", goArch: "amd64", npmOS: "win32", npmArch: "x64", suffix: "win32-x64", binaryFile: binaryName + ".exe"},
	{goOS: "windows", goArch: "arm64", npmOS: "win32", npmArch: "arm64", suffix: "win32-arm64", binaryFile: binaryName + ".exe"},
}

func platformPackageName(t target) string {
	return rootPackageName + "-" + t.suffix
}

func platformDirName(t target) string {
	return "cli-" + t.suffix
}
