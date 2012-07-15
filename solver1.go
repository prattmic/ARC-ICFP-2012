package main

import (
        "./icfp"
        "fmt"
)

func main() {
    fmt.Println("********************");
    fmt.Println("* Solver 1         *");
    fmt.Println("********************");

    mine := new(icfp.Mine)
    err := mine.FromFile("maps/contest1.map", 100)

    if err != nil {
        fmt.Printf("Map failed to load, Error: %s\n", err)
    }

    //Load map data
    mine.ParseLayout()

    //Print initial stats
    fmt.Printf("Water: %d\n", mine.Water)
    fmt.Printf("Flooding: %d\n", mine.Flooding)
    fmt.Printf("Waterproof: %d\n", mine.Robot.Waterproof)
    fmt.Printf("Trampolines: %v\n", mine.Trampolines)

    mine.Print()
    move(mine,'D')
    mine.Print()
    move(mine,'L')
    mine.Print()
    move(mine,'R')
    mine.Print()
    move(mine,'D')
    mine.Print()
    move(mine,'D')
    mine.Print()
}

func move(mine *icfp.Mine, move byte) bool {
        switch move{
        case 'L', 'l':
            move := icfp.Coord{mine.Robot.Coord[0], mine.Robot.Coord[1]-1}
            if mine.ValidMove(move) {
                mine.Update(move)
                return true
            } else {
                return false
            }
            return false
        case 'R', 'r':
            move := icfp.Coord{mine.Robot.Coord[0], mine.Robot.Coord[1]+1}
            if mine.ValidMove(move) {
                mine.Update(move)
                return true
            } else {
                return false
            }
            return false
        case 'U', 'u':
            move := icfp.Coord{mine.Robot.Coord[0]-1, mine.Robot.Coord[1]}
            if mine.ValidMove(move) {
                mine.Update(move)
                return true
            } else {
                return false
            }
            return false
        case 'D', 'd':
            move := icfp.Coord{mine.Robot.Coord[0]+1, mine.Robot.Coord[1]}
            if mine.ValidMove(move) {
                mine.Update(move)
                return true
            } else {
                return false
            }
            return false
        case 'W', 'w':
            move := icfp.Coord{mine.Robot.Coord[0], mine.Robot.Coord[1]}
            if mine.ValidMove(move) {
                mine.Update(move)
                return true
            } else {
                return false
            }
            return false
        }
    return false
}