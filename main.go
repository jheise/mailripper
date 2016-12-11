package main

import (
	// standard
	"bufio"
	"fmt"
	"io"
	"net/mail"
	"os"
	"strings"
)

func parseMail(emailBuf io.Reader) error {
	var dest *os.File
	var targetFile io.Writer
	inAttachment := false
	inHeader := false
	inBody := false
	inFile := false
	encoding := ""
	filename := ""

	emailMsg, err := mail.ReadMessage(emailBuf)
	if err != nil {
		return err
	}

	emailReader := bufio.NewReader(emailMsg.Body)
	for {
		line, err := emailReader.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "--==") && strings.HasSuffix(line, "==") {
			//fmt.Printf("LINE: %s\n", line)
			inAttachment = true
			inHeader = true
			continue
		}
		if inAttachment && inHeader {
			if len(line) == 0 {
				inHeader = false
				inBody = true
				continue
			} else {
				//fmt.Printf("LINE: %s\n", line)
				if inHeader {
					fields := strings.Split(line, " ")
					if fields[0] == "Content-Disposition:" {
						//fmt.Printf("%s\n", fields[2])
						filesplit := strings.Split(fields[2], "\"")
						filename = filesplit[1]
						//fmt.Printf("filename: %s\n", filename)

						dest, err = os.Create(filename)
						if err != nil {
							return err
						}
						targetFile = bufio.NewWriter(dest)

						inFile = true
					}
					if fields[0] == "Content-Transfer-Encoding:" {
						encoding = fields[1]
					}
				}
			}
		}

		if inAttachment && inBody {
			if len(line) == 0 {
				inBody = false
				inFile = false
				if filename != "" {
					fmt.Printf("Wrote %s %s encoded file\n", filename, encoding)
				}
				//targetFile.Flush()
				dest.Close()
				encoding = ""
				filename = ""
				continue
			} else {
				if inFile {
					//targetFile.WriteString(line)
					fmt.Fprintf(targetFile, "%s", line)
				}
			}
		}

		if strings.HasPrefix(line, "--==") && strings.HasSuffix(line, "==--") {
			inAttachment = false
			inBody = false
		}
	}

	return nil
}

func usage() {
	fmt.Println("usage: mailripper filename")
}

func main() {
	if len(os.Args) != 2 {
		usage()
		os.Exit(1)
	}

	sourceFile := os.Args[1]
	emailFile, err := os.Open(sourceFile)
	if err != nil {
		panic(err)
	}
	defer emailFile.Close()

	emailBuf := bufio.NewReader(emailFile)
	err = parseMail(emailBuf)
	if err != nil {
		panic(err)
	}
}
