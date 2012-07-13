package main

import (
        "os"
        "bufio"
        "fmt"
)

func main() {
    file, err := os.Open("maps/contest1.map")
    if err != nil {
        fmt.Printf("Error: %s", err)
    }
    fileinfo, err := file.Stat()
    
    r := bufio.NewReaderSize(file, int(fileinfo.Size()))

    data := make([][]byte, 10, 100)

    i := 0
    for ; ; i++ {
        line, _, err := r.ReadLine()
        if err != nil {
            break
        }
        data[i] = line;
    }

    j := i
    for i = 0; i < j; i++ {
        fmt.Println(string(data[i]))
    }
}
