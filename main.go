package main

import (
	"github.com/gofiber/fiber/v2/middleware/cors"

	"log"

	"github.com/teampui/pac/bundb"

	// 放在所有人前面，這樣就可以在一開始就讀取環境變數
	"github.com/joho/godotenv"
	// _ "github.com/joho/godotenv/autoload"

	// 載入 Pac
	"github.com/teampui/pac"
	"github.com/teampui/pac/redis"

	// 載入服務

	"login/handler"
	"login/pkg/customEorr/repository"
	// "github.com/teampui/comico-api/cmd/api/session-handler"
	// pkgHandler "github.com/teampui/comico-api/pkg/handler"
	// "github.com/teampui/comico-api/pkg/repository"
	// "github.com/teampui/comico-api/pkg/service"

	"os"
	"time"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	log.Println("print env: " + os.Getenv("DSN"))
	log.Println("Starting comico-api...")

	// 開始一個新的 App
	app := pac.NewApp(
		pac.ListenPortFromEnv(":3245"),    // 如果環境變數裡沒設定的話，預設 :7777
		pac.UseLogger(),                   // 使用請求記錄器
		bundb.ProvideDB(os.Getenv("DSN")), // 使用 BunDB 作為資料庫層
		redis.ProvideSession(redis.SessionConfig{
			ClientKeystore: "cookie:942",
			RedisURL:       os.Getenv("REDIS_DSN"),
			Expiration:     24 * time.Hour,
		}),
	)

	app.Router().Use(cors.New(cors.Config{
		AllowHeaders:     "Origin,Content-Type,Accept,Content-Length,Accept-Language,Accept-Encoding,Connection,Access-Control-Allow-Origin",
		AllowOrigins:     "*",
		AllowCredentials: true,
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
	}))

	app.Add(&repository.CustomerInMemoryRepo{})

	// API 路由
	// session handler
	app.Add(&handler.AuthHandler{})

	// 開始工作
	app.Start()
}
