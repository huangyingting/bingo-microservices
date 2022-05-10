package search

type Alias struct {
	Alias string   `json:"alias"`
	Oid   string   `json:"oid"`
	Tags  []string `json:"tags"`
}

type ISearch interface {
	Index(aliasTags Alias) error
	SearchAnd(oid string, tags []string) ([]Alias, error)
	SearchOr(oid string, tags []string) ([]Alias, error)
	Suggest(text string) ([]string, error)
	Delete(alias string, oid string) error
}
