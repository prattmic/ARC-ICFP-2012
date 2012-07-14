package icfp

import (
        "os"
        "bufio"
        "strconv"
        "regexp"
)

type Map        [][]byte
type Coord      [2]int
type Robot struct {
    Coord       Coord
    Waterproof  int
    Moves       int
    Watermoves  int
}
type Lift struct {
    Coord       Coord
    Open        bool
}
type Mine struct {
    Layout      Map
    Robot       Robot
    Lambda      []Coord
    Rocks       []Coord
    Lift        Lift
    Water       int
    Flooding    int
}

var RoboChar    byte = 'R'
var RockChar    byte = '*'
var WallChar    byte = '#'
var LambdaChar  byte = '\\'
var EarthChar   byte = '.'
var EmptyChar   byte = ' '
var CLiftChar   byte = 'L'
var OLiftChar   byte = 'O'

func (mine *Mine) ParseLayout() {
    mine.Lambda = make([]Coord, 0, 100)
    mine.Rocks = make([]Coord, 0, 100)

    for i := range mine.Layout {
        for j := range mine.Layout[i] {
            if mine.Layout[i][j] == LambdaChar {
                mine.Lambda = append(mine.Lambda, Coord{i,j})
            } else if mine.Layout[i][j] == RockChar {
                mine.Rocks = append(mine.Rocks, Coord{i,j})
            } else if mine.Layout[i][j] == CLiftChar {
                mine.Lift.Coord = Coord{i,j}
                mine.Lift.Open = false
            } else if mine.Layout[i][j] == RoboChar {
                mine.Robot.Coord = Coord{i,j}
            }
        }
    }
}

func (mine *Mine) ValidMove(move Coord) bool {
    y := Abs(mine.Robot.Coord[0]-move[0])
    x := Abs(mine.Robot.Coord[1]-move[1])
    tile := mine.Layout[move[0]][move[1]]

    if x != 0 && y != 0 {
        return false
    } else if x > 1 || y > 1 {
        return false
    } else if tile == EmptyChar || tile == EarthChar || tile == LambdaChar || tile == OLiftChar {
        return true
    }

    return false
}

func (mine *Mine) FromFile(name string, capacity uint32) (err error) {
    file, err := os.Open(name)
    if err != nil {
        return err
    }
    fileinfo, err := file.Stat()
    
    r := bufio.NewReaderSize(file, int(fileinfo.Size()))

    data := make([][]byte, 0, capacity)

    mine.Water = 0
    mine.Flooding = 0
    mine.Robot.Waterproof = 10

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
            mine.Water, _ = strconv.Atoi(match[1])
        } else if match := findSubmatch("Flooding\\s+([0-9]+)", string(line)); match != nil && len(match) == 2 {
            mine.Flooding, _ = strconv.Atoi(match[1])
        } else if match := findSubmatch("Waterproof\\s+([0-9]+)", string(line)); match != nil && len(match) == 2 {
            mine.Robot.Waterproof, _ = strconv.Atoi(match[1])
        } else {
            data = append(data, line)
        }
    }

    mine.Layout = data

    return nil
}
