package yasha

import (
	"compress/bzip2"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/dotabuff/yasha/dota"
	"github.com/stretchr/testify/assert"
)

type testCase struct {
	matchId int64
	url     string

	expectLastChatMessage string
}

// Secret vs Cloud 9 played on patch 6.83c
func TestEsportsPatch683b(t *testing.T) {
	c := &testCase{
		matchId: 1405240741,
		url:     "http://replay135.valve.net/570/1405240741_220241732.dem.bz2",
		expectLastChatMessage: "Gg",
	}

	testReplayCase(t, c)
}

// Navi vs Basically Unknown, played on patch 6.84p0
func TestEsportsPatch684p0(t *testing.T) {
	c := &testCase{
		matchId: 1450235906,
		url:     "http://replay136.valve.net/570/1450235906_1463120933.dem.bz2",
		expectLastChatMessage: "gg",
	}

	testReplayCase(t, c)
}

// No Respeta Funadores vs Who Needs Skill, played on patch 6.84p1
func TestEsportsPatch684p1(t *testing.T) {
	c := &testCase{
		matchId: 1458895412,
		url:     "http://replay123.valve.net/570/1458895412_140022944.dem.bz2",
		expectLastChatMessage: "gg",
	}

	testReplayCase(t, c)
}

// Manually scrutinised match, played on patch 6.84p1
func TestPublicMatchPatch684p1(t *testing.T) {
	assert := assert.New(t)

	data, err := getReplayData(1456774107, "http://s.tsai.co/replays/1456774107.dem")
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

	lastChatMessage := ""

	parser := NewParser(data)
	parser.OnSayText2 = func(n int, o *dota.CUserMsg_SayText2) {
		//t.Logf("OnSayText2: %+v", o)
		lastChatMessage = o.GetText()
	}

	parser.OnChatEvent = func(n int, o *dota.CDOTAUserMsg_ChatEvent) {
		//t.Logf("OnChatEvent: %+v", o)
	}
	parser.Parse()

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
