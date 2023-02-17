package main

import (
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

type HUD struct {
	ScoreDisplayFont  font.Face
	ResultDisplayFont font.Face
}

func newHUD() (*HUD, error) {
	tt, err := opentype.Parse(fonts.PressStart2P_ttf)
	if err != nil {
		return nil, err
	}

	const dpi = 72
	scoreDisplayFont, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    76,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return nil, err
	}

	resultDisplayFont, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    18,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return nil, err
	}

	return &HUD{
		ScoreDisplayFont:  scoreDisplayFont,
		ResultDisplayFont: resultDisplayFont,
	}, nil
}
