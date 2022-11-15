package repository

import (
	"context"
	_ "context"
	"fmt"
	"github.com/teampui/pac"
	"github.com/uptrace/bun"
	"login/modle"
)

type CustomerRepoInterface interface {
	pac.Service
	GetProfileByUsername(username string) (string, error)
	Createdaccountdate(username string, password string) (string, error)
}
type CustomerDataRepo struct {
	db *bun.DB
}

func (repo *CustomerDataRepo) Register(app *pac.App) {
	app.Repositories.Add("customerdata", repo)

	repo.db = pac.Must[*bun.DB](
		pac.Svc[*bun.DB](app, "db"),
		"repository/product: cannot start due to no valid nsmh db found",
	)
}
func (cu *CustomerDataRepo) GetProfileByUsername(username string) (string, error) {
	ctx := context.Background()
	Customerdata := new(model.Customerdata)
	err := cu.db.
		NewSelect().
		Model(Customerdata).
		Where("username=?", username).
		Scan(ctx)
	if err != nil {
		return "", err
	}
	return Customerdata.Password, nil
}
func (cu *CustomerDataRepo) Createdaccountdate(username string, password string) (string, error) {

	return "", nil
}

// 以下為舊版功能
type CustomerInMemoryRepo struct {
	data map[string]string
}

func (repo *CustomerInMemoryRepo) Register(app *pac.App) {
	app.Repositories.Add("customer", repo)

	data := make(map[string]string)
	data["admin"] = "password"
	data["shane"] = "123"

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
func (repo *CustomerInMemoryRepo) Createdaccountdate(username string, password string) (string, error) {
	_, ok := repo.data[username]
	if ok {
		return "", fmt.Errorf("該帳號已存在")
	}
	data := make(map[string]string)
	data[username] = password
	repo.data = data
	return "OK", nil

}
