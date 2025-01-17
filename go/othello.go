package othello
//package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"io/ioutil"
	"math"
	"net/http"

	//"strconv"
)



func main() {
	fakeBoard := makeFakeBoard()
	moves := fakeBoard.ValidMoves()
	if len(moves) < 1 {
		fmt.Println("PASS")
		return
	}
	fmt.Println(getBestMove(*fakeBoard))
}




func init() {
	http.HandleFunc("/", getMove)
}

type Game struct {
	Board Board `json:board`
}

// Provide a generic handler for move requests. If no board state is
// specified then a simple HTML form is provided to let users paste
// JSON state (which can be copy-pasted from a game running on
// http://step-reversi.appspot.com/ ).
func getMove(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	var js []byte
	defer r.Body.Close()
	js, _ = ioutil.ReadAll(r.Body)
	if len(js) < 1 {
		js = []byte(r.FormValue("json"))
	}
	if len(js) < 1 {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, `
		<body>
			<form method=get>
				Paste JSON here:<p/><textarea name=json cols=80 rows=24></textarea>
				<p/><input type=submit>
			</form>
		</body>`)
		return
	}
	var game Game
	err := json.Unmarshal(js, &game)
	if err != nil {
		fmt.Fprintf(w, "invalid json %v? %v", string(js), err)
		return
	}
	board := game.Board
	log.Infof(ctx, "got board: %v", board)
	moves := board.ValidMoves()
	if len(moves) < 1 {
		fmt.Fprintf(w, "PASS")
		return
	}
	// NOTE TO STUDENTS: This next line is the main line you'll want to
	// change.  Right now this is just picking a random move out of the
	// list of possible moves, but you'll want to make this choose a
	// better move (probably using some game tree traversal algorithm
	// like MinMax).

	//move := moves[rand.Intn(len(moves))] // random
	//move := greedy(board, moves)			// greedy
	move := getBestMove(board)
	fmt.Fprintf(w, "[%d,%d]", move.Where[0], move.Where[1])
}

// Just choose a move that gets most pieces.
func greedy(b Board, moves []Move) Move{
	max := 0
	var best Move
	var cnt int
	for _, move := range moves {
		nextBoard, _ := b.After(move)
		white, black := nextBoard.CountColors()
		switch move.As {
		case White:
			cnt = white
		case Black:
			cnt = black
		}
		if cnt > max {
			max = cnt
			best = move
		}
	}
	return best
}


// Depth is always set to 6 for now.
func getBestMove(b Board) Move {
	me := b.Next
	//_, bestMove := b.ScoreMM(6, me, b.CountEmpty())
	_, bestMove := b.ScoreAB(6, me, b.CountEmpty(), -math.MaxInt32, math.MaxInt32)

	/*
	best := -math.MaxInt32
	var bestMove Move
	for _, move := range b.ValidMoves() {
		if best <= b.ScoreAB(6, me, b.CountEmpty(), -math.MaxInt32, math.MaxInt32) {
			bestMove = move
		}
		fmt.Println("-----")
	}
	*/
	return bestMove
}


type Piece int8

const (
	Empty Piece = iota
	Black Piece = iota
	White Piece = iota

	// Red/Blue are aliases for Black/White
	Red  = Black
	Blue = White
)

func (p Piece) Opposite() Piece {
	switch p {
	case White:
		return Black
	case Black:
		return White
	default:
		return Empty
	}
}

type Board struct {
	// Pieces says what pieces are where.
	Pieces [8][8]Piece
	// Next says what the color of the next piece played must be.
	Next Piece
}

// Scoring by Mini-Max
/*
func (b Board) ScoreMM(depth int, myPiece Piece, emptySpace int) (int, Move) {

	if depth < 1 || emptySpace < 1 {
		//fmt.Println(b.String())

		// Not time left to recurse, just evaluate this board and return.

		if emptySpace < 12 {
			return b.EvalByPieceNum(myPiece), Move{}
		} else {
			return b.EvalByScore(myPiece), Move{}
		}


		return b.EvalByScore(myPiece), Move{}

	}

	bestScore := b.Player().MinScore(myPiece)
	var bestMove Move

	// Search each valid move and score them, choose the best one.
	// If the move is my turn, choose the maximum score. Otherwise, choose the minimum.
	for _, move := range b.ValidMoves() {
		nextBoard, _ := b.Clone().Exec(move)
		score, _ := (*nextBoard).ScoreMM(depth - 1, myPiece, emptySpace - 1)

		switch b.Player() {
		case myPiece:
			if score >= bestScore {
				bestScore = score
				bestMove = move
			}
		default:
			if score <= bestScore {
				bestScore = score
				bestMove = move
			}
		}
	}
	fmt.Print("player: ")
	fmt.Print(b.Player())
	fmt.Print(" best: ")
	fmt.Println(bestScore)
	fmt.Print("best move: ")
	fmt.Println(bestMove)
	return bestScore, bestMove
}
 */


// Scoring by Alpha-Beta
func (b Board) ScoreAB(depth int, myPiece Piece, emptySpace int, alpha int, beta int) (int, Move) {

	// Not recurse, just evaluate this board and return.
	if depth < 1 || emptySpace < 1 {
		//fmt.Print(b.String())
		//fmt.Println(emptySpace)

		if emptySpace < 10 {
			return b.EvalByPieceNum(myPiece), Move{}
		} else {
			return b.EvalByScore(myPiece, emptySpace), Move{}
		}


		// The number of available place for the next player.
		placeable := len(b.ValidMoves())
		if b.Next != myPiece {
			placeable *= -1
		}

		// Score of the board from the score table
		score := b.EvalByScore(myPiece, emptySpace)

		return score + placeable * 2, Move{}
	}

	// Search each valid move and score them, choose the best one.
	// If the move is my turn, choose the maximum score. Otherwise, choose the minimum.
	var best int
	var bestMove Move

	if b.Next == myPiece {
		best = -math.MaxInt32
		for _, move := range b.ValidMoves() {
			nextBoard, _ := b.Clone().Exec(move)
			score, _ := (*nextBoard).ScoreAB(depth - 1, myPiece, emptySpace - 1, alpha, beta)
			if score >= best {
				best = score
				bestMove = move
			}
			// alpha = max(alpha, best)
			if alpha < best {
				alpha = best
			}
			// beta cut-off
			if alpha >= beta {
				break
			}
		}
		fmt.Print("player: ")
		fmt.Print(b.Player())
		fmt.Print(" max: ")
		fmt.Println(best)
	} else {
		best = math.MaxInt32
		for _, move := range b.ValidMoves() {
			nextBoard, _ := b.Clone().Exec(move)
			score, _ := (*nextBoard).ScoreAB(depth - 1, myPiece, emptySpace - 1, alpha, beta)
			if score <= best {
				best = score
				bestMove = move
			}
			// beta = min(beta, best)
			if best < beta {
				beta = best
			}
			// alpha cut-off
			if alpha >= beta {
				break
			}
		}
		fmt.Print("player: ")
		fmt.Print(b.Player())
		fmt.Print(" min: ")
		fmt.Println(best)
	}
	return best, bestMove
}


func (b Board) EvalByPieceNum (myPiece Piece) int {
	white, black := b.CountColors()
	val := 0
	switch myPiece {
	case White:
		val = white
	case Black:
		val = black
	}

	fmt.Print(" count: ")
	fmt.Println(val)
	return val
}


func (b Board) EvalByScore(myPiece Piece, emptySpace int) int {
	wscore, bscore := b.ScoreColors(emptySpace)
	val := 0

	switch myPiece {
	case White:
		val = wscore
	case Black:
		val = bscore
	}

	fmt.Print(" score: ")
	fmt.Println(val)
	return val
}


// Count the number of white and black pieces.
func (b Board) CountColors() (int, int) {
	white := 0
	black := 0
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			switch b.Pieces[i][j] {
			case White:
				white++
			case Black:
				black++
			}
		}
	}
	return white, black
}

// Count the number of empty spaces.
func (b Board) CountEmpty() int {
	cnt := 0
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			if b.Pieces[i][j] == Empty {
				cnt++
			}
		}
	}
	return cnt
}


// Return the score of a position.
/*
func getScore1(p Position) int{
	scoreTable := [][]int {
		{ 383, -15, 5, 4, 4, 5, -15, 383 },
		{ -15, -112, -3, 0, 0, -3, -112, -15 },
		{ 5, -3, -2, 5, 5, -2, -3, 5 },
		{ 4, 0, 5, 10, 10, 5, 0, 4 },
		{ 4, 0, 5, 10, 10, 5, 0, 4 },
		{ 5, -3, -2, 5, 5, -2, -3, 5 },
		{ -15, -112, -3, 0, 0, -3, -112, -15 },
		{ 383, -15, 5, 4, 4, 5, -15, 383 },
	}
	return scoreTable[p[0]][p[1]]
}

func getScore2(p Position) int{
	scoreTable := [][]int {
		{ 383, -15, -2, -4, -4, -2, -15, 383 },
		{ -15, -112, -3, 0, 0, -3, -112, -15 },
		{ -2, -3, -2, 5, 5, -2, -3, -2 },
		{ -4, 0, 5, 10, 10, 5, 0, -4 },
		{ -4, 0, 5, 10, 10, 5, 0, -4 },
		{ -2, -3, -2, 5, 5, -2, -3, -2 },
		{ -15, -112, -3, 0, 0, -3, -112, -15 },
		{ 383, -15, -2, -4, -4, -2, -15, 383 },
	}
	return scoreTable[p[0]][p[1]]
}
*/


func getScore1(p Position) int{
	scoreTable := [][]int {
		{ 50,  -12,  5,  4,  4,  5, -12,  50 },
		{ -12, -20, -3, -3, -3, -3, -20, -12 },
		{ 5,   -3,   0, -1, -1,  0,  -3,   5 },
		{ 4,   -3,  -1, -1, -1, -1,  -3,   4 },
		{ 4,   -3,  -1, -1, -1, -1,  -3,   4 },
		{ 5,   -3,   0, -1, -1,  0,  -3,   5 },
		{ -12, -20, -3, -3, -3, -3, -20, -12 },
		{ 50,  -12,  5,  4,  4,  5, -12,  50 },
	}
	return scoreTable[p[0]][p[1]]
}


func getScore2(p Position) int{
	scoreTable := [][]int {
		{  50, -12,  0, -1, -1,  0, -12,  50 },
		{ -12, -15, -3, -3, -3, -3, -15, -12 },
		{   0,  -3,  0, -1, -1,  0,  -3,   0 },
		{  -1,  -3, -1, -1, -1, -1,  -3,  -1 },
		{  -1,  -3, -1, -1, -1, -1,  -3,  -1 },
		{   0,  -3,  0, -1, -1,  0,  -3,   0 },
		{ -12, -15, -3, -3, -3, -3, -15, -12 },
		{  50, -12,  0, -1, -1,  0, -12,  50 },
	}
	return scoreTable[p[0]][p[1]]
}


// Sum up the score of black in a situation.
func (b Board) ScoreColors(emptySpace int) (int, int) {
	wscore := 0
	bscore := 0

	if (emptySpace > 48) {
		for i := 0; i < 8; i++ {
			for j := 0; j < 8; j++ {
				switch b.Pieces[i][j] {
				case White:
					wscore += getScore1(Position{i, j})
				case Black:
					bscore += getScore1(Position{i, j})

				}
			}
		}
	} else {
		for i := 0; i < 8; i++ {
			for j := 0; j < 8; j++ {
				switch b.Pieces[i][j] {
				case White:
					wscore += getScore2(Position{i, j})
				case Black:
					bscore += getScore2(Position{i, j})

				}
			}
		}
	}
	return wscore, bscore
}


// Return the player who should put on the board.
func (b Board) Player() Piece {
	return b.Next
}

// ??
func (p Piece) MinScore(myPiece Piece) int {
	if p == myPiece {
		return -math.MaxInt32
	} else {
		return math.MaxInt32
	}
}


func makeFakeBoard() *Board {
	b := Board{}

	b.Pieces = [8][8]Piece{
		{ Empty, Empty, Black, Black, Black, Black, Empty, White },
		{ Empty, Empty, Black, White, Black, Black, Empty, White },
		{ Black, Black, Black, White, Black, Black, Black, White },
		{ Black, Black, Black, White, Black, Black, Black, White },
		{ Black, Black, White, Black, Black, White, Black, White },
		{ Black, Black, Black, Black, Black, White, White, White },
		{ Empty, Black, Black, Black, Black, Black, Empty, White },
		{ Empty, Empty, White, White, White, White, White, White },
	}

	b.Next = Black

	return &b
}



/*

func makeFakeBoard() *Board {
	b := Board{}

	b.Pieces = [8][8]Piece{
		{ Empty, Empty, Empty, White, Empty, Empty, Empty, Empty },
		{ Empty, Empty, Empty, White, Empty, Empty, Empty, Empty },
		{ Empty, Empty, Empty, White, Empty, Empty, Empty, Empty },
		{ Empty, Empty, Empty, White, Empty, Empty, Empty, Empty },
		{ Empty, Empty, Empty, Black, Empty, Empty, Empty, Empty },
		{ Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty },
		{ Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty },
		{ Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty },
	}

	b.Next = Black

	return &b
}

 */


//-------------------------------

// Position represents a position on the othello board. Valid board
// coordinates are 1-8 (not 0-7)!
type Position [2]int

// Valid returns true iff this is a valid board position.
func (p Position) Valid() bool {
	ok := func(i int) bool { return 1 <= i && i <= 8 }
	return ok(p[0]) && ok(p[1])
}

// Pass returns true iff this move position represents a pass.
func (p Position) Pass() bool {
	return !p.Valid()
}

// Move describes a move on an Othello board.
type Move struct {
	// Where a piece is going to be placed. If Where is zeros, or
	// another invalid coordinate, it indicates a pass.
	Where Position
	// As is the player taking the player taking the turn.
	As Piece
}

// At returns a pointer to the piece at a given position.
func (b *Board) At(p Position) *Piece {
	return &b.Pieces[p[1]-1][p[0]-1]
}

// Get returns the piece at a given position.
func (b *Board) Get(p Position) Piece {
	return *b.At(p)
}

// Exec runs a move on a given Board, updating the given board, and
// returning it. Returns error if the move is illegal.
func (b *Board) Exec(m Move) (*Board, error) {
	if !m.Where.Pass() {
		if _, err := b.realMove(m); err != nil {
			return b, err
		}
	} else {
		// Attempting to pass.
		valid := b.ValidMoves()
		if len(valid) > 0 {
			return nil, fmt.Errorf("%v illegal move: there are valid moves available: %v", m, valid)
		}
	}
	b.Next = b.Next.Opposite()
	return b, nil
}

// Clone makes a new identical copy of an existing board and returns a
// pointer to it.
func (b *Board) Clone() *Board {
	clone := *b
	return &clone
}

// Returns the state of a new board after the given move. Returns an
// unchanged board and an error if the move is illegal.
func (b Board) After(m Move) (Board, error) {
	if _, err := b.Exec(m); err != nil {
		return b, err
	}
	return b, nil
}

// realMove executes a move that isn't a PASS. Use Exec instead to
// execute any move (include PASS moves).
func (b *Board) realMove(m Move) (*Board, error) {
	captures, err := b.tryMove(m)
	if err != nil {
		return nil, err
	}

	for _, p := range append(captures, m.Where) {
		*b.At(p) = m.As
	}
	return b, nil
}

type direction Position

var dirs []direction

func init() {
	for x := -1; x <= 1; x++ {
		for y := -1; y <= 1; y++ {
			if x == 0 && y == 0 {
				continue
			}
			dirs = append(dirs, direction{x, y})
		}
	}
}

// tryMove tries a non-PASS move without actually executing it.
// Returns the list of captures that would happen.
func (b *Board) tryMove(m Move) ([]Position, error) {
	if b.Get(m.Where) != Empty {
		return nil, fmt.Errorf("%v illegal move: %v is occupied by %v", m, m.Where, b.Get(m.Where))
	}

	var captures []Position
	for _, dir := range dirs {
		captures = append(captures, b.findCaptures(m, dir)...)
	}

	if len(captures) < 1 {
		return nil, fmt.Errorf("%v illegal move: no pieces were captured", m)
	}
	return captures, nil
}

func translate(p Position, d direction) Position {
	return Position{p[0] + d[0], p[1] + d[1]}
}

func (b *Board) findCaptures(m Move, dir direction) []Position {
	var caps []Position
	for p := m.Where; true; caps = append(caps, p) {
		p = translate(p, dir)
		if !p.Valid() {
			// End of board.
			return []Position{}
		}
		switch *b.At(p) {
		case m.As:
			return caps
		case Empty:
			return []Position{}
		}
	}
	panic("impossible")
}

// Returns a slice of valid moves for the given Board.
func (b *Board) ValidMoves() []Move {
	var moves []Move
	for y := 1; y <= 8; y++ {
		for x := 1; x <= 8; x++ {
			m := Move{Where: Position{x, y}, As: b.Next}
			_, err := b.tryMove(m)
			if err == nil {
				moves = append(moves, m)
			}
		}
	}
	return moves
}

// Converts a Board into a human-readable ASCII art diagram.
func (b Board) String() string {
	buf := &bytes.Buffer{}
	buf.WriteString("\n")
	buf.WriteString(" |ABCDEFGH|\n")
	buf.WriteString("-+--------+\n")
	for y := 0; y < 8; y++ {
		fmt.Fprintf(buf, "%d|", y+1)
		for x := 0; x < 8; x++ {
			p := b.Pieces[y][x]
			switch p {
			case Red:
				buf.WriteString("X")
			case Blue:
				buf.WriteString("O")
			default:
				buf.WriteString(" ")
			}
		}
		fmt.Fprintf(buf, "|%d\n", y+1)
	}
	buf.WriteString("-+--------+\n")
	buf.WriteString(" |ABCDEFGH|\n")
	return buf.String()
}
