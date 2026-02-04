package commands

import (
	"fmt"

	"github.com/ahmedsat/erp-reports-cli/erp"
	"github.com/ahmedsat/erp-reports-cli/utils"
)

type LocationTable struct {
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}

type FarmApplication struct {
	Name           string          `json:"name"`
	Engineer       string          `json:"engineer_name"`
	UserName       string          `json:"user_name"`
	MapCoordinate  string          `json:"map_coordinate"`
	LocationTables []LocationTable `json:"location_table"`
}

func (f FarmApplication) DocTypeName() string {
	return "Farm Application"
}

func FarmApplications(args []string) (err error) {
	if len(args) == 0 {
		return fmt.Errorf("not enough arguments")
	}

	switch args[0] {
	case "info":
		if len(args) < 2 {
			return fmt.Errorf("not enough arguments")
		}
		app, err := erp.Get1[FarmApplication](args[1])
		if err != nil {
			return err
		}
		fmt.Println(app)

	case "create_map":
		if len(args) < 2 {
			return fmt.Errorf("not enough arguments")
		}
		return CreateMap(args[1])
	case "validate":
		if len(args) < 2 {
			return fmt.Errorf("not enough arguments")
		}
		err = ValidateFarmApplication(args[1])
		if err != nil {
			return
		}
	default:
		err = fmt.Errorf("unknown farm applications command: %s", args[0])
	}

	return err
}

func CreateMap(id string) (err error) {

	app, err := erp.Get1[FarmApplication](id)
	if err != nil {
		return err
	}

	if len(app.LocationTables) == 0 {
		return fmt.Errorf("no location table")
	}

	if len(app.LocationTables) > 1 {
		return fmt.Errorf("more than one location table")
	}

	if len(app.MapCoordinate) == 0 {
		return fmt.Errorf("no map coordinate")
	}

	res, err := erp.CallMethod("create_map_record_if_inside", map[string]any{
		"app":     app.Name,
		"polygon": app.MapCoordinate,
		"lat":     app.LocationTables[0].Latitude,
		"lng":     app.LocationTables[0].Longitude,
	})
	if err != nil {
		return err
	}
	if res != nil {
		fmt.Println(string(res))
	}
	return
}

func ValidateFarmApplication(id string) (err error) {

	apps, err := erp.Get[FarmApplication](utils.Filters{}, []string{"name"})
	if err != nil {
		return
	}

	fmt.Println(apps[0])

	return
}
