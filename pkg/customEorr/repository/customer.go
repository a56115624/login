package repository

import (
	"context"
	"github.com/teampui/pac"
	"github.com/uptrace/bun"
	"login/modle"
)

type CustomerRepoInterface interface {
	pac.Service
	GetProfileById(id int64) (*model.Customer, error)
	GetCustomerRecordById(id int64) (*model.Customer, error)
	GetProfileByUsername(username string) (*model.Customer, error)

	CreateNewUser(user *model.Customer) (int64, error)
	ChangePassword(user *model.Customer) error

	UpdateUser(user *model.Customer) error

	GetPurchasedByUserId(userId int64) ([]*model.Purchased, error)
	CheckPurchasedExistByUserIdAndChapterId(userId int64, chapterId int64) bool
}

type CustomerPostgresRepo struct {
	db *bun.DB
}

func (repo *CustomerPostgresRepo) Register(app *pac.App) {
	app.Repositories.Add("customer", repo)

	repo.db = pac.Must[*bun.DB](
		pac.Svc[*bun.DB](app, "db"),
		"repository/customer: cannot start due to no valid customer db found",
	)

	//repo.db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
}

func (repo *CustomerPostgresRepo) GetPurchasedByUserId(userId int64) ([]*model.Purchased, error) {
	purchaseds := make([]*model.Purchased, 0)
	ctx := context.Background()

	if err := repo.db.NewSelect().
		Model(&purchaseds).
		Relation("Chapter").
		Relation("Chapter.Item").
		Where("customer_id = ?", userId).
		Order("created_at desc").Scan(ctx); err != nil {
		return nil, err
	}

	return purchaseds, nil
}

func (repo *CustomerPostgresRepo) CheckPurchasedExistByUserIdAndChapterId(userId int64, chapterId int64) bool {
	ctx := context.Background()
	txn := new(model.Purchased)
	exist, _ := repo.db.NewSelect().Model(txn).Where("customer_id = ? and chapter_id = ?", userId, chapterId).Exists(ctx)

	return exist
}

func (repo *CustomerPostgresRepo) GetTxnsByUserId(userId int64) ([]*model.Txn, error) {
	txns := make([]*model.Txn, 0)
	ctx := context.Background()

	if err := repo.db.NewSelect().
		Model(&txns).
		Relation("Purchased").
		Relation("Purchased.Chapter").
		Relation("Purchased.Chapter.Item").
		Where("txn.customer_id = ?", userId).
		Order("created_at desc").Scan(ctx); err != nil {
		return nil, err
	}

	return txns, nil
}

func (repo *CustomerPostgresRepo) UpdateUser(user *model.Customer) error {
	ctx := context.Background()
	_, err := repo.db.NewUpdate().Model(user).OmitZero().WherePK().Exec(ctx)
	return err
}

func (repo *CustomerPostgresRepo) GetProfileById(id int64) (*model.Customer, error) {
	user := new(model.Customer)
	ctx := context.Background()

	err := repo.db.NewSelect().Model(user).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return user, nil
}
func (repo *CustomerPostgresRepo) GetCustomerRecordById(id int64) (*model.Customer, error) {
	user := new(model.Customer)
	ctx := context.Background()

	err := repo.db.NewSelect().
		Model(user).
		Relation("Txn").
		Relation("Order").
		Where("customer.id = ?", id).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (repo *CustomerPostgresRepo) GetProfileByUsername(username string) (*model.Customer, error) {
	user := new(model.Customer)
	ctx := context.Background()

	err := repo.db.NewSelect().Model(user).Where("username = ?", username).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (repo *CustomerPostgresRepo) CreateNewUser(user *model.Customer) (int64, error) {
	ctx := context.Background()
	_, err := repo.db.NewInsert().Model(user).Exec(ctx)
	if err != nil {
		return 0, err
	}
	return user.Id, err
}

func (repo *CustomerPostgresRepo) ChangePassword(user *model.Customer) error {
	ctx := context.Background()
	_, err := repo.db.NewUpdate().Model(user).Exec(ctx)
	return err
}
