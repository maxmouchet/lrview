package lib

import (
	"path/filepath"
)

const (
	TreeNodeKindCollection = "collections"
	TreeNodeKindFolder     = "folders"
)

type Collection struct {
	ID     int64 `db:"id_local"`
	Name   string
	Parent *int64
}

func (c *Collection) TreeNode() *TreeNode {
	return &TreeNode{
		ID:       c.ID,
		Name:     c.Name,
		Kind:     TreeNodeKindCollection,
		Children: []*TreeNode{},
	}
}

type Folder struct {
	ID           int64  `db:"id_local"`
	PathFromRoot string `db:"pathFromRoot"`
}

func (f *Folder) TreeNode() *TreeNode {
	return &TreeNode{
		ID:       f.ID,
		Name:     filepath.Base(f.PathFromRoot),
		Kind:     TreeNodeKindFolder,
		Children: []*TreeNode{},
	}
}

type Image struct {
	ID int64 `db:"id_local"`
}

type ImageCacheEntry struct {
	UUID   string
	Digest string
}

type TreeNode struct {
	ID       int64
	Name     string
	Kind     string
	Children []*TreeNode
}
