package web

import (
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/alexedwards/scs/v2"
	"github.com/gorilla/mux"
	svc "{{gomod.name}}/services"
	"{{gomod.name}}/web/auth"
	"golang.org/x/oauth2"
)

func DefaultGatewayAddress() string {
	gateway_addr := os.Getenv("{{app.classname | toupper}}_WEB_PORT")
	if gateway_addr != "" {
		return gateway_addr
	}
	return ":10080"
}

// We are just storing the tokens in mem - just for testing
// Move to a datastore later
type TokenInfo struct {
	Provider  string
	AuthType  string
	AuthToken *oauth2.Token
	UserInfo  map[string]any
}

type {{app.classname}}App struct {
	AllowLocalDev  bool
	Api            *{{app.classname}}Api
	Auth           *auth.{{app.classname}}Auth
	AuthMiddleware *auth.AuthMiddleware
	ClientMgr      *svc.ClientMgr

	router  *mux.Router
	Session *scs.SessionManager

	// Storing tokens in-mem for now to test oauth
	tokenLock sync.RWMutex
	allTokens []*TokenInfo
}

func New{{app.classname}}App(grpcAddr string) (app *{{app.classname}}App, err error) {
	clientMgr := svc.NewClientMgr(grpcAddr)
	session := scs.New() //scs.NewCookieManager("u46IpCV9y5Vlur8YvODJEhgOY8m9JVE4"),

	am := auth.AuthMiddleware{
		SessionGetter: func(r *http.Request, key string) any {
			return session.GetString(r.Context(), key)
		},
	}

	// Uncomment if using some kind of database session
	// Session.Store = NewGCDSessionStore(0)
	omauth := &auth.{{app.classname}}Auth{
		ClientMgr:      clientMgr,
		Session:        session,
		AuthMiddleware: &am,
	}

	omapi := &{{app.classname}}Api{
		GrpcAddr:       grpcAddr,
		AuthMiddleware: &am,
	}

	app = &{{app.classname}}App{
		AllowLocalDev:  true,
		ClientMgr:      clientMgr,
		Api:            omapi,
		Auth:           omauth,
		Session:        session,
		AuthMiddleware: &am,
	}
	omauth.UserStore = app
	return app, err
}

func (web *{{app.classname}}App) GetRouter() http.Handler {
	if web.router == nil {
		log.Println("Creating App router...")
		web.router = mux.NewRouter().StrictSlash(true)
		web.router.Use(web.AuthMiddleware.ExtractUser)

		// TODO - Also see how we can allow cross site auth cookies
		if web.AllowLocalDev {
			web.router.Use(cors)
		}

		authRouter := web.router.PathPrefix("/auth").Subrouter()
		web.Auth.SetupRoutes(authRouter)

		apiRouter := web.router.PathPrefix("/api").Subrouter()
		web.Api.SetupRoutes(apiRouter)
	}
	return web.Session.LoadAndSave(web.router)
}

func (web *{{app.classname}}App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	web.GetRouter().ServeHTTP(w, r)
}

type User struct {
	id string
}

func (u *User) Id() string { return u.id }

// Just to manage tokens in mem
func (o *{{app.classname}}App) GetUserByID(userId string) (auth.User, error) {
	// very hacky check here
	// TODO - Fix this with a) use a real storage and b) with some kind of proper ID mapping
	o.tokenLock.RLock()
	defer o.tokenLock.RUnlock()
	for _, ti := range o.allTokens {
		if ti.Provider == "google" {
			return &User{id: ti.UserInfo["email"].(string)}, nil
		}
		if ti.AuthType == "saml" {
			return &User{id: ti.UserInfo["email"].(string)}, nil
		}
	}
	return nil, nil
}

func (o *{{app.classname}}App) EnsureAuthUser(authtype string, provider string, token *oauth2.Token, userInfo map[string]any) (auth.User, error) {
	o.tokenLock.Lock()
	defer o.tokenLock.Unlock()
	userEmail := userInfo["email"].(string)
	for _, ti := range o.allTokens {
		if ti.AuthType == authtype && ti.Provider == provider && ti.UserInfo["email"].(string) == userEmail {
			// Tokenfor this user already found - so just update the user info and token
			ti.UserInfo = userInfo
			ti.AuthToken = token
			return &User{id: userEmail}, nil
		}
	}
	// none found so add it
	o.allTokens = append(o.allTokens, &TokenInfo{AuthType: authtype, Provider: provider, AuthToken: token, UserInfo: userInfo})
	return &User{id: userEmail}, nil
}
