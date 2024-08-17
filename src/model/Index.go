package model

// Index structure à passer en argument de toutes les templates "*-index" pour afficher une page entière
type Index struct {
	IsLogged bool
	IsAdmin  bool
	Query    string
	Data     any
}
