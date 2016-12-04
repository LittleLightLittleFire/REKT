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
	MedalSecKilled // Killed within two seconds of the previous

	// TODO: More to come
)

var medalMap = map[Medal]string{
	MedalLargestToday: "\U0001F396",
	MedalLargestWeek:  "\U0001F3C5",
	MedalLargestMonth: "\U0001F3C6",
	Medal100k:         "\U0001F4AF",
	MedalSecKilled:    "\U0001F525",
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
	if l.Quantity > scores.HighestDay {
		scores.HighestDay = l.Quantity
		medals = append(medals, MedalLargestToday)
	}

	if l.Quantity > scores.HighestWeek {
		scores.HighestWeek = l.Quantity
		medals = append(medals, MedalLargestWeek)
	}

	if l.Quantity > scores.HighestMonth {
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

	if now.Unix()-streak.UnixTime > 20 {
		streak.Count = 0
	}
	streak.Count++

	// Issue the medal for being Seckilled
	if now.Unix()-streak.UnixTime <= 2 {
		medals = append(medals, MedalSecKilled)
	}

	streak.UnixTime = now.Unix()
	s.HighScores.Kills[l.Symbol] = streak

	// Issue the snark
	// Because we have limited text, we will not be able to issue snark every single time.

	// USD value:    0 -------- 10k ---------- 50k-------------- 500k --------->
	// Snark prob:       0%-5%         5%-20%        20%-100%
	//
	// Each awarded medal boosts by 5%
	var issueSnark bool

	usdVal := l.USDValue()
	switch {
	case usdVal <= 10000:
		issueSnark = lerp(0, 10000, usdVal, 0.00, 0.05) > rand.Float64()
	case usdVal <= 50000:
		issueSnark = lerp(10000, 50000, usdVal, 0.05, 0.20) > rand.Float64()
	default:
		issueSnark = lerp(50000, 500000, usdVal, 0.20, 1.00) > rand.Float64()
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
		streakStrRaw = s.MultiKill[len(s.MultiKill)-1] + " x" + strconv.Itoa(streak.Count)
	} else {
		streakStrRaw = s.MultiKill[streak.Count]
	}
	streakStr := strings.Replace(streakStrRaw, "$SYMBOL", string(l.Symbol), -1)
	snarkStr := strings.Replace(snark, "$SYMBOL", string(l.Symbol), -1)

	return DecoratedLiquidation{
		Streak:      streakStr,
		Medals:      medals,
		Snark:       snarkStr,
		Liquidation: l,
	}
}

// String implements Stringer.
func (dl DecoratedLiquidation) String() string {
	base := dl.Liquidation.String()

	// Add medals
	if len(dl.Medals) > 0 {
		base += " "
		for _, medal := range dl.Medals {
			base += medalMap[medal]
		}
	}

	if dl.Streak != "" {
		base += " ~ " + dl.Streak
	}

	// If the text gets too long, don't bother writing the snark
	if len(base)+3+len(dl.Snark) < 140 {
		if dl.Snark != "" {
			base += " ~ " + dl.Snark
		}
	}

	return base
}
