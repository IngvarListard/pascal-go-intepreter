package main

import (
	"bufio"
	"fmt"
	"github.com/IngvarListard/pascal-go-intepreter/pkg/calc3"
	"log"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("calc> ")
		s, _ := reader.ReadString('\n')
		s = strings.Replace(s, "\n", "", -1)

		if s == "exit" {
			log.Fatal("end")
		}

		p, err := calc3.NewParser(s)
		if err != nil {
			fmt.Println("error: ", err)
			continue
		}
		i, err := calc3.NewInterpreter(p)
		if err != nil {
			fmt.Println("error: ", err)
			continue
		}

		r, err := i.Term()
		if err != nil {
			fmt.Println("Error: ", err)
		}
		fmt.Println(r)
	}
}
