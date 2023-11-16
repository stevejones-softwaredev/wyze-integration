package main

import (
  "github.com/go-resty/resty/v2"
  "time"
  "fmt"
  "log"
  "os"
  "strings"
  "strconv"
  "github.com/slack-go/slack"
  "sjones/wyze-go/wyze"
)

func main() {
  client := resty.New()

  environment := validateNeededInputs()

  wyzeToken := fmt.Sprintf("%s==", environment["WYZE_ACCESS_TOKEN"])

  accessToken := wyze.GetWyzeAccessToken(client, wyzeToken)

  lookback_seconds,_ := strconv.Atoi(environment["WYZE_LOOKBACK_SECONDS"])
  lookback_seconds *= -1

  end_time := time.Now()
  begin_time := end_time.Add(time.Second * time.Duration(lookback_seconds))

  files := wyze.GetWyzeCamThumbnails(client,
    getOptionalVar("WYZE_HOME", "./", &environment),
    accessToken,
    10,
    parseCamList(environment["WYZE_CAM_LIST"]),
    begin_time,
    end_time)
  deviceMap := getDeviceMacList(client, accessToken)

  for _,file := range files {
    msg := fmt.Sprintf("Recorded at %s from %s",time.UnixMilli(file.Timestamp).Format(time.RFC850), deviceMap[file.Mac].Nickname)
    fmt.Println(msg)
    api := slack.New(environment["SLACK_OAUTH_TOKEN"])
    fileInfo,_ := os.Stat(file.Path)
    uploadParams := slack.UploadFileV2Parameters{
      Channel: environment["SLACK_CHANNEL"],
      File: file.Path,
      Title: msg,
      InitialComment: msg,
      Filename: file.Path,
      FileSize: int(fileInfo.Size()),
    }
    file, err := api.UploadFileV2(uploadParams)

    if err != nil {
      fmt.Printf("%s\n", err)
      continue
    }
    fmt.Printf("ID: %s, Title: %s\n", file.ID, file.Title)
  }
}

func parseCamList(camList string) []string {
  return strings.Split(camList, ",")
}

func getDeviceMacList(client *resty.Client,
    accessToken string) map[string]wyze.WyzeDevice {
  deviceMap := make(map[string]wyze.WyzeDevice)

  devices := wyze.GetWyzeDeviceList(client, accessToken)

  for _,device := range devices {
    deviceMap[device.MAC] = device
  }

  return deviceMap
}

func validateNeededInputs() map[string]string{
  getenvironment := func(data []string, getkeyval func(item string) (key, val string)) map[string]string {
      items := make(map[string]string)
      for _, item := range data {
          key, val := getkeyval(item)
          items[key] = val
      }
      return items
    }

  environment := getenvironment(os.Environ(), func(item string) (key, val string) {
      splits := strings.Split(item, "=")
      key = splits[0]
      val = splits[1]
      return
  })

  checkRequiredVar("WYZE_ACCESS_TOKEN", environment)
  checkRequiredVar("SLACK_OAUTH_TOKEN", environment)
  checkRequiredVar("WYZE_CAM_LIST", environment)
  checkRequiredVar("SLACK_CHANNEL", environment)
  getOptionalVar("WYZE_LOOKBACK_SECONDS", "330", &environment)
  getOptionalVar("WYZE_HOME", "./", &environment)

  return environment
}

func checkRequiredVar(name string, env map[string]string) string {
  value, ok := env[name]
  if !ok {
    log.Fatalf("Missing required environment variable %s; exiting", name)
  }

  return value
}

func getOptionalVar(name string, defaultValue string, env *map[string]string)  string {
  value, ok := (*env)[name]
  if !ok {
    log.Printf("Missing optional environment variable %s; using default %s", name, defaultValue)
    return defaultValue
  } else {
    return value
  }
}

