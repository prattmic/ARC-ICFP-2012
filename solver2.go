package main

import (
        "os"
        "./icfp"
        "fmt"
        "container/list"
        "os/signal"
        "time"
)

const (
    C   int = 45
    D   int = 10
    E   int = 40
    F   int = 35
)

type AStar struct {
    Mine    *icfp.Mine
    D       int
    H       int
    E       *list.Element
}
func usage() {
    fmt.Fprintf(os.Stderr, "usage: %s [mapfile]\n", os.Args[0])
    os.Exit(2)
}

func main() {
    if len(os.Args) != 2 {
        usage()
    }

    fmt.Println("********************");
    fmt.Println("* Solver 2         *");
    fmt.Println("********************");

    mine := new(icfp.Mine)
    err := mine.FromFile(os.Args[1], 100, false)

    if err != nil {
        fmt.Printf("Map failed to load, Error: %s\n", err)
    }

    //Load map data
    mine.ParseLayout()
    mine.Print()

    //Catch SIGINT
    sig := make(chan os.Signal, 10)
    signal.Notify(sig, os.Interrupt)

    //Print initial stats
    fmt.Printf("Water: %d\n", mine.Water)
    fmt.Printf("Flooding: %d\n", mine.Flooding)
    fmt.Printf("Waterproof: %d\n", mine.Robot.Waterproof)
    fmt.Printf("Trampolines: %v\n", mine.Trampolines)

    mapQ := list.New()
    bestSol := new(AStar)
    bestSol.Mine = mine
    bestSol.D = bestSol.GetD()
    bestSol.H = bestSol.GetH()

    bestSol.E = mapQ.PushBack(bestSol)

    options := []byte{'U','D','L','R'}

    fmt.Printf("Distance to lift %d\n",bestSol.Mine.LiftDist())
    var counter = 0
    var Solved = false

    for i:=1;i<1500;i++ {
        
        if i%500==0 {
            bestSol.Mine.Print()
        }
        //Select Best Map
        tmpSol, ok := mapQ.Front().Value.(*AStar)
        if ok {
            bestSol = tmpSol
        } else {
            return
        }
        //fmt.Printf("Length: %d\n",mapQ.Len())
        // make children of Best map
        mapQ.Remove(bestSol.E)
        //fmt.Printf("Length: %d\n",mapQ.Len())

        for j:=0;j<4;j++ {
            newMine := bestSol.Mine.Copy()
            if move(newMine,options[j]) && !newMine.Robot.Dead {
                //if newMine.FloodFillRouteHome() {
                //    fmt.Println("Clear to go home")
                //}
                counter++
                //Create new sol
                tmpSol := new(AStar)
                tmpSol.Mine = newMine
                tmpSol.D = tmpSol.GetD()
                tmpSol.H = tmpSol.GetH()

                //Sort the new solution into the map queue
                for e:= mapQ.Front(); e!= nil; e=e.Next() {
                    stackSol ,ok := e.Value.(*AStar)
                    if ok {
                        if (tmpSol.H+tmpSol.D)<(stackSol.H+stackSol.D) {
                            tmpSol.E = mapQ.InsertBefore(tmpSol,e)
                            break;
                        }                     }
                }
                if tmpSol.E==nil {
                    tmpSol.E = mapQ.PushBack(tmpSol)
                }

                if newMine.Complete {// || newMine.Robot.Lambda >= 3{
                    bestSol = tmpSol
                    Solved = true
                    goto solved
                }
            }
        }

        select {
        case <-sig:
            fmt.Println("SIGINT")
            bestSol.Mine.Print()
            fmt.Printf("%+v\n", bestSol.Mine)
        default:
            if i%1000 == 0 {
                time.Sleep(1*time.Microsecond)
            }
        }
    }

solved:
        if Solved {
            bestSol.Mine.Print()
            fmt.Printf("%s\n",bestSol.Mine.Command)
            fmt.Printf("Score: %d\n",bestSol.Mine.Score())
            fmt.Printf("Counter: %d\n",counter)
        } else {
            fmt.Println("No solution found")
        }
}
func (sol *AStar) GetH() int {
    return len(sol.Mine.Lambda)*C+sol.Mine.LiftDist()*D
}

func (sol *AStar) GetD() int {
    return sol.Mine.Robot.Moves*E-sol.Mine.Robot.Lambda*F
}

func move(mine *icfp.Mine, command byte) bool {
        switch command{
        case 'L', 'l':
            move := icfp.Coord{mine.Robot.Coord[0], mine.Robot.Coord[1]-1}
            if mine.ValidMove(move, false) {
                mine.Update(move,command )
                return true
            } else {
                return false
            }
            return false
        case 'R', 'r':
            move := icfp.Coord{mine.Robot.Coord[0], mine.Robot.Coord[1]+1}
            if mine.ValidMove(move, false) {
                mine.Update(move, command)
                return true
            } else {
                return false
            }
            return false
        case 'U', 'u':
            move := icfp.Coord{mine.Robot.Coord[0]-1, mine.Robot.Coord[1]}
            if mine.ValidMove(move, false) {
                mine.Update(move, command)
                return true
            } else {
                return false
            }
            return false
        case 'D', 'd':
            move := icfp.Coord{mine.Robot.Coord[0]+1, mine.Robot.Coord[1]}
            if mine.ValidMove(move, false) {
                mine.Update(move, command)
                return true
            } else {
                return false
            }
            return false
        case 'W', 'w':
            move := icfp.Coord{mine.Robot.Coord[0], mine.Robot.Coord[1]}
            if mine.ValidMove(move, false) {
                mine.Update(move, command)
                return true
            } else {
                return false
            }
            return false
        case 'A', 'a':
            mine.Command = append(mine.Command,'A')
            mine.Robot.Abort = true
            return true
        case 'S', 's':
            move := icfp.Coord{mine.Robot.Coord[0], mine.Robot.Coord[1]}
            if mine.ValidMove(move, true) {
                mine.Update(move, command)
                return true
            } else {
                return false
            }
            return false
        }
    return false
}
