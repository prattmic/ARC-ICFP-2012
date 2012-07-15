package main

import (
        "os"
        "bufio"
        "fmt"
        "unsafe"
        "./icfp"
)

func main() {
    mine := new(icfp.Mine)
    err := mine.FromFile("maps/beard1_test.map", 100, false)

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
    fmt.Printf("Growth: %v\n", mine.Growth)
    fmt.Printf("Razors: %v\n", mine.Robot.Razors)

    fmt.Printf("\nMine struct:\n%+v\n\n", mine)

    fmt.Printf("Moving left allowed: %t\n", mine.ValidMove(icfp.Coord{mine.Robot.Coord[0], mine.Robot.Coord[1]-1}, false))
    fmt.Printf("Moving right allowed: %t\n", mine.ValidMove(icfp.Coord{mine.Robot.Coord[0], mine.Robot.Coord[1]+1}, false))
    fmt.Printf("Moving down allowed: %t\n", mine.ValidMove(icfp.Coord{mine.Robot.Coord[0]+1, mine.Robot.Coord[1]}, false))
    fmt.Printf("Moving up allowed: %t\n", mine.ValidMove(icfp.Coord{mine.Robot.Coord[0]-1, mine.Robot.Coord[1]}, false))

    size := unsafe.Sizeof(*mine) + unsafe.Sizeof(mine.Layout)
    for i := range mine.Layout {
        size += unsafe.Sizeof(mine.Layout[i])
    }

    fmt.Printf("Size: %d\n", size)


    //mine.Update(icfp.Coord{2,3})
    //for i := range mine.Layout {
    //    fmt.Println(string(mine.Layout[i]))
    //}
    //fmt.Printf("\nMine struct:\n%+v\n\n", mine)

    serve(mine)
}

func serve(mine *icfp.Mine) {
    r := bufio.NewReaderSize(os.Stdin, 64)

    for char, err := r.ReadByte() ; err == nil ; char, err = r.ReadByte() {
        switch char {
        case 'L', 'l':
            move := icfp.Coord{mine.Robot.Coord[0], mine.Robot.Coord[1]-1}
            if mine.ValidMove(move, false) {
                mine.Update(move, char)
            } else {
                fmt.Println("Invalid move")
            }
            mine.Print()
        case 'R', 'r':
            move := icfp.Coord{mine.Robot.Coord[0], mine.Robot.Coord[1]+1}
            if mine.ValidMove(move, false) {
                mine.Update(move, char)
            } else {
                fmt.Println("Invalid move")
            }
            mine.Print()
        case 'U', 'u':
            move := icfp.Coord{mine.Robot.Coord[0]-1, mine.Robot.Coord[1]}
            if mine.ValidMove(move, false) {
                mine.Update(move, char)
            } else {
                fmt.Println("Invalid move")
            }
            mine.Print()
        case 'D', 'd':
            move := icfp.Coord{mine.Robot.Coord[0]+1, mine.Robot.Coord[1]}
            if mine.ValidMove(move, false) {
                mine.Update(move, char)
            } else {
                fmt.Println("Invalid move")
            }
            mine.Print()
        case 'W', 'w':
            move := icfp.Coord{mine.Robot.Coord[0], mine.Robot.Coord[1]}
            if mine.ValidMove(move, false) {
                mine.Update(move, char)
            } else {
                fmt.Println("Invalid move")
            }
            mine.Print()
        case 'S', 's':
            move := icfp.Coord{mine.Robot.Coord[0], mine.Robot.Coord[1]}
            if mine.ValidMove(move, true) {
                mine.Update(move, char)
            } else {
                fmt.Println("Invalid move")
            }
            mine.Print()
        }
    }
}
