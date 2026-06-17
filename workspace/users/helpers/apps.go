package helpers

import "github.com/gold-gym/gymkit"

// App origin constants — identifies which app/client made the request via X-App-Origin header.
const AppWeb = "web"
const AppMobile = "mobile"
const AppAdmin = "admin"
const AppSystem = "system"
const AppWorker = "worker"

const LoginAsUser = 1
const LoginAsAdmin = 2
const LoginAsSystem = 3
const LoginAsWorker = 4

func CheckLoginAs(appOrigin string) int8 {
	switch appOrigin {
	case AppWeb, AppMobile:
		return LoginAsUser
	case AppAdmin:
		return LoginAsAdmin
	case AppWorker:
		return LoginAsWorker
	default:
		return LoginAsSystem
	}
}

func IsUserAppID(appOrigin string) bool {
	return map[string]bool{
		gymkit.GetEnv("APP_BUYER_NAME", AppWeb):    true,
		gymkit.GetEnv("APP_MOBILE_NAME", AppMobile): true,
	}[appOrigin]
}

func IsAdminAppID(appOrigin string) bool {
	return map[string]bool{
		gymkit.GetEnv("APP_ADMIN_NAME", AppAdmin): true,
	}[appOrigin]
}

func IsInternalAppID(appOrigin string) bool {
	return appOrigin == AppSystem || appOrigin == AppWorker
}
