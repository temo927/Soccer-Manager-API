package auth

import (
	"context"
	"math/rand"
	"time"

	"soccer-manager-api/internal/domain"
	"soccer-manager-api/internal/ports/repository"
	"soccer-manager-api/pkg/jwt"
	"soccer-manager-api/pkg/password"

	"github.com/google/uuid"
)


type AuthUseCase struct {
	userRepo   repository.UserRepository
	teamRepo   repository.TeamRepository
	playerRepo repository.PlayerRepository
	jwtSecret  string
	jwtExpHours int
}


func NewAuthUseCase(
	userRepo repository.UserRepository,
	teamRepo repository.TeamRepository,
	playerRepo repository.PlayerRepository,
	jwtSecret string,
	jwtExpHours int,
) *AuthUseCase {
	return &AuthUseCase{
		userRepo:    userRepo,
		teamRepo:    teamRepo,
		playerRepo:  playerRepo,
		jwtSecret:   jwtSecret,
		jwtExpHours: jwtExpHours,
	}
}


type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}


type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}


type AuthResponse struct {
	Token string      `json:"token"`
	User  *domain.User `json:"user"`
}


func (uc *AuthUseCase) Register(ctx context.Context, req RegisterRequest) (*AuthResponse, error) {

	existingUser, _ := uc.userRepo.GetByEmail(ctx, req.Email)
	if existingUser != nil {
		return nil, domain.ErrUserAlreadyExists
	}


	passwordHash, err := password.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}


	user := domain.NewUser(req.Email, passwordHash)
	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}


	teamName := generateDefaultTeamName()
	teamCountry := "Unknown"
	team := domain.NewTeam(user.ID, teamName, teamCountry)
	if err := uc.teamRepo.Create(ctx, team); err != nil {
		return nil, err
	}


	players := uc.generateInitialPlayers(team.ID)
	if err := uc.playerRepo.CreateBatch(ctx, players); err != nil {
		return nil, err
	}


	token, err := jwt.GenerateToken(user.ID, user.Email, uc.jwtSecret, uc.jwtExpHours)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		Token: token,
		User:  user,
	}, nil
}


func (uc *AuthUseCase) Login(ctx context.Context, req LoginRequest) (*AuthResponse, error) {

	user, err := uc.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, domain.ErrInvalidCredentials
	}


	if !password.CheckPasswordHash(req.Password, user.PasswordHash) {
		return nil, domain.ErrInvalidCredentials
	}


	token, err := jwt.GenerateToken(user.ID, user.Email, uc.jwtSecret, uc.jwtExpHours)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		Token: token,
		User:  user,
	}, nil
}


func (uc *AuthUseCase) generateInitialPlayers(teamID uuid.UUID) []*domain.Player {
	players := make([]*domain.Player, 0, 20)


	for i := 0; i < 3; i++ {
		players = append(players, domain.NewPlayer(
			&teamID,
			generateRandomFirstName(),
			generateRandomLastName(),
			generateRandomCountry(),
			domain.PositionGoalkeeper,
		))
	}


	for i := 0; i < 6; i++ {
		players = append(players, domain.NewPlayer(
			&teamID,
			generateRandomFirstName(),
			generateRandomLastName(),
			generateRandomCountry(),
			domain.PositionDefender,
		))
	}


	for i := 0; i < 6; i++ {
		players = append(players, domain.NewPlayer(
			&teamID,
			generateRandomFirstName(),
			generateRandomLastName(),
			generateRandomCountry(),
			domain.PositionMidfielder,
		))
	}


	for i := 0; i < 5; i++ {
		players = append(players, domain.NewPlayer(
			&teamID,
			generateRandomFirstName(),
			generateRandomLastName(),
			generateRandomCountry(),
			domain.PositionAttacker,
		))
	}

	return players
}

func generateDefaultTeamName() string {
	names := []string{
		"FC United", "City FC", "Athletic Club", "Sporting FC",
		"United FC", "City United", "Athletic United", "Sporting Club",
	}
	rand.Seed(time.Now().UnixNano())
	return names[rand.Intn(len(names))]
}

func generateRandomFirstName() string {
	names := []string{
		"John", "James", "Michael", "David", "Robert", "William", "Richard", "Joseph",
		"Thomas", "Charles", "Christopher", "Daniel", "Matthew", "Anthony", "Mark",
		"Donald", "Steven", "Paul", "Andrew", "Joshua", "Kenneth", "Kevin", "Brian",
	}
	rand.Seed(time.Now().UnixNano())
	return names[rand.Intn(len(names))]
}

func generateRandomLastName() string {
	names := []string{
		"Smith", "Johnson", "Williams", "Brown", "Jones", "Garcia", "Miller", "Davis",
		"Rodriguez", "Martinez", "Hernandez", "Lopez", "Wilson", "Anderson", "Thomas",
		"Taylor", "Moore", "Jackson", "Martin", "Lee", "Thompson", "White", "Harris",
	}
	rand.Seed(time.Now().UnixNano())
	return names[rand.Intn(len(names))]
}

func generateRandomCountry() string {
	countries := []string{
		"Brazil", "Argentina", "Spain", "Germany", "France", "Italy", "England",
		"Portugal", "Netherlands", "Belgium", "Croatia", "Uruguay", "Colombia",
		"Mexico", "Chile", "Poland", "Denmark", "Sweden", "Norway", "Greece",
	}
	rand.Seed(time.Now().UnixNano())
	return countries[rand.Intn(len(countries))]
}

