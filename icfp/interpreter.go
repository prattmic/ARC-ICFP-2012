package icfp

import (
        "os"
        "bufio"
        "strconv"
        "regexp"
        "fmt"
)

type Map        [][]byte
type Coord      [2]int
type CoordSlice []Coord
type Robot struct {
    Coord       Coord
    Waterproof  int
    Moves       int
    Watermoves  int
    Lambda      int
    Dead        bool
}
type Lift struct {
    Coord       Coord
    Open        bool
}
type Rock struct {
    Curr        Coord
    Prev        Coord
}
type RockSlice  []Rock
type Mine struct {
    Layout      Map
    Robot       Robot
    Lambda      CoordSlice
    Rocks       RockSlice
    Lift        Lift
    Water       int
    Flooding    int
    Complete    bool
}

var RoboChar    byte = 'R'
var RockChar    byte = '*'
var WallChar    byte = '#'
var LambdaChar  byte = '\\'
var EarthChar   byte = '.'
var EmptyChar   byte = ' '
var CLiftChar   byte = 'L'
var OLiftChar   byte = 'O'

func (mine *Mine) ParseLayout() {
    mine.Lambda = make([]Coord, 0, 100)
    mine.Rocks = make([]Rock, 0, 100)

    for i := range mine.Layout {
        for j := range mine.Layout[i] {
            if mine.Layout[i][j] == LambdaChar {
                mine.Lambda = append(mine.Lambda, Coord{i,j})
            } else if mine.Layout[i][j] == RockChar {
                mine.Rocks = append(mine.Rocks, Rock{Coord{i,j}, Coord{i,j}})
            } else if mine.Layout[i][j] == CLiftChar {
                mine.Lift.Coord = Coord{i,j}
                mine.Lift.Open = false
            } else if mine.Layout[i][j] == RoboChar {
                mine.Robot.Coord = Coord{i,j}
            }
        }
    }
}

func (mine *Mine) Update(move Coord) {
    var updatedRockPrev = false

    updated := make([][]byte, len(mine.Layout))

    // Create new blank map
    for i := range mine.Layout {
        updated[i] = make([]byte, len(mine.Layout[i]))
    }


    //Robot Movement
    //Get lambda
    if mine.Layout[move[0]][move[1]]==LambdaChar {   
        mine.Robot.Lambda++

        /* Get index in list */
        coordi, err := mine.Lambda.FindCoord(Coord{move[0], move[1]})
        if err != nil {
            fmt.Printf("Error: %s\n", err)
            return
        }

        /* Delete it */
        mine.Lambda = append(mine.Lambda[:coordi], mine.Lambda[coordi+1:]...)
    }

    //Move rock
    if mine.Layout[move[0]][move[1]] == RockChar {      
        switch {
        case mine.Robot.Coord[1]<move[1]:
            mine.Layout[move[0]][move[1]+1] = RockChar

            rock, err := mine.Rocks.FindRock(Coord{move[0],move[1]})
            if err != nil {
                fmt.Printf("Error: %s\n", err)
                return
            }

            mine.Rocks[rock].Prev = mine.Rocks[rock].Curr
            mine.Rocks[rock].Curr = Coord{move[0],move[1]+1}
            updatedRockPrev = true
        case mine.Robot.Coord[1]>move[1]:
            mine.Layout[move[0]][move[1]-1] = RockChar

            rock, err := mine.Rocks.FindRock(Coord{move[0],move[1]})
            if err != nil {
                fmt.Printf("Error: %s\n", err)
                return
            }

            mine.Rocks[rock].Prev = mine.Rocks[rock].Curr
            mine.Rocks[rock].Curr = Coord{move[0],move[1]-1}
            updatedRockPrev = true
        }
    }

    //Check for completion    
    if mine.Layout[move[0]][move[1]]==OLiftChar {
        mine.Complete = true
        return
    } else {
        mine.Complete = false
    }

    //Move the robot
    mine.Layout[mine.Robot.Coord[0]][mine.Robot.Coord[1]] = EmptyChar
    mine.Layout[move[0]][move[1]] = RoboChar
    mine.Robot.Coord = move

    // Loop through and update the level
    for i := range mine.Layout {
        for j := range mine.Layout[i] {
            switch mine.Layout[i][j] {
            case RoboChar:
                updated[i][j] = RoboChar
            case EmptyChar: 
                if updated[i][j] != RockChar {
                    updated[i][j] = EmptyChar
                }
            case LambdaChar:
                updated[i][j] = LambdaChar
            case EarthChar:
                updated[i][j] = EarthChar
            case WallChar:
                updated[i][j] = WallChar
            case RockChar:
                switch {
                case mine.Layout[i+1][j] == EmptyChar:
                    //Rule 1
                    updated[i][j] = EmptyChar
                    updated[i+1][j] = RockChar

                    rock, err := mine.Rocks.FindRock(Coord{i,j})
                    if err != nil {
                        fmt.Printf("Error: %s\n", err)
                        return
                    }

                    if updatedRockPrev == false {
                        mine.Rocks[rock].Prev = mine.Rocks[rock].Curr
                    }
                    mine.Rocks[rock].Curr = Coord{i+1,j}
                case (mine.Layout[i+1][j] == RockChar || mine.Layout[i+1][j] == LambdaChar) && mine.Layout[i][j+1] == EmptyChar && mine.Layout[i+1][j+1] == EmptyChar:
                    //Rule 2 and 4
                    updated[i][j] = EmptyChar
                    updated[i+1][j+1] = RockChar

                    rock, err := mine.Rocks.FindRock(Coord{i,j})
                    if err != nil {
                        fmt.Printf("Error: %s\n", err)
                        return
                    }
                    if updatedRockPrev == false {
                        mine.Rocks[rock].Prev = mine.Rocks[rock].Curr
                    }
                    mine.Rocks[rock].Curr = Coord{i+1,j+1}
                case mine.Layout[i+1][j] == RockChar && mine.Layout[i][j-1] == EmptyChar && mine.Layout[i+1][j-1] == EmptyChar:
                    //Rule 3
                    updated[i][j] = EmptyChar
                    updated[i+1][j-1] = RockChar

                    rock, err := mine.Rocks.FindRock(Coord{i,j})
                    if err != nil {
                        fmt.Printf("Error: %s\n", err)
                        return
                    }
                    if updatedRockPrev == false {
                        mine.Rocks[rock].Prev = mine.Rocks[rock].Curr
                    }
                    mine.Rocks[rock].Curr = Coord{i+1,j-1}
                default:
                    updated[i][j] = RockChar
                }
            }
        }
    }

    // Update State of the lift gate
    if mine.Lift.Open {
        updated[mine.Lift.Coord[0]][mine.Lift.Coord[1]] = OLiftChar
    } else {
        updated[mine.Lift.Coord[0]][mine.Lift.Coord[1]] = CLiftChar
    }

    //Update survival of the robot
    if mine.Layout[mine.Robot.Coord[0]][mine.Robot.Coord[1]] == RockChar {
        mine.Robot.Dead = true
    // Ending condition #3
    } else if mine.Layout[mine.Robot.Coord[0]-1][mine.Robot.Coord[1]] == RockChar {
        rock, err := mine.Rocks.FindRock(Coord{mine.Robot.Coord[0]-1,mine.Robot.Coord[1]})
        if err != nil {
            fmt.Printf("Error: %s\n", err)
            return
        }

        // Falling
        if mine.Rocks[rock].Curr[0] != mine.Rocks[rock].Prev[0] {
            mine.Robot.Dead = true
        }
    }

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

    mine.Layout = updated 
    mine.Robot.Moves++
}

func (mine *Mine) IsFlooded(loc Coord) bool {
    if (len(mine.Layout) - loc[0]) < (mine.Water + mine.Robot.Moves/mine.Flooding) {
        return true
    }

    return false
}

func (mine *Mine) ValidMove(move Coord) bool {
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

    switch tile {
    case RockChar:
        switch {
        case horz == -1 && mine.Layout[move[0]][move[1]-1] == EmptyChar:    // Pushable Rock
            return true
        case horz == 1 && mine.Layout[move[0]][move[1]+1] == EmptyChar:     // Pushable Rock
            return true
        default:
            return false
        }
    case EmptyChar, EarthChar, LambdaChar, OLiftChar:
        switch {
        case mine.Layout[move[0]-1][move[1]] == RockChar:                   // Rock above space
            rock, err := mine.Rocks.FindRock(Coord{move[0]-1,move[1]})
            if err != nil {
                fmt.Printf("Error: %s\n", err)
                return false
            }

            // Can move if rock above isn't falling
            return mine.Rocks[rock].Curr[0] == mine.Rocks[rock].Prev[0];
        case mine.Layout[move[0]-2][move[1]] == RockChar && mine.Layout[move[0]-1][move[1]] == EmptyChar:   // Rock up 2 with empty space between
            return false
        default:
            return true
        }
    }

    // If move was valid, you should have returned by now
    return false
}

func (mine *Mine) FromFile(name string, capacity uint32) (err error) {
    file, err := os.Open(name)
    if err != nil {
        return err
    }
    fileinfo, err := file.Stat()
    
    r := bufio.NewReaderSize(file, int(fileinfo.Size()))

    data := make([][]byte, 0, capacity)

    mine.Water = 0
    mine.Flooding = 0
    mine.Robot.Waterproof = 10
    mine.Robot.Lambda = 0
    mine.Robot.Dead = false
    mine.Complete = false

    i := 0
    for ; ; i++ {
        line, _, err := r.ReadLine()
        if err != nil {
            break
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
        } else {
            data = append(data, line)
        }
    }

    mine.Layout = data

    return nil
}
