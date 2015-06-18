package yasha

import (
	"compress/bzip2"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/dotabuff/yasha/dota"
	"github.com/stretchr/testify/assert"
)

func init() {
	debugMode = true
}

type testCase struct {
	matchId int64
	url     string

	expectLastChatMessage string
	expectHeroKillCount   map[string]int
	expectHeroDeathCount  map[string]int
}

// Esports match, played on patch 6.83c
func TestEsportsPatch683b(t *testing.T) {
	c := &testCase{
		matchId: 1405240741,
		url:     "https://s3-us-west-2.amazonaws.com/yasha.dotabuff/1405240741.dem",
		expectLastChatMessage: "Gg",
		expectHeroKillCount: map[string]int{
			"npc_dota_hero_ember_spirit": 0,
			"npc_dota_hero_broodmother":  5,
		},
		expectHeroDeathCount: map[string]int{
			"npc_dota_hero_chen":        2,
			"npc_dota_hero_broodmother": 0,
			"npc_dota_hero_sniper":      2,
			"npc_dota_hero_phoenix":     2,
		},
	}

	testReplayCase(t, c)
}

// Esports match, played on patch 6.84p0
func TestEsportsPatch684p0(t *testing.T) {
	c := &testCase{
		matchId: 1450235906,
		url:     "https://s3-us-west-2.amazonaws.com/yasha.dotabuff/1450235906.dem",
		expectLastChatMessage: "gg",
		expectHeroKillCount: map[string]int{
			"npc_dota_hero_broodmother": 3,
		},
		expectHeroDeathCount: map[string]int{
			"npc_dota_hero_broodmother": 7,
		},
	}

	testReplayCase(t, c)
}

// Esports match, played on patch 6.84p1
func TestEsportsPatch684p1(t *testing.T) {
	c := &testCase{
		matchId: 1458895412,
		url:     "https://s3-us-west-2.amazonaws.com/yasha.dotabuff/1458895412.dem",
		expectLastChatMessage: "gg",
		expectHeroKillCount: map[string]int{
			"npc_dota_hero_faceless_void": 3,
		},
		expectHeroDeathCount: map[string]int{
			"npc_dota_hero_faceless_void": 2,
		},
	}

	testReplayCase(t, c)
}

// Esports match, played on patch 6.84c
func TestEsportsPatch684c(t *testing.T) {
	c := &testCase{
		matchId: 1483980562,
		url:     "https://s3-us-west-2.amazonaws.com/yasha.dotabuff/1483980562.dem",
		expectLastChatMessage: "gg wp",
		expectHeroKillCount: map[string]int{
			"npc_dota_hero_dragon_knight": 5,
			"npc_dota_hero_bristleback":   1,
		},
		expectHeroDeathCount: map[string]int{
			"npc_dota_hero_earthshaker": 6,
			"npc_dota_hero_bristleback": 3,
		},
	}

	testReplayCase(t, c)
}

// Manually scrutinised match, played on patch 6.84p1
func TestPublicMatchPatch684p1(t *testing.T) {
	assert := assert.New(t)

	data, err := getReplayData(1456774107, "https://s3-us-west-2.amazonaws.com/yasha.dotabuff/1456774107.dem")
	if err != nil {
		t.Fatalf("unable to get replay: %s", err)
	}

	parser := NewParser(data)
	parser.OnSayText2 = func(n int, o *dota.CUserMsg_SayText2) {
	}

	earthshakerDeaths := 0
	spiritBreakerDeaths := 0
	parser.OnCombatLog = func(entry CombatLogEntry) {
		// t.Logf("OnCombatLog: %s: %+v", reflect.TypeOf(entry), entry)
		switch log := entry.(type) {
		case *CombatLogDeath:
			if log.Target == "npc_dota_hero_earthshaker" {
				earthshakerDeaths++
			}
			if log.Target == "npc_dota_hero_spirit_breaker" {
				spiritBreakerDeaths++
			}
		}
	}

	var now time.Duration
	var gameTime, preGameStarttime float64
	parser.OnEntityPreserved = func(pe *PacketEntity) {
		if pe.Name == "DT_DOTAGamerulesProxy" {
			gameTime = pe.Values["DT_DOTAGamerules.m_fGameTime"].(float64)
			preGameStarttime = pe.Values["DT_DOTAGamerules.m_flPreGameStartTime"].(float64)
			now = time.Duration(gameTime-preGameStarttime) * time.Second
		}
	}

	// entindex:3 order_type:1 units:349 position:<x:6953.3125 y:6920.8438 z:384 > queue:false
	unitOrderCount := 0
	unitOrderQueuedCount := 0
	specificUnitOrder := false
	parser.OnSpectatorPlayerUnitOrders = func(n int, o *dota.CDOTAUserMsg_SpectatorPlayerUnitOrders) {
		unitOrderCount++
		if *o.Queue == true {
			unitOrderQueuedCount++
		}
		if *o.Entindex == 3 && *o.OrderType == 1 && o.Units[0] == 349 && *o.Queue == false &&
			*o.Position.X == 6953.3125 && *o.Position.Y == 6920.8438 && *o.Position.Y == 384.0 {
			specificUnitOrder = true
		}
	}

	chatWheelMessagesCount := 0
	parser.OnChatWheel = func(n int, o *dota.CDOTAUserMsg_ChatWheel) {
		chatWheelMessagesCount++
	}

	parser.Parse()

	assert.Equal(8, earthshakerDeaths)
	assert.Equal(11, spiritBreakerDeaths)          // not actually right but verified in replay
	assert.Equal(55316, unitOrderCount)            // regression test
	assert.Equal(102, unitOrderQueuedCount)        // regression test
	assert.Equal(int64(2585000000000), int64(now)) // regression test
	assert.Equal(0, chatWheelMessagesCount)        // regression test
}

func testReplayCase(t *testing.T, c *testCase) {
	assert := assert.New(t)

	data, err := getReplayData(c.matchId, c.url)
	if err != nil {
		t.Fatalf("unable to get replay: %s", err)
	}

	worldMins := &Vector3{}
	worldMaxes := &Vector3{}
	lastChatMessage := ""
	heroKillCount := make(map[string]int)
	heroDeathCount := make(map[string]int)

	parser := NewParser(data)
	parser.OnSayText2 = func(n int, o *dota.CUserMsg_SayText2) {
		lastChatMessage = o.GetText()
	}

	parser.OnChatEvent = func(n int, o *dota.CDOTAUserMsg_ChatEvent) {
	}

	parser.OnCombatLog = func(entry CombatLogEntry) {
		switch log := entry.(type) {
		case *CombatLogDeath:
			if strings.HasPrefix(log.Target, "npc_dota_hero_") {
				if _, ok := heroKillCount[log.Source]; !ok {
					heroKillCount[log.Source] = 0
				}
				heroKillCount[log.Source] += 1
			}

			if _, ok := heroDeathCount[log.Target]; !ok {
				heroDeathCount[log.Target] = 0
			}
			heroDeathCount[log.Target] += 1
		}
	}

	parser.OnEntityCreated = func(ent *PacketEntity) {
		if ent.Tick == 0 && ent.Name == "DT_WORLD" {
			worldMins = ent.Values["DT_WORLD.m_WorldMins"].(*Vector3)
			worldMaxes = ent.Values["DT_WORLD.m_WorldMaxs"].(*Vector3)
		}
	}

	parser.Parse()

	/*
		for _, table := range parser.Stsh.current {
			if table.Name == "instancebaseline" {
				for _, i := range table.Items {
					classId := atoi(i.Str)
					className := parser.ClassInfosNameMapping[classId]
					_dump_fixture(_sprintf("instancebaseline/%d_%s", c.matchId, className), i.Data)
				}
			}
		}
	*/

	// Make sure we have found the death counts for specified heroes
	if c.expectHeroDeathCount != nil {
		for hero, count := range c.expectHeroDeathCount {
			assert.Equal(count, heroDeathCount[hero], "expected hero %s to have death count %d", hero, count)
		}
	}

	// Make sure we have found the kill counts for specified heroes.
	if c.expectHeroKillCount != nil {
		for hero, count := range c.expectHeroKillCount {
			assert.Equal(count, heroKillCount[hero], "expected hero %s to have kill count %d", hero, count)
		}
	}

	// Make sure we find the DT_WORLD entity and it has the correct min and max dimensions.
	// This serves to help ensure our Float and Vector3 parsing is correct.
	assert.Equal(&Vector3{X: -8576.0, Y: -7680.0, Z: -1536.0}, worldMins)
	assert.Equal(&Vector3{X: 9216.0, Y: 8192.0, Z: 256.0}, worldMaxes)

	// Make sure we found the chat messages and have properly found the last one
	assert.Equal(c.expectLastChatMessage, lastChatMessage)
}

func getReplayData(matchId int64, url string) ([]byte, error) {
	path := fmt.Sprintf("replays/%d.dem", matchId)
	if data, err := ioutil.ReadFile(path); err == nil {
		fmt.Printf("read replay %d from %s\n", matchId, path)
		return data, nil
	}

	fmt.Printf("downloading replay %d from %s...\n", matchId, url)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Return an error if we don't get a 200
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("invalid status %d", resp.StatusCode)
	}

	var data []byte
	if url[len(url)-3:] == "bz2" {
		data, err = ioutil.ReadAll(bzip2.NewReader(resp.Body))
	} else {
		data, err = ioutil.ReadAll(resp.Body)
	}

	if err != nil {
		return nil, err
	}

	if err := ioutil.WriteFile(path, data, 0644); err != nil {
		return nil, err
	}

	fmt.Printf("downloaded replay %d from %s to %s\n", matchId, url, path)

	return data, nil
}
