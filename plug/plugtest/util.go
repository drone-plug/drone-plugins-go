package plugtest

import (
	"encoding/json"
	"log"
)

func JSON(v interface{}) string {
	data, err := json.Marshal(&v)
	if err != nil {
		log.Fatal(err)
	}
	return string(data)
}
