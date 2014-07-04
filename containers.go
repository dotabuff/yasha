package core

import "github.com/davecgh/go-spew/spew"

var pp = spew.Dump

type AbilityTracker struct {
	HeroHandle int
	Level      int
	Tick       int
	Name       string
}

type Abilities []*AbilityTracker

func (p Abilities) Len() int           { return len(p) }
func (p Abilities) Less(i, j int) bool { return p[i].Tick < p[j].Tick }
func (p Abilities) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

type LastHitTracker struct {
	HeroHandle int
	Tick       int
	LastHit    int
}

type LastHits []*LastHitTracker

func (p LastHits) Len() int           { return len(p) }
func (p LastHits) Less(i, j int) bool { return p[i].Tick < p[j].Tick }
func (p LastHits) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
