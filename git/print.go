package git

import (
	"fmt"
	"io"
)

func (c *Commit) Print(w io.Writer) {
	fmt.Fprintf(w, "tree %s\n", c.Tree)
	fmt.Fprintf(w, "parent %s\n", c.Parent)
	fmt.Fprintf(w, "author %s\n", c.Author)
	fmt.Fprintf(w, "committer %s\n", c.Committer)
	fmt.Fprintf(w, "\n%s", c.Message)
}

func (tree *Tree) Print(w io.Writer) {
	for tree.Next() {
		entry := tree.Entry()
		fmt.Fprintf(w, "%08o %s %s\t%s\n", entry.Mode, entry.Type, entry.ID, entry.Name)
	}
}

func (b *Blob) Print(w io.Writer) {
	io.Copy(w, b)
}

func (t *Tag) Print(w io.Writer) {
	fmt.Fprintf(w, "object %s\n", t.Object)
	fmt.Fprintf(w, "type %s\n", t.ObjType)
	fmt.Fprintf(w, "tagger %s\n", t.Tagger)
	fmt.Fprintf(w, "\n%s", t.Message)
}

func (d *Delta) Print(w io.Writer) {
    
}
