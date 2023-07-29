package commands

import (
	"fmt"
	"log"
	"math/rand"
	"regexp"
	"strconv"

	"github.com/NekoFluff/discord"
	"github.com/bwmarrin/discordgo"
)

func Roll() discord.Command {
	command := "roll"

	return discord.Command{
		Command: discordgo.ApplicationCommand{
			Name:        command,
			Description: fmt.Sprintf("Roll some dice (e.g. `%s 3d4+(3d6-1d4)*1d3+6`)", command),
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "dice",
					Description: "The dice to roll",
					Required:    true,
				},
			},
		},
		Handler: func(s discord.Session, i *discordgo.InteractionCreate) {
			options := i.ApplicationCommandData().Options
			optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
			for _, opt := range options {
				optionMap[opt.Name] = opt
			}

			diceRoll := fmt.Sprint(optionMap["dice"].Value)

			calculations := calculateDiceRoll(diceRoll)
			if calculations == nil {
				err := respondToInteraction(s, i.Interaction, "Invalid dice roll")
				if err != nil {
					log.Println(err)
				}
				return
			}

			response := fmt.Sprintf("__Dice Roll: %s__\n\n", diceRoll)

			for _, diceRoll := range calculations.diceRolls {
				response += fmt.Sprintf("%s\n", diceRoll)
			}
			response += fmt.Sprintf("### Result: %d\n", calculations.result)
			err := respondToInteraction(s, i.Interaction, response)

			if err != nil {
				log.Println(err)
			}
		},
	}
}

func calculateDiceRoll(input string) *Calculations {
	tokens := tokenize(input)
	postfix := infixToPostfix(tokens)
	calculations := evaluatePostfix(postfix)
	return calculations
}

func rollDice(numDice int, sides int) int {
	result := 0
	for i := 0; i < numDice; i++ {
		result += rand.Intn(sides) + 1
	}
	return result
}

func evaluatePostfix(tokens []string) *Calculations {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Panic occurred: ", err)
		}
	}()

	stack := make([]int, 0)

	calculations := Calculations{}

	for _, token := range tokens {
		if isNumber(token) {
			num, _ := strconv.Atoi(token)
			stack = append(stack, num)
		} else if token == "d" {
			numDice := stack[len(stack)-2]
			sides := stack[len(stack)-1]
			stack = stack[:len(stack)-2]
			result := rollDice(numDice, sides)
			stack = append(stack, result)
			calculations.diceRolls = append(calculations.diceRolls, fmt.Sprintf("%dd%d: %d", numDice, sides, result))
		} else {
			b := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			a := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			switch token {
			case "+":
				stack = append(stack, a+b)
			case "-":
				stack = append(stack, a-b)
			case "*":
				stack = append(stack, a*b)
			case "/":
				stack = append(stack, a/b)
			}
		}
	}

	calculations.result = stack[0]
	return &calculations
}

func isNumber(token string) bool {
	_, err := strconv.Atoi(token)
	return err == nil
}

func tokenize(input string) []string {
	re := regexp.MustCompile(`\d+|[+*/\(\)\-d]`)
	tokens := re.FindAllString(input, -1)
	return tokens
}

func infixToPostfix(tokens []string) []string {
	precedence := map[string]int{
		"d": 4,
		"*": 3,
		"/": 3,
		"+": 2,
		"-": 2,
		"(": 1,
	}

	output := make([]string, 0)
	operators := make([]string, 0)

	for _, token := range tokens {
		if isNumber(token) {
			output = append(output, token)
		} else if token == "(" {
			operators = append(operators, token)
		} else if token == ")" {
			for len(operators) > 0 && operators[len(operators)-1] != "(" {
				output = append(output, operators[len(operators)-1])
				operators = operators[:len(operators)-1]
			}
			operators = operators[:len(operators)-1] // Discard the "("
		} else {
			for len(operators) > 0 && precedence[token] <= precedence[operators[len(operators)-1]] {
				output = append(output, operators[len(operators)-1])
				operators = operators[:len(operators)-1]
			}
			operators = append(operators, token)
		}
	}

	for len(operators) > 0 {
		output = append(output, operators[len(operators)-1])
		operators = operators[:len(operators)-1]
	}

	return output
}

type Calculations struct {
	diceRolls []string
	result    int
}
