package stormpathweb

import (
	"bytes"
	"time"

	"github.com/jarias/stormpath-sdk-go"
	"github.com/spf13/viper"
)

//Config is holds the load configuration from the different locations and files
var Config = &webConfig{}

type webConfig struct {
	Produces []string
	BasePath string
	//Login
	LoginURI     string
	LoginNextURI string
	LoginView    string
	LoginEnabled bool
	//Logout
	LogoutURI     string
	LogoutNextURI string
	LogoutEnabled bool
	//Register
	RegisterURI              string
	RegisterNextURI          string
	RegisterView             string
	RegisterEnabled          bool
	RegisterAutoLoginEnabled bool
	//Forgot password
	ForgotPasswordURI     string
	ForgotPasswordNextURI string
	ForgotPasswordView    string
	ForgotPasswordEnabled *bool
	//Verify
	VerifyURI     string
	VerifyNextURI string
	VerifyView    string
	VerifyEnabled *bool
	//Change password
	ChangePasswordURI              string
	ChangePasswordNextURI          string
	ChangePasswordAutoLoginEnabled bool
	ChangePasswordView             string
	ChangePasswordErrorURI         string
	ChangePasswordEnabled          *bool
	//Social
	FacebookCallbackURI string
	FacebookScope       string
	GoogleCallbackURI   string
	GoogleScope         string
	GithubCallbackURI   string
	GithubScope         string
	LinkedinCallbackURI string
	LinkedinScope       string
	//IDSite
	IDSiteEnabled     bool
	IDSiteLoginURI    string
	IDSiteForgotURI   string
	IDSiteRegisterURI string
	//Callback
	CallbackEnabled bool
	CallbackURI     string
	//Access Token Cookie
	AccessTokenCookieHTTPOnly bool
	AccessTokenCookieName     string
	AccessTokenCookieSecure   *bool
	AccessTokenCookiePath     string
	AccessTokenCookieDomain   string
	//Refresh Token Cookie
	RefreshTokenCookieHTTPOnly bool
	RefreshTokenCookieName     string
	RefreshTokenCookieSecure   *bool
	RefreshTokenCookiePath     string
	RefreshTokenCookieDomain   string
	//OAuth2
	OAuth2Enabled                                  bool
	OAuth2URI                                      string
	OAuth2ClientCredentialsGrantTypeEnabled        bool
	OAuth2ClientCredentialsGrantTypeAccessTokenTTL time.Duration
	OAuth2PasswordGrantTypeEnabled                 bool
	OAuth2PasswordGrantTypeValidationStrategy      string
	//Merge
	MeEnabled bool
	MeURI     string
	MeExpand  map[string]interface{}
}

func loadConfig() {
	stormpath.InitLog()

	viper.SetConfigType("yaml")

	//Load bundled default config
	defaultConfig, err := Asset("config/web.stormpath.yaml")
	if err != nil {
		stormpath.Logger.Panicf("[ERROR] Couldn't load default bundle configuration: %s", err)
	}

	viper.ReadConfig(bytes.NewBuffer(defaultConfig))

	//Merge users custom configuration
	viper.SetConfigFile("stormpath.yaml")
	viper.AddConfigPath("~/.stormpath/")
	viper.AddConfigPath(".")
	err = viper.MergeInConfig()
	if err != nil {
		stormpath.Logger.Println("[WARN] User didn't provide custom configuration")
	}

	Config.Produces = viper.GetStringSlice("stormpath.web.produces")
	Config.BasePath = viper.GetString("stormpath.web.basePath")

	loadSocialConfig()
	loadCookiesConfig()
	loadEndpointsConfig()
	loadOAuth2Config()
}

func loadOAuth2Config() {
	Config.OAuth2Enabled = viper.GetBool("stormpath.web.oauth2.enabled")
	Config.OAuth2URI = viper.GetString("stormpath.web.oauth2.uri")
	Config.OAuth2ClientCredentialsGrantTypeEnabled = viper.GetBool("stormpath.web.oauth2.client_credentials.enabled")
	Config.OAuth2ClientCredentialsGrantTypeAccessTokenTTL = time.Duration(viper.GetInt("stormpath.web.oauth2.client_credentials.accessToken.ttl")) * time.Second
	Config.OAuth2PasswordGrantTypeEnabled = viper.GetBool("stormpath.web.oauth2.password.enabled")
	Config.OAuth2PasswordGrantTypeValidationStrategy = viper.GetString("stormpath.web.oauth2.password.validationStrategy")
}

func loadSocialConfig() {
	Config.FacebookCallbackURI = viper.GetString("stormpath.web.social.facebook.uri")
	Config.FacebookScope = viper.GetString("stormpath.web.social.facebook.scope")
	Config.GoogleCallbackURI = viper.GetString("stormpath.web.social.google.uri")
	Config.GoogleScope = viper.GetString("stormpath.web.social.google.scope")
	Config.LinkedinCallbackURI = viper.GetString("stormpath.web.social.linkedin.uri")
	Config.LinkedinScope = viper.GetString("stormpath.web.social.linkedin.scope")
	Config.GithubCallbackURI = viper.GetString("stormpath.web.social.github.uri")
	Config.GithubScope = viper.GetString("stormpath.web.social.github.scope")
}

func loadCookiesConfig() {
	//AccessToken
	Config.AccessTokenCookieHTTPOnly = viper.GetBool("stormpath.web.accessTokenCookie.httpOnly")
	Config.AccessTokenCookieName = viper.GetString("stormpath.web.accessTokenCookie.name")
	Config.AccessTokenCookieSecure = loadBoolPtr("stormpath.web.accessTokenCookie.secure")
	Config.AccessTokenCookiePath = viper.GetString("stormpath.web.accessTokenCookie.path")
	Config.AccessTokenCookieDomain = viper.GetString("stormpath.web.accessTokenCookie.domain")
	//RefreshToken
	Config.RefreshTokenCookieHTTPOnly = viper.GetBool("stormpath.web.refreshTokenCookie.httpOnly")
	Config.RefreshTokenCookieName = viper.GetString("stormpath.web.refreshTokenCookie.name")
	Config.RefreshTokenCookieSecure = loadBoolPtr("stormpath.web.refreshTokenCookie.secure")
	Config.RefreshTokenCookiePath = viper.GetString("stormpath.web.refreshTokenCookie.path")
	Config.RefreshTokenCookieDomain = viper.GetString("stormpath.web.refreshTokenCookie.domain")
}

func loadEndpointsConfig() {
	//Login
	Config.LoginURI = viper.GetString("stormpath.web.login.uri")
	Config.LoginNextURI = viper.GetString("stormpath.web.login.nextUri")
	Config.LoginView = viper.GetString("stormpath.web.login.view")
	Config.LoginEnabled = viper.GetBool("stormpath.web.login.enabled")
	//Register
	Config.RegisterURI = viper.GetString("stormpath.web.register.uri")
	Config.RegisterView = viper.GetString("stormpath.web.register.view")
	Config.RegisterNextURI = viper.GetString("stormpath.web.register.uri")
	Config.RegisterEnabled = viper.GetBool("stormpath.web.register.enabled")
	Config.RegisterAutoLoginEnabled = viper.GetBool("stormpath.web.register.autoLogin")
	//Verify
	Config.VerifyURI = viper.GetString("stormpath.web.verifyEmail.uri")
	Config.VerifyEnabled = loadBoolPtr("stormpath.web.verifyEmail.enabled")
	Config.VerifyView = viper.GetString("stormpath.web.verifyEmail.view")
	Config.VerifyNextURI = viper.GetString("stormpath.web.verifyEmail.nextUri")
	//Forgot Password
	Config.ForgotPasswordURI = viper.GetString("stormpath.web.forgotPassword.uri")
	Config.ForgotPasswordNextURI = viper.GetString("stormpath.web.forgotPassword.nextUri")
	Config.ForgotPasswordView = viper.GetString("stormpath.web.forgotPassword.view")
	Config.ForgotPasswordEnabled = loadBoolPtr("stormpath.web.forgotPassword.enabled")
	//Change Password
	Config.ChangePasswordURI = viper.GetString("stormpath.web.changePassword.uri")
	Config.ChangePasswordNextURI = viper.GetString("stormpath.web.changePassword.nextUri")
	Config.ChangePasswordView = viper.GetString("stormpath.web.changePassword.view")
	Config.ChangePasswordEnabled = loadBoolPtr("stormpath.web.changePassword.enabled")
	Config.ChangePasswordAutoLoginEnabled = viper.GetBool("stormpath.web.changePassword.autoLogin")
	Config.ChangePasswordErrorURI = viper.GetString("stormpath.web.changePassword.errorUri")
	//Logout
	Config.LogoutURI = viper.GetString("stormpath.web.logout.uri")
	Config.LogoutNextURI = viper.GetString("stormpath.web.logout.nextUri")
	Config.LogoutEnabled = viper.GetBool("stormpath.web.logout.enabled")
	//IDSite
	Config.IDSiteEnabled = viper.GetBool("stormpath.web.idSite.enabled")
	Config.IDSiteLoginURI = viper.GetString("stormpath.web.idSite.loginUri")
	Config.IDSiteForgotURI = viper.GetString("stormpath.web.idSite.forgotUri")
	Config.IDSiteRegisterURI = viper.GetString("stormpath.web.idSite.registerUri")
	Config.CallbackEnabled = viper.GetBool("stormpath.web.callback.enabled")
	Config.CallbackURI = viper.GetString("stormpath.web.callback.uri")
	//Me
	Config.MeEnabled = viper.GetBool("stormpath.web.me.enabled")
	Config.MeURI = viper.GetString("stormpath.web.me.uri")
	Config.MeExpand = viper.GetStringMap("stormpath.web.me.expand")
}

func loadBoolPtr(key string) *bool {
	val := viper.Get(key)
	if val == nil {
		return nil
	}
	b := viper.GetBool(key)
	return &b
}

func IsForgotPasswordEnabled(application *stormpath.Application) bool {
	mapping := application.DefaultAccountStoreMapping

	if mapping != nil && mapping.IsAccountStoreDirectory() {
		directory, err := stormpath.GetDirectory(mapping.AccountStore.Href, stormpath.MakeDirectoryCriteria().WithAccountCreationPolicy().WithPasswordPolicy())
		if err != nil {
			return false
		}
		if Config.ForgotPasswordEnabled != nil {
			return stormpath.Enabled == directory.PasswordPolicy.ResetEmailStatus && *Config.ForgotPasswordEnabled
		}
		return stormpath.Enabled == directory.PasswordPolicy.ResetEmailStatus
	}

	return false
}

func IsVerifyEnabled(application *stormpath.Application) bool {
	mapping := application.DefaultAccountStoreMapping

	if mapping != nil && mapping.IsAccountStoreDirectory() {
		directory, err := stormpath.GetDirectory(mapping.AccountStore.Href, stormpath.MakeDirectoryCriteria().WithAccountCreationPolicy().WithPasswordPolicy())
		if err != nil {
			return false
		}
		if Config.VerifyEnabled != nil {
			return stormpath.Enabled == directory.AccountCreationPolicy.VerificationEmailStatus && *Config.VerifyEnabled
		}
		return stormpath.Enabled == directory.AccountCreationPolicy.VerificationEmailStatus
	}

	return false
}
