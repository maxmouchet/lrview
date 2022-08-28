package lib

import (
	"fmt"
	"github.com/aalpern/luminosity"
	"github.com/jmoiron/sqlx"
	"path/filepath"
)

type Previews struct {
	db   *sqlx.DB
	path string
}

func OpenPreviews(path string) (*Previews, error) {
	dbPath := filepath.Join(path, "previews.db")
	db, err := sqlx.Open("sqlite3", dbPath+"?mode=ro")
	if err != nil {
		return nil, err
	}
	previews := &Previews{
		db:   db,
		path: path,
	}
	return previews, nil
}

func (p *Previews) Close() error {
	return p.db.Close()
}

func (p *Previews) GetImageCacheEntry(id int64) (*ImageCacheEntry, error) {
	const query = `
	SELECT UUID, Digest
	FROM ImageCacheEntry
	WHERE imageId = ?
	`
	return QueryOne[ImageCacheEntry](p.db, query, id)
}

func (p *Previews) GetPreviewPath(id int64) (string, error) {
	ice, err := p.GetImageCacheEntry(id)
	if err != nil {
		return "", err
	}
	path := filepath.Join(
		p.path,
		ice.UUID[0:1],
		ice.UUID[0:4],
		fmt.Sprintf("%s-%s.lrprev", ice.UUID, ice.Digest),
	)
	return path, nil
}

func (p *Previews) GetPreview(id int64) ([]byte, error) {
	path, err := p.GetPreviewPath(id)
	if err != nil {
		return nil, err
	}
	pf, err := luminosity.OpenPreviewFile(path)
	if err != nil {
		return nil, err
	}
	defer pf.Close()
	if len(pf.Sections) < 2 {
		return nil, fmt.Errorf("no embedded previews available")
	}
	return pf.Sections[len(pf.Sections)-1].ReadData()
}
