package icfp

import (
        "regexp"
        "fmt"
        "errors"
        "bufio"
        "io"
)

func (coords CoordSlice) FindCoord(item Coord) (index int, err error) {
    for i := range coords {
        if coords[i] == item {
            return i, nil
        }
    }

    return -1, fmt.Errorf("Item not found in CoordSlice: %v", item)
}
func (mine *Mine) Copy() *Mine {
    tmp := new(Mine)
    *tmp = *mine

    tmp.Layout= make([][]byte, len(mine.Layout))
    for i := range mine.Layout {
        newSlice := make([]byte, len(mine.Layout[i]))
        copy(newSlice,mine.Layout[i])
        tmp.Layout[i] = newSlice
        fmt.Println(tmp.Layout[i])
    }

    return tmp

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
