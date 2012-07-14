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
    err := mine.FromFile("maps/pratttramp1.map", 100)

    if err != nil {
        fmt.Printf("Error: %s\n", err)
    }

    for i := range mine.Layout {
        fmt.Println(string(mine.Layout[i]))
    }
    mine.ParseLayout()
    fmt.Printf("Water: %d\n", mine.Water)
    fmt.Printf("Flooding: %d\n", mine.Flooding)
    fmt.Printf("Waterproof: %d\n", mine.Robot.Waterproof)
    fmt.Printf("Trampolines: %v\n", mine.Trampolines)

    fmt.Printf("\nMine struct:\n%+v\n\n", mine)

    fmt.Printf("Moving left allowed: %t\n", mine.ValidMove(icfp.Coord{mine.Robot.Coord[0], mine.Robot.Coord[1]-1}))
    fmt.Printf("Moving right allowed: %t\n", mine.ValidMove(icfp.Coord{mine.Robot.Coord[0], mine.Robot.Coord[1]+1}))
    fmt.Printf("Moving down allowed: %t\n", mine.ValidMove(icfp.Coord{mine.Robot.Coord[0]+1, mine.Robot.Coord[1]}))
    fmt.Printf("Moving up allowed: %t\n", mine.ValidMove(icfp.Coord{mine.Robot.Coord[0]-1, mine.Robot.Coord[1]}))


    mine.Update(icfp.Coord{2,3})
    for i := range mine.Layout {
        fmt.Println(string(mine.Layout[i]))
    }
    fmt.Printf("\nMine struct:\n%+v\n\n", mine)

    //serve(mine)
}

func serve(mine *icfp.Mine) {
    r := bufio.NewReaderSize(os.Stdin, 64)
    
    for char, err := r.ReadByte() ; err == nil ; char, err = r.ReadByte() {
	switch char {
	case 'L':
            fmt.Println(mine.ValidMove(icfp.Coord{mine.Robot.Coord[0], mine.Robot.Coord[1]-1}))
	case 'R':
            fmt.Println(mine.ValidMove(icfp.Coord{mine.Robot.Coord[0], mine.Robot.Coord[1]+1}))
	case 'U':
            fmt.Println(mine.ValidMove(icfp.Coord{mine.Robot.Coord[0]-1, mine.Robot.Coord[1]}))
	case 'D':
            fmt.Println(mine.ValidMove(icfp.Coord{mine.Robot.Coord[0]+1, mine.Robot.Coord[1]}))
        }
    }
}
