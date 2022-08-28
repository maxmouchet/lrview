package lib

import (
	"github.com/jmoiron/sqlx"
	"path/filepath"
	"strings"
)

type Catalog struct {
	db   *sqlx.DB
	path string
}

func OpenCatalog(path string) (*Catalog, error) {
	db, err := sqlx.Open("sqlite3", path+"?mode=ro")
	if err != nil {
		return nil, err
	}
	catalog := &Catalog{
		db:   db,
		path: path,
	}
	return catalog, nil
}

func (c *Catalog) Close() error {
	return c.db.Close()
}

func (c *Catalog) Name() string {
	return strings.TrimSuffix(filepath.Base(c.path), ".lrcat")
}

func (c *Catalog) GetCollections() ([]Collection, error) {
	const query = `
	SELECT id_local, name, parent
	FROM AgLibraryCollection
	ORDER BY name
	`
	return QueryAll[Collection](c.db, query)
}

func (c *Catalog) GetCollectionsImages(ids []int64) ([]Image, error) {
	query := `
	SELECT Adobe_images.id_local
	FROM Adobe_images
	JOIN AgLibraryCollectionImage
	ON Adobe_images.id_local = AgLibraryCollectionImage.image
	WHERE AgLibraryCollectionImage.collection IN (?)
	`
	query, args, err := sqlx.In(query, ids)
	if err != nil {
		return nil, err
	}
	return QueryAll[Image](c.db, query, args...)
}

func (c *Catalog) GetCollectionsTree() (*TreeNode, error) {
	collections, err := c.GetCollections()
	if err != nil {
		return nil, err
	}
	root := TreeNode{
		ID:       -1,
		Name:     "Collections",
		Kind:     TreeNodeKindCollection,
		Children: []*TreeNode{},
	}
	nodes := make(map[int64]*TreeNode)
	for _, collection := range collections {
		nodes[collection.ID] = collection.TreeNode()
	}
	for _, collection := range collections {
		parent := &root
		if collection.Parent != nil {
			parent = nodes[*collection.Parent]
		}
		parent.Children = append(parent.Children, nodes[collection.ID])
	}
	return &root, nil
}

func (c *Catalog) GetFolders() ([]Folder, error) {
	const query = `
	SELECT id_local, pathFromRoot
	FROM AgLibraryFolder
	ORDER BY pathFromRoot
	`
	return QueryAll[Folder](c.db, query)
}

func (c *Catalog) GetFoldersImages(ids []int64) ([]Image, error) {
	query := `
	SELECT Adobe_images.id_local
	FROM Adobe_images
	JOIN AgLibraryFile
	ON Adobe_images.rootFile = AgLibraryFile.id_local
	JOIN AgLibraryFolder
	ON AgLibraryFile.folder = AgLibraryFolder.id_local
	WHERE AgLibraryFolder.id_local = (?)
	`
	query, args, err := sqlx.In(query, ids)
	if err != nil {
		return nil, err
	}
	return QueryAll[Image](c.db, query, args...)
}

func (c *Catalog) GetFoldersTree() (*TreeNode, error) {
	folders, err := c.GetFolders()
	if err != nil {
		return nil, err
	}
	pathToId := make(map[string]int64)
	for _, folder := range folders {
		pathFromRoot := strings.TrimSuffix(folder.PathFromRoot, "/")
		pathToId[pathFromRoot] = folder.ID
	}
	root := TreeNode{
		ID:       -1,
		Name:     "Folders",
		Kind:     TreeNodeKindFolder,
		Children: []*TreeNode{},
	}
	nodes := make(map[int64]*TreeNode)
	for _, folder := range folders {
		nodes[folder.ID] = folder.TreeNode()
	}
	for _, folder := range folders {
		pathFromRoot := strings.TrimSuffix(folder.PathFromRoot, "/")
		pathComponents := strings.Split(pathFromRoot, "/")
		parent := &root
		if len(pathComponents) > 1 {
			parentPath := strings.Join(pathComponents[:len(pathComponents)-1], "/")
			parent = nodes[pathToId[parentPath]]
		}
		parent.Children = append(parent.Children, nodes[folder.ID])
	}
	return &root, nil
}
