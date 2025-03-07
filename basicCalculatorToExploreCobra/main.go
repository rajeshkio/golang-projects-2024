package main

import (
	"log"
	"strconv"

	"github.com/spf13/cobra"
)

func operation(op string, n1, n2 string) (float32, error) {

	num1, err := strconv.Atoi(n1)
	if err != nil {
		return 0, err
	}

	num2, err := strconv.Atoi(n2)
	if err != nil {
		return 0, err
	}

	fl1 := float32(num1)
	fl2 := float32(num2)
	switch op {
	case "add":
		return fl1 + fl2, nil
	case "sub":
		return fl1 - fl2, nil
	case "div":
		return fl1 / fl2, nil
	}
	return 0, nil
}

var cmdAdd = &cobra.Command{
	Use:   "add",
	Short: "Add two numbers",
	Long:  "Add two numbers: add <number1> <number2>",
	Run: func(cmd *cobra.Command, args []string) {
		result, err := operation("add", args[0], args[1])
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("result: %f\n", result)
	},
	Args: cobra.ExactArgs(2),
}

func main() {
	var calculator = &cobra.Command{
		Use:   "Calculator",
		Short: "A simple calculator",
	}
	calculator.AddCommand(cmdAdd)
	calculator.Execute()
}
