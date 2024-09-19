package elo

import (
	"math"
	"sort"
)

const (
	kFactor       = 15
	homeAdvantage = 50
)

type Team struct {
	Name   string
	Rating float64
}

type Game struct {
	Date      string
	HomeTeam  string
	AwayTeam  string
	HomeScore int
	AwayScore int
	HomeWin   bool
}

func CalculateElo(teams map[string]*Team, games []Game) {
	for _, game := range games {
		homeTeam := teams[game.HomeTeam]
		awayTeam := teams[game.AwayTeam]

		homeRating := homeTeam.Rating + homeAdvantage
		awayRating := awayTeam.Rating

		homeExpected := expectedScore(homeRating, awayRating)
		awayExpected := expectedScore(awayRating, homeRating)

		var homeActual, awayActual float64
		if game.HomeWin {
			homeActual = 1
			awayActual = 0
		} else {
			homeActual = 0
			awayActual = 1
		}

		homeTeam.Rating += kFactor * (homeActual - homeExpected)
		awayTeam.Rating += kFactor * (awayActual - awayExpected)
	}
}

func expectedScore(rating1, rating2 float64) float64 {
	return 1 / (1 + math.Pow(10, (rating2-rating1)/400))
}

func CalculateBrierScore(teams map[string]*Team, games []Game) float64 {
	var totalSquaredError float64
	for _, game := range games {
		homeTeam := teams[game.HomeTeam]
		awayTeam := teams[game.AwayTeam]

		homeRating := homeTeam.Rating + homeAdvantage
		awayRating := awayTeam.Rating

		homeProbability := expectedScore(homeRating, awayRating)

		var actualOutcome float64
		if game.HomeWin {
			actualOutcome = 1
		}

		squaredError := math.Pow(homeProbability-actualOutcome, 2)
		totalSquaredError += squaredError
	}

	return totalSquaredError / float64(len(games))
}

func GetSortedTeams(teams map[string]*Team) []*Team {
	var sortedTeams []*Team
	for _, team := range teams {
		sortedTeams = append(sortedTeams, team)
	}
	sort.Slice(sortedTeams, func(i, j int) bool {
		return sortedTeams[i].Rating > sortedTeams[j].Rating
	})
	return sortedTeams
}

// package elo
//
// import (
// 	"math"
// 	"sort"
// )
//
// const (
// 	startKFactor  = 32
// 	endKFactor    = 10
// 	decayRate     = 0.01 // Adjust this rate to control how fast the K-factor decays
// 	homeAdvantage = 70
// )
//
// type Team struct {
// 	Name   string
// 	Rating float64
// }
//
// type Game struct {
// 	Date      string
// 	HomeTeam  string
// 	AwayTeam  string
// 	HomeScore int
// 	AwayScore int
// 	HomeWin   bool
// }
//
// func CalculateElo(teams map[string]*Team, games []Game) {
// 	for i, game := range games {
// 		kFactor := decayedKFactor(i)
//
// 		homeTeam := teams[game.HomeTeam]
// 		awayTeam := teams[game.AwayTeam]
//
// 		homeRating := homeTeam.Rating + homeAdvantage
// 		awayRating := awayTeam.Rating
//
// 		homeExpected := expectedScore(homeRating, awayRating)
// 		awayExpected := expectedScore(awayRating, homeRating)
//
// 		var homeActual, awayActual float64
// 		if game.HomeWin {
// 			homeActual = 1
// 			awayActual = 0
// 		} else {
// 			homeActual = 0
// 			awayActual = 1
// 		}
//
// 		homeTeam.Rating += kFactor * (homeActual - homeExpected)
// 		awayTeam.Rating += kFactor * (awayActual - awayExpected)
// 	}
// }
//
// func decayedKFactor(gameIndex int) float64 {
// 	// Exponential decay formula
// 	return endKFactor + (startKFactor-endKFactor)*math.Exp(-decayRate*float64(gameIndex))
// }
//
// func expectedScore(rating1, rating2 float64) float64 {
// 	return 1 / (1 + math.Pow(10, (rating2-rating1)/400))
// }
//
// func CalculateBrierScore(teams map[string]*Team, games []Game) float64 {
// 	var totalSquaredError float64
// 	for _, game := range games {
// 		homeTeam := teams[game.HomeTeam]
// 		awayTeam := teams[game.AwayTeam]
//
// 		homeRating := homeTeam.Rating + homeAdvantage
// 		awayRating := awayTeam.Rating
//
// 		homeProbability := expectedScore(homeRating, awayRating)
//
// 		var actualOutcome float64
// 		if game.HomeWin {
// 			actualOutcome = 1
// 		}
//
// 		squaredError := math.Pow(homeProbability-actualOutcome, 2)
// 		totalSquaredError += squaredError
// 	}
//
// 	return totalSquaredError / float64(len(games))
// }
//
// func GetSortedTeams(teams map[string]*Team) []*Team {
// 	var sortedTeams []*Team
// 	for _, team := range teams {
// 		sortedTeams = append(sortedTeams, team)
// 	}
// 	sort.Slice(sortedTeams, func(i, j int) bool {
// 		return sortedTeams[i].Rating > sortedTeams[j].Rating
// 	})
// 	return sortedTeams
// }
