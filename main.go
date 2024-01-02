package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
)

var ACCIDENT = 0
var LOGO = 0
var FLY = 0
var C51 = 0

const OK = 0
const ERR = -1

var COLS, LINES int
var s tcell.Screen

func main() {
	tcell.SetEncodingFallback(tcell.EncodingFallbackASCII)
	var err error
	s, err = tcell.NewScreen()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	if err = s.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer s.Fini()

	s.SetStyle(tcell.StyleDefault)
	s.Clear()

	s.SetCursorStyle(tcell.CursorStyleDefault)

	COLS, LINES = s.Size()
	args := os.Args[1:]

	for _, arg := range args {
		if arg[0] == '-' {
			option(arg[1:])
		}
	}

	// Didn't work, the library internally handled
	// signal.Ignore(os.Interrupt)

	s.HideCursor()

	for x := COLS - 1; ; x-- {
		if LOGO == 1 {
			if add_sl(x) == ERR {
				break
			}
		} else if C51 == 1 {
			if add_C51(x) == ERR {
				break
			}
		} else {
			if add_D51(x) == ERR {
				break
			}
		}

		s.Show()
		// s.Sync()

		time.Sleep(40 * time.Millisecond)
	}
	s.ShowCursor(0, 0)
}

func my_mvaddstr(y, x int, str string) {
	// if y < 0 {
	// 	return
	// }
	for ; x < 0; x, str = x+1, str[1:] {
		if len(str) == 0 {
			return
		}
	}
	for _, char := range str {
		// s.SetCell(x, y, tcell.StyleDefault, char)
		s.SetContent(x, y, char, nil, tcell.StyleDefault)
		x++
	}
}

func option(str string) {
	for _, char := range str {
		switch char {
		case 'a':
			ACCIDENT = 1
		case 'F':
			FLY = 1
		case 'l':
			LOGO = 1
		case 'c':
			C51 = 1
		default:
		}
	}
}

var add_sl = func() func(x int) int {
	sl := [LOGOPATTERNS][LOGOPATTERNS + 1]string{
		{LOGO1, LOGO2, LOGO3, LOGO4, LWHL11, LWHL12, DELLN},
		{LOGO1, LOGO2, LOGO3, LOGO4, LWHL21, LWHL22, DELLN},
		{LOGO1, LOGO2, LOGO3, LOGO4, LWHL31, LWHL32, DELLN},
		{LOGO1, LOGO2, LOGO3, LOGO4, LWHL41, LWHL42, DELLN},
		{LOGO1, LOGO2, LOGO3, LOGO4, LWHL51, LWHL52, DELLN},
		{LOGO1, LOGO2, LOGO3, LOGO4, LWHL61, LWHL62, DELLN},
	}

	coal := [LOGOPATTERNS + 1]string{LCOAL1, LCOAL2, LCOAL3, LCOAL4, LCOAL5, LCOAL6, DELLN}
	car := [LOGOPATTERNS + 1]string{LCAR1, LCAR2, LCAR3, LCAR4, LCAR5, LCAR6, DELLN}

	return func(x int) int {
		y := 0
		py1, py2, py3 := 0, 0, 0

		if x < -LOGOLENGTH {
			return -1
		}
		y = LINES/2 - 3

		if FLY == 1 {
			y = (x / 6) + LINES - (COLS / 6) - LOGOHEIGHT
			py1, py2, py3 = 2, 4, 6
		}
		for i := 0; i <= LOGOHEIGHT; i++ {
			my_mvaddstr(y+i, x, sl[(LOGOLENGTH+x)/3%LOGOPATTERNS][i])
			my_mvaddstr(y+i+py1, x+21, coal[i])
			my_mvaddstr(y+i+py2, x+42, car[i])
			my_mvaddstr(y+i+py3, x+63, car[i])
		}
		if ACCIDENT == 1 {
			add_man(y+1, x+14)
			add_man(y+1+py2, x+45)
			add_man(y+1+py2, x+53)
			add_man(y+1+py3, x+66)
			add_man(y+1+py3, x+74)
		}
		add_smoke(y-1, x+LOGOFUNNEL)
		return 0
	}
}()

var add_D51 = func() func(x int) int {
	d51 := [D51PATTERNS][D51HEIGHT + 1]string{
		{D51STR1, D51STR2, D51STR3, D51STR4, D51STR5, D51STR6, D51STR7, D51WHL11, D51WHL12, D51WHL13, D51DEL},
		{D51STR1, D51STR2, D51STR3, D51STR4, D51STR5, D51STR6, D51STR7, D51WHL21, D51WHL22, D51WHL23, D51DEL},
		{D51STR1, D51STR2, D51STR3, D51STR4, D51STR5, D51STR6, D51STR7, D51WHL31, D51WHL32, D51WHL33, D51DEL},
		{D51STR1, D51STR2, D51STR3, D51STR4, D51STR5, D51STR6, D51STR7, D51WHL41, D51WHL42, D51WHL43, D51DEL},
		{D51STR1, D51STR2, D51STR3, D51STR4, D51STR5, D51STR6, D51STR7, D51WHL51, D51WHL52, D51WHL53, D51DEL},
		{D51STR1, D51STR2, D51STR3, D51STR4, D51STR5, D51STR6, D51STR7, D51WHL61, D51WHL62, D51WHL63, D51DEL},
	}
	coal := [D51HEIGHT + 1]string{COAL01, COAL02, COAL03, COAL04, COAL05, COAL06, COAL07, COAL08, COAL09, COAL10, COALDEL}

	return func(x int) int {
		y, dy := 0, 0

		if x < -D51LENGTH {
			return -1
		}
		y = LINES/2 - 5

		if FLY == 1 {
			y = (x / 7) + LINES - (COLS / 7) - D51HEIGHT
			dy = 1
		}
		for i := 0; i <= D51HEIGHT; i++ {
			my_mvaddstr(y+i, x, d51[(D51LENGTH+x)%D51PATTERNS][i])
			my_mvaddstr(y+i+dy, x+53, coal[i])
		}
		if ACCIDENT == 1 {
			add_man(y+2, x+43)
			add_man(y+2, x+47)
		}
		add_smoke(y-1, x+D51FUNNEL)
		return OK
	}
}()

var add_C51 = func() func(x int) int {
	c51 := [C51PATTERNS][C51LENGTH + 1]string{
		{C51STR1, C51STR2, C51STR3, C51STR4, C51STR5, C51STR6, C51STR7, C51WH11, C51WH12, C51WH13, C51WH14, C51DEL},
		{C51STR1, C51STR2, C51STR3, C51STR4, C51STR5, C51STR6, C51STR7, C51WH21, C51WH22, C51WH23, C51WH24, C51DEL},
		{C51STR1, C51STR2, C51STR3, C51STR4, C51STR5, C51STR6, C51STR7, C51WH31, C51WH32, C51WH33, C51WH34, C51DEL},
		{C51STR1, C51STR2, C51STR3, C51STR4, C51STR5, C51STR6, C51STR7, C51WH41, C51WH42, C51WH43, C51WH44, C51DEL},
		{C51STR1, C51STR2, C51STR3, C51STR4, C51STR5, C51STR6, C51STR7, C51WH51, C51WH52, C51WH53, C51WH54, C51DEL},
		{C51STR1, C51STR2, C51STR3, C51STR4, C51STR5, C51STR6, C51STR7, C51WH61, C51WH62, C51WH63, C51WH64, C51DEL},
	}
	coal := [C51LENGTH + 1]string{COALDEL, COAL01, COAL02, COAL03, COAL04, COAL05, COAL06, COAL07, COAL08, COAL09, COAL10, COALDEL}

	return func(x int) int {
		y, dy := 0, 0

		if x < -C51LENGTH {
			return -1
		}
		y = LINES/2 - 5

		if FLY == 1 {
			y = (x / 7) + LINES - (COLS / 7) - C51HEIGHT
			dy = 1
		}
		for i := 0; i <= C51HEIGHT; i++ {
			my_mvaddstr(y+i, x, c51[(C51LENGTH+x)%C51PATTERNS][i])
			my_mvaddstr(y+i+dy, x+55, coal[i])
		}
		if ACCIDENT == 1 {
			add_man(y+3, x+45)
			add_man(y+3, x+49)
		}
		add_smoke(y-1, x+C51FUNNEL)
		return OK
	}
}()

var add_man = func() func(y, x int) {
	man := [2][2]string{{"", "(O)"}, {"Help!", "\\O/"}}
	return func(y, x int) {

		for i := 0; i < 2; i++ {
			my_mvaddstr(y+i, x, man[(LOGOLENGTH+x)/12%2][i])
		}
	}
}()

var add_smoke = func() func(y, x int) {
	type smokes struct {
		y, x, ptrn, kind int
	}
	const SMOKEPTNS = 16
	S := make([]smokes, 1000)
	var sum int
	Smoke := [2][SMOKEPTNS]string{
		{"(   )", "(    )", "(    )", "(   )", "(  )",
			"(  )", "( )", "( )", "()", "()", "O", "O", "O", "O", "O", " "},
		{"(@@@)", "(@@@@)", "(@@@@)", "(@@@)", "(@@)",
			"(@@)", "(@)", "(@)", "@@", "@@", "@", "@", "@", "@", "@", " "},
	}
	Eraser := [SMOKEPTNS]string{"     ", "      ", "      ", "     ", "   ", "   ", "  ", "  ", "  ", " ", " ", " ", " ", " ", " ", " "}
	dy := []int{2, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	dx := []int{-2, -1, 0, 1, 1, 1, 1, 1, 2, 2, 2, 2, 2, 3, 3, 3}
	return func(y, x int) {
		if x%4 == 0 {
			for i := 0; i < sum; i++ {
				my_mvaddstr(S[i].y, S[i].x, Eraser[S[i].ptrn])
				S[i].y -= dy[S[i].ptrn]
				S[i].x += dx[S[i].ptrn]
				if S[i].ptrn < SMOKEPTNS-1 {
					S[i].ptrn++
				}
				my_mvaddstr(S[i].y, S[i].x, Smoke[S[i].kind][S[i].ptrn])
			}
			my_mvaddstr(y, x, Smoke[sum%2][0])
			S[sum].y = y
			S[sum].x = x
			S[sum].ptrn = 0
			S[sum].kind = sum % 2
			sum++
		}
	}
}()
