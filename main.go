package main

import (
	"github.com/BattlesnakeOfficial/rules/cli/commands"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

type Snake struct {
	Name string
	URL  string
}
type GameParams struct {
	Size    int     `json:"size"`
	Type    string  `json:"type"`
	Timeout int     `json:"timeout"`
	Snakes  []Snake `json:"snakes"`
}

func main() {
	e := echo.New()
	e.POST("/battlesnake", func(c echo.Context) error {
		// start game with given URLs, board size, and game type
		gameParams := GameParams{
			Size:    11,
			Type:    "standard",
			Timeout: 500,
		}
		if err := c.Bind(&gameParams); err != nil {
			return err
		}

		if err := RunGame(gameParams); err != nil {
			return err
		}
		return c.String(http.StatusOK, "Game finished")
	})
	e.Logger.Fatal(e.Start(":8999"))
}

func RunGame(gameParams GameParams) error {
	gameState := &commands.GameState{}
	gameState.Width = gameParams.Size
	gameState.Height = gameParams.Size
	gameState.Timeout = gameParams.Timeout
	gameState.GameType = gameParams.Type
	gameState.MapName = "standard"
	gameState.Seed = time.Now().UTC().UnixNano()
	gameState.FoodSpawnChance = 15 // taken from CLI
	gameState.MinimumFood = 1
	gameState.HazardDamagePerTurn = 14
	gameState.ShrinkEveryNTurns = 25
	for _, s := range gameParams.Snakes {
		gameState.Names = append(gameState.Names, s.Name)
		gameState.URLs = append(gameState.URLs, s.URL)
	}

	if err := gameState.Initialize(); err != nil {
		return err
	}
	return gameState.Run()
}
