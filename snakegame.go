package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

//#include <windows.h>
//#include <conio.h>
/*
//
void gotoxy(int x,int y)
{
    COORD c;
    c.X=x,c.Y=y;
    SetConsoleCursorPosition(GetStdHandle(STD_OUTPUT_HANDLE),c);
}

int onKeyboard()
{
    return _getch();
}
*/
import "C"

type loct struct {
	i, j int
}

func randLoct() loct {
	x := rand.Int() % 10000
	y := rand.Int() % 10000
	return loct{x % 20, y % 20}
}

var (
	area      = [20][20]byte{}
	food      bool
	direction byte
	head      loct
	tail      loct
	size      int
	headChar  = byte('#')
	scoreLoct = loct{22, 0}
	tipLoct   = loct{22, 2}
)

func draw(p loct, c byte) {
	C.gotoxy(C.int(toX(p.i)), C.int(toY(p.j)))
	fmt.Fprintf(os.Stdout, "%c", c)
}

func setScore(s int) {
	C.gotoxy(C.int(toX(scoreLoct.i)), C.int(toY(scoreLoct.j)))
	fmt.Fprintf(os.Stdout, "score:%d", s)
}

func setTip() {
	C.gotoxy(C.int(toX(tipLoct.i)), C.int(toY(tipLoct.j)))
	fmt.Fprint(os.Stdout, "Operation Tip: Press the arrow keys to change the moving direction")
}

func toX(i int) int {
	return i*2 + 4
}
func toY(j int) int {
	return j + 2
}

func init() {
	head, tail = loct{4, 4}, loct{4, 4}
	area[4][4] = 'H'
	direction, size = 'R', 1

	draw(head, headChar)

	setScore(0)
	setTip()

	rand.Seed(int64(time.Now().Unix()))

	C.gotoxy(0, 0)

	fmt.Fprintln(os.Stdout, `
  +-----------------------------------------+
  |                                         |
  |                                         |
  |                                         |
  |                                         |
  |                                         |
  |                                         |
  |                                         |
  |                                         |
  |                                         |
  |                                         |
  |                                         |
  |                                         |
  |                                         |
  |                                         |
  |                                         |
  |                                         |
  |                                         |
  |                                         |
  |                                         |
  |                                         |
  +-----------------------------------------+
`)
}

func main() {
	go onKeyboardEvent()

	for {
		time.Sleep(time.Millisecond * 400)

		if direction == 'P' {
			continue
		}

		if !food {
			for {

				randLoct := randLoct()
				if area[randLoct.i][randLoct.j] == 0 {
					area[randLoct.i][randLoct.j] = 'F'
					draw(randLoct, 'o')
					food = true
					break
				}
			}

		}

		area[head.i][head.j] = direction
		oldHead := head

		switch direction {
		case 'U':
			head.j--
		case 'L':
			head.i--
		case 'R':
			head.i++
		case 'D':
			head.j++
		}

		if head.i < 0 || head.i >= 20 || head.j < 0 || head.j >= 20 {
			dead()
			break
		}

		headVal := area[head.i][head.j]
		if headVal == 'F' {
			food = false

			draw(oldHead, '*')
			draw(head, headChar)

			size++
			setScore(size - 1)
		} else if headVal == 0 {
			draw(oldHead, '*')
			draw(head, headChar)

			dir := area[tail.i][tail.j]

			area[tail.i][tail.j] = 0
			draw(tail, ' ')

			switch dir {
			case 'U':
				tail.j--
			case 'L':
				tail.i--
			case 'R':
				tail.i++
			case 'D':
				tail.j++
			}
		} else {
			dead()

			break
		}

	}

	time.Sleep(60 * time.Second)
}

func onKeyboardEvent() {
	for {
		switch byte(C.onKeyboard()) {
		case 72:
			if direction == 'D' {
				continue
			}
			direction = 'U'
		case 75:
			if direction == 'R' {
				continue
			}
			direction = 'L'
		case 77:
			if direction == 'L' {
				continue
			}
			direction = 'R'
		case 80:
			if direction == 'U' {
				continue
			}
			direction = 'D'
		case 32:
			direction = 'P'
		}
	}
}

func dead() {
	C.gotoxy(0, 23)
	fmt.Fprintln(os.Stdout, "Game over!")
	//fmt.Fprintln(os.Stdout, "Press Enter key to play again; press Backspace key to exit:")s
	C.gotoxy(0, 0)
}
