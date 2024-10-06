package main

import (
	"Valeria/entities"
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct {
	//image and position variables for the player
	player      *entities.Player
	enemies     []*entities.Enemy
	meat        []*entities.Meat
	tilemapJSON *TilemapJSON
	tilemapImg  *ebiten.Image
	cam         *Camera
}

func (g *Game) Update() error {
	//react to key presses
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.player.X += 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.player.X -= 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		g.player.Y -= 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		g.player.Y += 2
	}

	//enemies follow the player
	for _, sprite := range g.enemies {
		if sprite.FollowsPlayer {
			if sprite.X < g.player.X {
				sprite.X += 0.5
			} else if sprite.X > g.player.X {
				sprite.X -= 0.5
			}
			if sprite.Y < g.player.Y {
				sprite.Y += 0.5
			} else if sprite.Y > g.player.Y {
				sprite.Y -= 0.5
			}
		}

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
	for _, layer := range g.tilemapJSON.Layers {
		//loop over tiles in the layer data
		for index, id := range layer.Data {
			//get the tile position
			x := index % layer.Width
			y := index / layer.Width

			//convert the tile position to pixel position
			x *= 16
			y *= 16

			//get the position on the image where the tile id is
			srcX := (id - 1) % 22
			srcY := (id - 1) / 22

			//convert the src tile position to pixel src position
			srcX *= 16
			srcY *= 16

			//set the drawImageOptions to draw the tile at x,y
			opts.GeoM.Translate(float64(x), float64(y))

			opts.GeoM.Translate(g.cam.X, g.cam.Y)

			//draw the tile
			screen.DrawImage(
				//cropping out the tile that we want from the spritesheet
				g.tilemapImg.SubImage(image.Rect(srcX, srcY, srcX+16, srcY+16)).(*ebiten.Image),
				&opts,
			)

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
		cam:         NewCamera(0.0, 0.0),
	}

	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
