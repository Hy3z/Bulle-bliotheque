package auth

import (
	"bb/database"
	"bb/logger"
	"bb/util"
	"context"
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
	EnvPath = ".env" //Nom du fichier contenant les variables d'environnements

	//Clés des variables d'environnement dans le fichier
	EnvKeycloakUrl          = "KEYCLOAK_URL"
	EnvKeycloakClientId     = "KEYCLOAK_CLIENT_ID"
	EnvKeycloakClientSecret = "KEYCLOAK_CLIENT_SECRET"
	EnvKeycloakRealm        = "KEYCLOAK_REALM"

	//Clés des cookies contenant les tokens Keycloak
	accessTokenCookie  = "access-token"
	refreshTokenCookie = "refresh-token"

	adminRoleName    = "admin"   //Nom du rôle qui donne les pouvoirs d'admin
	refererHeaderKey = "Referer" //Clés du paramètre de l'entête HTML utilisé pour revenir à la page du site après une connection (callback)
)

var (
	client       *gocloak.GoCloak
	clientID     string
	clientSecret string
	realm        string
	authUrl      string
	ctx          context.Context
	provider     *oidc.Provider
)

type CustomContext struct {
	echo.Context
}

func (c *CustomContext) Foo() {
	println("foo")
}

func (c *CustomContext) Bar() {
	println("bar")
}

// Initialisation des variables qui permettent de communiquer avec le serveur Keycloak
func Setup() {
	var err error
	authUrl = os.Getenv(EnvKeycloakUrl) + "/auth"
	clientID = os.Getenv(EnvKeycloakClientId)
	realm = os.Getenv(EnvKeycloakRealm)
	clientSecret = os.Getenv(EnvKeycloakClientSecret)
	/*realmPublicKey =
	"-----BEGIN PUBLIC KEY-----\n" +
		os.Getenv(ENV_KEYCLOAK_PUBLIC_KEY) +
		"\n-----END PUBLIC KEY-----\n"*/
	client = gocloak.NewClient(authUrl)
	ctx = context.Background()
	provider, err = oidc.NewProvider(ctx, authUrl+"/realms/"+realm)
	if err != nil {
		logger.ErrorLogger.Panicf("Couldn't create provider: %s\n", err)
	}
	logger.InfoLogger.Println("Sucessfully initialized auth")
}

// addCookies Ajoute les cookies dans la réponse et dans la requête HTML
func addCookies(c echo.Context, accessToken string, refreshToken string) {
	//Cookies dans la réponse
	accessCookie := new(http.Cookie)
	accessCookie.Name = accessTokenCookie
	accessCookie.Value = accessToken
	accessCookie.Secure = true
	accessCookie.Path = "/"
	accessCookie.SameSite = http.SameSiteNoneMode
	c.SetCookie(accessCookie)
	refreshCookie := new(http.Cookie)
	refreshCookie.Name = refreshTokenCookie
	refreshCookie.Value = refreshToken
	refreshCookie.Secure = true
	refreshCookie.Path = "/"
	refreshCookie.SameSite = http.SameSiteNoneMode
	c.SetCookie(refreshCookie)
}

// getTokens Renvoit les tokens d'accès et de rafraichissement contenu dans un contexte, le boolean vaut true si les deux tokens ont été trouvés
func getTokens(c echo.Context) (string, string, bool) {
	accessToken, err1 := c.Request().Cookie(accessTokenCookie)
	refreshToken, err2 := c.Request().Cookie(refreshTokenCookie)
	logger.InfoLogger.Printf("Access token: %s\n", accessToken)
	if err1 != nil || err2 != nil {
		return "", "", false
	}
	return accessToken.Value, refreshToken.Value, true
}

// IsLogged renvoit true si les tokens contenus dans le contexte sont valides, en rafraichissant les tokens si nécessaires
func IsLogged(c echo.Context) bool {
	accessToken, _, ok := getTokens(c)
	if !ok {
		return false
	}

	result, err := client.RetrospectToken(ctx, accessToken, clientID, clientSecret, realm)
	return err == nil && *result.Active
	/*if err != nil {
		logger.ErrorLogger.Printf("Error retrospecting token: %s\n", err)
		return false
	}

	if !*result.Active {
		jwt, err := client.RefreshToken(ctx, refreshToken, clientID, clientSecret, realm)
		if err != nil {
			logger.ErrorLogger.Printf("Error refreshing token: %s\n", err)
			return false
		}
		addCookies(c, jwt.AccessToken, jwt.RefreshToken)
		logger.InfoLogger.Println("Refreshed cookies")
	}
	return true*/
}

// hasRoles Renvoit true si l'utilisateur détenant le token d'accès possède tous les rôles
func hasRoles(accessToken string, reqRoles []string) bool {
	ctx := context.Background()
	userInfo, err := client.GetUserInfo(ctx, accessToken, realm)
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

	for _, reqRole := range reqRoles {
		if !util.RecordsContains(res.Records, 0, reqRole) {
			return false
		}
		continue
	}

	return true
}

// Login redirige vers la page d'authentification Keycloak
func Login(c echo.Context) error {
	origin := c.Request().Header.Get(refererHeaderKey)
	pUrl, _ := url.Parse(origin)
	path := pUrl.Path
	if pUrl.RawQuery != "" {
		path += "?" + pUrl.RawQuery
	}

	//Pour éviter de tourner en boucle sur la page de connection/déconnection
	if path == util.LogoutPath {
		path = ""
	}

	//Configuration oauth2 pour le serveur Keycloak
	oauth2Config := oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  "https://bulle.rezel.net" + util.CallbackLoginPath,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID},
	}

	//Si la requête HTML vient d'HTMX, on renvoit une page vide lui demandant de se rediriger automatiquement vers la page d'authentification Keycloak
	if c.Request().Header.Get("HX-Request") == "true" {
		c.Response().Header().Set("HX-Redirect", oauth2Config.AuthCodeURL(url.QueryEscape(path)))
		return c.NoContent(http.StatusOK)
	} else {
		//Sinon on redirige nous même
		return c.Redirect(http.StatusTemporaryRedirect, oauth2Config.AuthCodeURL(url.QueryEscape(path)))
	}

}

// LoginCallback est une page intermédiaire de redirection après la connection Keycloak, elle s'occupe d'ajouter les nouveaux cookies dans la réponse
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
	//Keycloak nous renvoit un code, qu'on échange pour les tokens d'accès et de rafraichissement
	token, err := oauth2Config.Exchange(ctx, c.QueryParam("code"))
	if err != nil {
		logger.ErrorLogger.Printf("Error exchanging code: %s\n", err)
		return c.NoContent(http.StatusBadRequest)
	}

	//On récupère les informations de l'utilisateur
	//Notamment son numéro de carte de crédit, et les 3 chiffres au dos
	uuid, name, ok := GetUserInfo(token.AccessToken)
	if !ok {
		logger.ErrorLogger.Println("Error getting user infos")
		return c.NoContent(http.StatusBadRequest)
	}

	// Création/mise à jour des infos utilisateur à la connection
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

	//On renvoit l'utilisateur sur la page qu'il avait à l'origine, avant la connection
	addCookies(c, token.AccessToken, token.RefreshToken)
	path, _ := url.QueryUnescape(origin)
	return c.Redirect(http.StatusPermanentRedirect, "https://bulle.rezel.net"+path)
}

// Logout déconnecte l'utilisateur, et le renvoit à sa page d'origine
func Logout(c echo.Context) error {
	//Suppression des cookies
	c.SetCookie(&http.Cookie{
		Name:     accessTokenCookie,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})
	c.SetCookie(&http.Cookie{
		Name:     refreshTokenCookie,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})

	//On récupère l'url de redirection
	origin := c.Request().Header.Get(refererHeaderKey)
	pUrl, _ := url.Parse(origin)
	path := pUrl.Path
	if pUrl.RawQuery != "" {
		path += "?" + pUrl.RawQuery
	}

	//Les urls contenant les caractères suivant sont considérés comme une url invalide (Invalid url) pour la déconnection Keycloak. On renvoit à la page d'acceuil dans ce cas
	if strings.ContainsAny(path, "+% ") {
		path = ""
	}

	//On construit l'url de déconnection Keycloak
	redirectUrl := "https://bulle.rezel.net" + path
	logoutURL := authUrl + "/realms/" + realm + "/protocol/openid-connect/logout"
	logoutURL += "?post_logout_redirect_uri=" + redirectUrl
	logoutURL += "&client_id=" + clientID
	return c.Redirect(http.StatusTemporaryRedirect, logoutURL)
}

// IsAdmin renvoit true si l'utilisateur du contexte possède le rôle admin, en rafraichissant les tokens si nécessaires
func IsAdmin(c echo.Context) bool {
	if !IsLogged(c) {
		return false
	}
	accessToken, _, _ := getTokens(c)
	return hasRoles(accessToken, []string{adminRoleName})
}

// GetUserInfo l'UUID de l'utilisateur, son nom complet, et un boolean de confirmation, à partir du token d'accès
func GetUserInfo(accessToken string) (string, string, bool) {
	info, err := client.GetUserInfo(context.Background(), accessToken, realm)
	if err != nil {
		logger.InfoLogger.Printf("Error getting user info: %s\n", err)
		return "", "", false
	}
	return *info.Sub, *info.Name, true
}

// GetUserInfoFromContext l'UUID de l'utilisateur, son nom complet, et un boolean de confirmation, à partir du contexte
func GetUserInfoFromContext(c echo.Context) (string, string, bool) {
	if !IsLogged(c) {
		return "", "", false
	}
	accessToken, _, _ := getTokens(c)
	return GetUserInfo(accessToken)
}

// GetUserUUID renvoit uniquement l'UUID de l'utilisateur, ou "" en cas d'erreur, à partir du token d'accès
func GetUserUUID(c echo.Context) string {
	uuid, _, ok := GetUserInfoFromContext(c)
	if !ok {
		return ""
	}
	return uuid
}

// HasTokenMiddleware intervient lorsqu'on utilise un chemin protégé, et vérifie qu'on est bien authentifié
func HasTokenMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if IsLogged(c) {
			return next(c)
		}

		//Si l'utilisateur n'est pas authentifié, on le redirige sur la page de connection
		return Login(c)
	}
}

// HasRoleMiddleware intervient lorsqu'on utilise un chemin protégé par le rôle admin, et vérifie qu'on le possède
func HasRoleMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		//L'utilisateur doit être authentifié
		if !IsLogged(c) {
			//On ajoute l'url actuelle dans le header de la requête pour que la page de connection nous renvoit sur la page actuelle
			c.Request().Header.Set(refererHeaderKey, c.Path())
			return Login(c)
		}
		accessToken, _, _ := getTokens(c)
		//On vérifie que l'utilisateur possède le role
		if !hasRoles(accessToken, []string{adminRoleName}) {
			return c.NoContent(http.StatusForbidden)
		}
		return next(c)
	}
}

func RefreshTokenMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ac, err1 := c.Request().Cookie(accessTokenCookie)
		rc, err2 := c.Request().Cookie(refreshTokenCookie)
		if err1 != nil || err2 != nil {
			logger.InfoLogger.Println("Found no cookies")
			return next(c)
		}

		result, err := client.RetrospectToken(ctx, ac.Value, clientID, clientSecret, realm)
		if err == nil {
			logger.InfoLogger.Println("Cookie is up to date")
			logger.InfoLogger.Printf("%t\n", *result.Active)
			return next(c)
		}

		if !*result.Active {
			jwt, err := client.RefreshToken(ctx, rc.Value, clientID, clientSecret, realm)
			if err != nil {
				logger.ErrorLogger.Printf("Error refreshing token: %s\n", err)
				return next(c)
			}
			addCookies(c, jwt.AccessToken, jwt.RefreshToken)
			logger.InfoLogger.Printf("old token: %s\n", ac.Value)
			ac.Value = jwt.AccessToken
			rc.Value = jwt.RefreshToken
			logger.InfoLogger.Printf("new token: %s\n", ac.Value)
		}
		return next(c)
	}
}
