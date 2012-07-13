package main

import (
        "os"
        "bufio"
        "fmt"
)

type Map [][]byte

func main() {
    mine, err := MapFromFile("maps/contest6.map", 100)
    if err != nil {
        fmt.Printf("Error: %s\n", err)
    }

    for i := 0; i < len(mine); i++ {
        fmt.Println(string(mine[i]))
    }
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
