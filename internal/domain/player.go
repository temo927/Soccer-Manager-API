package domain

import (
	"math/rand"
	"time"

	"github.com/google/uuid"
)


type Position string

const (
	PositionGoalkeeper Position = "goalkeeper"
	PositionDefender   Position = "defender"
	PositionMidfielder Position = "midfielder"
	PositionAttacker   Position = "attacker"
)


type Player struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	TeamID      *uuid.UUID `json:"team_id,omitempty" db:"team_id"`
	FirstName   string     `json:"first_name" db:"first_name"`
	LastName    string     `json:"last_name" db:"last_name"`
	Country     string     `json:"country" db:"country"`
	Age         int        `json:"age" db:"age"`
	Position    Position   `json:"position" db:"position"`
	MarketValue float64    `json:"market_value" db:"market_value"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}

const (
	InitialPlayerValue = 1000000.00
	MinAge             = 18
	MaxAge             = 40
)


func NewPlayer(teamID *uuid.UUID, firstName, lastName, country string, position Position) *Player {
	rand.Seed(time.Now().UnixNano())
	age := MinAge + rand.Intn(MaxAge-MinAge+1)

	return &Player{
		ID:          uuid.New(),
		TeamID:      teamID,
		FirstName:   firstName,
		LastName:    lastName,
		Country:     country,
		Age:         age,
		Position:    position,
		MarketValue: InitialPlayerValue,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}


func (p *Player) UpdateMarketValue() {
	rand.Seed(time.Now().UnixNano())
	increasePercent := 0.10 + rand.Float64()*0.90
	p.MarketValue = p.MarketValue * (1 + increasePercent)
	p.UpdatedAt = time.Now()
}


func (p *Player) Transfer(newTeamID uuid.UUID) {
	p.TeamID = &newTeamID
	p.UpdateMarketValue()
	p.UpdatedAt = time.Now()
}


func (p *Player) IsOwnedBy(teamID uuid.UUID) bool {
	return p.TeamID != nil && *p.TeamID == teamID
}
