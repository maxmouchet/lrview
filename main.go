package main

import (
	"embed"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"html/template"
	"log"
	"lrview/lib"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func exists(path string) {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		panic(fmt.Errorf("%s does not exists", path))
	}
}

//go:embed assets/* templates/*
var fs embed.FS

func main() {
	catalogPath := os.Getenv("LRVIEW_CATALOG_PATH")
	if catalogPath == "" {
		log.Fatal("Please set LRVIEW_CATALOG_PATH")
	}
	exists(catalogPath)
	catalog, err := lib.OpenCatalog(catalogPath)
	check(err)
	defer catalog.Close()

	previewsPath := os.Getenv("LRVIEW_PREVIEWS_PATH")
	if previewsPath == "" {
		previewsPath = strings.ReplaceAll(catalogPath, ".lrcat", " Previews.lrdata")
		log.Printf("LRVIEW_PREVIEWS_PATH is not set, using %s", previewsPath)
	}
	exists(previewsPath)
	previews, err := lib.OpenPreviews(previewsPath)
	check(err)
	defer previews.Close()

	r := gin.Default()
	r.SetHTMLTemplate(template.Must(template.New("").ParseFS(fs, "templates/*.html")))
	r.StaticFS("/static", http.FS(fs))

	r.GET("/", func(c *gin.Context) {
		collections, err := catalog.GetCollectionsTree()
		check(err)
		folders, err := catalog.GetFoldersTree()
		check(err)
		c.HTML(http.StatusOK, "index.html", gin.H{
			"catalog_name": catalog.Name(),
			"collections":  collections,
			"folders":      folders,
		})
	})

	r.GET("/collections/:id", func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		check(err)
		collections, err := catalog.GetCollectionsTree()
		check(err)
		folders, err := catalog.GetFoldersTree()
		check(err)
		images, err := catalog.GetCollectionsImages([]int64{id})
		check(err)
		c.HTML(http.StatusOK, "index.html", gin.H{
			"catalog_name": catalog.Name(),
			"collections":  collections,
			"folders":      folders,
			"images":       images,
		})
	})

	r.GET("/folders/:id", func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		check(err)
		collections, err := catalog.GetCollectionsTree()
		check(err)
		folders, err := catalog.GetFoldersTree()
		check(err)
		images, err := catalog.GetFoldersImages([]int64{id})
		check(err)
		c.HTML(http.StatusOK, "index.html", gin.H{
			"catalog_name": catalog.Name(),
			"collections":  collections,
			"folders":      folders,
			"images":       images,
		})
	})

	r.GET("/images/:id/preview", func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		check(err)
		data, err := previews.GetPreview(id)
		check(err)
		c.Data(http.StatusOK, "image/jpeg", data)
	})

	check(r.Run())
}
