package utils

import (
	"bytes"
	"fmt"
	"image"
	"net/http"
	"github.com/fogleman/gg"
)

type Color struct {
	R int
	G int
	B int
}

func NewColor(r, g, b int) *Color {
	return &Color {
		R: r,
		G: g,
		B: b,
	}
}

func MemberCard(color *Color, level, exp, expToLevelUp int, userID, displayName, avatar string) (*bytes.Buffer, error) {
	im, err := gg.LoadJPG("./assets/images/banner.jpg")

	if err != nil {
		return nil, err
	}

	width, height := im.Bounds().Dx(), im.Bounds().Dy()
	dc := gg.NewContext(width, height)

	dc.DrawImage(im, 0, 0)

	dc.SetRGBA255(33, 8, 24, 220)
	dc.DrawRectangle(0, 0, float64(width), float64(height))
	dc.Fill()


	dc.SetRGB255(226, 115, 204)

	x, y := 100.0, 110.0
	radius := 60.0
	dc.DrawCircle(x, y, radius)
	dc.Clip()

	url := fmt.Sprintf("https://cdn.discordapp.com/avatars/%s/%s.png", userID, avatar)

	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	remoteImage, _, err := image.Decode(resp.Body)

	if err != nil {
		return nil, err
	}

	dc.DrawImageAnchored(remoteImage, int(x), int(y), 0.5, 0.5)
	dc.ResetClip()

	if err := dc.LoadFontFace("./assets/fonts/robmax.ttf", 28); err != nil {
		return nil, err
	}

	dc.DrawStringAnchored(displayName, 275, 70, 0.5, 0.5)

	widthBar := 500

	dc.SetRGB255(84, 61, 72)
	dc.DrawRoundedRectangle(200, 120, float64(widthBar), 20, 10)
	dc.Fill()

	expWidth := float64(exp) *float64(widthBar) / float64(expToLevelUp)

	dc.SetRGB255(color.R, color.G, color.B)
	dc.SetRGB255(242, 140, 220)
	dc.DrawRoundedRectangle(200, 120, expWidth, 20, 10)
	dc.Fill()

	levelText := fmt.Sprintf("Lvl. %d", level)
	widthText, _ := dc.MeasureString(levelText)

	if err := dc.LoadFontFace("./assets/fonts/zen.ttf", 24); err != nil {
		return nil, err
	}

	dc.SetRGB255(242, 140, 220)
	dc.DrawString(levelText, 200 + 500 - widthText, 78)

	if err := dc.LoadFontFace("./assets/fonts/zen.ttf", 10); err != nil {
		return nil, err
	}

	expText := fmt.Sprintf("Exp. %d/%d", exp, expToLevelUp)
	widthText, _ = dc.MeasureString(expText)

	dc.SetRGBA255(255, 255, 255, 160)
	dc.DrawString(expText, 198+500-widthText, 165)

	buf := new(bytes.Buffer)
	if err := dc.EncodePNG(buf);err != nil {
		return nil, err
	}

	return buf, nil
}
