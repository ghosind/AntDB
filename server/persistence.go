package server

import (
	"path"

	"github.com/ghosind/antdb/persistence"
)

func (s *Server) saveRDB() error {
	filePath := path.Join(s.dir, s.dbFilename)

	return persistence.RDBSave(s.databases, filePath)
}
