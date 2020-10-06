package main

import (
	"encoding/json"
	"fmt"
	"github.com/rs/cors"
	"net/http"
	"sync"
)

//func beacon(res http.ResponseWriter, req *http.Request) {
//	if n, err := io.Copy(res, req.Body); err != nil {
//		fmt.Println(n, err)
//	}
//}

type Username string

type description struct {
	Type string `json:"type"`
	Sdp  string `json:"sdp"`
}

type candidate struct {
	Candidate        string `json:"candidate"`
	SdpMid           string `json:"sdpMid"`
	SdpMLineIndex    uint   `json:"sdpMLineIndex"`
	UsernameFragment string `json:"usernameFragment"`
}

type peerInfo struct {
	Description *description `json:"description"`
	Candidates  []candidate  `json:"candidates"`
}

var mt sync.Mutex
var desMap map[Username]description
var candidateMap map[Username][]candidate

func init() {
	desMap = map[Username]description{}
	candidateMap = map[Username][]candidate{}
}

func getRemoteInfo(res http.ResponseWriter, req *http.Request) {
	mt.Lock()
	defer mt.Unlock()
	query := req.URL.Query()
	username := Username(query.Get("username"))

	var result peerInfo

	description, ok1 := desMap[username]
	if ok1 {
		result.Description = &description
		delete(desMap, username)
	}

	candidates, ok2 := candidateMap[username]
	if ok1 || ok2 {
		result.Candidates = candidates
		delete(candidateMap, username)
	}

	fmt.Println("reading of", username, "is", result)
	bs, _ := json.Marshal(result)
	_, _ = res.Write(bs)
}

func registerDescription(res http.ResponseWriter, req *http.Request) {
	mt.Lock()
	defer mt.Unlock()
	query := req.URL.Query()
	username := Username(query.Get("username"))

	des := description{}
	_ = json.NewDecoder(req.Body).Decode(&des)

	fmt.Println("register description", username, des)
	desMap[username] = des

	_, _ = res.Write(nil)
}

func registerCandidate(res http.ResponseWriter, req *http.Request) {
	mt.Lock()
	defer mt.Unlock()
	query := req.URL.Query()
	username := Username(query.Get("username"))

	cadi := candidate{}
	_ = json.NewDecoder(req.Body).Decode(&cadi)

	fmt.Println("register candidate", username, candidate{})
	candidateMap[username] = append(candidateMap[username], cadi)

	_, _ = res.Write(nil)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/read", getRemoteInfo)
	mux.HandleFunc("/post/description", registerDescription)
	mux.HandleFunc("/post/candidate", registerCandidate)

	handler := cors.Default().Handler(mux)
	if err := http.ListenAndServe(":9999", handler); err != nil {
		panic(err)
	}
}
