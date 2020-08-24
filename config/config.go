package config

const (
	// Port for listenning
	Port          = ":8080"
	MathPlusN     = "+"
	MathPlusA     = "p"
	MathPlusSpace = " "
	MathNegative  = "-"
	MathMultiply  = "*"
	MathDevide    = "/"
)

type MathFormulaStruct struct {
	FirstValue  float64
	Operation   string
	SecondValue float64
	Answer      string
}

type PageVariables struct {
	PageTitle    string
	PageGameMove GameMove
	OtherGames   []string
}

type GameMove struct {
	CurrentGame     string
	GamePosition    string
	CurrentPosition string
	NextPosition    string
	GameTurn        string
	GameImage       string
	PossibleMoves   []string
}
