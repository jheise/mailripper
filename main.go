package main

import (
    "bufio"
    "fmt"
    "io"
    "github.com/mohamedattahri/mail"
    "os"
    )

func main() {
    testfile := "Invoices.eml"
    fmt.Printf("testing net/mail on %s\n", testfile)

    emailFile, err := os.Open(testfile)
    if err != nil {
        panic(err)
    }
    defer emailFile.Close()

    emailBuf := bufio.NewReader(emailFile)

    emailMsg, err := mail.ReadMessage(emailBuf)
    if err != nil {
        panic(err)
    }

    emailReader := bufio.NewReader(emailMsg.Body)
    for {
        line, err := emailReader.ReadString('\n')
        fmt.Printf("LINE: %s\n", line)
        if err == io.EOF {
            break
        } else if err != nil {
            panic(err)
        }
    }
}
