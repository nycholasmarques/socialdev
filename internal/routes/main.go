package routes

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/nycholasmarques/socialdev/internal/db"
	authHttp "github.com/nycholasmarques/socialdev/internal/modules/auth/http"
	comment_handlers "github.com/nycholasmarques/socialdev/internal/modules/comment/http"
	favorite_handlers "github.com/nycholasmarques/socialdev/internal/modules/favorite/http"
	file_handlers "github.com/nycholasmarques/socialdev/internal/modules/file/http"
	liker_handler "github.com/nycholasmarques/socialdev/internal/modules/like/http"
	post_handler "github.com/nycholasmarques/socialdev/internal/modules/posts/http"
	role_routes "github.com/nycholasmarques/socialdev/internal/modules/role/routes"
	theme_routes "github.com/nycholasmarques/socialdev/internal/modules/theme/routes"
	user_routes "github.com/nycholasmarques/socialdev/internal/modules/user/routes"
)

func SetupRoutes(queries *db.Queries) *http.ServeMux {
	r := http.NewServeMux()

	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	uploadsDir := filepath.Join(wd, "uploads")
	r.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir(uploadsDir))))

	apiRouter := http.NewServeMux()
	authHttp.SetupAuthRoutes(apiRouter, queries)
	user_routes.SetupUserRoutes(apiRouter, queries)
	role_routes.SetupUserRoutes(apiRouter, queries)
	post_handler.SetupPostRoutes(apiRouter, queries)
	liker_handler.SetupLikeRoutes(apiRouter, queries)
	comment_handlers.SetupCommentRoutes(apiRouter, queries)
	favorite_handlers.SetupFavoriteRoutes(apiRouter, queries)
	theme_routes.SetupThemeRoutes(apiRouter, queries)
	file_handlers.SetupFileRoutes(apiRouter, queries)

	r.Handle("/api/", http.StripPrefix("/api", apiRouter))

	return r
}
