package tests

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
	"github.com/nicholas-kebbas/cs686-blockchain-p3-nicholas-kebbas/p2"
)

func TestBlockChainBasic(t *testing.T) {
	jsonBlockChain := "[{\"hash\":\"2d4017fd7181ea93899f43f1fbc38a21543456a35b12a4600b689d947f19b089\",\"timestamp\":123,\"height\":2,\"parentHash\":\"49de64b587bac530c06762c0e7fe55785a8ae0d5227395bee0f59a5679b9ddb4\",\"size\":242,\"signature_p\":\"\",\"publickey\":\"\\ufffd3\\u0017QA\\u0015\\ufffdd\\ufffd_\\ufffd\\u0013\\ufffd\\ufffd)\\u0002T\\ufffd\\ufffdg\\ufffd\\ufffd\\ufffd\\ufffd0\\u0002\\ufffd\\u0018^\\ufffd\\ufffd\\u0003\\u0007\\ufffd\\ufffdQ\\u0017\\ufffd\\ufffd\\\"\\u0002\\ufffdsUNT\\ufffd(\\ufffd\\ufffd}\\ufffd\\ufffd\\u001a\\ufffd\\ufffd)Ly~\\ufffd8\\ufffd\",\"mpt\":{\"1\":\"Hilary Clinton\"}}, {\"hash\":\"2d4017fd7181ea93899f43f1fbc38a21543456a35b12a4600b689d947f19b089\",\"timestamp\":123,\"height\":2,\"parentHash\":\"49de64b587bac530c06762c0e7fe55785a8ae0d5227395bee0f59a5679b9ddb4\",\"size\":242,\"signature_p\":\"\",\"publickey\":\"\\ufffd3\\u0017QA\\u0015\\ufffdd\\ufffd_\\ufffd\\u0013\\ufffd\\ufffd)\\u0002T\\ufffd\\ufffdg\\ufffd\\ufffd\\ufffd\\ufffd0\\u0002\\ufffd\\u0018^\\ufffd\\ufffd\\u0003\\u0007\\ufffd\\ufffdQ\\u0017\\ufffd\\ufffd\\\"\\u0002\\ufffdsUNT\\ufffd(\\ufffd\\ufffd}\\ufffd\\ufffd\\u001a\\ufffd\\ufffd)Ly~\\ufffd8\\ufffd\",\"mpt\":{\"1\":\"Donald Trump\"}}]"
	bc := p2.NewBlockChain()
	p2.DecodeFromJSON(&bc, jsonBlockChain)
	jsonNew := bc.EncodeToJSON()
	var realValue []p2.JsonBlock
	var expectedValue []p2.JsonBlock
	err := json.Unmarshal([]byte(jsonNew), &realValue)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	err = json.Unmarshal([]byte(jsonBlockChain), &expectedValue)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	if !reflect.DeepEqual(realValue, expectedValue) {
		fmt.Println("=========Real=========")
		fmt.Println(realValue)
		fmt.Println("=========Expcected=========")
		fmt.Println(expectedValue)
		t.Fail()
	}
}

func TestBlockChainBasic2(t *testing.T) {
	jsonBlockChain := "[{\"height\":1,\"timeStamp\":1551025401,\"hash\":\"6c9aad47a370269746f172a464fa6745fb3891194da65e3ad508ccc79e9a771b\",\"parentHash\":\"genesis\",\"size\":2089,\"mpt\":{\"CS686\":\"BlockChain\",\"test1\":\"value1\",\"test2\":\"value2\",\"test3\":\"value3\",\"test4\":\"value4\"}},{\"height\":2,\"timeStamp\":1551025401,\"hash\":\"944eb943b05caba08e89a613097ac5ac7d373d863224d17b1958541088dc20e2\",\"parentHash\":\"6c9aad47a370269746f172a464fa6745fb3891194da65e3ad508ccc79e9a771b\",\"size\":2146,\"mpt\":{\"CS686\":\"BlockChain\",\"test1\":\"value1\",\"test2\":\"value2\",\"test3\":\"value3\",\"test4\":\"value4\"}},{\"height\":2,\"timeStamp\":1551025401,\"hash\":\"f8af68feadf25a635bc6e81c08f81c6740bbe1fb2514c1b4c56fe1d957c7448d\",\"parentHash\":\"6c9aad47a370269746f172a464fa6745fb3891194da65e3ad508ccc79e9a771b\",\"size\":707,\"mpt\":{\"ge\":\"Charles\"}},{\"height\":3,\"timeStamp\":1551025401,\"hash\":\"f367b7f59c651e69be7e756298aad62fb82fddbfeda26cb06bfd8adf9c8aa094\",\"parentHash\":\"f8af68feadf25a635bc6e81c08f81c6740bbe1fb2514c1b4c56fe1d957c7448d\",\"size\":707,\"mpt\":{\"ge\":\"Charles\"}},{\"height\":3,\"timeStamp\":1551025401,\"hash\":\"05ac44dd82b6cc398a5e9664add21856ae19d107d9035af5fc54c9b0ffdef336\",\"parentHash\":\"944eb943b05caba08e89a613097ac5ac7d373d863224d17b1958541088dc20e2\",\"size\":2146,\"mpt\":{\"CS686\":\"BlockChain\",\"test1\":\"value1\",\"test2\":\"value2\",\"test3\":\"value3\",\"test4\":\"value4\"}}]"
	bc := p2.NewBlockChain()
	p2.DecodeFromJSON(&bc, jsonBlockChain)
	jsonNew := bc.EncodeToJSON()
	var realValue []p2.JsonBlock
	var expectedValue []p2.JsonBlock
	err := json.Unmarshal([]byte(jsonNew), &realValue)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	err = json.Unmarshal([]byte(jsonBlockChain), &expectedValue)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	if !reflect.DeepEqual(realValue, expectedValue) {
		fmt.Println("=========Real=========")
		fmt.Println(realValue)
		fmt.Println("=========Expcected=========")
		fmt.Println(expectedValue)
		t.Fail()
	}

}

func TestBlock(t *testing.T) {

}