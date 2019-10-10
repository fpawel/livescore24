package livescore

type Champ struct {
	Name string `json:"name"`
	Games []Game `json:"games,omitempty"`
}

type Game struct {
	Team1 string `json:"team1,omitempty"`
	Team2 string `json:"team2,omitempty"`
	Time string `json:"time,omitempty"`
	Score string `json:"score,omitempty"`
	Score2 string `json:"score2,omitempty"`
	Win float64 `json:"win,omitempty"`
	Draw float64 `json:"draw,omitempty"`
	Lose float64 `json:"lose,omitempty"`
}
