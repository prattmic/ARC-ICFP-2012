package main

import (
        "os"
        "bufio"
        "fmt"
)

type Map [][]byte
type Coord [2]int
type Robot struct {
    coord Coord
    mine Map
}

var RoboChar byte = 'R'
var RockChar byte = '*'
var WallChar byte = '#'
var LambdaChar byte = '\\'
var EarthChar byte = '.'
var EmptyChar byte = ' '
var CLiftChar byte = 'L'
var OLiftChar byte = 'O'

func main() {
    mine, err := MapFromFile("maps/contest1.map", 100)
    if err != nil {
        fmt.Printf("Error: %s\n", err)
    }

    for i := range mine {
        fmt.Println(string(mine[i]))
    }

    robot := new(Robot)
    robot.mine = mine

    robot.coord = mine.currentLocation()
    fmt.Printf("You are at %d\n", robot.coord)

    fmt.Printf("Moving left is: %t\n", robot.validMove(Coord{robot.coord[0], robot.coord[1]-1}))
    fmt.Printf("Moving down is: %t\n", robot.validMove(Coord{robot.coord[0]+1, robot.coord[1]}))

    serve(robot)
}

func serve(robot *Robot) {
    r := bufio.NewReaderSize(os.Stdin, 64)
    
    var err error = nil

    for err == nil {
        char, err := r.ReadByte()

        if char == 'L' {
            fmt.Println(robot.validMove(Coord{robot.coord[0], robot.coord[1]-1}))
        } else if char == 'R' {
            fmt.Println(robot.validMove(Coord{robot.coord[0], robot.coord[1]+1}))
        } else if char == 'U' {
            fmt.Println(robot.validMove(Coord{robot.coord[0]-1, robot.coord[1]}))
        } else if char == 'D' {
            fmt.Println(robot.validMove(Coord{robot.coord[0]+1, robot.coord[1]}))
        }

        _ = err
    }
}

func (robo *Robot) validMove(move Coord) bool {
    y := Abs(robo.coord[0]-move[0])
    x := Abs(robo.coord[1]-move[1])
    tile := robo.mine[move[0]][move[1]]

    if x != 0 && y != 0 {
        return false
    } else if x > 1 || y > 1 {
        return false
    } else if tile == EmptyChar || tile == EarthChar || tile == LambdaChar || tile == OLiftChar {
        return true
    }

    return false
}

func (mine Map) currentLocation() (Coord) {
    for i := range mine {
        for j := range mine[i] {
            if mine[i][j] == 'R' {
                return Coord{i,j}
            }
        }
    }

    return Coord{-1,-1}
}

func MapFromFile(name string, capacity uint32) (mine Map, err error) {
    file, err := os.Open(name)
    if err != nil {
        return nil, err
    }
    fileinfo, err := file.Stat()
    
    r := bufio.NewReaderSize(file, int(fileinfo.Size()))

    data := make([][]byte, 0, capacity)

    i := 0
    for ; ; i++ {
        line, _, err := r.ReadLine()
        if err != nil {
            break
        }
        data = append(data, line)
    }

    return data, nil
}

func Abs(n int) int {
    if n < 0 {
        return -n
    }

    return n
}
