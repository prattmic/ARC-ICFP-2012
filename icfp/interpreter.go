package icfp

import (
        "os"
        "io"
        "bufio"
        "strconv"
        "regexp"
        "fmt"
)

type Cell interface {
    Parse(coord Coord, mine *Mine)
    Update(coord Coord, mine *Mine, updated Map)
    MergeRobot(coord Coord, mine *Mine) bool
    Byte() byte
}

type CellSlice      []Cell
type CellSliceSlice []CellSlice

type RobotCell  byte
type RockCell   byte
type WallCell   byte
type LambdaCell byte
type EarthCell  byte
type EmptyCell  byte
type LiftCell   byte
type TrampCell  byte
type TargCell   byte
type BeardCell  byte
type RazorCell  byte

type Map        CellSliceSlice
type Coord      [2]int
type CoordSlice []Coord
type Robot struct {
    Coord       Coord
    Waterproof  int
    Moves       int
    Watermoves  int
    Lambda      int
    Razors      int
    Dead        bool
    Abort       bool
}
type Lift struct {
    Coord       Coord
    Open        bool
}
type Target struct {
    Num         int
    Coord       Coord
    TrampCoord  Coord
}
type Tramp      map[string]Target
type Mine struct {
    Command     []byte
    Layout      Map
    Robot       Robot
    Lambda      CoordSlice
    Lift        Lift
    Trampolines Tramp
    Water       int
    Flooding    int
    Growth      int
    Gcount      int
    Complete    bool
}

const (
    ROBOT   RobotCell   = 'R'
    ROCK    RockCell    = '*'
    WALL    WallCell    = '#'
    LAMBDA  LambdaCell  = '\\'
    EARTH   EarthCell   = '.'
    EMPTY   EmptyCell   = ' '
    CLIFT   LiftCell    = 'L'
    OLIFT   LiftCell    = 'O'
    BEARD   BeardCell   = 'W'
    RAZOR   RazorCell   = '!'
)

func (mine *Mine) Init() {
    mine.Water = 0
    mine.Flooding = 0
    mine.Robot.Waterproof = 10
    mine.Robot.Lambda = 0
    mine.Robot.Dead = false
    mine.Robot.Abort = false
    mine.Complete = false
    mine.Trampolines = make(Tramp)
    mine.Command = make([]byte,0,100)
    mine.Growth = 25 - 1
    mine.Robot.Razors = 0
}

func (mine *Mine) ParseLayout() {
    mine.Lambda = make([]Coord, 0, 100)

    for i := range mine.Layout {
        for j := range mine.Layout[i] {
            mine.Layout[i][j].Parse(Coord{i,j}, mine)
        }
    }
}

func (mine *Mine) Update(move Coord, command byte) {
    mine.Command = append(mine.Command,command)

    shave := command=='S'||command=='s'

    updated := NewMap(mine.Layout)

    //Update moves counter
    mine.Robot.Moves++

    //Trim beards
    if shave {
        mine.shave()
    }

    //Robot Movement
    mine.Layout[mine.Robot.Coord[0]][mine.Robot.Coord[1]].MergeRobot(move, mine)

    // Loop through and update the level
    for i := len(mine.Layout)-1; i>=0; i-- {
        for j := range mine.Layout[i] {
            mine.Layout[i][j].Update(Coord{i,j}, mine, updated)
        }
    }

    // Update State of the lift gate
    if mine.Lift.Open {
        updated[mine.Lift.Coord[0]][mine.Lift.Coord[1]] = OLIFT
    } else {
        updated[mine.Lift.Coord[0]][mine.Lift.Coord[1]] = CLIFT
    }

    //Update survival of the robot
    // The robot is dead if there is a rock above it that was not there previously (it is falling on the robot!)
    mine.Robot.Dead =(mine.Layout[mine.Robot.Coord[0]-1][mine.Robot.Coord[1]] != ROCK)&&(updated[mine.Robot.Coord[0]-1][mine.Robot.Coord[1]] == ROCK)

    //Update water damage
    if mine.Flooding != 0 {
        if mine.IsFlooded(mine.Robot.Coord) {
            mine.Robot.Watermoves++
        } else {
            mine.Robot.Watermoves = 0
        }
        if mine.Robot.Watermoves>mine.Robot.Waterproof {
            mine.Robot.Dead = true
        }
    }

    if mine.Gcount == 0 {
        mine.Gcount = mine.Growth - 1
    } else {
        mine.Gcount--
    }
    mine.Layout = updated 
}


func (mine *Mine) IsFlooded(loc Coord) bool {
    if (len(mine.Layout) - loc[0]) < (mine.Water + mine.Robot.Moves/mine.Flooding) {
        return true
    }

    return false
}

func (mine *Mine) LiftDist() int {
    return Abs(mine.Robot.Coord[0]-mine.Lift.Coord[0]) + Abs(mine.Robot.Coord[1]-mine.Lift.Coord[1])
}

func (mine *Mine) ValidMove(move Coord, shave bool) bool {
    y := Abs(mine.Robot.Coord[0]-move[0])
    x := Abs(mine.Robot.Coord[1]-move[1])
    tile := mine.Layout[move[0]][move[1]]
    // -1 = Left 0 = No horz 1 = Right
    horz := move[1] - mine.Robot.Coord[1]

    if (x != 0 && y != 0) || x > 1 || y > 1  {                              // Wrong move distance
        return false
    }

    if mine.Flooding != 0 && mine.IsFlooded(move) && mine.Robot.Waterproof <= mine.Robot.Watermoves {   // Drowning
        return false
    }

    if shave {
        if mine.Robot.Razors == 0 {
            return false
        } else {
            found := false
            for i := move[0]-1; i <= move[0]+1; i++ {
                for j := move[1]-1; j <= move[1]+1; j++ {
                    if mine.Layout[i][j] == BEARD {
                        found = true
                    }
                }
            }
            if !found {
                return false
            }
        }
    }

    // Move down with rock above
    if (move[0]-mine.Robot.Coord[0] == 1) && mine.Layout[move[0]-2][move[1]] == ROCK {
        return false
    }

    switch {
    case 'A' <= tile.Byte() && tile.Byte() <= 'I':
        return true
    case tile == ROBOT:
        return true
    case tile == ROCK:
        switch {
        case horz == -1 && mine.Layout[move[0]][move[1]-1] == EMPTY:    // Pushable Rock
            return true
        case horz == 1 && mine.Layout[move[0]][move[1]+1] == EMPTY:     // Pushable Rock
            return true
        default:
            return false
        }
    case tile == EMPTY, tile == EARTH, tile == LAMBDA, tile == RAZOR:
        switch {
        case (move[0]-2) >= 0 && len(mine.Layout[move[0]-2]) > move[1] && mine.Layout[move[0]-2][move[1]] == ROCK && mine.Layout[move[0]-1][move[1]] == EMPTY:   // Rock up 2 with empty space between
            return false
        default:
            return true
        }
    case tile == OLIFT:
        return true
    }

    // If move was valid, you should have returned by now
    return false
}

func (mine *Mine) FromFile(name string, capacity uint32, printonread bool) (err error) {
    file, err := os.Open(name)
    if err != nil {
        return err
    }
    fileinfo, err := file.Stat()
    
    r := bufio.NewReaderSize(file, int(fileinfo.Size()))

    mine.Load(r, capacity, printonread)

    return err
}

func (mine *Mine) Load(r *bufio.Reader, capacity uint32, printonread bool) (err error) {
    data := make(Map, 0, capacity)

    mine.Init()
    
    i := 0
    for ; ; i++ {
        line, err := ReadLine(r)
        if err == io.EOF {
            break
        } else if err != nil {
            fmt.Printf("Error: %s\n", err)
        }

        if printonread {
            fmt.Println(string(line))
        }

        // Blank
        if matched, _ := regexp.Match("^(?:\\s|\\n)*$", line); matched == true {
            continue
        } else if match := findSubmatch("Water\\s+([0-9]+)", string(line)); match != nil && len(match) == 2 {
            mine.Water, _ = strconv.Atoi(match[1])
        } else if match := findSubmatch("Flooding\\s+([0-9]+)", string(line)); match != nil && len(match) == 2 {
            mine.Flooding, _ = strconv.Atoi(match[1])
        } else if match := findSubmatch("Waterproof\\s+([0-9]+)", string(line)); match != nil && len(match) == 2 {
            mine.Robot.Waterproof, _ = strconv.Atoi(match[1])
        } else if match := findSubmatch("Trampoline\\s+([A-Z]+)\\s+targets\\s+([0-9]+)", string(line)); match != nil && len(match) == 3 {
            num, _ := strconv.Atoi(match[2])
            mine.Trampolines[match[1]] = Target{num, Coord{-1,-1}, Coord{-1,-1}}
        } else if match := findSubmatch("Growth\\s+([0-9]+)", string(line)); match != nil && len(match) == 2 {
            mine.Growth, _ = strconv.Atoi(match[1])
            mine.Gcount = mine.Growth - 1
        } else if match := findSubmatch("Razors\\s+([0-9]+)", string(line)); match != nil && len(match) == 2 {
            mine.Robot.Razors, _ = strconv.Atoi(match[1])
        } else {
            cells := Bytes2Cells(line)
            data = append(data, cells)
        }
    }

    mine.Layout = data

    return nil
}
