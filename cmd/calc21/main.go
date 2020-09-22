package main

import (
	"bufio"
	"fmt"
	"github.com/IngvarListard/pascal-go-intepreter/pkg/calc21"
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

		interpreter := calc21.NewInterpreter(text)
		res := interpreter.Eval()
		fmt.Println(res)
	}
}
