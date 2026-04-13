package cli

const (
	ansiReset = "\033[0m"
	ansiBold  = "\033[1m"
	ansiDim   = "\033[2m"
)

func bold(s string) string {
	return ansiBold + s + ansiReset
}

func dim(s string) string {
	return ansiDim + s + ansiReset
}
