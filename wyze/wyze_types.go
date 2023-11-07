package wyze

type WyzeAccessTokenRequest struct {
  AppVer string           `json:"app_ver"`
  PhoneId string          `json:"phone_id"`
  RefreshToken string      `json:"refresh_token"`
  SC string               `json:"sc"`
  SV string               `json:"sv"`
  TS string               `json:"ts"`
}

type WyzeAccessTokenData struct {
  AccessToken  string`json:"access_token"`
}

type WyzeAccessTokenResponse struct {
  Data WyzeAccessTokenData `json:"data"`
}

type WyzeDevice struct {
  MAC string      `json:"mac"`
  Nickname string `json:"nickname"`
}

type WyzeDeviceData struct {
  DeviceList []WyzeDevice `json:"device_list"`
}

type WyzeDeviceListResponse struct {
  Timestamp int       `json:"ts"`
  Code string         `json:"code"`
  Message string      `json:"msg"`
  Data WyzeDeviceData `json:"data"`
}

type WyzeDeviceListRequest struct {
  AppVer string           `json:"app_ver"`
  PhoneId string          `json:"phone_id"`
  AccessToken string      `json:"access_token"`
  SC string               `json:"sc"`
  SV string               `json:"sv"`
  TS string               `json:"ts"`
}

type WyzeEventRequest struct {
  AppVer string           `json:"app_ver"`
  PhoneId string          `json:"phone_id"`
  AccessToken string      `json:"access_token"`
  SC string               `json:"sc"`
  SV string               `json:"sv"`
  Devices []string        `json:"device_mac_list"`
  Count int            `json:"count"`
  OrderBy string          `json:"order_by"`
  PhoneSystemType string  `json:"phone_system_type"`
  BeginTime string        `json:"begin_time"`
  EndTime string          `json:"end_time"`
  TS string               `json:"ts"`
}

type WyzeDownloadedFile struct {
  Path string
  Url string
  Mac string
  Timestamp int64
}

type WyzeFile struct {
  URL string `json:"url"`
}

type WyzeCamEvent struct {
  DeviceMac string    `json:"device_mac"`
  EventTime int       `json:"event_ts"`
  FileList []WyzeFile `json:"file_list"`
}

type WyzeEventResponseData struct {
  EventList []WyzeCamEvent `json:"event_list"`
}

type WyzeEventResponse struct {
  Timestamp int              `json:"ts"`
  Code string                `json:"code"`
  Message string             `json:"msg"`
  Data WyzeEventResponseData `json:"data"`
}
