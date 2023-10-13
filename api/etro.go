package api

import (
	"encoding/json"
	"net/http"
)

type EtroResponse struct {
	TotalParams []struct {
		ID    any    `json:"id,omitempty"`
		Name  string `json:"name"`
		Value any    `json:"value"`
		Units string `json:"units,omitempty"`
	} `json:"totalParams"`
}

type Stats struct {
	AverageItemLevel float64
	WeaponDamage     int
	MainStat         int
	Vitality         int
	CriticalHit      int
	DirectHit        int
	Determination    int
	SkillSpeed       int
	SpellSpeed       int
	Tenacity         int
	Piety            int
}

func GetFromEtro(id string) (*Stats, error) {
	// Request the data from the Etro API, URL: https://etro.gg/api/gearsets/ + id
	res, err := http.Get("https://etro.gg/api/gearsets/" + id)
	if err != nil {
		return &Stats{}, err
	}
	// Parse the JSON response into a struct
	var etroResponse EtroResponse
	err = json.NewDecoder(res.Body).Decode(&etroResponse)
	if err != nil {
		return &Stats{}, err
	}
	// Create a new Stats struct
	stats := &Stats{}
	// Loop through the totalParams array
	for _, param := range etroResponse.TotalParams {
		// Check the name of the param and set the value in the Stats struct
		switch param.Name {
		case "STR":
			stats.MainStat = int(param.Value.(float64))
		case "DEX":
			stats.MainStat = int(param.Value.(float64))
		case "INT":
			stats.MainStat = int(param.Value.(float64))
		case "MND":
			stats.MainStat = int(param.Value.(float64))
		case "VIT":
			stats.Vitality = int(param.Value.(float64))
		case "CRT":
			stats.CriticalHit = int(param.Value.(float64))
		case "DH":
			stats.DirectHit = int(param.Value.(float64))
		case "DET":
			stats.Determination = int(param.Value.(float64))
		case "SKS":
			stats.SkillSpeed = int(param.Value.(float64))
		case "SPS":
			stats.SpellSpeed = int(param.Value.(float64))
		case "TEN":
			stats.Tenacity = int(param.Value.(float64))
		case "PIE":
			stats.Piety = int(param.Value.(float64))
		case "Average Item Level":
			stats.AverageItemLevel = param.Value.(float64)
		case "Weapon Damage":
			stats.WeaponDamage = int(param.Value.(float64))
		}
	}
	// Return the struct
	return stats, nil
}
