package helpers

import "github.com/gold-gym/gymkit"

const AppWeb = "web"
const AppMobile = "mobile"
const AppAdmin = "admin"
const AppSystem = "system"
const AppWorker = "worker"

func IsUserAppID(appOrigin string) bool {
	return map[string]bool{
		gymkit.GetEnv("APP_BUYER_NAME", AppWeb):     true,
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
