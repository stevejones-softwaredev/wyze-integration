package wyze

const wyzeDeveloperApi string = "wyze_developer_api"
const wyzeRequestTimestamp string = "4070908800000" // January 1, 2099 12:00:00 AM UTC for unspecified reasons
const wyzeApiHost string = "api.wyzecam.com"
const wyzeAuthHost string = "auth-prod.api.wyze.com"
const wyzeContentType string = "application/json"

const wyzeAuthEndpoint string = "https://auth-prod.api.wyze.com/api/user/login"
const wyzeAccessTokenEndpoint string = "https://api.wyzecam.com/app/user/refresh_token"
const wyzeGetDeviceListEndpoint string = "https://api.wyzecam.com/app/v2/home_page/get_object_list"
const wyzeGetEventListEndpoint string = "https://api.wyzecam.com/app/v2/device/get_event_list"
