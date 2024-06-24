package auth

import (
	"bb/logger"
	"bb/util"
	"context"
	"errors"
	"fmt"
	"github.com/Nerzal/gocloak/v13"
	"github.com/coreos/go-oidc"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
	"net/http"
	"net/url"
	"os"
)

const (
	ENV_PATH                   = ".env"
	authHeaderKey              = "auth"
	ENV_KEYCLOAK_URL           = "KEYCLOAK_URL"
	ENV_KEYCLOAK_CLIENT_ID     = "KEYCLOAK_CLIENT_ID"
	ENV_KEYCLOAK_CLIENT_SECRET = "KEYCLOAK_CLIENT_SECRET"
	ENV_KEYCLOAK_REALM         = "KEYCLOAK_REALM"
	access_token_cookie        = "access-token"
	refresh_token_cookie       = "refresh-token"
	admin_role_name            = "my.role.dev"
	refererHeaderKey           = "Referer"
)

var (
	pathEndedError error = errors.New("JWT path ended")
	jwtKeyError    error = errors.New("Wrong key")
)

var (
	client       *gocloak.GoCloak
	clientID     string
	clientSecret string
	realm        string
	kcUrl        string
	ctx          context.Context
	provider     *oidc.Provider
)

func Setup() {
	var err error
	kcUrl = os.Getenv(ENV_KEYCLOAK_URL)
	clientID = os.Getenv(ENV_KEYCLOAK_CLIENT_ID)
	realm = os.Getenv(ENV_KEYCLOAK_REALM)
	clientSecret = os.Getenv(ENV_KEYCLOAK_CLIENT_SECRET)
	client = gocloak.NewClient(kcUrl)
	ctx = context.Background()
	provider, err = oidc.NewProvider(ctx, kcUrl+"/realms/"+realm)
	if err != nil {
		logger.ErrorLogger.Panicf("Couldn't create provider: %s\n", err)
	}
	logger.InfoLogger.Println("Sucessfully initialized auth")
}

func getTokens(c echo.Context) (string, string, bool) {
	access_token, err1 := c.Request().Cookie(access_token_cookie)
	refresh_token, err2 := c.Request().Cookie(refresh_token_cookie)
	if err1 != nil || err2 != nil {
		return "", "", false
	}
	return access_token.Value, refresh_token.Value, true
}

func hasToken(c echo.Context) (bool, *gocloak.JWT) {
	access_token, refresh_token, ok := getTokens(c)
	if !ok {
		return false, nil
	}

	result, err := client.RetrospectToken(ctx, access_token, clientID, clientSecret, realm)
	if err != nil {
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

func sliceContains(slice []interface{}, tofind string) bool {
	for _, inter := range slice {
		sinter, ok := inter.(string)
		if !ok {
			continue
		}
		if sinter == tofind {
			return true
		}
	}
	return false
}

func jwtWalk(jwt jwt.MapClaims, keys ...string) (interface{}, error) {
	if len(keys) == 0 {
		return jwt, nil
	}
	next, ok := jwt[keys[0]]
	if !ok {
		return nil, pathEndedError
	}
	for _, key := range keys[1:] {
		temp, ok := next.(map[string]interface{})
		if !ok {
			logger.InfoLogger.Printf("%s\n", key)
			return nil, pathEndedError
		}
		next, ok = temp[key]
		if !ok {
			return nil, jwtKeyError
		}
	}
	return next, nil
}

func hasRoles(c echo.Context, req_roles []string) bool {
	access_token, _, ok := getTokens(c)
	if !ok {
		return false
	}

	key :=
		"-----BEGIN PUBLIC KEY-----\n" +
			"MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA0nJP/hPocUHoginR37MwpF35vhl+icQjMvdLeQ78RYG5tV9TJb2pxowPe6CIxty3pvfMjrCfUJTY8gx/9GLm8qiqdVRmw0FEUVho58gg/+3wVjcJf3As08Y4wbtIa4rHoQH3VDC5O42rq8q/elH457vY8E0BeNpKy8W5VGrwELYP6FPpO8eTx6EFwsg6E6R4R5z8bouwQg4sudS8MN5bqt9FgYO/QahNQStMNoB8NVQ/EccCZWbadZApsTXz/0Sxe7LGnaOPc50CTMv4HvXaQdv4d5U+qug+6t7oqj0VqL0ZX9SE1TmoAIlCsclQui5RYx+r5MAGoXchnZ/TwgoKHwIDAQAB" +
			"\n-----END PUBLIC KEY-----\n"

	pk, err := jwt.ParseRSAPublicKeyFromPEM([]byte(key))
	if err != nil {
		logger.InfoLogger.Printf("Error parsing public key: %s\n", err)
		return false
	}

	token, err := jwt.ParseWithClaims(access_token, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return pk, nil
	})

	if err != nil {
		logger.InfoLogger.Printf("Error parsing token: %s\n", err)
		return false
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return false
	}

	iroles, err := jwtWalk(claims, "resource_access", clientID, "roles")
	if err != nil {
		logger.InfoLogger.Printf("Error walking jwt: %s\n", err)
	}

	roles, ok := iroles.([]interface{})
	if !ok {
		return false
	}

	for _, req_role := range req_roles {
		if !sliceContains(roles, req_role) {
			logger.InfoLogger.Printf("%s not in\n", req_role)
			return false
		}
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
		c.Request().Header.Set(refererHeaderKey, c.Path())
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

		if jwt != nil {
			addCookies(&c, jwt.AccessToken, jwt.RefreshToken)
		}

		if !hasRoles(c, []string{admin_role_name}) {
			return c.NoContent(http.StatusForbidden)
		}
		return next(c)
	}
}

func addCookies(c *echo.Context, access_token string, refresh_token string) {
	accessCookie := new(http.Cookie)
	accessCookie.Name = access_token_cookie
	accessCookie.Value = access_token
	accessCookie.Secure = true
	accessCookie.Path = "/"
	(*c).SetCookie(accessCookie)
	refreshCookie := new(http.Cookie)
	refreshCookie.Name = refresh_token_cookie
	refreshCookie.Value = refresh_token
	refreshCookie.Secure = true
	refreshCookie.Path = "/"
	(*c).SetCookie(refreshCookie)
}

func Login(c echo.Context) error {
	origin := c.Request().Header.Get(refererHeaderKey)
	pUrl, _ := url.Parse(origin)
	path := pUrl.Path
	if pUrl.RawQuery != "" {
		path += "?" + pUrl.RawQuery
	}
	oauth2Config := oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  "http://localhost:42069" + util.CallbackLoginPath,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID},
	}
	return c.Redirect(http.StatusTemporaryRedirect, oauth2Config.AuthCodeURL(url.QueryEscape(path)))
}

func LoginCallback(c echo.Context) error {
	origin := c.QueryParam("state")
	oauth2Config := oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     provider.Endpoint(),
		RedirectURL:  "http://localhost:42069" + util.CallbackLoginPath,
		Scopes:       []string{oidc.ScopeOpenID},
	}
	token, err := oauth2Config.Exchange(context.Background(), c.QueryParam("code"))
	if err != nil {
		logger.ErrorLogger.Printf("Error exchanging code: %s\n", err)
		return c.NoContent(http.StatusBadRequest)
	}

	addCookies(&c, token.AccessToken, token.RefreshToken)
	path, _ := url.QueryUnescape(origin)
	return c.Redirect(http.StatusPermanentRedirect, "http://localhost:42069"+path)
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
	redirectUrl := "http://localhost:42069" + path

	url := kcUrl + "/realms/" + realm + "/protocol/openid-connect/logout"
	url += "?post_logout_redirect_uri=" + redirectUrl
	url += "&client_id=" + clientID
	return c.Redirect(http.StatusTemporaryRedirect, url)
}
