package backend

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
)

func Persist(tabs Tables) {
	os.MkdirAll(baseDataPath, 0755)
	for cID, m := range tabs {
		jsonData, err := json.Marshal(m)
		if err != nil {
			log.Fatal(err)
		}
		filePath := baseDataPath + "/" + strconv.Itoa(int(cID)) + ".json"
		file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0644)
		defer file.Close()
		if err != nil {
			log.Fatal(err)
		}
		n, err := file.Write(jsonData)
		if err != nil {
			log.Fatal(err)
		}
		err = file.Truncate(int64(n))
	}
}
