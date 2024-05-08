package main

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
	"strconv"

	tmio "github.com/Johnnycyan/go-tmio-sdk"
)

type DDPlayer struct {
	Height      float64 `json:"height"`
	Name        string  `json:"name"`
	Rank        int     `json:"rank"`
	TS          int     `json:"ts"`
	UpdateCount int     `json:"update_count"`
	WSID        string  `json:"wsid"`
}

type Leaderboard []DDPlayer

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		fmt.Println("Usage: deepdip <port>")
		return
	}
	port := args[0]
	fmt.Println("Listening on http://localhost:" + port)
	http.HandleFunc("/pb", getPB)
	http.HandleFunc("/leaderboards", getLeaderboards)
	http.ListenAndServe(":"+port, nil)
}

func deepDipAPIPlayer(playerID string) (*DDPlayer, error) {
	url := fmt.Sprintf("https://dips-plus-plus.xk.io/leaderboard/%s", playerID)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var player DDPlayer
	if err := json.NewDecoder(resp.Body).Decode(&player); err != nil {
		return nil, err
	}
	return &player, nil
}

func deepDipAPILeaderboard() (*Leaderboard, error) {
	url := "https://dips-plus-plus.xk.io/leaderboard/global"
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var leaderboard Leaderboard
	if err := json.NewDecoder(resp.Body).Decode(&leaderboard); err != nil {
		return nil, err
	}
	return &leaderboard, nil
}

func getPB(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprint(w, "User not found")
		}
	}()
	username := r.URL.Query().Get("username")
	if username == "" {
		http.Error(w, "username is required", http.StatusBadRequest)
		return
	}
	playerID, err := tmio.GetPlayerID(username)
	if err != nil {
		fmt.Fprint(w, "Player not found")
		return
	}
	player, err := deepDipAPIPlayer(playerID)
	if err != nil {
		fmt.Fprint(w, "Player not found on DeepDip API")
		return
	}

	if player.Rank == 0 {
		fmt.Fprint(w, "Player not found on DeepDip API")
		return
	}

	roundedHeight := int(math.Round(player.Height))
	fmt.Fprint(w, player.Name+" is rank #"+strconv.Itoa(player.Rank)+" ("+strconv.Itoa(roundedHeight)+"m)")
}

func getLeaderboards(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprint(w, "User not found")
		}
	}()
	username := r.URL.Query().Get("username")
	var usernameExists bool
	if username == "" {
		usernameExists = false
	} else {
		usernameExists = true
	}

	leaderboard, err := deepDipAPILeaderboard()
	if err != nil {
		fmt.Fprint(w, "Leaderboard not found PANIC")
		return
	}

	var medal string
	var playersString string
	for i, player := range *leaderboard {
		if i >= 3 {
			break
		} else if i == 0 {
			medal = "🥇 "
		} else if i == 1 {
			medal = "🥈 "
		} else if i == 2 {
			medal = "🥉 "
		} else {
			medal = ""
		}
		roundedHeight := int(math.Round(player.Height))
		playersString += medal + player.Name + " (" + strconv.Itoa(roundedHeight) + "m) "
	}

	var player *DDPlayer
	if usernameExists {
		playerID, err := tmio.GetPlayerID(username)
		if err != nil {
			usernameExists = false
		}
		player, err = deepDipAPIPlayer(playerID)
		if err != nil {
			usernameExists = false
		}
	}

	var userString string
	if usernameExists {
		if player.Rank == 0 {
			userString = ""
		} else {
			roundedHeight := int(math.Round(player.Height))
			userString = "| " + player.Name + "'s PB is rank #" + strconv.Itoa(player.Rank) + " with a height of " + strconv.Itoa(roundedHeight) + "m"
		}
	} else {
		userString = ""
	}

	fullstring := "Current standing: " + playersString + userString
	fmt.Fprint(w, fullstring)
}
