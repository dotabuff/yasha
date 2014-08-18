package main

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/dotabuff/yasha"
	"github.com/dotabuff/yasha/dota"
)

func main() {
	parser := yasha.ParserFromFile("/home/manveru/Dropbox/Public/d2rp_test_replays/813230338.dem")
	parser.OnFileInfo = func(fileinfo *dota.CDemoFileInfo) { spew.Dump(fileinfo) }
	parser.OnSayText2 = func(tick int, obj *dota.CUserMsg_SayText2) { spew.Dump(tick, obj) }
	parser.OnChatEvent = func(tick int, obj *dota.CDOTAUserMsg_ChatEvent) { spew.Dump(tick, obj) }
	parser.Parse()
}
