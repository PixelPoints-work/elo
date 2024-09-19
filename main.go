package main

import (
	"fmt"
	"log"
	"nba-elo-rating-v-2/elo"
	"os"
	"path/filepath"
	"strings"

	"github.com/xuri/excelize/v2"
)

func main() {
	// Initialize teams and games
	teams := make(map[string]*elo.Team)
	var allGames []elo.Game

	// Read all Excel files in the data directory
	files, err := os.ReadDir("data")
	if err != nil {
		log.Fatalf("Unable to read data directory: %v", err)
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".xlsx") {
			filePath := filepath.Join("data", file.Name())
			games := readGamesFromFile(filePath, teams)
			allGames = append(allGames, games...)
		}
	}

	// Calculate Elo ratings
	elo.CalculateElo(teams, allGames)

	// Get sorted teams
	sortedTeams := elo.GetSortedTeams(teams)

	// Print Elo ratings
	fmt.Println("NBA Team Elo Ratings:")
	for _, team := range sortedTeams {
		fmt.Printf("%s: %.2f\n", team.Name, team.Rating)
	}

	// Calculate and print Brier score
	brierScore := elo.CalculateBrierScore(teams, allGames)
	fmt.Printf("\nBrier Score: %.4f\n", brierScore)
}

func readGamesFromFile(filePath string, teams map[string]*elo.Team) []elo.Game {
	var games []elo.Game

	// Open the Excel file
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		log.Printf("Unable to open file %s: %v", filePath, err)
		return games
	}
	defer f.Close()

	// Read the game data
	rows, err := f.GetRows("Sheet1")
	if err != nil {
		log.Printf("Unable to read rows from file %s: %v", filePath, err)
		return games
	}

	// Process each row (skipping the header)
	for _, row := range rows[1:] {
		if len(row) < 6 {
			continue // Skip rows with insufficient data
		}

		homeTeam := row[4]
		awayTeam := row[2]
		homeScore := parseInt(row[5])
		awayScore := parseInt(row[3])

		// Initialize teams if they don't exist
		if _, exists := teams[homeTeam]; !exists {
			teams[homeTeam] = &elo.Team{Name: homeTeam, Rating: 1500}
		}
		if _, exists := teams[awayTeam]; !exists {
			teams[awayTeam] = &elo.Team{Name: awayTeam, Rating: 1500}
		}

		// Add game to the list
		games = append(games, elo.Game{
			Date:      row[0],
			HomeTeam:  homeTeam,
			AwayTeam:  awayTeam,
			HomeScore: homeScore,
			AwayScore: awayScore,
			HomeWin:   homeScore > awayScore,
		})
	}
	return games
}

func parseInt(s string) int {
	var result int
	fmt.Sscanf(s, "%d", &result)
	return result
}
