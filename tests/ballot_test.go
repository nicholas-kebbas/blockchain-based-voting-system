package tests

import (
	"fmt"
	"github.com/nicholas-kebbas/cs686-blockchain-p3-nicholas-kebbas/voting"
	"reflect"
	"testing"
)

func TestBallot(t *testing.T) {
	ballotJson := "{\"Candidates\":{\"1\":\"Donald Trump\",\"2\":\"Hilary Clinton\"}}"
	jsonBallot:= voting.FromJson(ballotJson)
	realBallot := voting.Initial()
	realBallot.AddChoice("1", "Donald Trump")
	realBallot.AddChoice("2", "Hilary Clinton")
	realBallotAsJson := realBallot.ToJson()


	if !reflect.DeepEqual(jsonBallot, realBallotAsJson) {
		fmt.Println("=========Real=========")
		fmt.Println(jsonBallot)
		fmt.Println("=========Expected=========")
		fmt.Println(realBallotAsJson)
		t.Fail()
	}
}