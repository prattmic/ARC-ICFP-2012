package main

import (
        "os"
        "bufio"
        "strconv"
        "regexp"
        "fmt"
)

type Map        [][]byte
type Coord      [2]int
type Robot struct {
    coord       Coord
    waterproof  int
    moves       int
    watermoves  int
    lambda      int
}
type Lift struct {
    coord       Coord
    open        bool
}
type Mine struct {
    layout      Map
    robot       Robot
    lambda      []Coord
    rocks       []Coord
    lift        Lift
    water       int
    flooding    int
}

var RoboChar    byte = 'R'
var RockChar    byte = '*'
var WallChar    byte = '#'
var LambdaChar  byte = '\\'
var EarthChar   byte = '.'
var EmptyChar   byte = ' '
var CLiftChar   byte = 'L'
var OLiftChar   byte = 'O'

func main() {
    mine := new(Mine)
    err := mine.FromFile("maps/contest1.map", 100)

    if err != nil {
        fmt.Printf("Error: %s\n", err)
    }

    for i := range mine.layout {
        fmt.Println(string(mine.layout[i]))
    }
    fmt.Printf("Water: %d\n", mine.water)
    fmt.Printf("Flooding: %d\n", mine.flooding)
    fmt.Printf("Waterproof: %d\n", mine.robot.waterproof)

    mine.ParseLayout()
    fmt.Printf("\nMine struct:\n%+v\n\n", mine)

    fmt.Printf("Moving left is: %t\n", mine.validMove(Coord{mine.robot.coord[0], mine.robot.coord[1]-1}))
    fmt.Printf("Moving down is: %t\n", mine.validMove(Coord{mine.robot.coord[0]+1, mine.robot.coord[1]}))


    mine.update(Coord{2,4})
for i := range mine.layout {
        fmt.Println(string(mine.layout[i]))
    }



    serve(mine)
}

func (mine *Mine) ParseLayout() {
    mine.lambda = make([]Coord, 0, 100)
    mine.rocks = make([]Coord, 0, 100)

    for i := range mine.layout {
        for j := range mine.layout[i] {
            if mine.layout[i][j] == LambdaChar {
                mine.lambda = append(mine.lambda, Coord{i,j})
            } else if mine.layout[i][j] == RockChar {
                mine.rocks = append(mine.rocks, Coord{i,j})
            } else if mine.layout[i][j] == CLiftChar {
                mine.lift.coord = Coord{i,j}
                mine.lift.open = false
            } else if mine.layout[i][j] == RoboChar {
                mine.robot.coord = Coord{i,j}
            }
        }
    }
}

func (mine *Mine) update(move Coord) {
    updated := make([][]byte, len(mine.layout))

    mine.robot.coord = move
    for i := range mine.layout {
        updated[i] = make([]byte, len(mine.layout[i]))
	
        for j := range mine.layout[i] {
            if i==move[0] && j==move[1] {
		updated[i][j] = RoboChar
                mine.robot.lambda++
            } else if mine.layout[i][j] == RoboChar && (i==move[0] || j==move[1]) {
                updated[i][j] = EmptyChar
            } else if mine.layout[i][j] == EmptyChar{
		updated[i][j] = EmptyChar
            } else if mine.layout[i][j] == LambdaChar{
                updated[i][j] = LambdaChar
            } else if mine.layout[i][j] == EarthChar{
                updated[i][j] = EarthChar
            } else if mine.layout[i][j] == RockChar{
                //Rock logic goes here                
                updated[i][j] = RockChar
            } else if mine.layout[i][j] == WallChar {
                updated[i][j] = WallChar
            } else if mine.layout[i][j] == OLiftChar {
                updated[i][j] = OLiftChar
            } else if mine.layout[i][j] == CLiftChar {
                updated[i][j] = CLiftChar
            } else if mine.layout[i][j] == EarthChar {
                updated[i][j] = EarthChar
            }
        }
    }
    mine.layout = updated 
}

func serve(mine *Mine) {
    r := bufio.NewReaderSize(os.Stdin, 64)
    
    var err error = nil

    for err == nil {
        char, err := r.ReadByte()

        if char == 'L' {
            fmt.Println(mine.validMove(Coord{mine.robot.coord[0], mine.robot.coord[1]-1}))
        } else if char == 'R' {
            fmt.Println(mine.validMove(Coord{mine.robot.coord[0], mine.robot.coord[1]+1}))
        } else if char == 'U' {
            fmt.Println(mine.validMove(Coord{mine.robot.coord[0]-1, mine.robot.coord[1]}))
        } else if char == 'D' {
            fmt.Println(mine.validMove(Coord{mine.robot.coord[0]+1, mine.robot.coord[1]}))
        }

        _ = err
    }
}

func (mine *Mine) validMove(move Coord) bool {
    y := Abs(mine.robot.coord[0]-move[0])
    x := Abs(mine.robot.coord[1]-move[1])
    tile := mine.layout[move[0]][move[1]]

    if x != 0 && y != 0 {
        return false
    } else if x > 1 || y > 1 {
        return false
    } else if tile == EmptyChar || tile == EarthChar || tile == LambdaChar || tile == OLiftChar {
        return true
    }

    return false
}

func findSubmatch(reg string, line string) []string {
    re, err := regexp.Compile(reg)
    if err != nil {
        fmt.Printf("Error: %s", err)
        return nil
    }

    return re.FindStringSubmatch(line)
}

func (mine *Mine) FromFile(name string, capacity uint32) (err error) {
    file, err := os.Open(name)
    if err != nil {
        return err
    }
    fileinfo, err := file.Stat()
    
    r := bufio.NewReaderSize(file, int(fileinfo.Size()))

    data := make([][]byte, 0, capacity)

    mine.water = 0
    mine.flooding = 0
    mine.robot.waterproof = 10
    mine.robot.lambda = 0

    i := 0
    for ; ; i++ {
        line, _, err := r.ReadLine()
        if err != nil {
            break
        }

        // Blank
        if matched, _ := regexp.Match("^(?:\\s|\\n)*$", line); matched == true {
            continue
        } else if match := findSubmatch("Water\\s+([0-9]+)", string(line)); match != nil && len(match) == 2 {
            mine.water, _ = strconv.Atoi(match[1])
        } else if match := findSubmatch("Flooding\\s+([0-9]+)", string(line)); match != nil && len(match) == 2 {
            mine.flooding, _ = strconv.Atoi(match[1])
        } else if match := findSubmatch("Waterproof\\s+([0-9]+)", string(line)); match != nil && len(match) == 2 {
            mine.robot.waterproof, _ = strconv.Atoi(match[1])
        } else {
            data = append(data, line)
        }
    }

    mine.layout = data

    return nil
}

func Abs(n int) int {
    if n < 0 {
        return -n
    }

    return n
}
