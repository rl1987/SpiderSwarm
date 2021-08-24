package spsw

type MasterStore struct {
	UUID     string
	Backends []MasterStoreBackend
}
