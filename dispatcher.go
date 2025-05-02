// Copyright (c) 2025 Grigoriy Efimov
//
// Licensed under the MIT License. See LICENSE file in the project root for details.

package main

import (
	"image"
	"image/color"
	"time"

	ink "github.com/CatInBeard/inkview"
)

type DispatcherApp struct {
	font       *ink.Font
	fontH      int
	fontW      int
	app        ink.App
	screenSize image.Point
}

func (a *DispatcherApp) Init() error {
	ink.ClearScreen()
	ink.DrawTopPanel()
	ink.ShowHourglass()

	CreateInitDb()

	a.font = ink.OpenFont(ink.DefaultFontMono, a.fontH, true)
	a.font.SetActive(color.RGBA{0, 0, 0, 255})
	a.fontW = ink.CharWidth('a') // Work only for monospace font

	a.screenSize = ink.ScreenSize()

	ink.SetMessageDelay(time.Second * 10)

	ink.Infof(GetCurrentTranslation("app_name"), GetCurrentTranslation("welcome_message"))

	ink.QueryNetwork()
	err := ink.ConnectDefault()
	if err != nil {
		ink.Exit()
	}

	ink.HideHourglass()
	a.app = &ListApps{font: a.font, fontW: a.fontW, fontH: a.fontH, screenSize: a.screenSize}
	a.app.Init()

	return nil
}

func (a *DispatcherApp) Close() error {
	if a.app != nil {
		a.app.Close()
	}
	return nil
}

func (a *DispatcherApp) Draw() {
	if a.app != nil {
		a.app.Draw()
	}
}

func (a *DispatcherApp) Key(e ink.KeyEvent) bool {
	if a.app != nil {
		return a.app.Key(e)
	}
	return true
}

func (a *DispatcherApp) Pointer(e ink.PointerEvent) bool {

	if a.app != nil {
		return a.app.Pointer(e)
	}

	return true
}

func (a *DispatcherApp) Touch(e ink.TouchEvent) bool {
	if a.app != nil {
		return a.app.Touch(e)
	}
	return true
}

func (a *DispatcherApp) Orientation(o ink.Orientation) bool {
	if a.app != nil {
		return a.app.Orientation(o)
	}
	return true
}
