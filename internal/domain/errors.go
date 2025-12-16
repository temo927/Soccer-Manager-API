package domain

import "errors"


var (

	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")


	ErrTeamNotFound       = errors.New("team not found")
	ErrTeamAlreadyExists  = errors.New("team already exists")
	ErrTeamFull           = errors.New("team already has maximum number of players")
	ErrInsufficientBudget = errors.New("insufficient budget")
	ErrCannotBuyOwnPlayer = errors.New("cannot buy your own player")


	ErrPlayerNotFound          = errors.New("player not found")
	ErrPlayerNotOwned          = errors.New("player does not belong to your team")
	ErrPlayerAlreadyListed     = errors.New("player is already on transfer list")
	ErrPlayerNotOnTransferList = errors.New("player is not on transfer list")


	ErrTransferNotFound        = errors.New("transfer not found")
	ErrTransferListingNotFound = errors.New("transfer listing not found")
	ErrInvalidAskingPrice      = errors.New("invalid asking price")
)


type DomainError struct {
	Code    string
	Message string
	Err     error
}

func (e *DomainError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Message
}

func (e *DomainError) Unwrap() error {
	return e.Err
}


func NewDomainError(code, message string, err error) *DomainError {
	return &DomainError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}
