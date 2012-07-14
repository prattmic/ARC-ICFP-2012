package main

import (
        "os"
        "bufio"
        //"strconv"
        //"regexp"
        "fmt"
        "./icfp"
)

func main() {
    mine := new(icfp.Mine)
    err := mine.FromFile("maps/contest1.map", 100)

    if err != nil {
        fmt.Printf("Error: %s\n", err)
    }

    for i := range mine.Layout {
        fmt.Println(string(mine.Layout[i]))
    }
    fmt.Printf("Water: %d\n", mine.Water)
    fmt.Printf("Flooding: %d\n", mine.Flooding)
    fmt.Printf("Waterproof: %d\n", mine.Robot.Waterproof)

    mine.ParseLayout()
    fmt.Printf("\nMine struct:\n%+v\n\n", mine)

    fmt.Printf("Moving left is: %t\n", mine.ValidMove(icfp.Coord{mine.Robot.Coord[0], mine.Robot.Coord[1]-1}))
    fmt.Printf("Moving down is: %t\n", mine.ValidMove(icfp.Coord{mine.Robot.Coord[0]+1, mine.Robot.Coord[1]}))


    mine.Update(icfp.Coord{2,3})
    for i := range mine.Layout {
        fmt.Println(string(mine.Layout[i]))
    }
    fmt.Printf("\nMine struct:\n%+v\n\n", mine)

    serve(mine)
}

func serve(mine *icfp.Mine) {
    r := bufio.NewReaderSize(os.Stdin, 64)
    
    var err error = nil

    for err == nil {
        char, err := r.ReadByte()

        if char == 'L' {
            fmt.Println(mine.ValidMove(icfp.Coord{mine.Robot.Coord[0], mine.Robot.Coord[1]-1}))
        } else if char == 'R' {
            fmt.Println(mine.ValidMove(icfp.Coord{mine.Robot.Coord[0], mine.Robot.Coord[1]+1}))
        } else if char == 'U' {
            fmt.Println(mine.ValidMove(icfp.Coord{mine.Robot.Coord[0]-1, mine.Robot.Coord[1]}))
        } else if char == 'D' {
            fmt.Println(mine.ValidMove(icfp.Coord{mine.Robot.Coord[0]+1, mine.Robot.Coord[1]}))
        }

        _ = err
    }
}
