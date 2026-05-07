package kobo

import (
	"fmt"
	"slices"
	"strings"
	"unicode"

	"github.com/ahmedsat/ebda-cli/utils"
)

var regions = []string{
	"New Valley",
	"Sharqia",
	"Ismailia",
	"Qalyubia",
	"Daqahliya",
	"Galvina",
	"Gharbia",
	"North Sinai",
	"South Sinai",
	"Marsa Matrouh",
	"Siwa",
	"El Wahat AL Baharya",
	"Damietta",
	"Kafr El Sheikh",
	"Beni Suef",
	"Fayoum",
	"Al Behera",
	"Minya",
	"Qena",
	"Asyut",
	"Sohag",
	"Aswan",
	"Luxor",
	"Giza",
	"Menoufia",
	"Alexandria Desert",
}

type ValidationStatus struct {
	Timestamp int    `json:"timestamp"`
	Uid       string `json:"uid"`
	ByWhom    string `json:"by_whom"`
	Label     string `json:"label"`
}

type AtHouseAuditCommittee struct {
	AuditorName    string `json:"at_house/audit_committee_g/auditor_name"`
	AuditorNamePic string `json:"at_house/audit_committee_g/__014"`
}

type AtHouse001 struct {
	AtHouse001PlotName        string `json:"at_house/_001/plot_name"`
	AtHouse001PlotOwner       string `json:"at_house/_001/plot_owner"`
	AtHouse001PlotOwnerNumber string `json:"at_house/_001/plot_owner_number"`
	AtHouse001015             string `json:"at_house/_001/__015"`
}

type GroupBq2ft12GroupOz2rj74 struct {
	GroupBq2ft12GroupOz2rj74018 string `json:"group_bq2ft12/group_oz2rj74/__018"`
	GroupBq2ft12GroupOz2rj74030 string `json:"group_bq2ft12/group_oz2rj74/__030"`
	GroupBq2ft12GroupOz2rj74031 string `json:"group_bq2ft12/group_oz2rj74/__031"`
	GroupBq2ft12GroupOz2rj74032 string `json:"group_bq2ft12/group_oz2rj74/__032"`
	GroupBq2ft12GroupOz2rj74001 string `json:"group_bq2ft12/group_oz2rj74/__001"`
	GroupBq2ft12GroupOz2rj74035 string `json:"group_bq2ft12/group_oz2rj74/__035"`
	GroupBq2ft12GroupOz2rj74034 string `json:"group_bq2ft12/group_oz2rj74/__034"`
	GroupBq2ft12GroupOz2rj74036 string `json:"group_bq2ft12/group_oz2rj74/__036"`
	GroupBq2ft12GroupOz2rj74037 string `json:"group_bq2ft12/group_oz2rj74/__037"`
	GroupBq2ft12GroupOz2rj74038 string `json:"group_bq2ft12/group_oz2rj74/__038"`
	GroupBq2ft12GroupOz2rj74039 string `json:"group_bq2ft12/group_oz2rj74/__039"`
	GroupBq2ft12GroupOz2rj74029 string `json:"group_bq2ft12/group_oz2rj74/__029"`
	GroupBq2ft12GroupOz2rj74033 string `json:"group_bq2ft12/group_oz2rj74/__033"`
}

type GroupBq2ft12GroupNk73y42 struct {
	GroupBq2ft12GroupNk73y42040 string `json:"group_bq2ft12/group_nk73y42/__040"`
	GroupBq2ft12GroupNk73y42041 string `json:"group_bq2ft12/group_nk73y42/__041"`
	GroupBq2ft12GroupNk73y42042 string `json:"group_bq2ft12/group_nk73y42/__042"`
	GroupBq2ft12GroupNk73y42043 string `json:"group_bq2ft12/group_nk73y42/__043"`
	GroupBq2ft12GroupNk73y42044 string `json:"group_bq2ft12/group_nk73y42/__044"`
	GroupBq2ft12GroupNk73y42045 string `json:"group_bq2ft12/group_nk73y42/__045"`
	GroupBq2ft12GroupNk73y42046 string `json:"group_bq2ft12/group_nk73y42/__046"`
}

type GroupBq2ft12GroupOx7yv57 struct {
	GroupBq2ft12GroupOx7yv57FarmReportYesNo string `json:"group_bq2ft12/group_ox7yv57/farm_report_yes_no"`
	GroupBq2ft12GroupOx7yv57FarmReportPic   string `json:"group_bq2ft12/group_ox7yv57/farm_report_pic"`
}

type GroupCb7wg46 struct {
	GroupCb7wg46AreThereTreesOn string `json:"group_cb7wg46/_Are_there_trees_on"`
	GroupCb7wg46016             string `json:"group_cb7wg46/__016"`
	GroupCb7wg46028             string `json:"group_cb7wg46/__028"`
	GroupCb7wg46091             string `json:"group_cb7wg46/__091"`
	GroupCb7wg46093             string `json:"group_cb7wg46/__093"`
	GroupCb7wg46094             string `json:"group_cb7wg46/__094"`
	GroupCb7wg46095             string `json:"group_cb7wg46/__095"`
}

type GroupBo9ig41GroupEj02n94 struct {
	GroupBo9ig41GroupEj02n94081       string `json:"group_bo9ig41/group_ej02n94/__081"`
	GroupBo9ig41GroupEj02n94082       string `json:"group_bo9ig41/group_ej02n94/__082"`
	GroupBo9ig41GroupEj02n94T         string `json:"group_bo9ig41/group_ej02n94/T_"`
	GroupBo9ig41GroupEj02n94Source001 string `json:"group_bo9ig41/group_ej02n94/_Source_001"`
	GroupBo9ig41GroupEj02n94083       string `json:"group_bo9ig41/group_ej02n94/__083"`
	GroupBo9ig41GroupEj02n94084       string `json:"group_bo9ig41/group_ej02n94/__084"`
	GroupBo9ig41GroupEj02n94096       string `json:"group_bo9ig41/group_ej02n94/__096"`
}

type GroupBo9ig41 struct {
	GroupBo9ig41IsThereLivestock      string                     `json:"group_bo9ig41/_Is_there_livestock"`
	GroupBo9ig41Type                  string                     `json:"group_bo9ig41/_Type"`
	GroupBo9ig41077                   string                     `json:"group_bo9ig41/__077"`
	GroupBo9ig41078                   string                     `json:"group_bo9ig41/__078"`
	GroupBo9ig41080                   string                     `json:"group_bo9ig41/__080"`
	GroupBo9ig41079                   string                     `json:"group_bo9ig41/__079"`
	GroupBo9ig41GroupEj02n94G         []GroupBo9ig41GroupEj02n94 `json:"group_bo9ig41/group_ej02n94"`
	GroupBo9ig41AnimalPlaygroundYesNo string                     `json:"group_bo9ig41/animal_playground_yes_no"`
	GroupBo9ig41AnimalWellbeingYesNo  string                     `json:"group_bo9ig41/animal_wellbeing_yes_no"`
	GroupBo9ig41085                   string                     `json:"group_bo9ig41/__085"`
}

type GroupTm8rd70 struct {
	GroupTm8rd70029 string `json:"group_tm8rd70/__029"`
}

type GroupPh7bq74GroupZi8mj74 struct {
	GroupPh7bq74GroupZi8mj74055 string `json:"group_ph7bq74/group_zi8mj74/__055"`
	GroupPh7bq74GroupZi8mj74056 string `json:"group_ph7bq74/group_zi8mj74/__056"`
	GroupPh7bq74GroupZi8mj74057 string `json:"group_ph7bq74/group_zi8mj74/__057"`
	GroupPh7bq74GroupZi8mj74058 string `json:"group_ph7bq74/group_zi8mj74/__058"`
}

type GroupLm7bh83 struct {
	GroupLm7bh83Protection string `json:"group_lm7bh83/_Protection"`
	GroupLm7bh83051        string `json:"group_lm7bh83/__051"`
	GroupLm7bh83052        string `json:"group_lm7bh83/__052"`
	GroupLm7bh83053        string `json:"group_lm7bh83/__053"`
	GroupLm7bh83054        string `json:"group_lm7bh83/__054"`
}

type GroupWu36n20 struct {
	GroupWu36n20064                            string `json:"group_wu36n20/__064"`
	GroupWu36n20source                         string `json:"group_wu36n20/_Source"`
	GroupWu36n20066                            string `json:"group_wu36n20/__066"`
	GroupWu36n20067                            string `json:"group_wu36n20/__067"`
	GroupWu36n20068                            string `json:"group_wu36n20/__068"`
	GroupWu36n20069                            string `json:"group_wu36n20/__069"`
	GroupWu36n20070                            string `json:"group_wu36n20/__070"`
	GroupWu36n20097                            string `json:"group_wu36n20/__097"`
	GroupWu36n20GroupAl4ei07BiodynamicPrepa    string `json:"group_wu36n20/group_al4ei07/_Biodynamic_prepa"`
	GroupWu36n20GroupAl4ei07071                string `json:"group_wu36n20/group_al4ei07/__071"`
	GroupWu36n20GroupAl4ei07008                string `json:"group_wu36n20/group_al4ei07/__008"`
	GroupWu36n20GroupAl4ei07BiodynamicToolsPic string `json:"group_wu36n20/group_al4ei07/biodynamic_tools_pic"`
}

type GroupIa6hw30 struct {
	GroupIa6hw30SanitationFacilities    string `json:"group_ia6hw30/sanitation_facilities"`
	GroupIa6hw30SanitationFacilitiesPic string `json:"group_ia6hw30/sanitation_facilities_pic"`
}
type GroupWb0zy95 struct {
	GroupWb0zy95NonCompliance              string `json:"group_wb0zy95/non_compliance"`
	GroupWb0zy95092                        string `json:"group_wb0zy95/__092"`
	GroupWb0zy95CorrectiveMeasures         string `json:"group_wb0zy95/corrective_measures"`
	GroupWb0zy95099                        string `json:"group_wb0zy95/__099"`
	GroupWb0zy95CorrectiveMeasuresDeadline string `json:"group_wb0zy95/corrective_measures_deadline"`
	GroupWb0zy95AuditorsRecommendations    string `json:"group_wb0zy95/auditors_recommendations"`
	GroupWb0zy95004                        string `json:"group_wb0zy95/__004"`
}

type GroupPw9vq49 struct {
	GroupPw9vq49IsThereAStore string `json:"group_pw9vq49/_Is_there_a_store"`
	GroupPw9vq49072           string `json:"group_pw9vq49/__072"`
	GroupPw9vq49073           string `json:"group_pw9vq49/__073"`
	GroupPw9vq49074           string `json:"group_pw9vq49/__074"`
	GroupPw9vq49075           string `json:"group_pw9vq49/__075"`
	GroupPw9vq49076           string `json:"group_pw9vq49/__076"`
}

type EngineerData struct {
	EngineerName string `json:"engineer_data/engineer_name"`
	EngineerPic  string `json:"engineer_data/engineer_pic"`
	Region       string `json:"engineer_data/__002"`
}

type AtHouse struct {
	FarmNameAuto                          string                  `json:"at_house/farm_name_auto"`
	FarmId                                string                  `json:"at_house/farm_id"`
	AuditType                             string                  `json:"at_house/audit_type"`
	VisitDate                             string                  `json:"at_house/visit_date"`
	Season                                string                  `json:"at_house/__010"`
	FarmArea                              string                  `json:"at_house/farm_area"`
	FarmRepresentativeName                string                  `json:"at_house/__011"`
	FarmRepresentativePhoneNumber         string                  `json:"at_house/__012"`
	A_009                                 string                  `json:"at_house/__009"`
	A_001G                                []AtHouse001            `json:"at_house/_001"`
	A_018OwnershipmentDoc                 string                  `json:"at_house/_018"`
	GroupJs15s06PreviousCorrectiveActions string                  `json:"at_house/group_js15s06/previous_corrective_actions"`
	GroupJs15s06CorrectiveActionTaken     string                  `json:"at_house/group_js15s06/corrective_action_taken"`
	GroupJs15s06CorrectiveActionMethod    string                  `json:"at_house/group_js15s06/corrective_action_method"`
	EolPgsCriteriaYesNo                   string                  `json:"at_house/eol_pgs_criteria_yes_no"`
	A_017                                 string                  `json:"at_house/__017"`
	A_018                                 string                  `json:"at_house/__018"`
	AtHouseAuditCommitteeG                []AtHouseAuditCommittee `json:"at_house/audit_committee_g"`
}

type PGSNew struct {
	ID              int    `json:"_id"`
	Start           string `json:"start"`
	End             string `json:"end"`
	StartGeopoint   string `json:"start-geopoint"`
	Today           string `json:"today"`
	Username        string `json:"username"`
	Deviceid        string `json:"deviceid"`
	Phonenumber     string `json:"phonenumber"`
	RightsKnowledge string `json:"group_by3de59/__102"`
	EngineerData
	AtHouse
	GroupLl7vs36GPSOutlineOfFarm             string                     `json:"group_ll7vs36/GPS_Outline_of_Farm"`
	GroupLl7vs36IsThereASepar                string                     `json:"group_ll7vs36/_Is_there_a_separ"`
	GroupLl7vs36                             string                     `json:"group_ll7vs36/_"`
	GroupLl7vs36MethodOfSeparati             string                     `json:"group_ll7vs36/_Method_of_separati"`
	GroupLl7vs36023                          string                     `json:"group_ll7vs36/__023"`
	GroupLl7vs36019                          string                     `json:"group_ll7vs36/__019"`
	GroupLl7vs36020                          string                     `json:"group_ll7vs36/__020"`
	GroupLl7vs36021                          string                     `json:"group_ll7vs36/__021"`
	GroupLl7vs36022                          string                     `json:"group_ll7vs36/__022"`
	GroupLl7vs36024                          string                     `json:"group_ll7vs36/__024"`
	GroupLl7vs36025                          string                     `json:"group_ll7vs36/__025"`
	GroupLl7vs36026                          string                     `json:"group_ll7vs36/__026"`
	GroupLl7vs36027                          string                     `json:"group_ll7vs36/__027"`
	GroupLl7vs36005                          string                     `json:"group_ll7vs36/__005"`
	GroupLl7vs36007                          string                     `json:"group_ll7vs36/__007"`
	GroupLl7vs36Gmo                          string                     `json:"group_ll7vs36/_GMO"`
	GroupBq2ft12013                          string                     `json:"group_bq2ft12/__013"`
	GroupBq2ft12GroupOz2rj74G                []GroupBq2ft12GroupOz2rj74 `json:"group_bq2ft12/group_oz2rj74"`
	GroupBq2ft12GroupNk73y42G                []GroupBq2ft12GroupNk73y42 `json:"group_bq2ft12/group_nk73y42"`
	GroupBq2ft12GroupOx7yv57G                []GroupBq2ft12GroupOx7yv57 `json:"group_bq2ft12/group_ox7yv57"`
	GroupCb7wg46G                            []GroupCb7wg46             `json:"group_cb7wg46"`
	GroupBo9ig41G                            []GroupBo9ig41             `json:"group_bo9ig41"`
	GroupCm5lx00047                          string                     `json:"group_cm5lx00/__047"`
	GroupCm5lx00048                          string                     `json:"group_cm5lx00/__048"`
	GroupCm5lx00003                          string                     `json:"group_cm5lx00/__003"`
	GroupCm5lx00WellLicense                  string                     `json:"group_cm5lx00/well_license"`
	GroupCm5lx00049                          string                     `json:"group_cm5lx00/__049"`
	GroupCm5lx00IrrigationAmountFeddanYear   string                     `json:"group_cm5lx00/irrigation_amount_feddan_year"`
	GroupCm5lx00IsTillageCarriedOutOnThe     string                     `json:"group_cm5lx00/Is_tillage_carried_out_on_the_"`
	GroupCm5lx00IrrigationPollutantsYesNo    string                     `json:"group_cm5lx00/irrigation_pollutants_yes_no"`
	GroupCm5lx00029                          string                     `json:"group_cm5lx00/__029"`
	GroupCm5lx00050                          string                     `json:"group_cm5lx00/__050"`
	GroupPh7bq74086                          string                     `json:"group_ph7bq74/__086"`
	GroupPh7bq74062                          string                     `json:"group_ph7bq74/__062"`
	GroupPh7bq74GroupZi8mj74G                []GroupPh7bq74GroupZi8mj74 `json:"group_ph7bq74/group_zi8mj74"`
	GroupPh7bq74GroupRl6zg08059              string                     `json:"group_ph7bq74/group_rl6zg08/__059"`
	GroupPh7bq74GroupRl6zg08060              string                     `json:"group_ph7bq74/group_rl6zg08/__060"`
	GroupPh7bq74GroupRl6zg08061              string                     `json:"group_ph7bq74/group_rl6zg08/__061"`
	GroupYg03h10IsThereCropPlant             string                     `json:"group_yg03h10/_is_there_crop_plant_"`
	GroupYg03h10063                          string                     `json:"group_yg03h10/__063"`
	GroupTm8rd70G                            []GroupTm8rd70             `json:"group_tm8rd70"`
	GroupLm7bh83G                            []GroupLm7bh83             `json:"group_lm7bh83"`
	GroupSj6bq27006                          string                     `json:"group_sj6bq27/__006"`
	GroupSj6bq27buildingPreventionYesNo      string                     `json:"group_sj6bq27/building_prevention_yes_no"`
	GroupSj6bq27buildingPrevention           string                     `json:"group_sj6bq27/building_prevention"`
	GroupWu36n20G                            []GroupWu36n20             `json:"group_wu36n20"`
	GroupPw9vq49G                            []GroupPw9vq49             `json:"group_pw9vq49"`
	GroupIa6hw30G                            []GroupIa6hw30             `json:"group_ia6hw30"`
	GroupOt6ai14090                          string                     `json:"group_ot6ai14/__090"`
	GroupLk07o84ComplainsYesNo               string                     `json:"group_lk07o84/complains_yes_no"`
	GroupLk07o84ComplainsMethod              string                     `json:"group_lk07o84/complains_method"`
	GroupLk07o84ComplainProf                 string                     `json:"group_lk07o84/complain_prof"`
	GroupMl0jq70087                          string                     `json:"group_ml0jq70/__087"`
	GroupMl0jq70088                          string                     `json:"group_ml0jq70/__088"`
	GroupMl0jq70089                          string                     `json:"group_ml0jq70/__089"`
	GroupMl0jq70098                          string                     `json:"group_ml0jq70/__098"`
	GroupMl0jq70102                          string                     `json:"group_ml0jq70/__102"`
	GroupMl0jq70EssiveWorkingHours           string                     `json:"group_ml0jq70/_essive_working_hours"`
	GroupMl0jq70DoDocumentedContract         string                     `json:"group_ml0jq70/_Do_documented_contract"`
	GroupMl0jq70AttachAContract              string                     `json:"group_ml0jq70/_Attach_A_contract"`
	GroupMl0jq70IsSeasonalLaborEmployed      string                     `json:"group_ml0jq70/_Is_seasonal_labor_employed"`
	GroupMl0jq70AreTionateToTheCause         string                     `json:"group_ml0jq70/_Are_tionate_to_the_cause"`
	GroupMl0jq70488Eek8HoursPerDay           string                     `json:"group_ml0jq70/_48_8_eek_8_hours_per_day"`
	GroupMl0jq70DuringWorkingHours           string                     `json:"group_ml0jq70/_during_working_hours"`
	GroupMl0jq70AreDForOvertimeHours         string                     `json:"group_ml0jq70/_Are_d_for_overtime_hours"`
	GroupMl0jq70IsIntsAndSuggestions         string                     `json:"group_ml0jq70/_Is_ints_and_suggestions"`
	GroupMl0jq7015UnderTheAgeOf15            string                     `json:"group_ml0jq70/_15_under_the_age_of_15"`
	GroupMl0jq70SchoolingEducation           string                     `json:"group_ml0jq70/_schooling_education"`
	GroupMl0jq70LDangerousMachines           string                     `json:"group_ml0jq70/_l_dangerous_machines"`
	GroupMl0jq70ACcupationalInjuries         string                     `json:"group_ml0jq70/_A_ccupational_injuries"`
	GroupMl0jq70TiesAndCleanWater            string                     `json:"group_ml0jq70/_ties_and_clean_water"`
	GroupMl0jq70OccupationalSafety           string                     `json:"group_ml0jq70/_occupational_safety"`
	GroupMl0jq70NeedsAndChallenges           string                     `json:"group_ml0jq70/_needs_and_challenges"`
	GroupFy2dg78TStandardOfLiving            string                     `json:"group_fy2dg78/_t_standard_of_living"`
	GroupFy2dg78EWagesPaidOnTime             string                     `json:"group_fy2dg78/_e_wages_paid_on_time"`
	GroupFy2dg78ArePenaltie                  string                     `json:"group_fy2dg78/_Are_penaltie"`
	GroupFy2dg78PleaseSpeciETypesOfPenalties string                     `json:"group_fy2dg78/_Please_speci_e_types_of_penalties"`
	GroupFy2dg78HaveTheWorkersAtt            string                     `json:"group_fy2dg78/_Have_the_workers_att"`
	GroupFy2dg78KindlyIndNingsAndAttendance  string                     `json:"group_fy2dg78/_Kindly_ind_nings_and_attendance"`
	GroupYy5bi81OrSpiritualBeliefs           string                     `json:"group_yy5bi81/_or_spiritual_beliefs"`
	GroupYy5bi81103                          string                     `json:"group_yy5bi81/__103"`
	GroupYy5bi81ArOngTheFarmWorkers          string                     `json:"group_yy5bi81/_Ar_ong_the_farm_workers"`
	GroupYy5bi81DoCialPublicHolidays         string                     `json:"group_yy5bi81/_Do_cial_public_holidays"`
	GroupYy5bi81100                          string                     `json:"group_yy5bi81/__100"`
	GroupYy5bi81101                          string                     `json:"group_yy5bi81/__101"`
	GroupWb0zy95G                            []GroupWb0zy95             `json:"group_wb0zy95"`
	MetaInstanceID                           string                     `json:"meta/instanceID"`
	MetaRootUuid                             string                     `json:"meta/rootUuid"`
	MetaDeprecatedID                         string                     `json:"meta/deprecatedID"`
	FormhubUuid                              string                     `json:"formhub/uuid"`
	Version                                  string                     `json:"__version__"`
	XformIdString                            string                     `json:"_xform_id_string"`
	Uuid                                     string                     `json:"_uuid"`
	Attachments                              []Attachment               `json:"_attachments"`
	Status                                   string                     `json:"_status"`
	Geolocation                              []float64                  `json:"_geolocation"`
	SubmissionTime                           string                     `json:"_submission_time"`
	SubmittedBy                              string                     `json:"_submitted_by"`
	Tags                                     []struct{}                 `json:"_tags"`
	Notes                                    []struct{}                 `json:"_notes"`

	ValidationStatus `json:"_validation_status"`
	AttachmentMap    map[string]Attachment `json:"-"`
}

func (pgs PGSNew) GetFormID() string { return "aX4NJWgge6tooXjfSYXhrq" }

func (pgs PGSNew) Rate() (rate float64, issues []string) {

	pgs.AttachmentMap = make(map[string]Attachment)
	for i := range pgs.Attachments {
		_, ok := pgs.AttachmentMap[pgs.Attachments[i].MediaFileBasename]
		utils.Assert(!ok, fmt.Sprintf("duplicate attachment: %s in %d", pgs.Attachments[i].MediaFileBasename, pgs.ID))
		pgs.AttachmentMap[pgs.Attachments[i].MediaFileBasename] = pgs.Attachments[i]
	}

	checks := []utils.Check{
		pgs.check("ID"),
		pgs.check("Start"),
		pgs.check("End"),
		pgs.check("StartGeopoint"),
		pgs.check("Today"),
		pgs.check("Username"),
		pgs.check("Deviceid"),
		pgs.check("Phonenumber"),
		pgs.check("RightsKnowledge"),
		pgs.check("EngineerData.EngineerName"),
		pgs.check("EngineerData.EngineerPic"),
		pgs.check("EngineerData.Region"),
		pgs.check("AtHouse.FarmNameAuto"),
		pgs.check("AtHouse.FarmId"),
		pgs.check("AtHouse.AuditType"),
		pgs.check("AtHouse.VisitDate"),
		pgs.check("AtHouse.Season"),
		pgs.check("AtHouse.FarmArea"),
		pgs.check("AtHouse.FarmRepresentativeName"),
		pgs.check("AtHouse.FarmRepresentativePhoneNumber"),
		pgs.check("AtHouse.009"),
		pgs.check("AtHouse.001G"),
		pgs.check("AtHouse.018OwnershipmentDoc"),
		pgs.check("AtHouse.GroupJs15s06PreviousCorrectiveActions"),
		pgs.check("AtHouse.GroupJs15s06CorrectiveActionTaken"),
		pgs.check("AtHouse.GroupJs15s06CorrectiveActionMethod"),
		pgs.check("AtHouse.EolPgsCriteriaYesNo"),
		pgs.check("AtHouse.017"),
		pgs.check("AtHouse.018"),
		pgs.check("AtHouse.AuditCommitteeG"),
		pgs.check("GroupLl7vs36GPSOutlineOfFarm"),
		pgs.check("GroupLl7vs36IsThereASepar"),
		pgs.check("GroupLl7vs36"),
		pgs.check("GroupLl7vs36MethodOfSeparati"),
		pgs.check("GroupLl7vs36023"),
		pgs.check("GroupLl7vs36019"),
		pgs.check("GroupLl7vs36020"),
		pgs.check("GroupLl7vs36021"),
		pgs.check("GroupLl7vs36022"),
		pgs.check("GroupLl7vs36024"),
		pgs.check("GroupLl7vs36025"),
		pgs.check("GroupLl7vs36026"),
		pgs.check("GroupLl7vs36027"),
		pgs.check("GroupLl7vs36005"),
		pgs.check("GroupLl7vs36007"),
		pgs.check("GroupLl7vs36Gmo"),
		pgs.check("GroupBq2ft12013"),
		pgs.check("GroupBq2ft12GroupOz2rj74G"),
		pgs.check("GroupBq2ft12GroupNk73y42G"),
		pgs.check("GroupBq2ft12GroupOx7yv57G"),
		pgs.check("GroupCb7wg46G"),
		pgs.check("GroupBo9ig41G"),
		pgs.check("GroupCm5lx00047"),
		pgs.check("GroupCm5lx00048"),
		pgs.check("GroupCm5lx00003"),
		pgs.check("GroupCm5lx00WellLicense"),
		pgs.check("GroupCm5lx00049"),
		pgs.check("GroupCm5lx00IrrigationAmountFeddanYear"),
		pgs.check("GroupCm5lx00IsTillageCarriedOutOnThe"),
		pgs.check("GroupCm5lx00IrrigationPollutantsYesNo"),
		pgs.check("GroupCm5lx00029"),
		pgs.check("GroupCm5lx00050"),
		pgs.check("GroupPh7bq74086"),
		pgs.check("GroupPh7bq74062"),
		pgs.check("GroupPh7bq74GroupZi8mj74G"),
		pgs.check("GroupPh7bq74GroupRl6zg08059"),
		pgs.check("GroupPh7bq74GroupRl6zg08060"),
		pgs.check("GroupPh7bq74GroupRl6zg08061"),
		pgs.check("GroupYg03h10IsThereCropPlant"),
		pgs.check("GroupYg03h10063"),
		pgs.check("GroupTm8rd70G"),
		pgs.check("GroupLm7bh83G"),
		pgs.check("GroupSj6bq27006"),
		pgs.check("GroupSj6bq27buildingPreventionYesNo"),
		pgs.check("GroupSj6bq27buildingPrevention"),
		pgs.check("GroupWu36n20G"),
		pgs.check("GroupPw9vq49G"),
		pgs.check("GroupIa6hw30G"),
		pgs.check("GroupOt6ai14090"),
		pgs.check("GroupLk07o84ComplainsYesNo"),
		pgs.check("GroupLk07o84ComplainsMethod"),
		pgs.check("GroupLk07o84ComplainProf"),
		pgs.check("GroupMl0jq70087"),
		pgs.check("GroupMl0jq70088"),
		pgs.check("GroupMl0jq70089"),
		pgs.check("GroupMl0jq70098"),
		pgs.check("GroupMl0jq70102"),
		pgs.check("GroupMl0jq70EssiveWorkingHours"),
		pgs.check("GroupMl0jq70DoDocumentedContract"),
		pgs.check("GroupMl0jq70AttachAContract"),
		pgs.check("GroupMl0jq70IsSeasonalLaborEmployed"),
		pgs.check("GroupMl0jq70AreTionateToTheCause"),
		pgs.check("GroupMl0jq70488Eek8HoursPerDay"),
		pgs.check("GroupMl0jq70DuringWorkingHours"),
		pgs.check("GroupMl0jq70AreDForOvertimeHours"),
		pgs.check("GroupMl0jq70IsIntsAndSuggestions"),
		pgs.check("GroupMl0jq7015UnderTheAgeOf15"),
		pgs.check("GroupMl0jq70SchoolingEducation"),
		pgs.check("GroupMl0jq70LDangerousMachines"),
		pgs.check("GroupMl0jq70ACcupationalInjuries"),
		pgs.check("GroupMl0jq70TiesAndCleanWater"),
		pgs.check("GroupMl0jq70OccupationalSafety"),
		pgs.check("GroupMl0jq70NeedsAndChallenges"),
		pgs.check("GroupFy2dg78TStandardOfLiving"),
		pgs.check("GroupFy2dg78EWagesPaidOnTime"),
		pgs.check("GroupFy2dg78ArePenaltie"),
		pgs.check("GroupFy2dg78PleaseSpeciETypesOfPenalties"),
		pgs.check("GroupFy2dg78HaveTheWorkersAtt"),
		pgs.check("GroupFy2dg78KindlyIndNingsAndAttendance"),
		pgs.check("GroupYy5bi81OrSpiritualBeliefs"),
		pgs.check("GroupYy5bi81103"),
		pgs.check("GroupYy5bi81ArOngTheFarmWorkers"),
		pgs.check("GroupYy5bi81DoCialPublicHolidays"),
		pgs.check("GroupYy5bi81100"),
		pgs.check("GroupYy5bi81101"),
		pgs.check("GroupWb0zy95G"),
		pgs.check("MetaInstanceID"),
		pgs.check("MetaRootUuid"),
		pgs.check("MetaDeprecatedID"),
		pgs.check("FormhubUuid"),
		pgs.check("Version"),
		pgs.check("XformIdString"),
		pgs.check("Uuid"),
		pgs.check("Attachments"),
		pgs.check("Status"),
		pgs.check("Geolocation"),
		pgs.check("SubmissionTime"),
		pgs.check("SubmittedBy"),
		pgs.check("Tags"),
		pgs.check("Notes"),
		pgs.check("ValidationStatus"),
	}

	result := utils.EvaluateChecks(checks, nil)

	return result.Score, result.Issues
}

func (pgs PGSNew) check(s string) utils.Check {
	switch s {
	case "ID", "Start", "End", "Today", "Deviceid", "AtHouse.FarmNameAuto", "AtHouse.FarmArea":
		return utils.Check{Name: "ignore", Ok: true, Weight: 0}
	case "StartGeopoint":
		return utils.Check{Name: "يجب فتح تتبع الموقع", Ok: pgs.StartGeopoint != "", Weight: 1}
	case "Username":
		return utils.Check{Name: "تم تقديم الاستمارة من حساب مجهول", Ok: pgs.Username != "", Weight: 0}
	case "Phonenumber":
		issue, ok := validatePhoneNumber(pgs.Phonenumber)
		return utils.Check{Name: issue, Ok: ok, Weight: 0}
	case "RightsKnowledge":
		return utils.Check{Name: "المزارع ليس علي دراية بحقوق العمال", Ok: pgs.RightsKnowledge != "__1", Weight: 1}
	case "EngineerData.EngineerName":
		return utils.Check{Name: "اسم المهندس فارغ", Ok: pgs.EngineerData.EngineerName != "", Weight: 1}
	case "EngineerData.EngineerPic":
		att, ok := pgs.AttachmentMap[pgs.EngineerData.EngineerPic]
		if !ok {
			return utils.Check{Name: "صورة المهندس غير موجودة", Ok: false, Weight: 1}
		}
		_, err := DownloadAttach(att)
		if err != nil {
			return utils.Check{Name: "صورة المهندس غير قابلة للتحميل", Ok: false, Weight: 1}
		}
		return utils.Check{Name: "صورة المهندس بحاجة الي المراجعة اليدوية", Ok: false, Weight: 0}
	case "EngineerData.Region":
		return utils.Check{Name: "المنطقة غير صحيحة", Ok: slices.Contains(regions, pgs.EngineerData.Region), Weight: 1}
	case "AtHouse.FarmId":
		return utils.Check{Name: "كود المزرع فارغ", Ok: pgs.AtHouse.FarmId != "", Weight: 1}
	case "AtHouse.AuditType":
		return utils.Check{Name: "نوع الزيارة غير صحيح", Ok: slices.Contains([]string{"First_visit", "Follow_up_/_yearly"}, pgs.EngineerData.Region), Weight: 1}
	case "AtHouse.VisitDate":
		return utils.Check{Name: "تاريخ الزيارة غير صحيح", Ok: pgs.AtHouse.VisitDate != "", Weight: 1}
	case "AtHouse.Season":
		return utils.Check{Name: "الموسم الزراعي فارغ", Ok: pgs.AtHouse.Season != "", Weight: 1}

	default:
		return utils.Check{Name: fmt.Sprintf("Unknown field: %s", s), Ok: false, Weight: 1}
	}
}

func validatePhoneNumber(s string) (issue string, ok bool) {
	if strings.TrimSpace(s) == "" {
		return "رقم الهاتف فارغ", false
	}

	// Normalize: keep digits and optional leading '+'
	var b strings.Builder
	b.Grow(len(s))

	hasPlus := false

	for i, r := range s {
		switch {
		case r == '+':
			if i != 0 {
				return "علامة '+' في مكان غير صحيح", false
			}
			hasPlus = true
		case unicode.IsDigit(r):
			b.WriteRune(r)
		case r == ' ' || r == '-' || r == '(' || r == ')':
			continue
		default:
			return "يحتوي على أحرف غير صالحة", false
		}
	}

	num := b.String()

	switch {
	// Local: 01XXXXXXXXX
	case len(num) == 11 && strings.HasPrefix(num, "01"):
		if !validMobilePrefix(num[0:3]) {
			return "مقدمة رقم غير صحيحة", false
		}
		return "", true

	// International without '+': 201XXXXXXXXX
	case len(num) == 12 && strings.HasPrefix(num, "20"):
		local := "0" + num[2:]
		if !validMobilePrefix(local[0:3]) {
			return "مقدمة رقم غير صحيحة", false
		}
		return "", true

	// International with '+': +201XXXXXXXXX
	case hasPlus && len(num) == 12 && strings.HasPrefix(num, "20"):
		local := "0" + num[2:]
		if !validMobilePrefix(local[0:3]) {
			return "مقدمة رقم غير صحيحة", false
		}
		return "", true
	}

	return "صيغة رقم الهاتف غير صحيحة", false
}

func validMobilePrefix(prefix string) bool {
	switch prefix {
	case "010", "011", "012", "015":
		return true
	default:
		return false
	}
}
