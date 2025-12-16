package domain

import (
	"time"

	"github.com/google/uuid"
)


type Team struct {
	ID        uuid.UUID `json:"id" db:"id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	Name      string    `json:"name" db:"name"`
	Country   string    `json:"country" db:"country"`
	Budget    float64   `json:"budget" db:"budget"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}


type TeamWithValue struct {
	Team
	TotalValue float64 `json:"total_value"`
}

const (
	InitialBudget = 5000000.00
	MaxPlayers    = 20
)


func NewTeam(userID uuid.UUID, name, country string) *Team {
	return &Team{
		ID:        uuid.New(),
		UserID:    userID,
		Name:      name,
		Country:   country,
		Budget:    InitialBudget,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}


func (t *Team) CanAfford(price float64) bool {
	return t.Budget >= price
}


func (t *Team) DeductBudget(amount float64) {
	t.Budget -= amount
	t.UpdatedAt = time.Now()
}


func (t *Team) AddBudget(amount float64) {
	t.Budget += amount
	t.UpdatedAt = time.Now()
}
