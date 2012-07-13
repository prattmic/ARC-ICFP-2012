package main

import (
        "os"
        "bufio"
        "fmt"
)

type Map [][]byte

func main() {
    mine, err := MapFromFile("maps/contest1.map", 100)
    if err != nil {
        fmt.Printf("Error: %s\n", err)
    }

    for i := 0; i < len(mine); i++ {
        fmt.Println(string(mine[i]))
    }

    robot := mine.currentLocation()
    fmt.Printf("You are at %d\n", robot)
}

func (mine Map) currentLocation() ([2]int) {
    for i := 0; i < len(mine); i++ {
        for j := 0; j < len(mine[i]); j++ {
            if mine[i][j] == 'R' {
                return [2]int{i,j}
            }
        }
    }

    return [2]int{-1,-1}
}

func MapFromFile(name string, capacity uint32) (mine Map, err error) {
    file, err := os.Open(name)
    if err != nil {
        return nil, err
    }
    fileinfo, err := file.Stat()
    
    r := bufio.NewReaderSize(file, int(fileinfo.Size()))

    data := make([][]byte, 0, capacity)

    i := 0
    for ; ; i++ {
        line, _, err := r.ReadLine()
        if err != nil {
            break
        }
        data = append(data, line)
    }

    return data, nil
}
