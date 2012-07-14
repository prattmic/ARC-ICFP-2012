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


    mine.Update(icfp.Coord{2,4})
    for i := range mine.Layout {
        fmt.Println(string(mine.Layout[i]))
    }



    serve(mine)
}


func (mine *Mine) Update(move Coord) {
    updated := make([][]byte, len(mine.Layout))

    mine.Robot.Coord = move
    for i := range mine.Layout {
        updated[i] = make([]byte, len(mine.Layout[i]))
	
        for j := range mine.Layout[i] {
            if i==move[0] && j==move[1] {
                updated[i][j] = RoboChar
                mine.Robot.Lambda++
            } else if mine.Layout[i][j] == RoboChar && (i==move[0] || j==move[1]) {
                updated[i][j] = EmptyChar
            } else if mine.Layout[i][j] == EmptyChar{
		updated[i][j] = EmptyChar
            } else if mine.Layout[i][j] == LambdaChar{
                updated[i][j] = LambdaChar
            } else if mine.Layout[i][j] == EarthChar{
                updated[i][j] = EarthChar
            } else if mine.layout[i][j] == RockChar{
                //Rock logic goes here                
                updated[i][j] = RockChar
            } else if mine.Layout[i][j] == WallChar {
                updated[i][j] = WallChar
            } else if mine.Layout[i][j] == OLiftChar {
                updated[i][j] = OLiftChar
            } else if mine.Layout[i][j] == CLiftChar {
                updated[i][j] = CLiftChar
            } else if mine.Layout[i][j] == EarthChar {
                updated[i][j] = EarthChar
            }
        }
    }
    mine.Layout = updated 
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
