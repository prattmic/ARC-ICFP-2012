package icfp

import (
        "regexp"
        "fmt"
        "errors"
)

func IndexCoord(slice []Coord, item Coord) (index int, err error) {
    for i := range slice {
        if slice[i] == item {
            return i, nil
        }
    }

    return -1, errors.New("Item not found in []Coord")
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
