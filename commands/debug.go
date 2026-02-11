package commands

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/ahmedsat/ebda-cli/frappe"
	"github.com/ahmedsat/ebda-cli/frappe/types"
)

type Debug struct{}

// Description implements [main.subcommand].
func (d *Debug) Description() string {
	panic("unimplemented")
}

// Name implements [main.subcommand].
func (d *Debug) Name() string {
	return "debug"
}

// Result implements [main.subcommand].
func (d *Debug) Result() any {
	return nil
}

// Run implements [main.subcommand].
func (d *Debug) Run(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("not enough arguments")
	}
	debugCommand := args[0]
	switch args[0] {
	case "test":
		return d.Test(args[1:])
	case "update-farmer-name":
		return d.UpdateFarmerName(args[1:])
	default:
		fmt.Fprintf(os.Stderr, "unavailable commands: %s\n", debugCommand)
	}

	return nil
}

func (d *Debug) UpdateFarmerName(args []string) error {

	if len(args) < 3 {
		return fmt.Errorf("not enough arguments")
	}

	fameName := args[0]
	farmerNumber, err := strconv.Atoi(args[1])
	if err != nil {
		return err
	}
	farmerNewName := args[2]
	args = args[3:]

	farm, err := frappe.Get1[types.Farm](fameName)
	if err != nil {
		return err
	}

	if farmerNumber > len(farm.Farmers) {
		return fmt.Errorf("farmer number is too big: %d - max: %d", farmerNumber, len(farm.Farmers))
	}

	oldFarmer, err := frappe.Get1[types.Farmer](farm.Farmers[farmerNumber-1].FarmerName)
	if err != nil {
		return err
	}

	new := types.Farmer{
		FarmerName: farmerNewName,
		Phone:      oldFarmer.Phone,
		Gender:     oldFarmer.Gender,
	}

	_, err = frappe.Create(new)
	if err != nil && !errors.Is(err, frappe.ErrDedicated) {
		return err
	}

	trainings, err := frappe.Get[types.EbdaTraining](frappe.Filters{frappe.NewFilter("farm", frappe.Eq, farm.Name)}, []string{"name", "farm"})
	if err != nil {
		return err
	}

	for _, training := range trainings {
		training, err = frappe.Get1[types.EbdaTraining](training.Name)
		if err != nil {
			return err
		}
		for i := range training.Farmers {
			if training.Farmers[i].FarmerName == farm.Farmers[farmerNumber-1].FarmerName {
				training.Farmers[i].FarmerName = farmerNewName
				break
			}
		}
		_, err = frappe.UpdateDoc(training)
		if err != nil {
			return err
		}
	}

	followUps, err := frappe.Get[types.FarmFollowUp](frappe.Filters{frappe.NewFilter("farm", frappe.Eq, farm.Name)}, frappe.List{"name", "farm"})
	if err != nil {
		return err
	}

	for _, followUp := range followUps {
		followUp, err = frappe.Get1[types.FarmFollowUp](followUp.Name)
		if err != nil {
			return err
		}
		for i := range followUp.FarmersNames {
			if followUp.FarmersNames[i].Farmer == farm.Farmers[farmerNumber-1].FarmerName {
				followUp.FarmersNames[i].Farmer = farmerNewName
				break
			}
		}
		_, err = frappe.UpdateDoc(followUp)
		if err != nil {
			return err
		}
	}

	farm.Farmers[farmerNumber-1].FarmerName = farmerNewName
	_, err = frappe.UpdateDoc(farm)
	if err != nil {
		return err
	}

	return nil
}

func (d *Debug) Test(args []string) error {

	if len(args) < 1 {
		return fmt.Errorf("not enough arguments")
	}

	testModule := args[0]
	args = args[1:]

	switch testModule {
	case "frappe":
		if len(args) < 1 {
			return fmt.Errorf("not enough arguments")
		}
		resp, err := frappe.TestUrl(args[0])
		if err != nil {
			return err
		}
		fmt.Println(resp.Status)
	case "test":
		return Test(nil)
	default:
		fmt.Fprintf(os.Stderr, "unavailable commands: %s\n", testModule)
	}

	return nil
}

// Usage implements [main.subcommand].
func (d *Debug) Usage() string {
	panic("unimplemented")
}

func Test(args []string) error {

	return nil
}
