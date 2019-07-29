package api

import (
	"net/http"
	"time"

	"github.com/cakebin/smush/server/env"
	"github.com/cakebin/smush/server/routers/api/auth"
	"github.com/cakebin/smush/server/routers/api/match"
	"github.com/cakebin/smush/server/routers/api/user"
	"github.com/cakebin/smush/server/util/routing"
)

// Router is responsible for serving "/api"
// or delegating to the appropriate sub api-router
type Router struct {
	SysUtils    *env.SysUtils
	AuthRouter  *auth.Router
	MatchRouter *match.Router
	UserRouter  *user.Router
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

	if err != nil {
		http.Error(res, "Not logged in", http.StatusUnauthorized)
		return
	}

	_, err = r.SysUtils.Authenticator.CheckJWTToken(accessCookie.Value)
	if err != nil {
		// Check the refresh token
		refreshCookie, err := req.Cookie("smush-refresh-token")
		_, err = r.SysUtils.Authenticator.CheckJWTToken(refreshCookie.Value)
		if err != nil {
			// Refresh token is also invalid
			http.Error(res, "Session expired. Please log in again", http.StatusUnauthorized)
			return
		}
		// If refresh token is still valid then refresh the access token
		newExpirationTime := time.Now().Add(5 * time.Minute)
		newAccessToken, err := r.SysUtils.Authenticator.RefreshJWTAccessToken(accessCookie.Value, newExpirationTime)
		http.SetCookie(
			res,
			&http.Cookie{
				Name:    "smush-access-token",
				Value:   newAccessToken,
				Expires: newExpirationTime,
			},
		)
	}

	// Once we're authorized, allow api requests
	switch head {
	case "match":
		r.MatchRouter.ServeHTTP(res, req)
	case "user":
		r.UserRouter.ServeHTTP(res, req)
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
	return router
}
