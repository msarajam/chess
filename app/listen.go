package app

import (
	"errors"
	"fmt"
	"html/template"
	"image/color"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"webServer/config"

	guuid "github.com/google/uuid"
	"github.com/notnil/chess"
	"github.com/notnil/chess/image"
)

var (
	counter = 0
	mutex   = &sync.Mutex{}
)

func BasicHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "counter\n")
	fmt.Fprintf(w, "math\n")
	fmt.Fprintf(w, "chess\n")
}

func ChessHandler(w http.ResponseWriter, r *http.Request) {
	cp := ""
	np := ""
	cg := fmt.Sprint(guuid.New())
	gp := ""
	gt := ""
	existingGame := false

	keys, ok := r.URL.Query()["CurrentGame"]
	if ok {
		cg = keys[0]
		existingGame = true
	}
	if existingGame {
		keys, ok = r.URL.Query()["CurrentPosition"]
		if ok {
			cp = strings.ToLower(keys[0])
		}
		keys, ok = r.URL.Query()["GameTurn"]
		if ok {
			gt = keys[0]
		}

		keys, ok = r.URL.Query()["NextPosition"]
		if ok {
			np = strings.ToLower(keys[0])
		}
		data, err := ioutil.ReadFile("./games/" + cg + ".txt")
		if err != nil {
			fmt.Println("ERROR - Cant open the File : ", err.Error())
		} else {
			gp = string(data)
		}
	} else {
		_, err := os.Create("./games/" + cg + ".txt")
		if err != nil {
			fmt.Println("ERROR - Cant create a File", err.Error())
		}
		g := chess.NewGame()
		gp = fmt.Sprint(g.Position())
	}

	//fmt.Printf("CP : %v\nNP : %v\nCG : %v\nGP : %v\nGT : %v\n", cp, np, cg, gp, gt)

	fen, _ := chess.FEN(gp)
	g := chess.NewGame(fen)
	if cp != "" && np != "" {
		for _, v := range g.ValidMoves() {
			a := cp + np
			b := fmt.Sprint(v)
			if a == b {
				g.Move(v)
				fileName := "./games/" + cg + ".txt"
				if err := ioutil.WriteFile(fileName, []byte(fmt.Sprint(g.Position())), 0644); err != nil {
					fmt.Println("Error - ", err)
				}

				break
			}
		}
		fmt.Println("")
	}

	f, err := os.Create("./images/" + cg + ".svg")
	if err != nil {
		fmt.Println("Error Here line 59")
	}
	defer f.Close()

	fenStr := fmt.Sprintln(g.Position())
	pos := &chess.Position{}
	if err := pos.UnmarshalText([]byte(fenStr)); err != nil {
		fmt.Println("Error Here line 66")
	}

	if err := image.SVG(f, pos.Board()); err != nil {
		fmt.Println("Error Here line 73")
	}

	var pm []string
	for _, v := range g.ValidMoves() {
		pm = append(pm, fmt.Sprint(v))
	}

	gt = fmt.Sprint(pos.Turn())
	Title := "My Game"
	gm := config.GameMove{
		CurrentGame:     cg,
		CurrentPosition: cp,
		NextPosition:    np,
		GamePosition:    gp,
		GameTurn:        gt,
		PossibleMoves:   pm,
		GameImage:       "/images/" + cg + ".svg",
	}

	var files []string

	root := "./games/"
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {

		a := strings.Index(path, "es/")
		if a > 0 {
			path = path[a+len("es/"):]
		}
		b := strings.Index(path, ".txt")
		if b > 0 {
			path = path[:b]
		}
		fmt.Println(path)
		if len(path) > 1 {

			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		fmt.Println("ERROR - ", err)
	}
	/*
		for _, file := range files {
				fmt.Println(file)
		}
	*/

	MyPageVariables := config.PageVariables{
		PageTitle:    Title,
		PageGameMove: gm,
		OtherGames:   files,
	}

	t, err := template.ParseFiles("pages/homepage.html") //parse the html file homepage.html
	if err != nil {                                      // if there is an error
		log.Print("template parsing error: ", err) // log it
	}

	err = t.Execute(w, MyPageVariables)
	if err != nil {
		log.Print("template executing error: ", err)
	}
	//	fmt.Println(g.Position())
	fmt.Println(gm)
	fmt.Println("DONE")
}

func TempHandler(w http.ResponseWriter, r *http.Request) {
	g := chess.NewGame()

	for g.Outcome() == chess.NoOutcome {
		ms := g.ValidMoves()
		m := ms[rand.Intn(len(ms))]
		fmt.Println(m)
		g.Move(m)
		break
	}
	fmt.Println(g.Position().Board().Draw())
	//fmt.Println(g.Position())
	//fmt.Printf("Game Completed. %s by %s. \n", g.Outcome(), g.Method())
	fmt.Println(g.String())

	//	g.MoveStr("f3")
	/*
		g.MoveStr("e6")
		//	g.MoveStr("g4")
		//	g.MoveStr("Qh4")
		fmt.Println(g.Outcome())
		fmt.Println(g.Method())
		fmt.Fprintln(w, g.Outcome())
		fmt.Fprintln(w, g.Method())

		fmt.Println(g.Position().Board().Draw())
		fmt.Fprintf(w, g.Position().Board().Draw())
		g = chess.NewGame()
		fmt.Fprintf(w, g.Position().Board().Draw())
	*/

	f, err := os.Create("./images/TestChess.svg")
	if err != nil {
		fmt.Println("Error Here line 59")
	}
	defer f.Close()

	//fenStr := "rnbqkbnr/pppppppp/8/8/3P4/8/PPP1PPPP/RNBQKBNR b KQkq - 0 1"
	fenStr := fmt.Sprintln(g.Position())
	pos := &chess.Position{}
	if err := pos.UnmarshalText([]byte(fenStr)); err != nil {
		fmt.Println("Error Here line 66")
	}

	// write board SVG to file
	yellow := color.RGBA{255, 255, 0, 1}
	mark := image.MarkSquares(yellow, chess.D3, chess.D4)
	if err := image.SVG(f, pos.Board(), mark); err != nil {
		fmt.Println("Error Here line 73")
	}
}

func CounterHandler(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	counter++
	req := "Number of Request : " + strconv.Itoa(counter)
	fmt.Fprintf(w, req)
	mutex.Unlock()
}

func MathHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Query())
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("Error - %v", err)
	}
	req := string(b)
	if len(b) == 0 {
		for k, _ := range r.URL.Query() {
			req = k
			break
		}
	}
	fmt.Println("Got Request : ", req)

	var f config.MathFormulaStruct
	err = findOperation(&req, &f)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		fmt.Println(err.Error())
		return
	}
	formula := strings.Split(req, f.Operation)
	if checkFormula(formula, &f) {
		if err = calculation(&f); err != nil {
			fmt.Fprintf(w, err.Error())
			return
		}
	}
	fmt.Printf("Formula : %v\n", f)
	fmt.Fprintf(w, fmt.Sprintf("Answer : %v", f.Answer))
}

func findOperation(s *string, f *config.MathFormulaStruct) error {
	if strings.Contains(*s, config.MathPlusA) || strings.Contains(*s, config.MathPlusN) || strings.Contains(*s, config.MathPlusSpace) {
		*s = strings.Replace(*s, "p", "+", -1)
		*s = strings.Replace(*s, " ", "+", -1)
		f.Operation = config.MathPlusN
	} else if strings.Contains(*s, config.MathNegative) {
		f.Operation = config.MathNegative
	} else if strings.Contains(*s, config.MathMultiply) {
		f.Operation = config.MathMultiply
	} else if strings.Contains(*s, config.MathDevide) {
		f.Operation = config.MathDevide
	} else {
		return errors.New("Wrong Operation : " + *s)
	}
	return nil
}

func checkFormula(s []string, f *config.MathFormulaStruct) bool {
	if len(s) == 2 {
		for _, v := range s {
			if len(v) < 1 {
				return false
			}
		}
	} else {
		return false
	}
	var err error
	if f.FirstValue, err = strconv.ParseFloat(s[0], 64); err != nil {
		return false
	}

	if f.SecondValue, err = strconv.ParseFloat(s[1], 64); err != nil {
		return false
	}
	return true
}

func calculation(f *config.MathFormulaStruct) error {
	switch f.Operation {
	case config.MathPlusN:
		f.Answer = fmt.Sprintf("%.2f", f.FirstValue+f.SecondValue)
	case config.MathNegative:
		f.Answer = fmt.Sprintf("%.2f", f.FirstValue-f.SecondValue)
	case config.MathMultiply:
		f.Answer = fmt.Sprintf("%.2f", f.FirstValue*f.SecondValue)
	case config.MathDevide:
		f.Answer = fmt.Sprintf("%.2f", f.FirstValue/f.SecondValue)
	default:
		return errors.New(fmt.Sprintf("Wrong formula operation : %v", f.Operation))
	}
	return nil
}
