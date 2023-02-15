package main

import (
	"bytes"
	_ "embed"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
)

//go:embed wall.ogg
var wallOgg []byte

//go:embed bounce.ogg
var bounceOgg []byte

//go:embed score.ogg
var scoreOgg []byte

type Sound struct {
	player *audio.Player
}

func (s *Sound) Play() error {
	if !s.player.IsPlaying() {
		err := s.player.Rewind()
		if err != nil {
			return err
		}
		s.player.Play()
	}
	return nil
}

func LoadSounds() (map[string]*Sound, error) {
	sounds := map[string]*Sound{}
	for name, file := range soundFiles {
		reader := bytes.NewReader(file)
		decoded, err := vorbis.DecodeWithSampleRate(44100, reader)
		if err != nil {
			return nil, err
		}
		// player, err := audio.NewPlayer(audioContext, decoded)
		player, err := audioContext.NewPlayer(decoded)
		if err != nil {
			return nil, err
		}
		sounds[name] = &Sound{player: player}
	}
	return sounds, nil
}

var soundFiles = map[string][]byte{
	"wall":   wallOgg,
	"bounce": bounceOgg,
	"score":  scoreOgg,
}

// Initialize audio context
var audioContext = audio.NewContext(44100)
