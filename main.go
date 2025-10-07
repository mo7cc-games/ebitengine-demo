package main

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/yohamta/donburi"
)

// PositionData stores the world position of an entity in pixels.
type PositionData struct {
	X, Y float64
}

var Position = donburi.NewComponentType[PositionData]()

type Game struct {
	world     donburi.World
	playerEnt donburi.Entity

	tileImg   *ebiten.Image
	playerImg *ebiten.Image

	screenW int
	screenH int

	// camera position in world coordinates
	cameraX float64
	cameraY float64

	// camera speed in pixels per tick when buttons pressed
	camSpeed float64
}

func NewGame(w, h int) *Game {
	world := donburi.NewWorld()

	// create player entity with Position component
	e := world.Create(Position)
	entry := world.Entry(e)
	Position.SetValue(entry, PositionData{X: 0, Y: 0})

	// 1x1 white image used for tiles and player (scaled)
	tile := ebiten.NewImage(1, 1)
	tile.Fill(color.NRGBA{R: 200, G: 200, B: 200, A: 255})
	player := ebiten.NewImage(1, 1)
	player.Fill(color.NRGBA{R: 50, G: 160, B: 255, A: 255})

	return &Game{
		world:     world,
		playerEnt: e,
		tileImg:   tile,
		playerImg: player,
		screenW:   w,
		screenH:   h,
		cameraX:   0,
		cameraY:   0,
		camSpeed:  8.0,
	}
}

// Update handles input and updates player position.
func (g *Game) Update() error {
	entry := g.world.Entry(g.playerEnt)
	posPtr := Position.Get(entry)
	if posPtr == nil {
		// ensure component exists
		Position.SetValue(entry, PositionData{X: 0, Y: 0})
		posPtr = Position.Get(entry)
	}
	var pos PositionData
	if posPtr == nil {
		pos = PositionData{X: 0, Y: 0}
	} else {
		pos = *posPtr
	}

	speed := 200.0 / 60.0 // pixels per tick (approx 200 px/s)

	// arrow keys or WASD
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
		pos.Y -= speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
		pos.Y += speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		pos.X -= speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		pos.X += speed
	}

	Position.SetValue(entry, pos)

	// exit on ESC
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}

	// mouse-based camera buttons: if left mouse pressed and cursor over a button, move camera
	mx, my := ebiten.CursorPosition()
	leftPressed := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)

	// layout for button cross in left-bottom
	btnSize := 48
	spacing := 8
	margin := 12
	// anchor for left-bottom area
	baseX := margin
	baseY := g.screenH - margin

	// compute centers
	centerX := baseX + btnSize + spacing
	downY := baseY - btnSize
	midY := baseY - btnSize*2 - spacing
	upY := baseY - btnSize*3 - spacing*2

	// rectangles for buttons
	upRect := imageRect(centerX, upY, btnSize, btnSize)
	downRect := imageRect(centerX, downY, btnSize, btnSize)
	leftRect := imageRect(baseX, midY, btnSize, btnSize)
	rightRect := imageRect(baseX+(btnSize+spacing)*2, midY, btnSize, btnSize)

	if leftPressed {
		if pointInRect(mx, my, upRect) {
			g.cameraY -= g.camSpeed
		}
		if pointInRect(mx, my, downRect) {
			g.cameraY += g.camSpeed
		}
		if pointInRect(mx, my, leftRect) {
			g.cameraX -= g.camSpeed
		}
		if pointInRect(mx, my, rightRect) {
			g.cameraX += g.camSpeed
		}
	}

	return nil
}

// Draw renders a simple infinite grid and the player at the center of the screen.
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Clear()

	entry := g.world.Entry(g.playerEnt)
	posPtr := Position.Get(entry)
	var pos PositionData
	if posPtr == nil {
		pos = PositionData{X: 0, Y: 0}
	} else {
		pos = *posPtr
	}

	tileSize := 48
	halfW := g.screenW / 2
	halfH := g.screenH / 2

	// compute visible tile range (a small range for performance)
	viewTilesX := int(math.Ceil(float64(g.screenW)/float64(tileSize))) + 4
	viewTilesY := int(math.Ceil(float64(g.screenH)/float64(tileSize))) + 4

	// center the world on the camera (cameraX, cameraY)
	for dx := -viewTilesX; dx <= viewTilesX; dx++ {
		for dy := -viewTilesY; dy <= viewTilesY; dy++ {
			// world tile coordinate
			tileX := math.Floor((g.cameraX)/float64(tileSize)) + float64(dx)
			tileY := math.Floor((g.cameraY)/float64(tileSize)) + float64(dy)

			// convert tile to screen position
			screenX := (tileX*float64(tileSize) - g.cameraX) + float64(halfW)
			screenY := (tileY*float64(tileSize) - g.cameraY) + float64(halfH)

			op := &ebiten.DrawImageOptions{}
			op.GeoM.Scale(float64(tileSize), float64(tileSize))
			op.GeoM.Translate(screenX, screenY)

			// alternate colors slightly for a grid look
			if int(tileX+tileY)%2 == 0 {
				g.tileImg.Fill(color.NRGBA{R: 230, G: 230, B: 230, A: 255})
			} else {
				g.tileImg.Fill(color.NRGBA{R: 210, G: 210, B: 210, A: 255})
			}

			screen.DrawImage(g.tileImg, op)
		}
	}

	// draw player at its world position relative to camera
	playerSize := 28.0
	pop := &ebiten.DrawImageOptions{}
	pop.GeoM.Scale(playerSize, playerSize)
	// compute screen pos from world pos
	playerScreenX := (pos.X - g.cameraX) + float64(halfW)
	playerScreenY := (pos.Y - g.cameraY) + float64(halfH)
	pop.GeoM.Translate(playerScreenX-playerSize/2, playerScreenY-playerSize/2)
	screen.DrawImage(g.playerImg, pop)

	// debug text: world position
	ebitenutil.DebugPrint(screen, fmt.Sprintf("Player: %s, %s  Camera: %s, %s  (Arrows/WASD to move player, Esc to exit)", formatFloat(pos.X), formatFloat(pos.Y), formatFloat(g.cameraX), formatFloat(g.cameraY)))

	// draw camera control buttons in left-bottom
	btnSize := 48
	spacing := 8
	margin := 12
	baseX := margin
	baseY := g.screenH - margin
	centerX := baseX + btnSize + spacing
	downY := baseY - btnSize
	midY := baseY - btnSize*2 - spacing
	upY := baseY - btnSize*3 - spacing*2

	upRect := imageRect(centerX, upY, btnSize, btnSize)
	downRect := imageRect(centerX, downY, btnSize, btnSize)
	leftRect := imageRect(baseX, midY, btnSize, btnSize)
	rightRect := imageRect(baseX+(btnSize+spacing)*2, midY, btnSize, btnSize)

	// detect hover/press for color
	mx, my := ebiten.CursorPosition()
	leftPressed := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)

	drawButton := func(r imageRectT, label string) {
		rectOp := &ebiten.DrawImageOptions{}
		// solid image reuse: fill tileImg with button color then scale
		col := color.NRGBA{R: 100, G: 100, B: 100, A: 200}
		if pointInRect(mx, my, r) && leftPressed {
			col = color.NRGBA{R: 60, G: 160, B: 60, A: 220}
		} else if pointInRect(mx, my, r) {
			col = color.NRGBA{R: 80, G: 140, B: 80, A: 220}
		}
		g.tileImg.Fill(col)
		rectOp.GeoM.Scale(float64(r.w), float64(r.h))
		rectOp.GeoM.Translate(float64(r.x), float64(r.y))
		screen.DrawImage(g.tileImg, rectOp)
		ebitenutil.DebugPrintAt(screen, label, r.x+6, r.y+6)
	}

	drawButton(upRect, "↑")
	drawButton(downRect, "↓")
	drawButton(leftRect, "←")
	drawButton(rightRect, "→")
}

// simple rect type for ui
type imageRectT struct{ x, y, w, h int }

func imageRect(cx, cy, w, h int) imageRectT {
	return imageRectT{x: cx, y: cy, w: w, h: h}
}

func pointInRect(px, py int, r imageRectT) bool {
	return px >= r.x && px <= r.x+r.w && py >= r.y && py <= r.y+r.h
}

func formatFloat(v float64) string {
	// keep 1 decimal
	return strconv.FormatFloat(v, 'f', 1, 64)
}

// Layout returns logical screen size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.screenW, g.screenH
}

func main() {
	// window size
	w, h := 800, 600
	ebiten.SetWindowSize(w, h)
	ebiten.SetWindowTitle("Ebitengine + Donburi Infinite Map Demo")

	game := NewGame(w, h)

	if err := ebiten.RunGame(game); err != nil && err != ebiten.Termination {
		log.Fatal(err)
	}
}
