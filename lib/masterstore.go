package spiderswarm

type MasterStore struct {
	UUID     string
	Backends []MasterStoreBackend
}
