package main

import (
	"bufio"
	"fmt"
	"github.com/IngvarListard/pascal-go-intepreter/pkg/calc1"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	var text string
	for {
		fmt.Print("calc> ")
		text, _ = reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)

		if text == "exit" {
			fmt.Println("exiting...")
			break
		}

		if text == "" {
			continue
		}

		interpreter := calc1.NewInterpreter(text)
		res, err := interpreter.Eval()
		if err != nil {
			fmt.Printf("error occured: %v\n", err)
		}
		fmt.Println(res)
	}
}
