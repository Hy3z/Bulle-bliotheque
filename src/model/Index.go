package model

type Index struct {
	IsLogged bool
	IsAdmin  bool
	Query    string
	Data     any
}
