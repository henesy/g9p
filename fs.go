package g9p

import (
trie	"github.com/henesy/fstrie"
)

// FS represents a trie filesystem
type FS trie.Trie

// Make a new Filesystem
func (srv *Srv) MkFs(root File) FS {
	fs := FS(trie.New())
	//root.FileName = "/"
	fs.Root.Data = root
	srv.Fs = fs
	return srv.Fs
}

// Get a File by name and extract it
func (fs *FS) Get(name string) File {
	// TODO -- name sanitization

	return (*trie.Trie)(fs).Find(name).Data.(File)
}
