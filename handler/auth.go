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
	"strconv"
	_ "strings"
)

const LoggedUserIdStoredKey = "logged_user_id"
const LoggedUserTcStoredKey = "logged_user_tc"
const LoggedUserTtStoredKey = "logged_user_tt"

type AuthHandler struct {
	userRepo  repository.CustomerRepoInterface
	redisRepo repository.RedisRepoInterface
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
	// 登出
	r.Post("/logout", h.checkAuthn(), h.handleLogout)
	r.Get("/check", h.checkAuthn(), h.handleCheck)
	h.userRepo = pac.Must[repository.CustomerRepoInterface](
		pac.Repo[repository.CustomerRepoInterface](app, "customer"),
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

// handleLogout 處理使用者登出的情形
func (h *AuthHandler) handleLogout(c *fiber.Ctx) error {
	// tt := c.Query("tt", "")
	// tc := getTcBySession(c)

	// 先取得目前的登入使用者
	loggedUser, err := h.getLogged(c)
	if err != nil {
		return customError.New(c, 500, "cannot get current logged user")
	}
	// 從 Redis 裡移除指定的登入記錄
	err = h.removeSpecificConcurrent(loggedUser, h.getSessionID(c))
	if err != nil {
		return customError.New(c, 500, "cannot get current logged user")
	}

	// 把目前的 session 給取消掉
	err = h.destroyConnectionSession(c)
	if err != nil {
		return customError.New(c, 500, "cannot update session for logout")
	}

	// 回傳給使用者
	return c.SendString("我在這!")
	// return c.Status(200).JSON(fiber.Map{
	// 	"status":  200,
	// 	"message": "successfully logout",
	// 	"error":   false,
	// })
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

// checkAuthn 會檢查使用者是否有登入，如果沒登入的話回傳 401
func (h *AuthHandler) checkAuthn(opts ...any) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		if !h.isLogged(c) {
			return customError.New(c, 401, "not logged")
		}

		user, err := h.getLogged(c)

		if err != nil {
			return customError.New(c, 500, "cannot get logged user info")
		}
		// 如果目前這個使用者名下的線上可登入使用者沒有這個 session
		if !h.checkSpecificConcurrent(user, h.getSessionID(c)) {
			//
			err = h.destroyConnectionSession(c)
			if err != nil {
				return customError.New(c, 500, "cannot loggout")
			}

			return customError.New(c, 401, "not logged")
		}
		return c.Next()
	}
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
func (h *AuthHandler) checkSpecificConcurrent(loggedUserId int64, sessionId string) bool {
	members := h.redisRepo.GetSortedSetAllMembers(strconv.FormatInt(loggedUserId, 10))
	for _, member := range members {
		if member == sessionId {
			return true
		}
	}
	return false
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

func (h *AuthHandler) removeSpecificConcurrent(loggedUserId int64, sessionId string) error {
	return h.redisRepo.RemoveSortedSetMember(strconv.FormatInt(loggedUserId, 10), sessionId)
}
