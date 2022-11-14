package custom_error

import "github.com/gofiber/fiber/v2"

type ComicoError struct {
	Status    int    `json:"status"`
	Message   string `json:"message"`
	ErrorBool bool   `json:"error"`
	ErrorMsg  string `json:"-"`
}

func (e ComicoError) Error() string {
	return e.ErrorMsg
}

func (e ComicoError) ErrorCode() int {
	return e.Status
}

func (e ComicoError) ErrorMessage() string {
	return e.Message
}

func (e ComicoError) Return(c *fiber.Ctx) error {
	return c.Status(e.Status).JSON(e)
}

func NewAPIError(code int, responseMsg string, errorMsg string, errBool bool) *ComicoError {
	return &ComicoError{
		Status:    code,
		Message:   responseMsg,
		ErrorBool: errBool,
		ErrorMsg:  errorMsg,
	}
}

func New(c *fiber.Ctx, statusCode int, message string) error {
	return c.Status(statusCode).JSON(NewAPIError(statusCode, message, message, true))
}
