package main

import (
        "./icfp"
        "fmt"
        "container/list"
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

    mapQ := list.New()
    mapQ.PushBack(mine)

    options := []byte{'U','D','L','R'}

    var solved = false;

    for i:=1;i<20;i++ {
        for e:= mapQ.Front(); e!= nil; e=e.Next() {
            tmpMine ,ok := e.Value.(*icfp.Mine)
            if ok {
                for j:=0;j<4;j++ {
                    newMine := tmpMine.Copy()
                    if move(newMine,options[j]) {
                        mapQ.PushFront(newMine)
              //          fmt.Printf("%+v\n",newMine)
                        if newMine.Complete {
                            solved = true
                        }
                    }
                    if solved {
                        break
                    }
                }
            }
            if solved {
                break
            }
        }
        if solved {
            break
        }
    }
    if solved {
        tmpMine, ok := mapQ.Front().Value.(*icfp.Mine)
        if ok {
            tmpMine.Print()
        }
    }
/*
        for e:= mapQ.Front(); e!= nil; e=e.Next() {
            tmpMine ,ok := e.Value.(*icfp.Mine)
            if ok {
                tmpMine.Print()
                }
            }
  */      
    /*
    mine.Print()
    move(mine,'D')
    mine.Print()

    newMine := mine.Copy()
    move(newMine,'U')
    newMine.Print()

    move(mine,'L')
    mine.Print()

    //Illegal move
    if(!move(mine,'U')) {
        fmt.Println("Illegal move")
    } else {
    mine.Print()
    }

    move(mine,'R')
    mine.Print()
    move(mine,'D')
    mine.Print()
    move(mine,'D')
    mine.Print()
    fmt.Printf("%+v\n",mine)
    fmt.Printf("%+v\n",newMine)
    */
}

func move(mine *icfp.Mine, move byte) bool {
        switch move{
        case 'L', 'l':
            move := icfp.Coord{mine.Robot.Coord[0], mine.Robot.Coord[1]-1}
            if mine.ValidMove(move, false) {
                mine.Update(move, false)
                return true
            } else {
                return false
            }
            return false
        case 'R', 'r':
            move := icfp.Coord{mine.Robot.Coord[0], mine.Robot.Coord[1]+1}
            if mine.ValidMove(move, false) {
                mine.Update(move, false)
                return true
            } else {
                return false
            }
            return false
        case 'U', 'u':
            move := icfp.Coord{mine.Robot.Coord[0]-1, mine.Robot.Coord[1]}
            if mine.ValidMove(move, false) {
                mine.Update(move, false)
                return true
            } else {
                return false
            }
            return false
        case 'D', 'd':
            move := icfp.Coord{mine.Robot.Coord[0]+1, mine.Robot.Coord[1]}
            if mine.ValidMove(move, false) {
                mine.Update(move, false)
                return true
            } else {
                return false
            }
            return false
        case 'W', 'w':
            move := icfp.Coord{mine.Robot.Coord[0], mine.Robot.Coord[1]}
            if mine.ValidMove(move, false) {
                mine.Update(move, false)
                return true
            } else {
                return false
            }
            return false
        case 'S', 's':
            move := icfp.Coord{mine.Robot.Coord[0], mine.Robot.Coord[1]}
            if mine.ValidMove(move, true) {
                mine.Update(move, true)
                return true
            } else {
                return false
            }
            return false
        }
    return false
}
