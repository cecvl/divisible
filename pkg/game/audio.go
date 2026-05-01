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

	files := findMusicFiles("assets")
	if len(files) == 0 {
		return
	}
	g.musicFiles = files
	g.musicIndex = 0
	g.loadMusicAt(g.musicIndex)
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
		// advance to next track when current finishes
		// Use a small epsilon to avoid float issues
		if rl.GetMusicTimePlayed(g.music) >= rl.GetMusicTimeLength(g.music)-0.05 {
			g.nextMusic()
		}
	}
}

func findMusicFile(root string) (string, bool) {
	// deprecated: kept for compatibility
	files := findMusicFiles(root)
	if len(files) == 0 {
		return "", false
	}
	return files[0], true
}

func findMusicFiles(root string) []string {
	var files []string
	_ = filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d == nil || d.IsDir() {
			return nil
		}
		ext := strings.ToLower(filepath.Ext(path))
		if ext == ".mp3" || ext == ".ogg" || ext == ".wav" {
			files = append(files, path)
		}
		return nil
	})
	return files
}

func (g *Game) loadMusicAt(index int) {
	if index < 0 || index >= len(g.musicFiles) {
		return
	}
	// unload previous
	if rl.IsMusicValid(g.music) {
		rl.StopMusicStream(g.music)
		rl.UnloadMusicStream(g.music)
	}
	path := g.musicFiles[index]
	g.music = rl.LoadMusicStream(path)
	if rl.IsMusicValid(g.music) {
		// don't loop single file; we'll advance to next when finished
		g.music.Looping = false
		rl.SetMusicVolume(g.music, 0.35)
		rl.PlayMusicStream(g.music)
	}
}

func (g *Game) nextMusic() {
	if len(g.musicFiles) == 0 {
		return
	}
	g.musicIndex = (g.musicIndex + 1) % len(g.musicFiles)
	g.loadMusicAt(g.musicIndex)
}
