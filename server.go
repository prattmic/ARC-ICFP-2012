package main

import (
        "os"
        "bufio"
        "fmt"
        "./icfp"
)

func main() {
    mine := new(icfp.Mine)
    err := mine.FromFile("maps/beard1_test.map", 100, true)

    if err != nil {
        fmt.Printf("Error: %s\n", err)
    }

    mine.ParseLayout()

    serve(mine)

    mine.Print()
    fmt.Println(mine.Score())
}

func serve(mine *icfp.Mine) {
    r := bufio.NewReaderSize(os.Stdin, 64)

    for char, err := r.ReadByte() ; err == nil ; char, err = r.ReadByte() {
        switch char {
        case 'L', 'l':
            move := icfp.Coord{mine.Robot.Coord[0], mine.Robot.Coord[1]-1}
            if !mine.ValidMove(move, false) {
                //Wait
                move = icfp.Coord{mine.Robot.Coord[0], mine.Robot.Coord[1]}
            }
            mine.Update(move, false)
            //mine.Print()
        case 'R', 'r':
            move := icfp.Coord{mine.Robot.Coord[0], mine.Robot.Coord[1]+1}
            if !mine.ValidMove(move, false) {
                //Wait
                move = icfp.Coord{mine.Robot.Coord[0], mine.Robot.Coord[1]}
            }
            mine.Update(move, false)
            //mine.Print()
        case 'U', 'u':
            move := icfp.Coord{mine.Robot.Coord[0]-1, mine.Robot.Coord[1]}
            if !mine.ValidMove(move, false) {
                //Wait
                move = icfp.Coord{mine.Robot.Coord[0], mine.Robot.Coord[1]}
            }
            mine.Update(move, false)
            //mine.Print()
        case 'D', 'd':
            move := icfp.Coord{mine.Robot.Coord[0]+1, mine.Robot.Coord[1]}
            if !mine.ValidMove(move, false) {
                //Wait
                move = icfp.Coord{mine.Robot.Coord[0], mine.Robot.Coord[1]}
            }
            mine.Update(move, false)
            //mine.Print()
        case 'W', 'w':
            move := icfp.Coord{mine.Robot.Coord[0], mine.Robot.Coord[1]}
            if !mine.ValidMove(move, false) {
                //Wait
                move = icfp.Coord{mine.Robot.Coord[0], mine.Robot.Coord[1]}
            }
            mine.Update(move, false)
            //mine.Print()
        case 'S', 's':
            move := icfp.Coord{mine.Robot.Coord[0], mine.Robot.Coord[1]}
            if !mine.ValidMove(move, false) {
                //Wait
                move = icfp.Coord{mine.Robot.Coord[0], mine.Robot.Coord[1]}
            }
            mine.Update(move, false)
            //mine.Print()
        }
    }
}
