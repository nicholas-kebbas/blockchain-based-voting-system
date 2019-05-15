package voting

import (
	"encoding/json"
	"fmt"
)

type Ballot struct {
	Candidates map[string]string
}

func Initial() Ballot {
	ballot := Ballot{}
	ballot.Candidates = make(map[string]string)
	return ballot
}

func (ballot *Ballot) AddChoice(key string, value string) {
	ballot.Candidates[key] = value
}

func (ballot *Ballot) ToJson() string {
	encodedString, err := json.Marshal(ballot)
	if err != nil {
		fmt.Println("Error in Encode to Json in Block")
		fmt.Println(err)
		return ""
	}
	fmt.Println("Ballot To JSON")
	fmt.Println(string(encodedString))
	return string(encodedString)
}

func FromJson(jsonString string) Ballot {
	/* Convert the JSON string to bytes */
	bytes := []byte(jsonString)
	ballot := Initial()
	err := json.Unmarshal(bytes, &ballot)
	if err != nil {
		fmt.Println("Error in DecodeFromJson in Block")
	}
	return ballot
}

