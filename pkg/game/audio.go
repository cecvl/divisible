package game

import (
	"io/fs"
	"path/filepath"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func (g *Game) InitAudio() {
	if !rl.IsAudioDeviceReady() {
		return
	}

	musicPath, ok := findMusicFile("assets")
	if !ok {
		return
	}

	g.music = rl.LoadMusicStream(musicPath)
	if rl.IsMusicValid(g.music) {
		g.music.Looping = true
		rl.SetMusicVolume(g.music, 0.35)
		rl.PlayMusicStream(g.music)
	}
}

func (g *Game) CloseAudio() {
	if rl.IsMusicValid(g.music) {
		rl.StopMusicStream(g.music)
		rl.UnloadMusicStream(g.music)
	}
}

func (g *Game) updateAudio() {
	if rl.IsMusicValid(g.music) && g.State != StateFinished {
		rl.UpdateMusicStream(g.music)
	}
}

func findMusicFile(root string) (string, bool) {
	var found string
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d == nil || d.IsDir() || found != "" {
			return err
		}
		if strings.EqualFold(filepath.Ext(path), ".mp3") {
			found = path
		}
		return nil
	})
	if err != nil || found == "" {
		return "", false
	}
	return found, true
}