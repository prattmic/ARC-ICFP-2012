package icfp

import (
        "regexp"
        "fmt"
        "errors"
)

func (coords CoordSlice) FindCoord(item Coord) (index int, err error) {
    for i := range coords {
        if coords[i] == item {
            return i, nil
        }
    }

    return -1, errors.New("Item not found in CoordSlice")
}

func (rocks RockSlice) FindRock(curr Coord) (int, error) {
    for i := range rocks {
        if rocks[i].Curr == curr {
            return i, nil
        }
    }

    return -1, errors.New("Item not found in RockSlice")
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

func (mine *Mine) Print() {
    fmt.Printf("Current Score: %d\n",mine.score());
    for i := range mine.Layout {
        fmt.Println(string(mine.Layout[i]))
    }
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
