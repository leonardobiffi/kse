package version

var Version = "0.0.0-dev"
var template = `kse v{{printf "%s" .Version}}
`

func String() string {
	return Version
}

func Template() string {
	return template
}
