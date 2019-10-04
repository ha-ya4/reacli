package cmd


type cliContexter interface {
	String(name string) string
	Bool(name string) bool
}