package main

import (
	"encoding/json"
	"os"
)

type TilemapLayerJSON struct {
	Data   []int `json:"data"`
	Width  int   `json:"width"`
	Height int   `json:"height"`
}

type TilemapJSON struct {
	Layers []TilemapLayerJSON `json:"layers"`
}

func NewTilemapJSON(filepath string) (*TilemapJSON, error) {

	//Read the contents
	contents, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	//Unmarshal the contents
	var tilemapJSON TilemapJSON
	err = json.Unmarshal(contents, &tilemapJSON)
	if err != nil {
		return nil, err
	}

	return &tilemapJSON, nil
}
