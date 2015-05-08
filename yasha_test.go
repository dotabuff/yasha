package yasha

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/dotabuff/yasha/dota"
	"github.com/stretchr/testify/assert"
)

func TestRead(t *testing.T) {
	parser := ParserFromFile("1451081871_958173873.dem")
	parser.Parse()
	parser.OnSayText2 = func(n int, o *dota.CUserMsg_SayText2) { spew.Dump(o) }
	parser.OnChatEvent = func(n int, o *dota.CDOTAUserMsg_ChatEvent) { spew.Dump(o) }
	assert.EqualValues(t, 1, 4)
}
