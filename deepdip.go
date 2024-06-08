package main

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
	"strconv"
	"time"

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

type DDPlayerCurrent struct {
	DisplayName string       `json:"display_name"`
	UserID      string       `json:"user_id"`
	Last5Points [][2]float64 `json:"last_5_points"`
}

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
	http.HandleFunc("/current", getCurrentHeight)
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

func deepDipAPIPlayerCurrent(playerID string) (*DDPlayerCurrent, error) {
	url := fmt.Sprintf("https://dips-plus-plus.xk.io/live_heights/%s", playerID)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var player DDPlayerCurrent
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

func getTimeSince(ts int) string {
	// takes a unix timestamp and returns a string of the time since that timestamp
	// e.g. "1d 2h 3m 4s"

	// get the current time
	now := int(time.Now().Unix())

	// calculate the difference
	diff := now - ts

	// calculate the days, hours, minutes, and seconds
	days := diff / 86400
	diff = diff % 86400
	hours := diff / 3600
	diff = diff % 3600
	minutes := diff / 60
	seconds := diff % 60

	// build the string. Smallest possible format. So if there are any minutes we don't include seconds, if there are any hours we don't include minutes or seconds, etc. Example output: "3 hours ago", "A day ago", "2 days ago", "3 minutes ago", "A minute ago"
	var result string
	if days > 0 {
		if days == 1 {
			result = "A day ago"
		} else {
			result = strconv.Itoa(days) + " days ago"
		}
	} else if hours > 0 {
		if hours == 1 {
			result = "An hour ago"
		} else {
			result = strconv.Itoa(hours) + " hours ago"
		}
	} else if minutes > 0 {
		if minutes == 1 {
			result = "A minute ago"
		} else {
			result = strconv.Itoa(minutes) + " minutes ago"
		}
	} else {
		if seconds == 1 {
			result = "A second ago"
		} else {
			result = strconv.Itoa(seconds) + " seconds ago"
		}
	}

	return result
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
	clean := r.URL.Query().Get("clean")
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
	timeSince := getTimeSince(player.TS)
	if clean == "true" {
		fmt.Fprint(w, strconv.Itoa(roundedHeight))
		return
	}
	fmt.Fprint(w, player.Name+" is rank #"+strconv.Itoa(player.Rank)+" ("+strconv.Itoa(roundedHeight)+"m) ["+timeSince+"]")
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

	// define number of players to ignore
	ignore := r.URL.Query().Get("ignore")
	ignoreInt, err := strconv.Atoi(ignore)
	if err != nil {
		ignoreInt = 0
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
		if i < ignoreInt {
			continue
		}
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
			timeSince := getTimeSince(player.TS)
			userString = "| " + player.Name + "'s PB is rank #" + strconv.Itoa(player.Rank) + " with a height of " + strconv.Itoa(roundedHeight) + "m" + " [" + timeSince + "]"
		}
	} else {
		userString = ""
	}

	fullstring := playersString + userString
	fmt.Fprint(w, fullstring)
}

func getCurrentHeight(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	clean := r.URL.Query().Get("clean")
	playerID, err := tmio.GetPlayerID(username)
	if err != nil {
		fmt.Fprint(w, "Player not found")
		return
	}
	player, err := deepDipAPIPlayerCurrent(playerID)
	if err != nil {
		fmt.Fprint(w, "Player not found on DeepDip API")
		return
	}

	height := player.Last5Points[0][0]
	roundedHeight := int(math.Round(height))
	if clean == "true" {
		fmt.Fprint(w, strconv.Itoa(roundedHeight))
		return
	}
	fmt.Fprint(w, player.DisplayName+" is currently at "+strconv.Itoa(roundedHeight)+"m")
}
