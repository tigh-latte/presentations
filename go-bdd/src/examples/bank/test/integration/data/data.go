package data

import (
	"embed"
	"path"
)

type DataStore interface {
	Load(parts ...string) ([]byte, error)
}

type DataDir struct {
	Prefix string
	FS     embed.FS
}

func (d *DataDir) Load(parts ...string) ([]byte, error) {
	return d.FS.ReadFile(path.Join(append([]string{d.Prefix}, parts...)...))
}
