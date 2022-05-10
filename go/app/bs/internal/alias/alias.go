package alias

type IAlias interface {
	Next() (string, error)
	Validate(alias string) bool
	Close()
}
