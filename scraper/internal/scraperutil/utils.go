package scraperutil

import (
	"encoding/json"
	"os"
	"strconv"
	"time"

	"fightbettr.com/scraper/pkg/logger"
	"fightbettr.com/scraper/pkg/model"
)

func CreateNewCollection(c model.FightersCollection) {
	l := logger.Get()

	file, err := os.Create("./collection/fighters.json")
	if err != nil {
		l.Error("File creation error:", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(c); err != nil {
		l.Error("Error encoding JSON:", err)
		return
	}
}

func AddToExistedCollection(c model.FightersCollection) {
	l := logger.Get()

	file, err := os.Open("./collection/fighters.json")
	if err != nil {
		l.Error("File opening error:", err)
		return
	}
	defer file.Close()

	var existingFighters model.FightersCollection
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&existingFighters); err != nil {
		l.Error("Error decoding JSON:", err)
		return
	}

	existingFighters.Fighters = append(existingFighters.Fighters, c.Fighters...)
	collection := getUniqueCollection(existingFighters)

	file, err = os.Create("./collection/fighters.json")
	if err != nil {
		l.Error("File creation error:", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(collection); err != nil {
		l.Error("Error encoding JSON:", err)
		return
	}
}

func getUniqueCollection(c model.FightersCollection) model.FightersCollection {
	uniqueFightersMap := make(map[string]model.Fighter)
	uniqueFighters := make([]model.Fighter, 0, 500)

	for _, fighter := range c.Fighters {
		key := fighter.Name + fighter.NickName + strconv.Itoa(fighter.DebutTimestamp)
		if _, exists := uniqueFightersMap[key]; !exists {
			uniqueFightersMap[key] = fighter
			uniqueFighters = append(uniqueFighters, fighter)
		}
	}

	return model.FightersCollection{
		Fighters: uniqueFighters,
	}
}

func GetLoggerFlag(toAdd bool) int {
	if toAdd {
		return os.O_APPEND
	} else {
		return os.O_TRUNC
	}
}

// getDebutTimestamp parses the octagon debut date provided in the format "Jan. 2, 2006"
// and returns its Unix timestamp.This function is used to convert the debut date of a fighter in the octagon
// into a Unix timestamp for further processing.
func GetDebutTimestamp(octagonDebut string) int {
	l := logger.Get()

	layout := "Jan. 2, 2006"

	parsedTime, err := time.Parse(layout, octagonDebut)
	if err != nil {
		l.Error("Error while date parsing:", err)
		return 0
	}

	return int(parsedTime.Unix())
}
