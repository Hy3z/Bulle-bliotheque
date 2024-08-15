package auth

import (
	"bb/database"
	"bb/logger"
	"bb/util"
	"context"
	"errors"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/Nerzal/gocloak/v13"
	"github.com/coreos/go-oidc"
	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
)

const (
	ENV_PATH                   = ".env"
	authHeaderKey              = "auth"
	ENV_KEYCLOAK_URL           = "KEYCLOAK_URL"
	ENV_KEYCLOAK_CLIENT_ID     = "KEYCLOAK_CLIENT_ID"
	ENV_KEYCLOAK_CLIENT_SECRET = "KEYCLOAK_CLIENT_SECRET"
	ENV_KEYCLOAK_REALM         = "KEYCLOAK_REALM"
	ENV_KEYCLOAK_PUBLIC_KEY    = "KEYCLOACK_PUBLIC_KEY"
	access_token_cookie        = "access-token"
	refresh_token_cookie       = "refresh-token"
	admin_role_name            = "admin"
	refererHeaderKey           = "Referer"
)

var (
	pathEndedError error = errors.New("JWT path ended")
	jwtKeyError    error = errors.New("Wrong key")
)

var (
	client         *gocloak.GoCloak
	clientID       string
	clientSecret   string
	realm          string
	authUrl        string
	ctx            context.Context
	provider       *oidc.Provider
	realmPublicKey string
)

func Setup() {
	var err error
	authUrl = os.Getenv(ENV_KEYCLOAK_URL) + "/auth"
	clientID = os.Getenv(ENV_KEYCLOAK_CLIENT_ID)
	realm = os.Getenv(ENV_KEYCLOAK_REALM)
	clientSecret = os.Getenv(ENV_KEYCLOAK_CLIENT_SECRET)
	realmPublicKey =
		"-----BEGIN PUBLIC KEY-----\n" +
			os.Getenv(ENV_KEYCLOAK_PUBLIC_KEY) +
			"\n-----END PUBLIC KEY-----\n"
	client = gocloak.NewClient(authUrl)
	ctx = context.Background()
	provider, err = oidc.NewProvider(ctx, authUrl+"/realms/"+realm)
	if err != nil {
		logger.ErrorLogger.Panicf("Couldn't create provider: %s\n", err)
	}
	logger.InfoLogger.Println("Sucessfully initialized auth")
}

func getTokens(c echo.Context) (string, string, bool) {
	access_token, err1 := c.Request().Cookie(access_token_cookie)
	refresh_token, err2 := c.Request().Cookie(refresh_token_cookie)
	if err1 != nil || err2 != nil {
		//logger.InfoLogger.Println(access_token)
		//logger.InfoLogger.Println(refresh_token)
		return "", "", false
	}
	return access_token.Value, refresh_token.Value, true
}

func hasToken(c echo.Context) (bool, *gocloak.JWT) {
	access_token, refresh_token, ok := getTokens(c)
	if !ok {
		//logger.InfoLogger.Println("Not ok")
		return false, nil
	}
	if !ok {
		//logger.InfoLogger.Println("Not ok")
		return false, nil
	}

	result, err := client.RetrospectToken(ctx, access_token, clientID, clientSecret, realm)
	if err != nil {
		logger.ErrorLogger.Printf("Error retrospecting token: %s\n", err)
		return false, nil
	}

	if !*result.Active {
		newJWT, err := client.RefreshToken(ctx, refresh_token, clientID, clientSecret, realm)
		if err != nil {
			return false, nil
		}
		return true, newJWT
	}
	return true, nil
}

func hasRoles(c echo.Context, access_token string, req_roles []string) bool {
	ctx := context.Background()
	userInfo, err := client.GetUserInfo(ctx, access_token, realm)
	if err != nil {
		logger.ErrorLogger.Printf("Error getting user info: %s\n", err)
		return false
	}
	query, err := util.ReadCypherScript(util.CypherScriptDirectory + "/auth/getRolesByUUID.cypher")
	if err != nil {
		logger.ErrorLogger.Printf("Error reading script: %s\n", err)
		return false
	}

	res, err := database.Query(ctx, query, map[string]any{
		"uuid": *userInfo.Sub,
	})
	if err != nil {
		logger.ErrorLogger.Printf("Error on query %s: %s\n", query, err)
		return false
	}

	for _, req_role := range req_roles {
		if !util.RecordsContains(res.Records, 0, req_role) {
			return false
		}
		continue
	}

	return true
}

func HasTokenMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenPresent, jwt := hasToken(c)
		if tokenPresent {
			if jwt != nil {
				addCookies(&c, jwt.AccessToken, jwt.RefreshToken)
			}
			return next(c)
		}

		//On devra donc refaire l'action si on est pas encore connecté
		//Dans le futur on passera un paramètre si on veut être quand même redirigé
		//c.Request().Header.Set(refererHeaderKey, c.Path())

		return Login(c)
	}
}

func HasRoleMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenPresent, jwt := hasToken(c)

		if !tokenPresent {
			c.Request().Header.Set(refererHeaderKey, c.Path())
			return Login(c)
		}
		access_token, _, _ := getTokens(c)
		if jwt != nil {
			addCookies(&c, jwt.AccessToken, jwt.RefreshToken)
			access_token = jwt.AccessToken
		}

		if !hasRoles(c, access_token, []string{admin_role_name}) {
			return c.NoContent(http.StatusForbidden)
		}
		return next(c)
	}
}

func addCookies(c *echo.Context, access_token string, refresh_token string) {
	//Set cookies to response
	accessCookie := new(http.Cookie)
	accessCookie.Name = access_token_cookie
	accessCookie.Value = access_token
	accessCookie.Secure = true
	accessCookie.Path = "/"
	accessCookie.SameSite = http.SameSiteNoneMode
	(*c).SetCookie(accessCookie)
	refreshCookie := new(http.Cookie)
	refreshCookie.Name = refresh_token_cookie
	refreshCookie.Value = refresh_token
	refreshCookie.Secure = true
	refreshCookie.Path = "/"
	refreshCookie.SameSite = http.SameSiteNoneMode
	(*c).SetCookie(refreshCookie)

	//Set cookies to request so the new cookies are immediatly valid
	if ac, err := (*c).Request().Cookie(access_token_cookie); err == nil {
		ac.Value = access_token
	} else {
		(*c).Request().AddCookie(accessCookie)
	}
	if rc, err := (*c).Request().Cookie(refresh_token_cookie); err == nil {
		rc.Value = refresh_token
	} else {
		(*c).Request().AddCookie(refreshCookie)
	}
}

func Login(c echo.Context) error {
	origin := c.Request().Header.Get(refererHeaderKey)
	pUrl, _ := url.Parse(origin)
	path := pUrl.Path
	if pUrl.RawQuery != "" {
		path += "?" + pUrl.RawQuery
	}

	//Prevent login-logout loop, just in case
	if path == util.LogoutPath {
		path = ""
	}

	oauth2Config := oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  "https://bulle.rezel.net" + util.CallbackLoginPath,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID},
	}

	if c.Request().Header.Get("HX-Request") == "true" {
		c.Response().Header().Set("HX-Redirect", oauth2Config.AuthCodeURL(url.QueryEscape(path)))
		return c.NoContent(http.StatusOK)
	} else {
		return c.Redirect(http.StatusTemporaryRedirect, oauth2Config.AuthCodeURL(url.QueryEscape(path)))
	}

}

func LoginCallback(c echo.Context) error {
	origin := c.QueryParam("state")
	ctx := context.Background()
	oauth2Config := oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     provider.Endpoint(),
		RedirectURL:  "https://bulle.rezel.net" + util.CallbackLoginPath,
		Scopes:       []string{oidc.ScopeOpenID},
	}
	token, err := oauth2Config.Exchange(ctx, c.QueryParam("code"))
	if err != nil {
		logger.ErrorLogger.Printf("Error exchanging code: %s\n", err)
		return c.NoContent(http.StatusBadRequest)
	}

	uuid, name, err := GetUserInfo(token.AccessToken)
	if err != nil {
		logger.ErrorLogger.Printf("Error getting user infos: %s\n", err)
		return c.NoContent(http.StatusBadRequest)
	}

	//Création & Mise à jour des info utilisateur à la connection
	query, err := util.ReadCypherScript(util.CypherScriptDirectory + "/auth/createUser.cypher")
	if err != nil {
		logger.ErrorLogger.Printf("Error reading script: %s\n", err)
		return c.NoContent(http.StatusBadRequest)
	}
	_, err = database.Query(ctx, query, map[string]any{
		"uuid": uuid,
		"name": name,
	})
	if err != nil {
		logger.ErrorLogger.Printf("Error creating user: %s\n", err)
		return c.NoContent(http.StatusBadRequest)
	}

	addCookies(&c, token.AccessToken, token.RefreshToken)
	path, _ := url.QueryUnescape(origin)
	return c.Redirect(http.StatusPermanentRedirect, "https://bulle.rezel.net"+path)
}

func Logout(c echo.Context) error {
	//Delete previous cookies
	c.SetCookie(&http.Cookie{
		Name:     access_token_cookie,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})
	c.SetCookie(&http.Cookie{
		Name:     refresh_token_cookie,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})

	//Get redirect url
	origin := c.Request().Header.Get(refererHeaderKey)
	pUrl, _ := url.Parse(origin)
	path := pUrl.Path
	if pUrl.RawQuery != "" {
		path += "?" + pUrl.RawQuery
	}

	//Urls containing such characters will result in 'Invalid url' from keycloak logout for some reason. There may be other unallowed characters
	// So we default to root url
	if strings.ContainsAny(path, "+% ") {
		path = ""
	}

	redirectUrl := "https://bulle.rezel.net" + path

	url := authUrl + "/realms/" + realm + "/protocol/openid-connect/logout"
	url += "?post_logout_redirect_uri=" + redirectUrl
	url += "&client_id=" + clientID
	return c.Redirect(http.StatusTemporaryRedirect, url)
}

func IsLogged(c echo.Context) bool {
	ok, jwt := hasToken(c)
	if !ok {
		return false
	}

	if jwt != nil {
		logger.InfoLogger.Println("Changed cookies")
		addCookies(&c, jwt.AccessToken, jwt.RefreshToken)
	}
	return true
}

func IsAdmin(c echo.Context) bool {
	tokenPresent, jwt := hasToken(c)
	if !tokenPresent {
		return false
	}
	access_token, _, _ := getTokens(c)
	if jwt != nil {
		addCookies(&c, jwt.AccessToken, jwt.RefreshToken)
		access_token = jwt.AccessToken
	}
	return hasRoles(c, access_token, []string{admin_role_name})
}

func GetUserInfo(access_token string) (string, string, error) {
	info, err := client.GetUserInfo(context.Background(), access_token, realm)
	if err != nil {
		return "", "", err
	}

	return *info.Sub, *info.Name, nil
}

func GetUserInfoFromContext(c echo.Context) (string, string, bool) {
	access_token, _, ok := getTokens(c)
	if !ok {
		return "", "", false
	}
	uuid, name, err := GetUserInfo(access_token)
	return uuid, name, err == nil
}

// GetUserUUID returns user's UUID, empty if no user
func GetUserUUID(access_token string) string {
	uuid, _, err := GetUserInfo(access_token)
	if err != nil {
		return ""
	}
	return uuid
}

// GetUserUUIDFromContext GetUserUUID returns user's UUID, empty if no user
func GetUserUUIDFromContext(c echo.Context) string {
	access_token, _, ok := getTokens(c)
	if !ok {
		return ""
	}

	return GetUserUUID(access_token)
}

func GetUserName(access_token string) string {
	_, name, err := GetUserInfo(access_token)
	if err != nil {
		return ""
	}
	return name
}

// GetUserNameFromContext GetUserName returns user's UUID, empty if no user
func GetUserNameFromContext(c echo.Context) string {
	access_token, _, ok := getTokens(c)
	if !ok {
		return ""
	}

	info, err := client.GetUserInfo(context.Background(), access_token, realm)
	if err != nil {
		return ""
	}

	return *info.Sub
}
