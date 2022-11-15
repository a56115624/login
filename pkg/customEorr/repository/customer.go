package repository

import (
	_ "context"
	"fmt"

	"github.com/teampui/pac"
)

type CustomerRepoInterface interface {
	pac.Service
	GetProfileByUsername(username string) (string, error)
}

type CustomerInMemoryRepo struct {
	data map[string]string
}

func (repo *CustomerInMemoryRepo) Register(app *pac.App) {
	app.Repositories.Add("customer", repo)

	data := make(map[string]string)
	data["admin"] = "password"

	repo.data = data
	//repo.db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
}
func (repo *CustomerInMemoryRepo) GetProfileByUsername(username string) (string, error) {
	password, ok := repo.data[username]
	if !ok {
		return "", fmt.Errorf("沒有該帳號")
	}
	return password, nil
}
