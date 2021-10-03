package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/valyala/fastrand"
	"image"
	"log"
	"strconv"
)

const (
	screenWidth         = 1000
	screenHeight        = 1000
	balancingIterations = 10
	pixelSize           = 10
)

type Game struct {
	i          int
	pixSlice   [screenHeight / pixelSize][screenWidth / pixelSize]uint32
	pixSlice2  [screenHeight / pixelSize][screenWidth / pixelSize]uint32
	noiseImage *image.RGBA
}

func (g *Game) init() error {
	return nil
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		//fill matrix with random numbers
		for i1 := 0; i1 < screenHeight/pixelSize; i1++ {
			for i2 := 0; i2 < screenWidth/pixelSize; i2++ {
				x := uint32(uint8(fastrand.Uint32()))
				g.pixSlice2[i1][i2] = x
			}
		}
		//balancing {balancingIterations} times
		for j := 0; j < balancingIterations; j++ {
			for i1 := 0; i1 < screenHeight/pixelSize; i1++ {
				for i2 := 0; i2 < screenWidth/pixelSize; i2++ {
					g.pixSlice[i1][i2] = g.pixSlice2[i1][i2]
				}
			}
			for i1 := 0; i1 < screenHeight/pixelSize; i1++ {
				for i2 := 0; i2 < screenWidth/pixelSize; i2++ {
					if i1 == 0 {
						if i2 == 0 {
							g.pixSlice2[i1][i2] = (g.pixSlice[i1][i2] + g.pixSlice[i1][i2+1] + g.pixSlice[i1+1][i2] + g.pixSlice[i1+1][i2+1]) / 4
						} else if i2 == screenWidth/pixelSize-1 {
							g.pixSlice2[i1][i2] = (g.pixSlice[i1][i2-1] + g.pixSlice[i1][i2] + g.pixSlice[i1+1][i2-1] + g.pixSlice[i1+1][i2]) / 4
						} else {
							g.pixSlice2[i1][i2] = (g.pixSlice[i1][i2-1] + g.pixSlice[i1][i2] + g.pixSlice[i1][i2+1] + g.pixSlice[i1+1][i2-1] + g.pixSlice[i1+1][i2] + g.pixSlice[i1+1][i2+1]) / 6
						}
					} else if i1 == screenHeight/pixelSize-1 {
						if i2 == 0 {
							g.pixSlice2[i1][i2] = (g.pixSlice[i1-1][i2] + g.pixSlice[i1-1][i2+1] + g.pixSlice[i1][i2] + g.pixSlice[i1][i2+1]) / 4
						} else if i2 == screenWidth/pixelSize-1 {
							g.pixSlice2[i1][i2] = (g.pixSlice[i1-1][i2-1] + g.pixSlice[i1-1][i2] + g.pixSlice[i1][i2-1] + g.pixSlice[i1][i2]) / 4
						} else {
							g.pixSlice2[i1][i2] = (g.pixSlice[i1-1][i2-1] + g.pixSlice[i1-1][i2] + g.pixSlice[i1-1][i2+1] + g.pixSlice[i1][i2-1] + g.pixSlice[i1][i2] + g.pixSlice[i1][i2+1]) / 6
						}

					} else if i2 == 0 {
						g.pixSlice2[i1][i2] = (g.pixSlice[i1-1][i2] + g.pixSlice[i1-1][i2+1] + g.pixSlice[i1][i2] + g.pixSlice[i1][i2+1] + g.pixSlice[i1+1][i2] + g.pixSlice[i1+1][i2+1]) / 6
					} else if i2 == screenWidth/pixelSize-1 {
						g.pixSlice2[i1][i2] = (g.pixSlice[i1-1][i2-1] + g.pixSlice[i1-1][i2] + g.pixSlice[i1][i2-1] + g.pixSlice[i1][i2] + g.pixSlice[i1+1][i2-1] + g.pixSlice[i1+1][i2]) / 6
					} else {
						g.pixSlice2[i1][i2] = (g.pixSlice[i1-1][i2-1] + g.pixSlice[i1-1][i2] + g.pixSlice[i1-1][i2+1] + g.pixSlice[i1][i2-1] + g.pixSlice[i1][i2] + g.pixSlice[i1][i2+1] + g.pixSlice[i1+1][i2-1] + g.pixSlice[i1+1][i2] + g.pixSlice[i1+1][i2+1]) / 9
					}
				}
			}
		}
		//filling screen pixels
		i := 0
		var rc, gc, bc uint8
		for i1 := 0; i1 < screenHeight/pixelSize; i1++ {
			for i2 := 0; i2 < screenWidth/pixelSize; i2++ {
				n := uint8(g.pixSlice2[i1][i2])
				if i1 == 0 || i2 == 0 || i1 == screenHeight/pixelSize-1 || i2 == screenWidth/pixelSize-1 { //edge
					rc, gc, bc = 255, 0, 0
				} else if n > 137 { //mountains
					rc, gc, bc = 20, 15, 11
				} else if n > 132 { //hills
					rc, gc, bc = 159, 130, 0
				} else if n > 127 { //upper meadows
					rc, gc, bc = 50, 153, 0
				} else if n > 122 { //lower meadows
					rc, gc, bc = 86, 222, 71
				} else if n > 117 { //shallow water
					rc, gc, bc = 0, 127, 255
				} else { //deep water
					rc, gc, bc = 51, 51, 255
				}
				for i3 := 0; i3 < pixelSize; i3++ {
					for i4 := 0; i4 < pixelSize; i4++ {
						g.noiseImage.Pix[4*((i1*pixelSize+i3)*screenHeight+(i2*pixelSize+i4))] = rc
						g.noiseImage.Pix[4*((i1*pixelSize+i3)*screenHeight+(i2*pixelSize+i4))+1] = gc
						g.noiseImage.Pix[4*((i1*pixelSize+i3)*screenHeight+(i2*pixelSize+i4))+2] = bc
						g.noiseImage.Pix[4*((i1*pixelSize+i3)*screenHeight+(i2*pixelSize+i4))+3] = 0xff
						i++
					}
				}
			}
		}
		g.i++
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.ReplacePixels(g.noiseImage.Pix)
	txt := "Iteration " + strconv.Itoa(g.i)
	msg := fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f\n%s\nPress Space for a new random generated map", ebiten.CurrentTPS(), ebiten.CurrentFPS(), txt)
	ebitenutil.DebugPrint(screen, msg)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("RandomMaps")
	g := &Game{
		noiseImage: image.NewRGBA(image.Rect(0, 0, screenWidth, screenHeight)),
	}
	err := g.init()
	if err != nil {
		log.Fatal(err)
	}
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
