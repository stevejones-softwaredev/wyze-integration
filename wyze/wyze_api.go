package wyze

import (
  "github.com/go-resty/resty/v2"
  "fmt"
  "log"
  "strconv"
  "strings"
  "net/http"
  "time"
  "io"
  "os"
)

func GetWyzeRefreshToken(client *resty.Client, username string, password string) string {
  var refreshTokenResponse WyzeRefreshTokenResponse
  
  payload := WyzeRefreshTokenRequest{
    Name: username,
    Password: password,
  }

  _, err := client.R().
    SetHeader("Content-Type", wyzeContentType).
    SetHeader("Host", wyzeAuthHost).
    SetHeader("Keyid", "b4700b6b-9120-4622-8a6c-af5c67406510").
    SetHeader("Apikey", "5wYWKpbxRpWpnG4t3Sf5zMzgqXQ1rjf7AooPYlGx6aDHIQa3rs8x1YYZGpSI").
    SetBody(&payload).
    SetResult(&refreshTokenResponse).
    Post(wyzeAuthEndpoint)

    if err != nil {
      fmt.Println(err)
      return ""
    } else {
      return refreshTokenResponse.RefreshToken
    }
}

func GetWyzeAccessToken(client *resty.Client, refreshToken string) string {
  var accessTokenResponse WyzeAccessTokenResponse

  payload := WyzeAccessTokenRequest{
    AppVer: wyzeDeveloperApi,
    PhoneId: wyzeDeveloperApi,
    RefreshToken: refreshToken,
    SC: wyzeDeveloperApi,
    SV: wyzeDeveloperApi,
    TS: wyzeRequestTimestamp,
  }

  _, err := client.R().
    SetHeader("Content-Type", wyzeContentType).
    SetHeader("Host", wyzeApiHost).
    SetBody(&payload).
    SetResult(&accessTokenResponse).
    Post(wyzeAccessTokenEndpoint)

    if err != nil {
      fmt.Println(err)
      return ""
    } else {
      return accessTokenResponse.Data.AccessToken
    }
}

func GetWyzeCamThumbnails(client *resty.Client,
    downloadDirectory string,
    accessToken string,
    count int,
    devices []string,
    begin_time time.Time,
    end_time time.Time) []WyzeDownloadedFile {
  var thumbnailPaths []WyzeDownloadedFile
  var eventResponse WyzeEventResponse

  payload := WyzeEventRequest{
    AppVer: wyzeDeveloperApi,
    PhoneId: wyzeDeveloperApi,
    AccessToken: accessToken,
    SC: wyzeDeveloperApi,
    SV: wyzeDeveloperApi,
    Devices: devices,
    Count: count,
    OrderBy: "1",
    PhoneSystemType: "1",
    BeginTime: strconv.FormatInt(begin_time.UnixMilli(), 10),
    EndTime: strconv.FormatInt(end_time.UnixMilli(), 10),
    TS: wyzeRequestTimestamp,
  }

  _, err := client.R().
    SetHeader("Content-Type", wyzeContentType).
    SetHeader("Host", wyzeApiHost).
    SetBody(&payload).
    SetResult(&eventResponse).
    Post(wyzeGetEventListEndpoint)

  if err != nil {
    fmt.Println(err)
  } else {
    for _, event := range eventResponse.Data.EventList {
      for _, 	file := range event.FileList {
        fileName := fmt.Sprintf("%s%d.jpg", downloadDirectory, event.EventTime)

        _,err := os.Stat(fileName)

        if err == nil {
          log.Println("File", fileName, "already exists")
        } else {
          saveFile(file.URL, fileName)
          log.Println("File", fileName, "downloaded")
          downloadedFile := WyzeDownloadedFile{
            Path: fileName,
            Url: file.URL,
            Mac: event.DeviceMac,
            Timestamp: getTimestampFromFile(fileName),
          }
          thumbnailPaths = append(thumbnailPaths, downloadedFile)
        }
      }
    }
  }

  return thumbnailPaths
}

func getTimestampFromFile(filePath string) int64 {
  ts := strings.Split(filePath, ".")[0]
  slice := strings.Split(ts, "/")
  ts = strings.Split(ts, "/")[len(slice) - 1]
  value,_ := strconv.ParseInt(ts, 10, 64)

  return value
}

func GetWyzeDeviceList(client *resty.Client,
    accessToken string) []WyzeDevice {
  var deviceListResponse WyzeDeviceListResponse

  payload := WyzeDeviceListRequest{
    AppVer: wyzeDeveloperApi,
    PhoneId: wyzeDeveloperApi,
    AccessToken: accessToken,
    SC: wyzeDeveloperApi,
    SV: wyzeDeveloperApi,
    TS: wyzeRequestTimestamp,
  }

  _, err := client.R().
    SetHeader("Content-Type", wyzeContentType).
    SetHeader("Host", wyzeApiHost).
    SetBody(&payload).
    SetResult(&deviceListResponse).
    Post(wyzeGetDeviceListEndpoint)

  if err != nil {
    fmt.Println(err)
    return []WyzeDevice{}
  } else {
    return deviceListResponse.Data.DeviceList
  }
}

func saveFile(url string, fileName string) {
  response, err := http.Get(url)
  if err != nil {
    return
  }
  defer response.Body.Close()

  file, err := os.Create(fileName)
  if err != nil {
    return
  }
  defer file.Close()

  _, err = io.Copy(file, response.Body)
  if err != nil {
    return
  }
}

