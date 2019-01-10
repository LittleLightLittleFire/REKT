package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

type (
	// State tracks of the largest liquidations as well as kill streaks.
	State struct {
		SaveFile   string
		HighScores HighScores

		Snark      []string
		SnarkIndex int

		MultiKill []string
	}

	// Scores for a particular symbol.
	Scores struct {
		HighestDay   int64 `json:"highest_day"`
		HighestWeek  int64 `json:"highest_week"`
		HighestMonth int64 `json:"highest_month"`

		LastDay   int        `json:"last_day"`
		LastWeek  int        `json:"last_week"`
		LastMonth time.Month `json:"last_month"`
	}

	// Kill stores the last time a position was liquidated on a symbol.
	// It also stores the last time it was updated
	Kill struct {
		Count    int   `json:"count"`
		UnixTime int64 `json:"unix_time"`
	}

	// HighScores defines a data structure that store high scores.
	HighScores struct {
		Scores map[Symbol]Scores `json:"scores"`
		Kills  map[Symbol]Kill   `json:"kills"`
	}

	// A Medal is awarded to the liquidation if it breaks a high score.
	Medal int32

	// DecoratedLiquidation gives liqudation extra properties based on its timing and size.
	DecoratedLiquidation struct {
		Streak      string      // Multikills
		Medals      []Medal     // Medals
		Snark       string      // Snarky meme text to salt the wound
		Liquidation Liquidation // Actual liquidiation
	}
)

// Medals a liqudiation can win.
const (
	MedalLargestToday Medal = iota
	MedalLargestWeek
	MedalLargestMonth

	Medal100k      // Awarded for every 100k
	MedalStreak    // Killed as part of a kill streak
	MedalSecKilled // Killed within two seconds of the previous

	// TODO: More to come
)

// Twitter has extended the length limit.
const twitterLengthLimit = 280

var medalMap = map[Medal]string{
	MedalLargestToday: "", // Disabled since liquidations are pretty rare
	MedalLargestWeek:  "\U0001F3C5",
	MedalLargestMonth: "\U0001F3C6",
	Medal100k:         "\U0001F4AF",
	MedalStreak:       "\U0001F525",
	MedalSecKilled:    "\U000026A1",
}

// NewState returns a new state object.
func NewState() (*State, error) {
	// TODO: move hardcoded files out of here.
	highScoresFile := "high_scores.json"
	snarkFile := "text/memes.txt"
	multiKillFile := "text/kill_streaks.txt"

	var state State

	// Load high scores
	if f, err := os.Open(highScoresFile); err != nil {
		state.HighScores = HighScores{
			make(map[Symbol]Scores),
			make(map[Symbol]Kill),
		}
	} else {
		defer f.Close()

		if err := json.NewDecoder(f).Decode(&state.HighScores); err != nil {
			return nil, err
		}
	}
	state.SaveFile = highScoresFile

	// Load memes
	snarkText, err := ioutil.ReadFile(snarkFile)
	if err != nil {
		return nil, err
	}
	state.Snark = strings.Split(strings.TrimSpace(string(snarkText)), "\n")

	// Shuffle
	state.resetSnark()

	// Load multi-kill
	multiKillText, err := ioutil.ReadFile(multiKillFile)
	if err != nil {
		return nil, err
	}
	state.MultiKill = strings.Split(strings.TrimSpace(string(multiKillText)), "\n")

	return &state, nil
}

// resetSnark shuffles the snark array and resets the counter.
func (s *State) resetSnark() {
	s.SnarkIndex = 0
	for i := range s.Snark {
		j := rand.Intn(i + 1)
		s.Snark[i], s.Snark[j] = s.Snark[j], s.Snark[i]
	}

	log.Println("Banter order:")
	for _, v := range s.Snark {
		log.Println("    ", v)
	}
}

// Save stores the high scores back to disk.
func (s *State) Save() error {
	f, err := os.OpenFile(s.SaveFile, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := json.NewEncoder(f).Encode(s.HighScores); err != nil {
		return err
	}

	return nil
}

// Linear interpolation
func lerp(x, y, z int64, start, end float64) float64 {
	return start + ((float64(z)-float64(x))/(float64(y)-float64(x)))*(end-start)
}

// Decorate a new liqudation.
func (s *State) Decorate(l Liquidation) DecoratedLiquidation {
	// Hand out medals
	var medals []Medal
	scores := s.HighScores.Scores[l.Symbol]

	// Expire the scores if their time has reached
	now := time.Now()
	if now.Day() != scores.LastDay {
		scores.LastDay = now.Day()
		scores.HighestDay = 0
	}

	_, week := now.ISOWeek()
	if week != scores.LastWeek {
		scores.LastWeek = week
		scores.HighestWeek = 0
	}

	if now.Month() != scores.LastMonth {
		scores.LastMonth = now.Month()
		scores.HighestMonth = 0
	}

	// Issue medal for each of the periods
	if l.Quantity >= scores.HighestWeek {
		scores.HighestWeek = l.Quantity
		medals = append(medals, MedalLargestWeek)
	}

	if l.Quantity >= scores.HighestMonth {
		scores.HighestMonth = l.Quantity
		medals = append(medals, MedalLargestMonth)
	}

	// Award the 100k medals
	for i := int64(0); i < l.Quantity/100000; i++ {
		medals = append(medals, Medal100k)
	}

	s.HighScores.Scores[l.Symbol] = scores

	// Issue the streak
	streak := s.HighScores.Kills[l.Symbol]

	if now.Unix()-streak.UnixTime > 60 {
		streak.Count = 0
	}
	streak.Count++
	if streak.Count >= 2 {
		medals = append(medals, MedalStreak)
	}

	// Issue the medal for being Seckilled
	if now.Unix()-streak.UnixTime <= 10 {
		medals = append(medals, MedalSecKilled)
	}

	streak.UnixTime = now.Unix()
	s.HighScores.Kills[l.Symbol] = streak

	// Issue the snark
	// Because we have limited text, we will not be able to issue snark every single time.

	// USD value:    0 -------- 100k ------------ 500k ------------ 2m --------->
	// Snark prob:        0%            5%-10%           10%-30%
	var issueSnark bool

	usdVal := l.USDValue()
	switch {
	case usdVal <= 100000:
		issueSnark = false
	case usdVal <= 500000:
		issueSnark = lerp(50000, 500000, usdVal, 0.05, 0.10) > rand.Float64()
	default:
		issueSnark = lerp(500000, 2000000, usdVal, 0.10, 0.30) > rand.Float64()
	}

	var snark string

	if issueSnark {
		s.SnarkIndex = (s.SnarkIndex + 1) % len(s.Snark)
		// Check if we've wrapped around now
		if s.SnarkIndex == 0 {
			s.resetSnark()
		}
		snark = s.Snark[s.SnarkIndex]
	}

	// TODO: refactor this

	var streakStrRaw string
	streak.Count -= 2
	if streak.Count < 0 {
		// No streak
	} else if streak.Count >= len(s.MultiKill) {
		streakStrRaw = s.MultiKill[len(s.MultiKill)-1] + " x" + strconv.Itoa(streak.Count+2)
	} else {
		streakStrRaw = s.MultiKill[streak.Count]
	}
	streakStr := strings.Replace(streakStrRaw, "$SYMBOL", string(l.Symbol), -1)
	snarkStr := strings.Replace(snark, "$SYMBOL", string(l.Symbol), -1)

	dl := DecoratedLiquidation{
		Streak:      streakStr,
		Medals:      medals,
		Snark:       snarkStr,
		Liquidation: l,
	}

	return dl
}

func (dl DecoratedLiquidation) hasMedals() bool {
	return len(dl.Medals) > 0
}

func (dl DecoratedLiquidation) hasSnark() bool {
	return dl.Snark != ""
}

func (dl DecoratedLiquidation) hasStreak() bool {
	return dl.Streak != ""
}

func (dl DecoratedLiquidation) medalsRunes() (result []rune) {
	if !dl.hasMedals() {
		return
	}

	result = append(result, ' ')
	for _, medal := range dl.Medals {
		result = append(result, []rune(medalMap[medal])...)
	}
	return
}

func (dl DecoratedLiquidation) streakRunes() []rune {
	if !dl.hasStreak() {
		return nil
	}

	return append([]rune(" ~ "), []rune(dl.Streak)...)
}

func (dl DecoratedLiquidation) snarkRunes() []rune {
	if !dl.hasSnark() {
		return nil
	}

	return append([]rune(" ~ "), []rune(dl.Snark)...)
}

// String implements Stringer.
func (dl DecoratedLiquidation) String() string {
	// We need to fit our string into the Twitter length
	// However Twitter documentation is full of shit
	//     https://developer.twitter.com/en/docs/basics/counting-characters.html
	// They do not count emojis which are a single unicode codepoint as a single character, they count it as two
	// The fact is, they've complicated this so much it requires a library (twitter-text) to figure out exactly what length they'll calculate this to be
	// So erring on the side of safety, we'll count all text in medals as two characters
	// This leave us with a safety margin of three characters created by ` ~ ` for any emojis in the snark itself
	base := []rune(dl.Liquidation.String())

	if len(base)+len(dl.medalsRunes())*2+len(dl.streakRunes())+len(dl.snarkRunes()) <= twitterLengthLimit {
		// It just works
		base = append(base, dl.medalsRunes()...)
		base = append(base, dl.streakRunes()...)
		base = append(base, dl.snarkRunes()...)
		return string(base)
	}

	if len(base)+len(dl.medalsRunes())*2+len(dl.snarkRunes()) <= twitterLengthLimit {
		// We'll do without the streak then
		base = append(base, dl.medalsRunes()...)
		base = append(base, dl.snarkRunes()...)
		return string(base)
	}

	if len(base)+len(dl.snarkRunes()) <= twitterLengthLimit {
		medalLength := (twitterLengthLimit - (len(base) + len(dl.snarkRunes()))) / 2

		// We'll trim the medals so that we use up the entire text, unless the medals get trimmed to nothing
		if medalLength > 3 {
			base = append(base, dl.medalsRunes()[:medalLength]...)
		}
		base = append(base, dl.snarkRunes()...)
		return string(base)
	}

	return string(base)
}
