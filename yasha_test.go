package yasha

import (
	"compress/bzip2"
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

	reader := bzip2.NewReader(resp.Body)
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	if err := ioutil.WriteFile(path, data, 0644); err != nil {
		return nil, err
	}

	fmt.Printf("downloaded replay %d from %s to %s\n", matchId, url, path)

	return data, nil
}
