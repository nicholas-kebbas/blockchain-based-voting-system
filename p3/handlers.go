package p3

import (
	//"../p2"
	"./data"
	//"encoding/json"
	// "errors"
	"fmt"
	//"github.com/gorilla/mux"
	//"io"
	//"io/ioutil"
	//"math/rand"
	"net/http"
	//"os"
	//"strconv"
	//"strings"
	//"time"
)

var TA_SERVER = "http://localhost:6688"
var REGISTER_SERVER = TA_SERVER + "/peer"
var BC_DOWNLOAD_SERVER = TA_SERVER + "/upload"
var SELF_ADDR = "http://localhost:6686"

var SBC data.SyncBlockChain
var Peers data.PeerList
var ifStarted bool

func init() {
	// This function will be executed before everything else.
	// Do some initialization here.
}

// Register ID, download BlockChain, start HeartBeat
func Start(w http.ResponseWriter, r *http.Request) {}

// Display peerList and sbc
func Show(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s\n%s", Peers.Show(), SBC.Show())
}

// Register to TA's server, get an ID
func Register() {}

// Download blockchain from TA server
func Download() {}

// Upload blockchain to whoever called this method, return jsonStr
func Upload(w http.ResponseWriter, r *http.Request) {
	blockChainJson, err := SBC.BlockChainToJson()
	if err != nil {
		// data.PrintError(err, "Upload")
	}
	fmt.Fprint(w, blockChainJson)
}

// Upload a block to whoever called this method, return jsonStr
func UploadBlock(w http.ResponseWriter, r *http.Request) {}

// Received a heartbeat
func HeartBeatReceive(w http.ResponseWriter, r *http.Request) {}

// Ask another server to return a block of certain height and hash
func AskForBlock(height int32, hash string) {}

func ForwardHeartBeat(heartBeatData data.HeartBeatData) {}

func StartHeartBeat() {}