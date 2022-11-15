package handler

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/teampui/pac"
	_ "go/printer"
	_ "log"
	"login/modle"
	customError "login/pkg/customEorr"
	"login/pkg/customEorr/repository"
	_ "strings"
)

const LoggedUserIdStoredKey = "logged_user_id"
const LoggedUserTcStoredKey = "logged_user_tc"
const LoggedUserTtStoredKey = "logged_user_tt"

type AuthHandler struct {
	userRepo repository.CustomerRepoInterface
}

// type colors struct {
// 	user password
// }

func (h *AuthHandler) Register(app *pac.App) {
	r := app.Router().Group("/api/v1/auth")
	// 帳號登入
	r.Post("/login", h.handleLogin)
	// 帳號登入
	r.Post("/register", h.handleRegister)
	h.userRepo = pac.Must[repository.CustomerRepoInterface](
		pac.Repo[repository.CustomerRepoInterface](app, "customerdata"),
		"service/manager: cannot start due to no valid manager repo found")

}

// handleLogin 處理使用者登入的流程
func (h *AuthHandler) handleRegister(c *fiber.Ctx) error {

	// 取得目前登入的資訊, 如果已經登入則回傳 403
	if h.isLogged(c) {
		return customError.New(c, 403, "already logged")
	}

	// 開始收集登入資訊
	form := new(model.AuthUserLoginForm)
	if err := c.BodyParser(form); err != nil {
		return customError.New(c, 400, "cannot parse request")
	}
	user := form.Username
	password := form.Password
	check, err := h.userRepo.Createdaccountdate(user, password)
	if err != nil {
		return c.SendString(err.Error())
	}
	fmt.Println("註冊成功", check)

	return c.SendString("成功註冊")

}
func (h *AuthHandler) handleLogin(c *fiber.Ctx) error {

	// 取得目前登入的資訊, 如果已經登入則回傳 403
	if h.isLogged(c) {
		return customError.New(c, 403, "already logged")
	}

	// 開始收集登入資訊
	form := new(model.AuthUserLoginForm)
	if err := c.BodyParser(form); err != nil {
		return customError.New(c, 400, "cannot parse request")
	}
	user := form.Username
	password := form.Password
	password, err := h.userRepo.GetProfileByUsername(user)
	if err != nil {
		return c.SendString(err.Error())
	}
	if password != form.Password {
		return c.SendString("密碼錯誤")
	}
	return c.SendString("成功登入")
}

func (h *AuthHandler) handleCheck(c *fiber.Ctx) error {
	// 取得目前登入的資訊, 如果已經登入則回傳 403
	if h.isLogged(c) {
		return customError.New(c, 403, "already logged")
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  200,
		"message": "not logged",
		"error":   false,
	})
}

// checkLogged 會快速進入 session 檢查是否有指定的值
func (h *AuthHandler) isLogged(c *fiber.Ctx) bool {
	session := pac.NewSessionOperator(c)
	return session.Get(LoggedUserIdStoredKey, int64(0)).(int64) > 0
}

// getLogged 取得目前 Session 中登入的使用者 Id
func (h *AuthHandler) getLogged(c *fiber.Ctx) (int64, error) {
	session := pac.NewSessionOperator(c)
	userID, ok := session.Get(LoggedUserIdStoredKey, 0).(int64)

	if !ok {
		return 0, fmt.Errorf("cannot parse logged user Id")
	}

	return userID, nil
}

// getSessionID 取得目前 Session 的 Id (唯一識別符)
func (h *AuthHandler) getSessionID(c *fiber.Ctx) string {
	session, ok := c.Locals("session").(*session.Session)

	if !ok {
		return ""
	}

	return session.ID()
}

// destroyConnectionSession 會把目前登入中的 session 給整個刪除，若使用者需要請求則需要新給一個
func (h *AuthHandler) destroyConnectionSession(c *fiber.Ctx) error {
	session := pac.NewSessionOperator(c)
	return session.Destroy()
}

// 取出tc，如果有沒有則預設是0
func getTcBySession(c *fiber.Ctx) int64 {
	session := pac.NewSessionOperator(c)
	tc := session.Get(LoggedUserTcStoredKey, 0)

	referrer, ok := tc.(int64)
	if !ok {
		return 0
	}

	return referrer
}
