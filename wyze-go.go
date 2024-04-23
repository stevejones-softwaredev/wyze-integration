package main

import (
  "github.com/go-resty/resty/v2"
  "time"
  "fmt"
  "io/fs"
  "io/ioutil"
  "log"
  "os"
  "strings"
  "strconv"
  "github.com/slack-go/slack"
  "github.com/stevejones-softwaredev/wyze-api"
)

func main() {
  client := resty.New()

  environment := validateNeededInputs()

  accessToken := wyze.GetWyzeAccessToken(client, environment)

  start, end := getTimeBounds(environment)

  filteredFiles := wyze.GetWyzeCamThumbnails(client,
    getOptionalVar("WYZE_HOME", "./", &environment),
    accessToken,
     10,
    parseCamList(environment["WYZE_FILTERED_CAM_LIST"]),
    parseFilterList(environment["WYZE_FILTER_VALUES"]),
    start,
    end)

  files := wyze.GetWyzeCamThumbnails(client,
    getOptionalVar("WYZE_HOME", "./", &environment),
    accessToken,
    10,
    parseCamList(environment["WYZE_CAM_LIST"]),
    []int {},
    start,
    end)
  files = append(files, filteredFiles...)

  deviceMap := getDeviceMacList(client, accessToken)

  botApi := slack.New(environment["SLACK_OAUTH_BOT_TOKEN"])
  userApi := slack.New(environment["SLACK_OAUTH_USER_TOKEN"])

  for _,file := range files {
    createSlackMessageWithFile(file, deviceMap, environment["SLACK_CHANNEL"], botApi, userApi)
  }
}

func parseCamList(camList string) []string {
  return strings.Split(camList, ":")
}

func parseFilterList(filterList string) []int {
  var filters []int
  names := strings.Split(filterList, ":")
  filterMap := wyze.GetFilterNameToValueMap()

  for _,name := range names {
    filterValue, filterOk := filterMap[name]

    if (filterOk) {
      filters = append(filters, filterValue)
    } else {
      log.Println("Unknown filter specified: ", name)
    }
  }
  
  return filters
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

  checkRequiredVar("WYZE_USERNAME", environment)
  checkRequiredVar("WYZE_PASSWORD_HASH", environment)
  checkRequiredVar("WYZE_KEY_ID", environment)
  checkRequiredVar("WYZE_API_KEY", environment)
  checkRequiredVar("SLACK_OAUTH_BOT_TOKEN", environment)
  checkRequiredVar("SLACK_OAUTH_USER_TOKEN", environment)
  checkRequiredVar("WYZE_CAM_LIST", environment)
  checkRequiredVar("SLACK_CHANNEL", environment)
  getOptionalVar("WYZE_LOOKBACK_SECONDS", "330", &environment)
  getOptionalVar("WYZE_FILTERED_CAM_LIST", "", &environment)
  getOptionalVar("WYZE_FILTER_VALUES", "", &environment)
  getOptionalVar("WYZE_HOME", "./", &environment)

  return environment
}

func getTimeBounds(environment map[string]string) (time.Time, time.Time) {
  dirName := getOptionalVar("WYZE_HOME", "./", &environment)
  files, _ := ioutil.ReadDir(dirName)
  end_time := time.Now()
  lookback_seconds,_ := strconv.Atoi(environment["WYZE_LOOKBACK_SECONDS"])
  lookback_seconds *= -1

  var matchingFiles []fs.FileInfo

  for _,file := range files {
    if (strings.Contains(file.Name(), ".jpg")) {
      matchingFiles = append(matchingFiles, file)
    }
  }

  defer deleteExpiredThumbnails(dirName, matchingFiles)

  if len(matchingFiles) == 0 {
    return end_time.Add(time.Second * time.Duration(lookback_seconds)), end_time
  } else {
    lastFileName := matchingFiles[len(matchingFiles)-1].Name()
    beginStampString := strings.Split(lastFileName,".")[0]
    i, err := strconv.ParseInt(beginStampString, 10, 64)
    if err != nil {
      log.Print(err)
      return end_time.Add(time.Second * time.Duration(lookback_seconds)), end_time
    }

    return time.UnixMilli(i + 1000), end_time
  }
}

func deleteExpiredThumbnails(dir string, files []fs.FileInfo) {
  for _,file := range files {
    diff := time.Since(file.ModTime())
    if diff.Hours() > 24 {
      deleteFile := fmt.Sprintf("%s%s", dir, file.Name())
      err := os.Remove(deleteFile)

      if (err != nil) {
        fmt.Println(err)
      }
    }
  }
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

