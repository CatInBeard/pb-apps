// Copyright (c) 2025 Grigoriy Efimov
//
// Licensed under the MIT License. See LICENSE file in the project root for details.

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	ink "github.com/CatInBeard/inkview"
)

const repoConfig = "https://raw.githubusercontent.com/CatInBeard/pb-apps/main/repo.json"

type App struct {
	Name        string `json:"name"`
	BinaryName  string `json:"binary-name"`
	URL         string `json:"url"`
	Description string `json:"description"`
	License     string `json:"license"`
}

type Release struct {
	TagName string `json:"tag_name"`
}

func GetRemoteAppList() (map[string]App, error) {
	resp, err := http.Get(repoConfig)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data struct {
		Repositories []App `json:"repositories"`
	}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	apps := make(map[string]App)
	for _, app := range data.Repositories {
		apps[app.Name] = app
	}

	return apps, nil
}

func GetReleases(app App) ([]string, error) {
	urlParts := strings.Split(app.URL, "/")
	repoName := urlParts[len(urlParts)-1]
	owner := urlParts[len(urlParts)-2]

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases", owner, repoName)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var releases []Release
	err = json.Unmarshal(body, &releases)
	if err != nil {
		return nil, err
	}

	var releaseTags []string
	for _, release := range releases {
		releaseTags = append(releaseTags, release.TagName)
	}

	return releaseTags, nil
}

func DownloadAndExtract(app App, releaseName string) error {
	releaseURL := getReleaseURL(app, releaseName)
	fmt.Println("URL: \"" + releaseURL + "\"")

	tmpDir, err := os.MkdirTemp(ink.TempDir, "release-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	releaseZipPath := filepath.Join(tmpDir, "release.zip")
	err = downloadFile(releaseURL, releaseZipPath)
	if err != nil {
		return err
	}

	err = unpackRelease(releaseZipPath, tmpDir)
	if err != nil {
		return err
	}

	downloadsDir := filepath.Join(ink.GamePath)
	err = copyFile(filepath.Join(tmpDir, app.BinaryName), downloadsDir)
	if err != nil {
		return err
	}

	return nil
}

func getReleaseURL(app App, releaseName string) string {
	urlParts := strings.Split(app.URL, "/")
	repoName := urlParts[len(urlParts)-1]
	owner := urlParts[len(urlParts)-2]

	url := fmt.Sprintf("https://github.com/%s/%s/releases/download/%s/release.zip", owner, repoName, releaseName)
	return url

}

func downloadFile(url string, filePath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func unpackRelease(zipPath string, destDir string) error {
	cmd := exec.Command("unzip", zipPath, "-d", destDir)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(string(output))
		return err
	}

	return nil
}

func copyFile(src string, dest string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(filepath.Join(dest, filepath.Base(src)))
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return err
	}

	return nil
}
