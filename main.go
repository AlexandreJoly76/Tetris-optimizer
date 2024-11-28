package main

import (
	"fmt"
	"io"
	"math"
	"os"
	"strings"
	"time"
)

type Tetrominos struct {
	Shape   []Coords
	IsValid bool
	Symbole string
}

type Coords struct {
	X int
	Y int
}

const Emoji = "🌏💥🌛🌸🪐🌊🔥💧💨🍇🍎🥭#####################################"
const VoidSymbole = "🌑"

var EmojiRune = []rune(Emoji)

func main() {
	// Check if the user has entered a file name
	if len(os.Args) > 2 {
		fmt.Println("Please enter only one file name")
		os.Exit(0)
	}
	FileName := os.Args[1]
	file := CheckFile(FileName)
	if file != nil {
		defer file.Close()
	}
	// Read the file
	FileContent := ReadFile(file)
	ArrayOfTetrominos := CheckTetrominoSent(FileContent)
	for i := 0; i < len(ArrayOfTetrominos); i++ {
		if !ArrayOfTetrominos[i].IsValid {
			fmt.Println("Error")
			os.Exit(0)
		}
	}
	execTime := time.Now()
	grid := Solve(ArrayOfTetrominos)
	PrintPropres(grid)
	fmt.Println("Execution time: ", time.Since(execTime))
}

func CheckFile(FileName string) *os.File {
	// Check if the file exists
	if len(os.Args) == 1 {
		fmt.Println("Please enter a file name")
	}
	file, err := os.Open("./File/" + FileName)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(0)
	}
	return file
}

func ReadFile(file *os.File) string {
	FileContent := ""
	// Read the file
	buf := make([]byte, 1024)
	for {
		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			fmt.Println("Error: ", err)
			os.Exit(0)
		}
		if n == 0 {
			break
		}
		FileContent += string(buf[:n])
	}
	return FileContent
}

func ValidTetrominos() []Tetrominos {
	TetrominosValid := []Tetrominos{}
	filecontent, _ := os.ReadFile("./File/goodtetrominos.txt")
	TertrominosSplit := strings.Split(string(filecontent), "\n\n")
	for i := 0; i < len(TertrominosSplit); i++ {
		TetrominoSeparated := Tetrominos{}
		TetrominoSeparated.IsValid = true
		TetrominoSeparated.Shape = GetCoords(TertrominosSplit[i])
		TetrominosValid = append(TetrominosValid, TetrominoSeparated)

	}
	return TetrominosValid
}

func GetCoords(str string) []Coords {
	// Get the coordinates of the tetrominos
	coords := []Coords{}
	Lines := strings.Split(str, "\n")
	y := -1
	x := -1
	for i := 0; i < len(Lines); i++ {
		x = -1
		y++
		for j := 0; j < len(Lines[i]); j++ {
			x++
			if Lines[i][j] == '#' {
				coords = append(coords, Coords{x, y})
			}
		}
	}
	return coords
}

func CheckTetrominoSent(str string) []Tetrominos {
	// Check if the tetrominos are valid
	ArrayOfTetrominos := []Tetrominos{}
	TetrominoSeparated := strings.Split(str, "\n\n")
	for i := 0; i < len(TetrominoSeparated); i++ {
		Tetromino := Tetrominos{}
		Tetromino.Shape = GetCoords(TetrominoSeparated[i])
		if len(Tetromino.Shape) > 0 {
			ShiftTetrominos(&Tetromino)
			CheckIsValidTetrominos(&Tetromino)
			ArrayOfTetrominos = append(ArrayOfTetrominos, Tetromino)
		}

	}
	return ArrayOfTetrominos
}

func ShiftTetrominos(Tetromino *Tetrominos) {
	// Shift the tetrominos to the left and top
	minX := 100
	minY := 100
	for _, coords := range Tetromino.Shape {
		if coords.X < minX {
			minX = coords.X
		}
		if coords.Y < minY {
			minY = coords.Y
		}

	}
	for minX > 0 {
		for i := 0; i < len(Tetromino.Shape); i++ {
			Tetromino.Shape[i].X--
		}
		minX--
	}
	for minY > 0 {
		for i := 0; i < len(Tetromino.Shape); i++ {
			Tetromino.Shape[i].Y--
		}
		minY--
	}
}

func CheckIsValidTetrominos(tetrominos *Tetrominos) {
	// Check if the tetrominos are valid
	ValideTetrominos := ValidTetrominos()
	for _, ValidTetromino := range ValideTetrominos {
		tetrominos.IsValid = true
		for i, coords := range tetrominos.Shape {
			if coords != ValidTetromino.Shape[i] {
				tetrominos.IsValid = false
				break
			}

		}
		if tetrominos.IsValid {
			break
		}
	}
}

func PlacedTetrominos(grid [][]string, tetrominos []Tetrominos, index int) bool {
	//place the tetrominos on the grid
	if index == len(tetrominos) {
		return true
	}
	for y := 0; y < len(grid); y++ {
		for x := 0; x < len(grid[y]); x++ {
			if CanPlace(grid, tetrominos[index], x, y) {
				tetrominos[index].Symbole = string(EmojiRune[index])
				Place(grid, tetrominos[index], x, y)
				if PlacedTetrominos(grid, tetrominos, index+1) {
					return true
				}
				Remove(grid, tetrominos[index], x, y)
			}

		}
	}
	return false
}

func Solve(tetrominos []Tetrominos) [][]string {
	// Make the grid size the smallest possible.
	GridSize := int(math.Ceil(math.Sqrt(float64(len(tetrominos) * 4))))
	for {
		grid := make([][]string, GridSize)
		for Line := range grid {
			grid[Line] = make([]string, GridSize)
			for char := range grid[Line] {
				grid[Line][char] = VoidSymbole
			}
		}
		if PlacedTetrominos(grid, tetrominos, 0) {
			return grid
		}
		GridSize++
	}
}

func CanPlace(grid [][]string, tetrominos Tetrominos, x, y int) bool {
	//Find if the tetrominos can be placed on the grid
	for _, coords := range tetrominos.Shape {
		if coords.X+x >= len(grid) || coords.Y+y >= len(grid) || grid[coords.Y+y][coords.X+x] != VoidSymbole {
			return false
		}
	}
	return true
}

func Place(grid [][]string, tetrominos Tetrominos, x, y int) {
	//Place the tetrominos on the grid
	for _, coords := range tetrominos.Shape {
		grid[coords.Y+y][coords.X+x] = tetrominos.Symbole
	}
}

func Remove(grid [][]string, tetrominos Tetrominos, x, y int) {
	// Remove the tetrominos from the grid for testing
	for _, coords := range tetrominos.Shape {
		grid[coords.Y+y][coords.X+x] = VoidSymbole
	}
}

func PrintPropres(grid [][]string) {
	// Print the grid
	var GridLines []string
	for _, Line := range grid {
		GridLines = append(GridLines, strings.Join(Line, ""))
	}
	fmt.Println(strings.Join(GridLines, "\n"))
}
