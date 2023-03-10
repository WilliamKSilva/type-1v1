package api

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
)

type mockGameRepository struct {
    mock.Mock
}

type mockTextService struct {
    mock.Mock
}

func (m *mockGameRepository) Create (game *Game) error {
   args := m.Called(game)
   
   game.ID = 1 

   return args.Error(1) 
}

func (m *mockGameRepository) Update (id uint, data UpdateGameData) (*Game, error) {
   args := m.Called(id, data)

   game := args.Get(0).(*Game)

   return game, args.Error(1) 
}

func (m *mockGameRepository) Find (id uint) (*Game, error) {
    args := m.Called(id)

    if args.Get(0) != nil {
        game := args.Get(0).(*Game) 

        return game, nil
    }

    return nil, args.Error(1)
}

func (m *mockTextService) GetRandomText (trigger string) (string, error) {
   args := m.Called(trigger)

   return args.String(0), args.Error(1)
}

func TestShouldThrowIfPlayerOneNameIsMissing(t *testing.T) {
    repo := new(mockGameRepository)
    textService := new(mockTextService)

    repo.On("Create", mock.Anything).Return(nil, nil)
    textService.On("GetRandomText", mock.Anything).Return("bla bla bla bla", nil)

    g := gameService{repo, textService}

    gameData := NewGameData{}

    want := errors.New("Player one name is required")
    _, err := g.NewGame(gameData)

    if err.Error() != want.Error() {
        t.Errorf("Expected: %s, got %s", want, err)
    } 
}

func TestShouldThrowIfTextServiceThrows(t *testing.T) {
    repo := new(mockGameRepository)
    textService := new(mockTextService)

    repo.On("Create", mock.Anything).Return(nil, nil)
    textService.On("GetRandomText", mock.Anything).Return("", errors.New("Internal Server Error"))

    g := gameService{repo, textService}

    gameData := NewGameData{
        PlayerOne: "test",
    }

    want := errors.New("Internal Server Error")
    _, err := g.NewGame(gameData)

    if err.Error() != want.Error() {
        t.Errorf("Expected: %s, got %s", want, err)
    }
}

func TestShouldReturnAGameOnSuccess(t *testing.T) {
    repo := new(mockGameRepository)
    textService := new(mockTextService)

    repo.On("Create", mock.Anything).Return(nil, nil)
    textService.On("GetRandomText", mock.Anything).Return("bla bla bla bla", nil)

    g := gameService{repo, textService}

    gameData := NewGameData{
        PlayerOne: "test",
    }

    game, _ := g.NewGame(gameData)
    want := &Game{
        ID: 1,
        PlayerOne: "test",
        Status: Waiting,
        Text: "bla bla bla bla",
        Winner: "",
    }

    if game.ID != want.ID {
        t.Errorf("Expected: %d, got %d", want.ID, game.ID)
    } 
}

func TestUpdateShouldThrowIfNoDataIsProvided(t *testing.T) {
    repo := new(mockGameRepository)
    textService := new(mockTextService)

    repo.On("Create", mock.Anything).Return(nil, nil)
    textService.On("GetRandomText", mock.Anything).Return("bla bla bla bla", nil)

    g := gameService{repo, textService}

    id := uint(1) 
    updateData := UpdateGameData{}

    want := errors.New("Provide valid data to update an Game")
    _, err := g.UpdateGame(id, updateData)

    if err.Error() != want.Error() {
        t.Errorf("Expected: %s, got %s", want.Error(), err.Error())
    }
}

func TestFindShouldThrowIfNoIdIsProvided(t *testing.T) {
    repo := new(mockGameRepository)
    textService := new(mockTextService)

    repo.On("Find", mock.Anything).Return(nil, nil)
    textService.On("GetRandomText", mock.Anything).Return("bla bla bla bla", nil)

    g := gameService{repo, textService}

    want := errors.New("Invalid ID value") 
    id := 0
    _, err := g.FindGame(uint(id))

    if err.Error() != want.Error() {
        t.Errorf("Expected: %s, got %s", want.Error(), err.Error())
    }
}

func TestFindShouldThrowIfUserNotFound(t *testing.T) {
    repo := new(mockGameRepository)
    textService := new(mockTextService)

    repo.On("Find", mock.Anything).Return(nil, nil)
    textService.On("GetRandomText", mock.Anything).Return(nil, nil)

    g := gameService{repo, textService}
    
    want := errors.New("User not found")
    id := 1
    _, err := g.FindGame(uint(id))

    if err.Error() != want.Error() {
        t.Errorf("Expected: %s, got %s", want.Error(), err.Error())
    }
}

func TestFindShouldReturnGameOnSuccess(t *testing.T) {
    repo := new(mockGameRepository)
    textService := new(mockTextService)

    want := &Game{ID: 1, Text: "bla bla bla bla", PlayerOne: "test", PlayerTwo: "", Status: Waiting, Winner: ""}

    repo.On("Find", mock.Anything).Return(want, nil)
    textService.On("GetRandomText", mock.Anything).Return(nil, nil)

    g := gameService{repo, textService}

    id := 1
    game, _ := g.FindGame(uint(id)) 

    if game.ID != want.ID {
        t.Errorf("Expected: %d, got %d", want.ID, game.ID)
    }
}
