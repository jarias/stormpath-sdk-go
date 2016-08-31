package stormpathweb

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jarias/stormpath-sdk-go"
	"github.com/spf13/viper"
)

//Config is holds the load configuration from the different locations and files
var Config = &webConfig{}

type webConfig struct {
	//Application
	ApplicationName string
	ApplicationHref string
	//General web
	Produces []string
	BasePath string
	//Login
	LoginURI     string
	LoginNextURI string
	LoginView    string
	LoginEnabled bool
	LoginForm    form
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
	RegisterForm             form
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

	v := viper.New()

	v.SetConfigType("yaml")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	//Load bundled default config
	defaultConfig, err := Asset("config/web.stormpath.yaml")
	if err != nil {
		stormpath.Logger.Panicf("[ERROR] Couldn't load default bundle configuration: %s", err)
	}

	v.ReadConfig(bytes.NewBuffer(defaultConfig))

	//Merge users custom configuration
	v.SetConfigName("stormpath")
	v.AddConfigPath(os.Getenv("HOME") + "/.stormpath")
	v.AddConfigPath(".")
	err = v.MergeInConfig()
	if err != nil {
		stormpath.Logger.Println("[WARN] User didn't provide custom configuration")
	}

	Config.Produces = v.GetStringSlice("stormpath.web.produces")
	Config.BasePath = v.GetString("stormpath.web.basePath")

	Config.ApplicationHref = v.GetString("stormpath.application.href")
	Config.ApplicationName = v.GetString("stormpath.application.name")

	loadSocialConfig(v)
	loadCookiesConfig(v)
	loadEndpointsConfig(v)
	loadOAuth2Config(v)
}

func loadOAuth2Config(v *viper.Viper) {
	Config.OAuth2Enabled = v.GetBool("stormpath.web.oauth2.enabled")
	Config.OAuth2URI = v.GetString("stormpath.web.oauth2.uri")
	Config.OAuth2ClientCredentialsGrantTypeEnabled = v.GetBool("stormpath.web.oauth2.client_credentials.enabled")
	Config.OAuth2ClientCredentialsGrantTypeAccessTokenTTL = time.Duration(v.GetInt("stormpath.web.oauth2.client_credentials.accessToken.ttl")) * time.Second
	Config.OAuth2PasswordGrantTypeEnabled = v.GetBool("stormpath.web.oauth2.password.enabled")
	Config.OAuth2PasswordGrantTypeValidationStrategy = v.GetString("stormpath.web.oauth2.password.validationStrategy")
}

func loadSocialConfig(v *viper.Viper) {
	Config.FacebookCallbackURI = v.GetString("stormpath.web.social.facebook.uri")
	Config.FacebookScope = v.GetString("stormpath.web.social.facebook.scope")
	Config.GoogleCallbackURI = v.GetString("stormpath.web.social.google.uri")
	Config.GoogleScope = v.GetString("stormpath.web.social.google.scope")
	Config.LinkedinCallbackURI = v.GetString("stormpath.web.social.linkedin.uri")
	Config.LinkedinScope = v.GetString("stormpath.web.social.linkedin.scope")
	Config.GithubCallbackURI = v.GetString("stormpath.web.social.github.uri")
	Config.GithubScope = v.GetString("stormpath.web.social.github.scope")
}

func loadCookiesConfig(v *viper.Viper) {
	//AccessToken
	Config.AccessTokenCookieHTTPOnly = v.GetBool("stormpath.web.accessTokenCookie.httpOnly")
	Config.AccessTokenCookieName = v.GetString("stormpath.web.accessTokenCookie.name")
	Config.AccessTokenCookieSecure = loadBoolPtr("stormpath.web.accessTokenCookie.secure", v)
	Config.AccessTokenCookiePath = v.GetString("stormpath.web.accessTokenCookie.path")
	Config.AccessTokenCookieDomain = v.GetString("stormpath.web.accessTokenCookie.domain")
	//RefreshToken
	Config.RefreshTokenCookieHTTPOnly = v.GetBool("stormpath.web.refreshTokenCookie.httpOnly")
	Config.RefreshTokenCookieName = v.GetString("stormpath.web.refreshTokenCookie.name")
	Config.RefreshTokenCookieSecure = loadBoolPtr("stormpath.web.refreshTokenCookie.secure", v)
	Config.RefreshTokenCookiePath = v.GetString("stormpath.web.refreshTokenCookie.path")
	Config.RefreshTokenCookieDomain = v.GetString("stormpath.web.refreshTokenCookie.domain")
}

func loadEndpointsConfig(v *viper.Viper) {
	//Login
	Config.LoginURI = v.GetString("stormpath.web.login.uri")
	Config.LoginNextURI = v.GetString("stormpath.web.login.nextUri")
	Config.LoginView = v.GetString("stormpath.web.login.view")
	Config.LoginEnabled = v.GetBool("stormpath.web.login.enabled")
	Config.LoginForm = buildForm("login", v)
	//Register
	Config.RegisterURI = v.GetString("stormpath.web.register.uri")
	Config.RegisterView = v.GetString("stormpath.web.register.view")
	Config.RegisterNextURI = v.GetString("stormpath.web.register.uri")
	Config.RegisterEnabled = v.GetBool("stormpath.web.register.enabled")
	Config.RegisterAutoLoginEnabled = v.GetBool("stormpath.web.register.autoLogin")
	Config.RegisterForm = buildForm("register", v)
	//Verify
	Config.VerifyURI = v.GetString("stormpath.web.verifyEmail.uri")
	Config.VerifyEnabled = loadBoolPtr("stormpath.web.verifyEmail.enabled", v)
	Config.VerifyView = v.GetString("stormpath.web.verifyEmail.view")
	Config.VerifyNextURI = v.GetString("stormpath.web.verifyEmail.nextUri")
	//Forgot Password
	Config.ForgotPasswordURI = v.GetString("stormpath.web.forgotPassword.uri")
	Config.ForgotPasswordNextURI = v.GetString("stormpath.web.forgotPassword.nextUri")
	Config.ForgotPasswordView = v.GetString("stormpath.web.forgotPassword.view")
	Config.ForgotPasswordEnabled = loadBoolPtr("stormpath.web.forgotPassword.enabled", v)
	//Change Password
	Config.ChangePasswordURI = v.GetString("stormpath.web.changePassword.uri")
	Config.ChangePasswordNextURI = v.GetString("stormpath.web.changePassword.nextUri")
	Config.ChangePasswordView = v.GetString("stormpath.web.changePassword.view")
	Config.ChangePasswordEnabled = loadBoolPtr("stormpath.web.changePassword.enabled", v)
	Config.ChangePasswordAutoLoginEnabled = v.GetBool("stormpath.web.changePassword.autoLogin")
	Config.ChangePasswordErrorURI = v.GetString("stormpath.web.changePassword.errorUri")
	//Logout
	Config.LogoutURI = v.GetString("stormpath.web.logout.uri")
	Config.LogoutNextURI = v.GetString("stormpath.web.logout.nextUri")
	Config.LogoutEnabled = v.GetBool("stormpath.web.logout.enabled")
	//IDSite
	Config.IDSiteEnabled = v.GetBool("stormpath.web.idSite.enabled")
	Config.IDSiteLoginURI = v.GetString("stormpath.web.idSite.loginUri")
	Config.IDSiteForgotURI = v.GetString("stormpath.web.idSite.forgotUri")
	Config.IDSiteRegisterURI = v.GetString("stormpath.web.idSite.registerUri")
	Config.CallbackEnabled = v.GetBool("stormpath.web.callback.enabled")
	Config.CallbackURI = v.GetString("stormpath.web.callback.uri")
	//Me
	Config.MeEnabled = v.GetBool("stormpath.web.me.enabled")
	Config.MeURI = v.GetString("stormpath.web.me.uri")
	Config.MeExpand = v.GetStringMap("stormpath.web.me.expand")
}

func loadBoolPtr(key string, v *viper.Viper) *bool {
	val := v.Get(key)
	if val == nil {
		return nil
	}
	b := v.GetBool(key)
	return &b
}

func isForgotPasswordEnabled(application *stormpath.Application) bool {
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

func isVerifyEnabled(application *stormpath.Application) bool {
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

func buildForm(formName string, v *viper.Viper) form {
	form := form{}

	for _, fieldName := range getConfiguredFormFieldNames(formName, v) {
		field := field{
			Name:        fieldName,
			Label:       v.GetString("stormpath.web." + formName + ".form.fields." + fieldName + ".label"),
			PlaceHolder: v.GetString("stormpath.web." + formName + ".form.fields." + fieldName + ".placeHolder"),
			Visible:     v.GetBool("stormpath.web." + formName + ".form.fields." + fieldName + ".visible"),
			Enabled:     v.GetBool("stormpath.web." + formName + ".form.fields." + fieldName + ".enabled"),
			Required:    v.GetBool("stormpath.web." + formName + ".form.fields." + fieldName + ".required"),
			Type:        v.GetString("stormpath.web." + formName + ".form.fields." + fieldName + ".type"),
		}
		if field.Enabled {
			form.Fields = append(form.Fields, field)
		}
	}

	return form
}

func getConfiguredFormFieldNames(formName string, v *viper.Viper) []string {
	configuredFields := v.GetStringMapString("stormpath.web." + formName + ".form.fields")
	fieldOrder := v.GetStringSlice("stormpath.web." + formName + ".form.fieldOrder")

	for fieldName := range configuredFields {
		if !contains(fieldOrder, fieldName) {
			fieldOrder = append(fieldOrder, fieldName)
		}
	}
	return fieldOrder
}

func resolveAccountStores(application *stormpath.Application) {
	//see https://github.com/stormpath/stormpath-framework-spec/blob/master/configuration.md
	mappings, err := application.GetAccountStoreMappings(stormpath.MakeApplicationAccountStoreMappingsCriteria())
	if err != nil || len(mappings.Items) == 0 {
		panic(fmt.Errorf("No account stores are mapped to the specified application. Account stores are required for login and registration. \n"))
	}

	if application.DefaultAccountStoreMapping == nil && Config.RegisterEnabled {
		panic(fmt.Errorf("No default account store is mapped to the specified application. A default account store is required for registration. \n"))
	}
}

func resolveApplication() *stormpath.Application {
	//see https://github.com/stormpath/stormpath-framework-spec/blob/master/configuration.md
	applicationHref := Config.ApplicationHref
	applicationName := Config.ApplicationName

	tenant, err := stormpath.CurrentTenant()
	if err != nil {
		panic(fmt.Errorf("Fatal couldn't get current tenant: %s \n", err))
	}

	if applicationHref != "" {
		if !strings.Contains(applicationHref, "/applications/") {
			panic(fmt.Errorf("(%s) is not a valid Stormpath Application href \n", applicationHref))
		}

		application, err := stormpath.GetApplication(applicationHref, stormpath.MakeApplicationCriteria().WithDefaultAccountStoreMapping())
		if err != nil {
			panic(fmt.Errorf("The provided application could not be found. The provided application href was: %s \n", applicationHref))
		}
		return application
	}

	if applicationName != "" {
		applications, err := tenant.GetApplications(stormpath.MakeApplicationsCriteria().NameEq(applicationName).WithDefaultAccountStoreMapping())
		if err != nil || len(applications.Items) == 0 {
			panic(fmt.Errorf("The provided application could not be found. The provided application name was: %s \n", applicationName))
		}

		return &applications.Items[0]
	}

	//Get all apps if size > 1 && <= 2 return the one that's not name "Stormpath" else error

	applications, err := tenant.GetApplications(stormpath.MakeApplicationsCriteria().WithDefaultAccountStoreMapping())

	if len(applications.Items) > 2 || len(applications.Items) == 1 {
		panic(fmt.Errorf("Could not automatically resolve a Stormpath Application. Please specify your Stormpath Application in your configuration \n"))
	}

	var application stormpath.Application

	for _, app := range applications.Items {
		if app.Name != "Stormpath" {
			application = app
			break
		}
	}

	return &application
}
