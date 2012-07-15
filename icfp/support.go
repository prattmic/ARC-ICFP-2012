package icfp

import (
        "regexp"
        "fmt"
        "errors"
        "bufio"
        "io"
)

func NewMap(ref Map) Map {
    updated := make([][]byte, len(ref))

    // Create new blank map
    for i := range ref {
        updated[i] = make([]byte, len(ref[i]))
    }

    return updated
}

func (mine *Mine) Copy() *Mine {
    tmp := new(Mine)
    *tmp = *mine

    tmp.Layout= make([][]byte, len(mine.Layout))
    for i := range mine.Layout {
        newSlice := make([]byte, len(mine.Layout[i]))
        copy(newSlice,mine.Layout[i])
        tmp.Layout[i] = newSlice
    }
    newSlice := make(CoordSlice, len(mine.Lambda))
    copy(newSlice,mine.Lambda)
    tmp.Lambda = newSlice

    newSlice2 := make([]byte, len(mine.Command))
    copy(newSlice2,mine.Command)
    tmp.Command = newSlice2



    return tmp
}

func (mine *Mine) shave() bool {
    if mine.Robot.Razors > 0 {
        mine.Robot.Razors--

        for i := mine.Robot.Coord[0]-1; i <= mine.Robot.Coord[0]+1; i++ {
            for j := mine.Robot.Coord[1]-1; j <= mine.Robot.Coord[1]+1; j++ {
                if mine.Layout[i][j] == BeardChar {
                    mine.Layout[i][j] = EmptyChar
                }
            }
        }

        return true
    }

    return false
}

func (mine *Mine) eatLambda(move Coord) error {
    mine.Robot.Lambda++

    /* Get index in list */
    coordi, err := mine.Lambda.FindCoord(Coord{move[0], move[1]})
    if err != nil {
        fmt.Printf("Error: %s\n", err)
        mine.Robot.Lambda--
        return err
    }

    /* Delete it */
    mine.Lambda = append(mine.Lambda[:coordi], mine.Lambda[coordi+1:]...)

    if len(mine.Lambda) == 0 {
        mine.Lift.Open = true
    }

    return nil
}

func (mine *Mine) takejump(move Coord) {
    targ := mine.Trampolines[string(mine.Layout[move[0]][move[1]])]
    mine.RemoveTramps(targ)

    mine.Layout[mine.Robot.Coord[0]][mine.Robot.Coord[1]] = EmptyChar
    mine.Layout[move[0]][move[1]] = EmptyChar
    mine.Layout[targ.Coord[0]][targ.Coord[1]] = RoboChar
    mine.Robot.Coord = targ.Coord
}

func (coords CoordSlice) FindCoord(item Coord) (index int, err error) {
    for i := range coords {
        if coords[i] == item {
            return i, nil
        }
    }

    return -1, fmt.Errorf("Item not found in CoordSlice: %v", item)
}

func (mine *Mine) TargetCoord(n int, coord Coord) error {
    found := false
    for key, value := range mine.Trampolines {
        if value.Num == n {
            found = true
            value.Coord = coord
            mine.Trampolines[key] = value
        }
    }

    if found {
        return error(nil)
    }

    return errors.New("Item not found in Trampolines")
}

func (mine *Mine) RemoveTramps(targ Target) {
    for key, value := range mine.Trampolines {
        if value.Num == targ.Num {
            mine.Layout[value.TrampCoord[0]][value.TrampCoord[1]] = EmptyChar
            mine.Layout[value.Coord[0]][value.Coord[1]] = EmptyChar
            delete(mine.Trampolines, key)
        }
    }
}

func (mine *Mine) Score() int {
    if(mine.Robot.Dead) {
        return 0
    } else if(mine.Robot.Abort) {
        return mine.Robot.Lambda*50-mine.Robot.Moves
    } else if(mine.Complete) {
        return mine.Robot.Lambda*75-mine.Robot.Moves
    } else {
        return mine.Robot.Lambda*25-mine.Robot.Moves
    }
    return 0
}

func (mine *Mine) Print() {
    //fmt.Printf("Current Score: %d\n",mine.score());
    for i := range mine.Layout {
        fmt.Println(string(mine.Layout[i]))
    }
}

// Inspired by Alex Ray
func ReadLine(r *bufio.Reader) ([]byte, error) {
    l := make([]byte, 0, 4096)
    
    for {
        line, isPrefix, err := r.ReadLine()

        if err != nil && err != io.EOF {
            return nil, err
        }

        l = append(l, line...)

        if err == io.EOF {
            return l, err
        }
        if !isPrefix {
            break
        }
    }
    return l, nil
}

func Abs(n int) int {
    if n < 0 {
        return -n
    }

    return n
}

func findSubmatch(reg string, line string) []string {
    re, err := regexp.Compile(reg)
    if err != nil {
        fmt.Printf("Error: %s", err)
        return nil
    }

    return re.FindStringSubmatch(line)
}
