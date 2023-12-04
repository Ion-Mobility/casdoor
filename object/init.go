// Copyright 2021 The Casdoor Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package object

import (
	"fmt"
	"os"

	"github.com/casdoor/casdoor/util"
)

func InitDb() {
	existed := initBuiltInOrganization()
	if !existed {
		// initBuiltInPermission()
		initBuiltInProvider()
		initBuiltInUser()
		initBuiltInApplication()
		initBuiltInCert()
		// initBuiltInLdap()
	}

	existed = initBuiltInApiModel()
	if !existed {
		initBuiltInApiAdapter()
		initBuiltInApiEnforcer()
		initBuiltInUserModel()
		initBuiltInUserAdapter()
		initBuiltInUserEnforcer()
	}

	// initWebAuthn()
}

func getBuiltInAccountItems() []*AccountItem {
	return []*AccountItem{
		{Name: "Organization", Visible: true, ViewRule: "Public", ModifyRule: "Admin"},
		{Name: "ID", Visible: true, ViewRule: "Public", ModifyRule: "Immutable"},
		{Name: "Name", Visible: true, ViewRule: "Public", ModifyRule: "Admin"},
		{Name: "Display name", Visible: true, ViewRule: "Public", ModifyRule: "Self"},
		{Name: "Avatar", Visible: true, ViewRule: "Public", ModifyRule: "Self"},
		{Name: "User type", Visible: true, ViewRule: "Public", ModifyRule: "Admin"},
		{Name: "Password", Visible: true, ViewRule: "Self", ModifyRule: "Self"},
		{Name: "Email", Visible: true, ViewRule: "Public", ModifyRule: "Self"},
		{Name: "Phone", Visible: true, ViewRule: "Public", ModifyRule: "Self"},
		{Name: "Country code", Visible: true, ViewRule: "Public", ModifyRule: "Admin"},
		{Name: "Country/Region", Visible: true, ViewRule: "Public", ModifyRule: "Self"},
		{Name: "Location", Visible: true, ViewRule: "Public", ModifyRule: "Self"},
		{Name: "Affiliation", Visible: true, ViewRule: "Public", ModifyRule: "Self"},
		{Name: "Title", Visible: true, ViewRule: "Public", ModifyRule: "Self"},
		{Name: "Homepage", Visible: true, ViewRule: "Public", ModifyRule: "Self"},
		{Name: "Bio", Visible: true, ViewRule: "Public", ModifyRule: "Self"},
		{Name: "Tag", Visible: true, ViewRule: "Public", ModifyRule: "Admin"},
		{Name: "Signup application", Visible: true, ViewRule: "Public", ModifyRule: "Admin"},
		{Name: "Roles", Visible: true, ViewRule: "Public", ModifyRule: "Immutable"},
		{Name: "Permissions", Visible: true, ViewRule: "Public", ModifyRule: "Immutable"},
		{Name: "Groups", Visible: true, ViewRule: "Public", ModifyRule: "Admin"},
		{Name: "3rd-party logins", Visible: true, ViewRule: "Self", ModifyRule: "Self"},
		{Name: "Properties", Visible: false, ViewRule: "Admin", ModifyRule: "Admin"},
		{Name: "Is admin", Visible: true, ViewRule: "Admin", ModifyRule: "Admin"},
		{Name: "Is forbidden", Visible: true, ViewRule: "Admin", ModifyRule: "Admin"},
		{Name: "Is deleted", Visible: true, ViewRule: "Admin", ModifyRule: "Admin"},
		{Name: "Multi-factor authentication", Visible: true, ViewRule: "Self", ModifyRule: "Self"},
		{Name: "WebAuthn credentials", Visible: true, ViewRule: "Self", ModifyRule: "Self"},
		{Name: "Managed accounts", Visible: true, ViewRule: "Self", ModifyRule: "Self"},
	}
}

func initBuiltInOrganization() bool {
	organization, err := getOrganization("admin", "built-in")
	if err != nil {
		panic(err)
	}

	if organization != nil {
		return true
	}

	organization = &Organization{
		Owner:              "admin",
		Name:               "built-in",
		CreatedTime:        util.GetCurrentTime(),
		DisplayName:        os.Getenv("ORGANIZATION_DISPLAY_NAME"),
		WebsiteUrl:         os.Getenv("ORGANIZATION_URL"),
		Favicon:            os.Getenv("ORGANIZATION_IMAGE"),
		PasswordType:       "md5-salt",
		PasswordOptions:    []string{"AtLeast6"},
		CountryCodes:       []string{"VN", "ID", "SG", "CN"},
		DefaultAvatar:      os.Getenv("ORGANIZATION_IMAGE"),
		Tags:               []string{},
		Languages:          []string{"vi", "id", "en"},
		InitScore:          2000,
		AccountItems:       getBuiltInAccountItems(),
		EnableSoftDeletion: false,
		IsProfilePublic:    false,
	}
	_, err = AddOrganization(organization)
	if err != nil {
		panic(err)
	}

	return false
}

func initBuiltInUser() {
	user, err := getUser("built-in", "admin")
	if err != nil {
		panic(err)
	}
	if user != nil {
		return
	}

	user = &User{
		Owner:             "built-in",
		Name:              "admin",
		CreatedTime:       util.GetCurrentTime(),
		Id:                util.GenerateId(),
		Type:              "normal-user",
		Password:          os.Getenv("ADMIN_PASSWORD"),
		DisplayName:       "Admin",
		Email:             os.Getenv("ADMIN_EMAIL"),
		Phone:             os.Getenv("ADMIN_PHONE"),
		CountryCode:       "SG",
		Address:           []string{},
		Tag:               "staff",
		Score:             2000,
		Ranking:           1,
		IsAdmin:           true,
		IsForbidden:       false,
		IsDeleted:         false,
		SignupApplication: "app-built-in",
		CreatedIp:         os.Getenv("IP_ADDRESS"),
		Properties:        make(map[string]string),
	}
	_, err = AddUser(user)
	if err != nil {
		panic(err)
	}
}

func initBuiltInApplication() {
	application, err := getApplication("admin", "app-built-in")
	if err != nil {
		panic(err)
	}

	if application != nil {
		return
	}

	application = &Application{
		Owner:            "admin",
		Name:             "app-built-in",
		CreatedTime:      util.GetCurrentTime(),
		DisplayName:      os.Getenv("APP_DISPLAY_NAME"),
		Logo:             os.Getenv("ORGANIZATION_IMAGE"),
		Organization:     "built-in",
		Cert:             os.Getenv("CERT_NAME"),
		EnablePassword:   true,
		EnableCodeSignin: true,
		Providers: []*ProviderItem{
			{Name: os.Getenv("SMS_TWILIO_PROVIDER_NAME"), CanSignUp: true, CanSignIn: true, CanUnlink: true, Prompted: false, SignupGroup: "", Rule: "None", Provider: nil},
		},
		SignupItems: []*SignupItem{
			{Name: "ID", Visible: false, Required: true, Prompted: false, Rule: "Random"},
			{Name: "Phone", Visible: true, Required: true, Prompted: false, Rule: "Normal"},
			{Name: "Agreement", Visible: true, Required: true, Prompted: false, Rule: "None"},
		},
		Tags:          []string{},
		RedirectUris:  []string{},
		ExpireInHours: 168,
		FormOffset:    2,
	}
	_, err = AddApplication(application)
	if err != nil {
		panic(err)
	}
}

func initBuiltInCert() {
	cert, err := getCert("admin", os.Getenv("CERT_NAME"))
	if err != nil {
		panic(err)
	}

	if cert != nil {
		return
	}

	cert = &Cert{
		Owner:           "admin",
		Name:            os.Getenv("CERT_NAME"),
		CreatedTime:     util.GetCurrentTime(),
		DisplayName:     os.Getenv("CERT_DISPLAY_NAME"),
		Scope:           "JWT",
		Type:            "x509",
		CryptoAlgorithm: "RS256",
		BitSize:         4096,
		ExpireInYears:   20,
		Certificate:     os.Getenv("JWT_CERTIFICATE"),
		PrivateKey:      os.Getenv("JWT_PRIVATE_KEY"),
	}
	_, err = AddCert(cert)
	if err != nil {
		panic(err)
	}
}

// func initBuiltInLdap() {
// 	ldap, err := GetLdap("ldap-built-in")
// 	if err != nil {
// 		panic(err)
// 	}

// 	if ldap != nil {
// 		return
// 	}

// 	ldap = &Ldap{
// 		Id:         "ldap-built-in",
// 		Owner:      "built-in",
// 		ServerName: "BuildIn LDAP Server",
// 		Host:       "example.com",
// 		Port:       389,
// 		Username:   "cn=buildin,dc=example,dc=com",
// 		Password:   "123",
// 		BaseDn:     "ou=BuildIn,dc=example,dc=com",
// 		AutoSync:   0,
// 		LastSync:   "",
// 	}
// 	_, err = AddLdap(ldap)
// 	if err != nil {
// 		panic(err)
// 	}
// }

func initBuiltInProvider() {
	provider, err := GetProvider(util.GetId("admin", os.Getenv("SMS_TWILIO_PROVIDER_NAME")))
	if err != nil {
		panic(err)
	}

	if provider != nil {
		return
	}

	provider = &Provider{
		Owner:        "admin",
		Name:         os.Getenv("SMS_TWILIO_PROVIDER_NAME"),
		CreatedTime:  util.GetCurrentTime(),
		DisplayName:  os.Getenv("SMS_TWILIO_PROVIDER_NAME"),
		Category:     os.Getenv("SMS_TWILIO_PROVIDER_NAME"),
		Type:         "Twilio SMS",
		Method:       "Normal",
		ClientId:     os.Getenv("SMS_TWILIO_ACCOUNT_SID"),
		ClientSecret: os.Getenv("SMS_TWILIO_AUTH_TOKEN"),
		TemplateCode: os.Getenv("SMS_TEMPLATE"),
		AppId:        os.Getenv("SMS_TWILIO_MESSAGE_SERVICE_SID"),
	}
	_, err = AddProvider(provider)
	if err != nil {
		panic(err)
	}
}

// func initWebAuthn() {
// 	gob.Register(webauthn.SessionData{})
// }

func initBuiltInUserModel() {
	model, err := GetModel(fmt.Sprintf("%s/user-model-built-in", "built-in"))
	if err != nil {
		panic(err)
	}

	if model != nil {
		return
	}

	model = &Model{
		Owner:       "built-in",
		Name:        "user-model-built-in",
		CreatedTime: util.GetCurrentTime(),
		DisplayName: "Built-in Model",
		ModelText: `[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act`,
	}
	_, err = AddModel(model)
	if err != nil {
		panic(err)
	}
}

func initBuiltInApiModel() bool {
	model, err := GetModel(fmt.Sprintf("%s/api-model-built-in", "built-in"))
	if err != nil {
		panic(err)
	}

	if model != nil {
		return true
	}

	modelText := `[request_definition]
r = subOwner, subName, method, urlPath, objOwner, objName

[policy_definition]
p = subOwner, subName, method, urlPath, objOwner, objName

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = (r.subOwner == p.subOwner || p.subOwner == "*") && \
    (r.subName == p.subName || p.subName == "*" || r.subName != "anonymous" && p.subName == "!anonymous") && \
    (r.method == p.method || p.method == "*") && \
    (r.urlPath == p.urlPath || p.urlPath == "*") && \
    (r.objOwner == p.objOwner || p.objOwner == "*") && \
    (r.objName == p.objName || p.objName == "*") || \
    (r.subOwner == r.objOwner && r.subName == r.objName)`

	model = &Model{
		Owner:       "built-in",
		Name:        "api-model-built-in",
		CreatedTime: util.GetCurrentTime(),
		DisplayName: "API Model",
		ModelText:   modelText,
	}
	_, err = AddModel(model)
	if err != nil {
		panic(err)
	}
	return false
}

// func initBuiltInPermission() {
// 	permission, err := GetPermission("built-in/permission-built-in")
// 	if err != nil {
// 		panic(err)
// 	}
// 	if permission != nil {
// 		return
// 	}

// 	permission = &Permission{
// 		Owner:        "built-in",
// 		Name:         "permission-built-in",
// 		CreatedTime:  util.GetCurrentTime(),
// 		DisplayName:  "Built-in Permission",
// 		Description:  "Built-in Permission",
// 		Users:        []string{"built-in/*"},
// 		Groups:       []string{},
// 		Roles:        []string{},
// 		Domains:      []string{},
// 		Model:        "model-built-in",
// 		Adapter:      "",
// 		ResourceType: "Application",
// 		Resources:    []string{"app-built-in"},
// 		Actions:      []string{"Read", "Write", "Admin"},
// 		Effect:       "Allow",
// 		IsEnabled:    true,
// 		Submitter:    "admin",
// 		Approver:     "admin",
// 		ApproveTime:  util.GetCurrentTime(),
// 		State:        "Approved",
// 	}
// 	_, err = AddPermission(permission)
// 	if err != nil {
// 		panic(err)
// 	}
// }

func initBuiltInUserAdapter() {
	adapter, err := GetAdapter("built-in/user-adapter-built-in")
	if err != nil {
		panic(err)
	}

	if adapter != nil {
		return
	}

	adapter = &Adapter{
		Owner:       "built-in",
		Name:        "user-adapter-built-in",
		CreatedTime: util.GetCurrentTime(),
		Table:       "casbin_user_rule",
		UseSameDb:   true,
	}
	_, err = AddAdapter(adapter)
	if err != nil {
		panic(err)
	}
}

func initBuiltInApiAdapter() {
	adapter, err := GetAdapter("built-in/api-adapter-built-in")
	if err != nil {
		panic(err)
	}

	if adapter != nil {
		return
	}

	adapter = &Adapter{
		Owner:       "built-in",
		Name:        "api-adapter-built-in",
		CreatedTime: util.GetCurrentTime(),
		Table:       "casbin_api_rule",
		UseSameDb:   true,
	}
	_, err = AddAdapter(adapter)
	if err != nil {
		panic(err)
	}
}

func initBuiltInUserEnforcer() {
	enforcer, err := GetEnforcer("built-in/user-enforcer-built-in")
	if err != nil {
		panic(err)
	}

	if enforcer != nil {
		return
	}

	enforcer = &Enforcer{
		Owner:       "built-in",
		Name:        "user-enforcer-built-in",
		CreatedTime: util.GetCurrentTime(),
		DisplayName: "User Enforcer",
		Model:       "built-in/user-model-built-in",
		Adapter:     "built-in/user-adapter-built-in",
	}

	_, err = AddEnforcer(enforcer)
	if err != nil {
		panic(err)
	}
}

func initBuiltInApiEnforcer() {
	enforcer, err := GetEnforcer("built-in/api-enforcer-built-in")
	if err != nil {
		panic(err)
	}

	if enforcer != nil {
		return
	}

	enforcer = &Enforcer{
		Owner:       "built-in",
		Name:        "api-enforcer-built-in",
		CreatedTime: util.GetCurrentTime(),
		DisplayName: "API Enforcer",
		Model:       "built-in/api-model-built-in",
		Adapter:     "built-in/api-adapter-built-in",
	}

	_, err = AddEnforcer(enforcer)
	if err != nil {
		panic(err)
	}
}
