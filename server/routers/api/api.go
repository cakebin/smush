package api

import (
	"net/http"

	"github.com/cakebin/smush/server/env"
	"github.com/cakebin/smush/server/routers/api/auth"
	"github.com/cakebin/smush/server/routers/api/character"
	"github.com/cakebin/smush/server/routers/api/match"
	"github.com/cakebin/smush/server/routers/api/user"
	"github.com/cakebin/smush/server/util/routing"
)

// Router is responsible for serving "/api"
// or delegating to the appropriate sub api-router
type Router struct {
	SysUtils        *env.SysUtils
	AuthRouter      *auth.Router
	MatchRouter     *match.Router
	UserRouter      *user.Router
	CharacterRouter *character.Router
}

func (r *Router) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	var head string
	head, req.URL.Path = routing.ShiftPath(req.URL.Path)

	if head == "auth" {
		r.AuthRouter.ServeHTTP(res, req)
		return
	}

	/*-----------------------------------------
	       API Route Authentication
	------------------------------------------*/
	// Check the access token
	accessCookie, err := req.Cookie("smush-access-token")

	// If the access token exists, check the token value
	if err == nil {
		// If this is run, and there is an error (token is invalid),
		// err will not be nil in the next if block
		_, err = r.SysUtils.Authenticator.CheckJWTToken(accessCookie.Value)
	}

	// Access token is invalid, or there wasn't an access cookie in the first place
	// I do not like reusing the same err here... the cascading logic is confusing. But it'll do for now.
	if err != nil {
		http.Error(res, "Session expired. Please log in again", http.StatusUnauthorized)
		return
	}

	// Once we're authorized, allow api requests
	switch head {
	case "match":
		r.MatchRouter.ServeHTTP(res, req)
	case "user":
		r.UserRouter.ServeHTTP(res, req)
	case "character":
		r.CharacterRouter.ServeHTTP(res, req)
	default:
		http.Error(res, "404 Not Found", http.StatusNotFound)
	}
}

// NewRouter makes a new api router and sets up its children
// routers with access to the "SysUtils" environment object
func NewRouter(sysUtils *env.SysUtils) *Router {
	router := new(Router)
	router.SysUtils = sysUtils
	router.AuthRouter = auth.NewRouter(sysUtils)
	router.MatchRouter = match.NewRouter(sysUtils)
	router.UserRouter = user.NewRouter(sysUtils)
	router.CharacterRouter = character.NewRouter(sysUtils)
	return router
}
