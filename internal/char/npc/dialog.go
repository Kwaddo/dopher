package npc

import (
    DM "doom/internal/model"
    "github.com/veandco/go-sdl2/sdl"
    "github.com/veandco/go-sdl2/ttf"
)

type DialogRenderer DM.DialogRenderer

func NewDialogRenderer() (*DialogRenderer, error) {
    if err := ttf.Init(); err != nil {
        return nil, err
    }

    font, err := ttf.OpenFont("assets/dogicapixel.ttf", 24)
    if err != nil {
        ttf.Quit()
        return nil, err
    }

    return &DialogRenderer{
        Font:   font,
        Loaded: true,
    }, nil
}

func (dr *DialogRenderer) Close() {
    if dr.Loaded {
        dr.Font.Close()
        ttf.Quit()
        dr.Loaded = false
    }
}

func (dr *DialogRenderer) RenderDialog(renderer *sdl.Renderer, text string) error {
    surface, err := dr.Font.RenderUTF8Solid(
        text,
        sdl.Color{R: 255, G: 255, B: 255, A: 255},
    )
    if err != nil {
        return err
    }
    defer surface.Free()

    texture, err := renderer.CreateTextureFromSurface(surface)
    if err != nil {
        return err
    }
    defer texture.Destroy()

    textW := surface.W
    textH := surface.H

    return renderer.Copy(texture, nil, &sdl.Rect{
        X: int32(DM.ScreenWidth/2 - textW/2),
        Y: int32(DM.ScreenHeight - textH - 20),
        W: int32(textW),
        H: int32(textH),
    })
}
