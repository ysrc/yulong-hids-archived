package models

// Server model definiton.
type Server struct {
	baseModel
}

// NewServer init collectionName
func NewServer() Server {
	mdl := Server{}
	mdl.collectionName = "server"
	return mdl
}
