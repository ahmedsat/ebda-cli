package main

import (
	_ "embed"
	"fmt"

	"github.com/ahmedsat/ebda-cli/config"
	"github.com/ahmedsat/ebda-cli/frappe"
	"github.com/ahmedsat/ebda-cli/frappe/types"
)

var farms = []string{}

func main() {

	err := config.Configure()
	if err != nil {
		panic(err)
	}

	_, err = frappe.Login()
	if err != nil {
		panic(err)
	}

	for _, farmName := range farms {
		farm, err := frappe.Get1[types.Farm](farmName)
		if err != nil {
			panic(err)
		}

		if farm.FarmApplication == "" {
			fmt.Println(farmName, " : has no farm application")
			continue
		}

		farmApp, err := frappe.Get1[types.FarmApplication](farm.FarmApplication)
		if err != nil {
			fmt.Println(farmName, " : ", err)
			continue
		}

		farm.Farmers = make([]types.FarmFarmer, len(farmApp.Farmers))
		for i, farmer := range farmApp.Farmers {
			farm.Farmers[i] = types.FarmFarmer{
				FarmerName:            farmer.Farmer,
				TotalArea:             farmer.OwnedAreaInFeddan,
				AreaUnit:              "Feddan",
				FarmerNationalIdImage: farmer.FarmerNationalId,
				NationalIdNumber:      farmer.FarmerNationalIdNo,
				Phone:                 fmt.Sprintf("%d", farmer.FarmerPhoneNumber),
			}
		}

		err = farm.Update()
		if err != nil {
			fmt.Println(farmName, " : ", err)
			continue
		}

		fmt.Println(farmName, " : ", farm.TotalFarmers == farmApp.TotalFarmers)
	}
}
