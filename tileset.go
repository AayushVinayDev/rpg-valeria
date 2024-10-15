package main

import (
	"encoding/json"
	"image"
	"os"
	"path/filepath"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Tileset interface {
	Img(id int) *ebiten.Image
}

type UniformTilesetJSON struct {
	Path string `json:"image"`
}

type UniformTileset struct {
	img *ebiten.Image
	gid int
}

func (u *UniformTileset) Img(id int) *ebiten.Image {

	id -= u.gid

	//get the position on the Image where the tile id is
	srcX := id % 22
	srcY := id / 22
	//convert the src tile position to pixel src position
	srcX *= 16
	srcY *= 16

	//draw the tile
	return u.img.SubImage(image.Rect(srcX, srcY, srcX+16, srcY+16)).(*ebiten.Image)
}

type TileJSON struct {
	Id     int    `json:"id"`
	Path   string `json:"image"`
	Width  int    `json:"imagewidth"`
	Height int    `json:"imageheight"`
}

type DynamicTilesetJSON struct {
	Tiles []*TileJSON `json:"tiles"`
}

type DynamicTileset struct {
	imgs []*ebiten.Image
	gid  int
}

func (d *DynamicTileset) Img(id int) *ebiten.Image {
	id -= d.gid
	if id < 0 || id >= len(d.imgs) {
		// return a default image or handle the error
		return nil
	}
	return d.imgs[id]
}

func NewTileset(path string, gid int) (Tileset, error) {

	contents, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if strings.Contains(path, "buildings") {
		var dynamicTilesetJSON DynamicTilesetJSON
		err = json.Unmarshal(contents, &dynamicTilesetJSON)
		if err != nil {
			return nil, err
		}

		dynamicTileset := DynamicTileset{}
		dynamicTileset.gid = gid
		dynamicTileset.imgs = make([]*ebiten.Image, 0)

		for _, tileJSON := range dynamicTilesetJSON.Tiles {

			tileJSONPath := tileJSON.Path
			tileJSONPath = filepath.Clean(tileJSONPath)
			tileJSONPath = strings.ReplaceAll(tileJSONPath, "\\", "/")
			tileJSONPath = strings.TrimPrefix(tileJSONPath, "../")
			tileJSONPath = strings.TrimPrefix(tileJSONPath, "../")
			tileJSONPath = filepath.Join("assets/", tileJSONPath)

			img, _, err := ebitenutil.NewImageFromFile(tileJSONPath)
			if err != nil {
				return nil, err
			}
			dynamicTileset.imgs = append(dynamicTileset.imgs, img)
		}
		return &dynamicTileset, nil
	}

	//return uniform tileset
	var uniformTilesetJSON UniformTilesetJSON
	err = json.Unmarshal(contents, &uniformTilesetJSON)
	if err != nil {
		return nil, err
	}

	uniformTileset := UniformTileset{}

	tileJSONPath := uniformTilesetJSON.Path
	tileJSONPath = filepath.Clean(tileJSONPath)
	tileJSONPath = strings.ReplaceAll(tileJSONPath, "\\", "/")
	tileJSONPath = strings.TrimPrefix(tileJSONPath, "../")
	tileJSONPath = strings.TrimPrefix(tileJSONPath, "../")
	tileJSONPath = filepath.Join("assets/", tileJSONPath)

	img, _, err := ebitenutil.NewImageFromFile(tileJSONPath)
	if err != nil {
		return nil, err
	}
	uniformTileset.img = img
	uniformTileset.gid = gid

	return &uniformTileset, nil

}
