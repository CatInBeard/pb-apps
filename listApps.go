// Copyright (c) 2025 Grigoriy Efimov
//
// Licensed under the MIT License. See LICENSE file in the project root for details.

package main

import (
	"fmt"
	"image"
	"image/color"
	"math"

	ink "github.com/CatInBeard/inkview"
)

type ListApps struct {
	font                  *ink.Font
	screenSize            image.Point
	fontW                 int
	fontH                 int
	appList               map[string]App
	appKeys               []string
	page                  int
	currentSelectedAppKey string
	shouldRepaint         bool
}

const pageSize = 5

func (a *ListApps) Init() error {
	var err error
	a.appList, err = GetRemoteAppList()

	if err != nil {
		ink.Warningf(GetCurrentTranslation("app_name"), GetCurrentTranslation("get_list_error_message"))
		return err
	}

	a.appKeys = make([]string, 0, len(a.appList))
	for k := range a.appList {
		a.appKeys = append(a.appKeys, k)
	}

	a.page = 1
	a.shouldRepaint = true

	a.Draw()

	return nil
}

func (a *ListApps) Close() error {
	return nil
}

func (a *ListApps) Draw() {
	if !a.shouldRepaint {
		return
	}
	ink.ClearScreen()
	ink.DrawTopPanel()

	appsCount := len(a.appList)
	currentPageAppsCount := int(math.Min(float64(a.page*pageSize), float64(appsCount))) - (a.page-1)*pageSize

	for i := 0; i < currentPageAppsCount; i++ {
		appIndex := (a.page-1)*pageSize + i
		appName := a.appKeys[appIndex]
		app := a.appList[appName]
		var maxCharsOnLine int = a.screenSize.X/a.fontW - 10
		var appText string

		if countRealChar(app.Description) > maxCharsOnLine {
			appText = firstNRunes(app.Description, maxCharsOnLine-5) + "..."
		} else {
			appText = app.Description
		}

		if appName == a.currentSelectedAppKey {
			ink.FillArea(image.Rectangle{
				Min: image.Point{X: 0, Y: a.screenSize.Y / (pageSize + 1) * i},
				Max: image.Point{X: a.screenSize.X, Y: a.screenSize.Y / (pageSize + 1) * (i + 1)},
			}, color.RGBA{R: 173, G: 255, B: 173, A: 255})
		}

		ink.DrawString(image.Point{X: a.fontW * 5, Y: a.screenSize.Y/(pageSize+1)*i + 1}, app.Name)
		ink.DrawLine(image.Point{X: 10, Y: a.screenSize.Y / (pageSize + 1) * (i + 1)}, image.Point{X: a.screenSize.X - 10, Y: a.screenSize.Y / (pageSize + 1) * (i + 1)}, color.Black)
		ink.DrawString(image.Point{X: a.fontW * 5, Y: a.screenSize.Y/(pageSize+1)*i + 1 + int(float64(a.fontH)*1.5)}, appText)

	}

	ink.DrawString(image.Point{a.fontW * 5, a.screenSize.Y - a.fontH*8}, "<")

	ink.DrawString(image.Point{a.screenSize.X - a.fontW*5, a.screenSize.Y - a.fontH*8}, ">")

	pageHint := fmt.Sprintf("%v: %d", GetCurrentTranslation("page_hint"), a.page)

	x := (a.screenSize.X - countRealChar(pageHint)*a.fontW) / 2
	ink.DrawString(image.Point{x, a.screenSize.Y - a.fontH*8}, pageHint)

	ink.PartialUpdate(image.Rectangle{image.Point{0, 0}, a.screenSize})
}

func (a *ListApps) Key(e ink.KeyEvent) bool {
	return true
}

func (a *ListApps) Pointer(e ink.PointerEvent) bool {

	if e.State == ink.PointerDown {
		appsCount := len(a.appList)
		index := a.getItemOnPage(e.Point)
		currentPageAppsCount := int(math.Min(float64(a.page*pageSize), float64(appsCount))) - (a.page-1)*pageSize
		if index < currentPageAppsCount {
			a.currentSelectedAppKey = a.appKeys[a.getAppIndex(index)]
			ink.Repaint()
		}
	} else if e.State == ink.PointerUp {
		app := a.appList[a.currentSelectedAppKey]
		a.shouldRepaint = false
		ink.SetDialogHandler(a.handleInstallDialog)
		ink.Dialog(ink.Question, GetCurrentTranslation("install_question"), app.Name+"\n"+app.Description, GetCurrentTranslation("install_button"), GetCurrentTranslation("cancel_button"))
	}

	return true
}

func (a *ListApps) Touch(e ink.TouchEvent) bool {
	return true
}

func (a *ListApps) Orientation(o ink.Orientation) bool {
	return true
}

func (a *ListApps) getAppIndex(itemNumber int) int {
	return (a.page-1)*pageSize + itemNumber
}

func (a *ListApps) getItemOnPage(point image.Point) int {
	return point.Y / (a.screenSize.Y / (pageSize + 1))
}

func (a *ListApps) handleInstallDialog(button int) {

	if button == 1 {

		ink.ShowHourglass()
		app := a.appList[a.currentSelectedAppKey]
		releases, _ := GetReleases(app)
		err := DownloadAndExtract(app, releases[0])
		ink.HideHourglass()

		if err == nil {
			ink.Infof(GetCurrentTranslation("app_name"), GetCurrentTranslation("success_install"))
		} else {
			ink.Errorf(GetCurrentTranslation("app_name"), GetCurrentTranslation("failed_install"))
		}
	}

	a.shouldRepaint = true
	a.currentSelectedAppKey = ""
	ink.Repaint()
}
