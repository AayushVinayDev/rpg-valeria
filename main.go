package main

import (
	"Valeria/entities"
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func CheckCollisionHorizontal(sprite *entities.Sprite, colliders []image.Rectangle) {
	for _, collider := range colliders {
		if collider.Overlaps(
			image.Rect(
				int(sprite.X),
				int(sprite.Y),
				int(sprite.X)+16.0,
				int(sprite.Y)+16.0,
			),
		) {
			if sprite.Dx > 0.0 {
				sprite.X = float64(collider.Min.X) - 16.0
			} else if sprite.Dx < 0.0 {
				sprite.X = float64(collider.Max.X)
			}
		}
	}
}

func CheckCollisionVertical(sprite *entities.Sprite, colliders []image.Rectangle) {
	for _, collider := range colliders {
		if collider.Overlaps(
			image.Rect(
				int(sprite.X),
				int(sprite.Y),
				int(sprite.X)+16.0,
				int(sprite.Y)+16.0,
			),
		) {
			if sprite.Dy > 0.0 {
				sprite.Y = float64(collider.Min.Y) - 16.0
			} else if sprite.Dy < 0.0 {
				sprite.Y = float64(collider.Max.X)
			}
		}
	}
}

type Game struct {
	//image and position variables for the player
	player      *entities.Player
	enemies     []*entities.Enemy
	meat        []*entities.Meat
	tilemapJSON *TilemapJSON
	tilesets    []Tileset
	tilemapImg  *ebiten.Image
	cam         *Camera
	colliders   []image.Rectangle
}

func (g *Game) Update() error {
	//react to key presses
	g.player.Dx = 0.0
	g.player.Dy = 0.0

	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.player.Dx = -2
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.player.Dx = 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		g.player.Dy = -2
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		g.player.Dy = 2
	}

	g.player.X += g.player.Dx

	CheckCollisionHorizontal(g.player.Sprite, g.colliders)

	g.player.Y += g.player.Dy

	CheckCollisionVertical(g.player.Sprite, g.colliders)

	//enemies follow the player
	for _, sprite := range g.enemies {

		sprite.Dx = 0.0
		sprite.Dy = 0.0

		if sprite.FollowsPlayer {
			if sprite.X < g.player.X {
				sprite.Dx += 1
			} else if sprite.X > g.player.X {
				sprite.Dx -= 1
			}
			if sprite.Y < g.player.Y {
				sprite.Dy += 1
			} else if sprite.Y > g.player.Y {
				sprite.Dy -= 1
			}
		}

		sprite.X += sprite.Dx

		CheckCollisionHorizontal(sprite.Sprite, g.colliders)

		sprite.Y += sprite.Dy

		CheckCollisionVertical(sprite.Sprite, g.colliders)
	}

	for _, heal := range g.meat {
		if g.player.X > heal.X {
			g.player.Health += heal.AmtHeal
			//fmt.Printf("DOGGO ATE MEAT!! Health: %d\n", g.player.Health)
		}
	}

	g.cam.FollowTarget(g.player.X+8, g.player.Y+8, 320, 240)
	g.cam.Constraint(
		float64(g.tilemapJSON.Layers[0].Width)*16.0,
		float64(g.tilemapJSON.Layers[0].Height)*16.0,
		320,
		240,
	)

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	//fill the screen with a nice sky color
	screen.Fill(color.RGBA{120, 180, 255, 255})

	opts := ebiten.DrawImageOptions{}

	//loop over the layers
	for layerIndex, layer := range g.tilemapJSON.Layers {
		//loop over tiles in the layer data
		for index, id := range layer.Data {

			if id <= 0 {
				continue
			}
			// if id < 0 {
			// 	log.Println("invalid tile id: ", id)
			// 	continue
			// }

			//get the tile position
			x := index % layer.Width
			y := index / layer.Width

			//convert the tile position to pixel position
			x *= 16
			y *= 16

			img := g.tilesets[layerIndex].Img(id)

			opts.GeoM.Translate(float64(x), float64(y))

			opts.GeoM.Translate(0.0, -(float64(img.Bounds().Dy()) + 16))

			opts.GeoM.Translate(g.cam.X, g.cam.Y)

			screen.DrawImage(img, &opts)

			// //get the position on the image where the tile id is
			// srcX := (id - 1) % 22
			// srcY := (id - 1) / 22

			// //convert the src tile position to pixel src position
			// srcX *= 16
			// srcY *= 16

			// //set the drawImageOptions to draw the tile at x,y

			// //draw the tile
			// screen.DrawImage(
			// 	//cropping out the tile that we want from the spritesheet
			// 	g.tilemapImg.SubImage(image.Rect(srcX, srcY, srcX+16, srcY+16)).(*ebiten.Image),
			// 	&opts,
			// )

			//reset the opts for next tile
			opts.GeoM.Reset()

		}
	}

	//set the transalation of our drawImageOptions to the player's position
	opts.GeoM.Translate(g.player.X, g.player.Y)
	opts.GeoM.Translate(g.cam.X, g.cam.Y)

	//draw our player
	screen.DrawImage(
		g.player.Img.SubImage(
			image.Rect(0, 0, 16, 16),
		).(*ebiten.Image),
		&opts,
	)

	opts.GeoM.Reset()

	for _, sprite := range g.enemies {
		opts.GeoM.Translate(sprite.X, sprite.Y)
		opts.GeoM.Translate(g.cam.X, g.cam.Y)

		screen.DrawImage(
			sprite.Img.SubImage(
				image.Rect(0, 0, 16, 16),
			).(*ebiten.Image),
			&opts,
		)

		opts.GeoM.Reset()

	}
	for _, sprite := range g.meat {
		opts.GeoM.Translate(sprite.X, sprite.Y)
		opts.GeoM.Translate(g.cam.X, g.cam.Y)

		screen.DrawImage(
			sprite.Img.SubImage(
				image.Rect(0, 0, 16, 16),
			).(*ebiten.Image),
			&opts,
		)

		opts.GeoM.Reset()

	}

	for _, collider := range g.colliders {
		vector.StrokeRect(
			screen,
			float32(collider.Min.X)+float32(g.cam.X),
			float32(collider.Min.Y)+float32(g.cam.Y),
			float32(collider.Dx()),
			float32(collider.Dy()),
			1.0,
			color.RGBA{255, 0, 0, 255},
			true,
		)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Valeria")

	playerImg, _, err := ebitenutil.NewImageFromFile("assets/image/ninja.png")
	if err != nil {
		log.Fatal(err)
	}
	skeletonImg, _, err := ebitenutil.NewImageFromFile("assets/image/skeleton.png")
	if err != nil {
		log.Fatal(err)
	}
	meatImg, _, err := ebitenutil.NewImageFromFile("assets/image/Beaf.png")
	if err != nil {
		log.Fatal(err)
	}
	tilemapImg, _, err := ebitenutil.NewImageFromFile("assets/image/TilesetFloor.png")
	if err != nil {
		log.Fatal(err)
	}

	tilemapJSON, err := NewTilemapJSON("assets/maps/spawn.json")
	if err != nil {
		log.Fatal(err)
	}

	tilesets, err := tilemapJSON.GenerateTileset()
	if err != nil {
		log.Fatal(err)
	}

	game := Game{
		player: &entities.Player{
			Sprite: &entities.Sprite{
				Img: playerImg,
				X:   50.0,
				Y:   50.0,
			},
			Health: 50,
		},
		enemies: []*entities.Enemy{
			{
				Sprite: &entities.Sprite{
					Img: skeletonImg,
					X:   100.0,
					Y:   100.0,
				},
				FollowsPlayer: true,
			},
			{
				Sprite: &entities.Sprite{
					Img: skeletonImg,
					X:   10.0,
					Y:   180.0,
				},
				FollowsPlayer: false,
			},
			{
				Sprite: &entities.Sprite{
					Img: skeletonImg,
					X:   75.0,
					Y:   10.0,
				},
				FollowsPlayer: false,
			},
		},
		meat: []*entities.Meat{
			{
				Sprite: &entities.Sprite{
					Img: meatImg,
					X:   205.0,
					Y:   100.0,
				},
				AmtHeal: 5.0,
			},
		},
		tilemapJSON: tilemapJSON,
		tilemapImg:  tilemapImg,
		tilesets:    tilesets,
		cam:         NewCamera(0.0, 0.0),
		colliders: []image.Rectangle{
			image.Rect(100, 100, 116, 116),
		},
	}

	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
