package routes

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/nycholasmarques/socialdev/internal/db"
	authHttp "github.com/nycholasmarques/socialdev/internal/modules/auth/http"
	comment_handlers "github.com/nycholasmarques/socialdev/internal/modules/comment/http"
	favorite_handlers "github.com/nycholasmarques/socialdev/internal/modules/favorite/http"
	liker_handler "github.com/nycholasmarques/socialdev/internal/modules/like/http"
	post_handler "github.com/nycholasmarques/socialdev/internal/modules/posts/http"
	role_handler "github.com/nycholasmarques/socialdev/internal/modules/role/http"
	userHttp "github.com/nycholasmarques/socialdev/internal/modules/user/http"
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
	userHttp.SetupUserRoutes(apiRouter, queries)
	role_handler.SetupUserRoutes(apiRouter, queries)
	post_handler.SetupPostRoutes(apiRouter, queries)
	liker_handler.SetupLikeRoutes(apiRouter, queries)
	comment_handlers.SetupCommentRoutes(apiRouter, queries)
	favorite_handlers.SetupFavoriteRoutes(apiRouter, queries)
	SetupThemeRoutes(apiRouter, queries)
	SetupFileRoutes(apiRouter, queries)

	r.Handle("/api/", http.StripPrefix("/api", apiRouter))

	return r
}
