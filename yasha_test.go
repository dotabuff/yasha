package yasha

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

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
		url:     "http://s.tsai.co/replays/1405240741.dem",
		expectLastChatMessage: "Gg",
	}

	testReplayCase(t, c)
}

// Navi vs Basically Unknown, played on patch 6.84p0
func TestEsportsPatch684p0(t *testing.T) {
	c := &testCase{
		matchId: 1450235906,
		url:     "http://s.tsai.co/replays/1450235906.dem",
		expectLastChatMessage: "gg",
	}

	testReplayCase(t, c)
}

// No Respeta Funadores vs Who Needs Skill, played on patch 6.84p1
func TestEsportsPatch684p1(t *testing.T) {
	c := &testCase{
		matchId: 1458895412,
		url:     "http://s.tsai.co/replays/1458895412.dem",
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
		//t.Logf("OnSayText2: %+v", o)
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
				t.Logf("OnCombatLog: %+v", log)
				spiritBreakerDeaths++
			}
		}
	}

	parser.Parse()

	assert.Equal(8, earthshakerDeaths)
	assert.Equal(11, spiritBreakerDeaths) // not actually right but verified in replay
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
		t.Logf("OnSayText2: %+v", o)
		lastChatMessage = o.GetText()
	}

	parser.OnChatEvent = func(n int, o *dota.CDOTAUserMsg_ChatEvent) {
		t.Logf("OnChatEvent: %+v", o)
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

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if err := ioutil.WriteFile(path, data, 0644); err != nil {
		return nil, err
	}

	fmt.Printf("downloaded replay %d from %s to %s\n", matchId, url, path)

	return data, nil
}
