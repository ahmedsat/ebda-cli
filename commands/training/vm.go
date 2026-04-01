package training

import (
	"errors"
	"fmt"
	"slices"
)

type OptType int

const (
	Any OptType = iota
	Region
	Farm
	Code
	Module
	Not
	Or
	And
	Xor

	OptCount
)

func (o OptType) String() string {
	if o < 0 || o >= OptCount {
		return "Unknown"
	}

	return [...]string{
		"Any", "Region", "Farm", "Code", "Module", "Not", "Or", "And", "Xor",
	}[o]

}

type Opt struct {
	OptType
	Operand string
}

func (o Opt) String() string {
	return fmt.Sprintf("%s %s", o.OptType, o.Operand)
}

func Evaluate(v TrainingEntry, opts []Opt) (bool, error) {
	stack := []bool{}

	for _, opt := range opts {
		switch opt.OptType {
		case Region:
			stack = append(stack, slices.Contains(v.Regions, opt.Operand))
		case Farm:
			stack = append(stack, slices.Contains(v.Farms, opt.Operand))
		case Code:
			stack = append(stack, slices.Contains(v.Codes, opt.Operand))
		case Module:
			stack = append(stack, v.Modules[opt.Operand] > 0)
		case Not:
			if len(stack) == 0 {
				return false, errors.New("unexpected not")
			}
			stack[len(stack)-1] = !stack[len(stack)-1]
		case Or:
			if len(stack) < 2 {
				return false, errors.New("unexpected or")
			}
			stack[len(stack)-2] = stack[len(stack)-2] || stack[len(stack)-1]
			stack = stack[:len(stack)-1]
		case And:
			if len(stack) < 2 {
				return false, errors.New("unexpected and")
			}
			stack[len(stack)-2] = stack[len(stack)-2] && stack[len(stack)-1]
			stack = stack[:len(stack)-1]
		case Xor:
			if len(stack) < 2 {
				return false, errors.New("unexpected xor")
			}
			stack[len(stack)-2] = stack[len(stack)-2] != stack[len(stack)-1]
			stack = stack[:len(stack)-1]
		default:
			return false, fmt.Errorf("unexpected opt (%s)", opt)
		}
	}

	if len(stack) != 1 {
		return false, errors.New("unexpected end of program")
	}

	return stack[0], nil

}
