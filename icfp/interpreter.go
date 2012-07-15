package icfp

import (
        "os"
        "io"
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

var RoboChar    byte = 'R'
var RockChar    byte = '*'
var WallChar    byte = '#'
var LambdaChar  byte = '\\'
var EarthChar   byte = '.'
var EmptyChar   byte = ' '
var CLiftChar   byte = 'L'
var OLiftChar   byte = 'O'
var BeardChar   byte = 'W'
var RazorChar   byte = '!'

func (mine *Mine) Init() {
    mine.Water = 0
    mine.Flooding = 0
    mine.Robot.Waterproof = 10
    mine.Robot.Lambda = 0
    mine.Robot.Dead = false
    mine.Robot.Abort = false
    mine.Complete = false
    mine.Trampolines = make(Tramp)
    mine.Growth = 25 - 1
    mine.Robot.Razors = 0
}

func (mine *Mine) ParseLayout() {
    mine.Lambda = make([]Coord, 0, 100)

    for i := range mine.Layout {
        for j := range mine.Layout[i] {
            switch {
            case mine.Layout[i][j] == LambdaChar:
                mine.Lambda = append(mine.Lambda, Coord{i,j})
            case mine.Layout[i][j] == CLiftChar:
                mine.Lift.Coord = Coord{i,j}
                mine.Lift.Open = false
            case mine.Layout[i][j] == RoboChar:
                mine.Robot.Coord = Coord{i,j}
            case '0' <= mine.Layout[i][j] && mine.Layout[i][j] <= '9':
                num, _ := strconv.Atoi(string(mine.Layout[i][j]))
                mine.TargetCoord(num, Coord{i, j})
            case 'A' <= mine.Layout[i][j] && mine.Layout[i][j] <= 'I':
                targ := mine.Trampolines[string(mine.Layout[i][j])]
                targ.TrampCoord = Coord{i,j}
                mine.Trampolines[string(mine.Layout[i][j])] = targ
            }
        }
    }
}

func (mine *Mine) Update(move Coord, trim bool) {
    var trampjump = false

    updated := make([][]byte, len(mine.Layout))

    // Create new blank map
    for i := range mine.Layout {
        updated[i] = make([]byte, len(mine.Layout[i]))
    }

    //Update moves counter
    mine.Robot.Moves++

    //Trim beards
    if trim {
        if mine.Robot.Razors > 0 {
            i := mine.Robot.Coord[0]
            j := mine.Robot.Coord[1]
            mine.Robot.Razors--
            if mine.Layout[i+1][j-1] == BeardChar {
                mine.Layout[i+1][j-1] = EmptyChar
            }
            if mine.Layout[i+1][j] == BeardChar {
                mine.Layout[i+1][j] = EmptyChar
            }
            if mine.Layout[i+1][j+1] == BeardChar {
                mine.Layout[i+1][j+1] = EmptyChar
            }
            if mine.Layout[i][j-1] == BeardChar {
                mine.Layout[i][j-1] = EmptyChar
            }
            if mine.Layout[i][j+1] == BeardChar {
                mine.Layout[i][j+1] = EmptyChar
            }
            if mine.Layout[i-1][j-1] == BeardChar {
                mine.Layout[i-1][j-1] = EmptyChar
            }
            if mine.Layout[i-1][j] == BeardChar {
                mine.Layout[i-1][j] = EmptyChar
            }
            if mine.Layout[i-1][j+1] == BeardChar {
                mine.Layout[i-1][j+1] = EmptyChar
            }
        }
    }

    //Robot Movement
    switch {
    //Get lambda
    case mine.Layout[move[0]][move[1]] == LambdaChar:
        mine.Robot.Lambda++

        /* Get index in list */
        coordi, err := mine.Lambda.FindCoord(Coord{move[0], move[1]})
        if err != nil {
            fmt.Printf("Error: %s\n", err)
            return
        }

        /* Delete it */
        mine.Lambda = append(mine.Lambda[:coordi], mine.Lambda[coordi+1:]...)

        if len(mine.Lambda) == 0 {
            mine.Lift.Open = true
        }
    //Get razor
    case mine.Layout[move[0]][move[1]] == RazorChar:
        mine.Robot.Razors++
    //Move rock
    case mine.Layout[move[0]][move[1]] == RockChar:
        switch {
        case mine.Robot.Coord[1]<move[1]:
            mine.Layout[move[0]][move[1]+1] = RockChar

        case mine.Robot.Coord[1]>move[1]:
            mine.Layout[move[0]][move[1]-1] = RockChar
        }
    //Trampoline
    case 'A' <= mine.Layout[move[0]][move[1]] && mine.Layout[move[0]][move[1]] <= 'I':
        trampjump = true
        targ := mine.Trampolines[string(mine.Layout[move[0]][move[1]])]
        mine.RemoveTramps(targ)

        mine.Layout[mine.Robot.Coord[0]][mine.Robot.Coord[1]] = EmptyChar
        mine.Layout[move[0]][move[1]] = EmptyChar
        mine.Layout[targ.Coord[0]][targ.Coord[1]] = RoboChar
        mine.Robot.Coord = targ.Coord

        //Check for completion    
    case mine.Layout[move[0]][move[1]] == OLiftChar:
        mine.Complete = true
        return
    }

    mine.Complete = false

    //Move the robot
    if !trampjump {
        mine.Layout[mine.Robot.Coord[0]][mine.Robot.Coord[1]] = EmptyChar
        mine.Layout[move[0]][move[1]] = RoboChar
        mine.Robot.Coord = move
    }

    // Loop through and update the level
    for i := len(mine.Layout)-1; i>=0; i-- {
        for j := range mine.Layout[i] {
            switch mine.Layout[i][j] {
            default:
                updated[i][j] = mine.Layout[i][j]
            case EmptyChar: 
                if updated[i][j] != RockChar && updated[i][j] != BeardChar {
                    updated[i][j] = EmptyChar
                }
            case BeardChar:
                updated[i][j] = mine.Layout[i][j]
                if mine.Gcount == 0 {
                    if updated[i+1][j-1] == EmptyChar {
                        updated[i+1][j-1] = BeardChar
                    }
                    if updated[i+1][j] == EmptyChar {
                        updated[i+1][j] = BeardChar
                    }
                    if updated[i+1][j+1] == EmptyChar {
                        updated[i+1][j+1] = BeardChar
                    }
                    if updated[i][j-1] == EmptyChar {
                        updated[i][j-1] = BeardChar
                    }
                    if mine.Layout[i][j+1] == EmptyChar {
                        updated[i][j+1] = BeardChar
                    }
                    if mine.Layout[i-1][j-1] == EmptyChar {
                        updated[i-1][j-1] = BeardChar
                    }
                    if mine.Layout[i-1][j] == EmptyChar {
                        updated[i-1][j] = BeardChar
                    }
                    if mine.Layout[i-1][j+1] == EmptyChar {
                        updated[i-1][j+1] = BeardChar
                    }
                }
            case RockChar:
                switch {
                case mine.Layout[i+1][j] == EmptyChar && updated[i+1][j] != BeardChar:
                    //Rule 1
                    updated[i][j] = EmptyChar
                    updated[i+1][j] = RockChar

                case (mine.Layout[i+1][j] == RockChar || mine.Layout[i+1][j] == LambdaChar) && updated[i+1][j] != BeardChar && mine.Layout[i][j+1] == EmptyChar && updated[i][j+1] != BeardChar && mine.Layout[i+1][j+1] == EmptyChar && updated[i+1][j+1] != BeardChar:
                    //Rule 2 and 4
                    updated[i][j] = EmptyChar
                    updated[i+1][j+1] = RockChar
                case mine.Layout[i+1][j] == RockChar && updated[i+1][j] != BeardChar && mine.Layout[i][j-1] == EmptyChar && updated[i][j-1] != BeardChar && mine.Layout[i+1][j-1] == EmptyChar && updated[i+1][j-1] != BeardChar:
                    //Rule 3
                    updated[i][j] = EmptyChar
                    updated[i+1][j-1] = RockChar
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
    // The robot is dead if there is a rock above it that was not there previously (it is falling on the robot!)
    mine.Robot.Dead =(mine.Layout[mine.Robot.Coord[0]-1][mine.Robot.Coord[1]] != RockChar)&&(updated[mine.Robot.Coord[0]-1][mine.Robot.Coord[1]] == RockChar)

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

func (mine *Mine) ValidMove(move Coord, trim bool) bool {
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

    if trim {
        if mine.Robot.Razors == 0 {
            return false
        } else {
            found := false
            for i := move[0]-1; i <= move[0]+1; i++ {
                for j := move[1]-1; j <= move[1]+1; j++ {
                    if mine.Layout[i][j] == BeardChar {
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
    if (move[0]-mine.Robot.Coord[0] == 1) && mine.Layout[move[0]-2][move[1]] == RockChar {
        return false
    }

    switch {
    case 'A' <= tile && tile <= 'I':
        return true
    case tile == RoboChar:
        return true
    case tile == RockChar:
        switch {
        case horz == -1 && mine.Layout[move[0]][move[1]-1] == EmptyChar:    // Pushable Rock
            return true
        case horz == 1 && mine.Layout[move[0]][move[1]+1] == EmptyChar:     // Pushable Rock
            return true
        default:
            return false
        }
    case tile == EmptyChar, tile == EarthChar, tile == LambdaChar:
        switch {
        case (move[0]-2) >= 0 && len(mine.Layout[move[0]-2]) > move[1] && mine.Layout[move[0]-2][move[1]] == RockChar && mine.Layout[move[0]-1][move[1]] == EmptyChar:   // Rock up 2 with empty space between
            return false
        default:
            return true
        }
    case tile == OLiftChar:
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

    data := make([][]byte, 0, capacity)

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
            data = append(data, line)
        }
    }

    mine.Layout = data

    return nil
}
