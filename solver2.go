package main

import (
        "os"
        "./icfp"
        "fmt"
        "container/list"
        "os/signal"
        "time"
        "bufio"
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
    //if len(os.Args) != 2 {
    //    usage()
    //}

    //fmt.Println("********************");
    //fmt.Println("* Solver 2         *");
    //fmt.Println("********************");

    mine := new(icfp.Mine)
    //err := mine.FromFile(os.Args[1], 100, false)
    r := bufio.NewReaderSize(os.Stdin, 64)
    err := mine.Load(r, 100, false)

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
    fmt.Printf("Growth: %v\n", mine.Growth)

    bestScore := new(AStar)
    bestScore.Mine = mine
    bestScore.D = 0
    bestScore.H = 0

    mapQ := list.New()

    bestSol := new(AStar)
    bestSol.Mine = mine
    bestSol.D = bestSol.GetD()
    bestSol.H = bestSol.GetH()

    bestSol.E = mapQ.PushBack(bestSol)

    options := []byte{'U','D','L','R','S'}

    //fmt.Printf("Distance to lift %d\n",bestSol.Mine.LiftDist())
    var counter = 0
    var Solved = false

    for i:=1;i<15000000;i++ {

        //Select Best Map
        tmpSol, ok := mapQ.Front().Value.(*AStar)
        if ok {
            bestSol = tmpSol
        } else {
            return
        }

        //Copy off highest score
        if bestSol.Mine.Score() > bestScore.Mine.Score() {
            bestScore = bestSol
        }

        // make children of Best map
        mapQ.Remove(bestSol.E)
        for j:=0;j<5;j++ {
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
            fmt.Println("SIGINT\n")
            fmt.Printf("%sA\n",bestScore.Mine.Command)
            bestScore.Mine.Robot.Abort = true
            fmt.Printf("Score: %d\n", bestScore.Mine.Score())
            fmt.Printf("\n%+v\n", bestScore.Mine)
            bestScore.Mine.Print()
            return
        default:
            if i%20 == 0 {
                time.Sleep(1*time.Microsecond)
            }
        }
    }

solved:
        if Solved {
            fmt.Printf("%s",bestScore.Mine.Command)
            fmt.Printf("Score: %d\n", bestScore.Mine.Score())
            fmt.Printf("\n%+v\n", bestScore.Mine)
            bestScore.Mine.Print()
        } else {
            //fmt.Println("No solution found")
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
