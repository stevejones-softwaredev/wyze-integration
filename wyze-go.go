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
  "stevejones.softwaredev/wyze-go/wyze"
)

func main() {
  client := resty.New()

  environment := validateNeededInputs()
  
  wyzeToken := environment["WYZE_REFRESH_TOKEN"]

  accessToken := wyze.GetWyzeAccessToken(client, wyzeToken)

  start, end := getTimeBounds(environment)

  filteredFiles := wyze.GetWyzeCamThumbnails(client,
    getOptionalVar("WYZE_HOME", "./", &environment),
    accessToken,
    10,
    parseCamList(environment["WYZE_FILTERED_CAM_LIST"]),
    []int {102, 103}, // restrict to events tagged with "pet" of "vehicle" (for Syd's dumptruck)
    start,
    end)
  files := wyze.GetWyzeCamThumbnails(client,
    getOptionalVar("WYZE_HOME", "./", &environment),
    accessToken,
    10,
    parseCamList(environment["WYZE_CAM_LIST"]),
    []int {}, // restrict to events tagged with "pet" of "vehicle" (for Syd's dumptruck)
    start,
    end)
  files = append(files, filteredFiles...)

  deviceMap := getDeviceMacList(client, accessToken)

  botApi := slack.New(environment["SLACK_OAUTH_BOT_TOKEN"])
  userApi := slack.New(environment["SLACK_OAUTH_USER_TOKEN"])

  catNameSectionBlock, catActivitySectionBlock := createConstantSectionBlocks()
  
  for _,file := range files {
    msg := fmt.Sprintf("Recorded at %s from %s",time.UnixMilli(file.Timestamp).Format(time.RFC850), deviceMap[file.Mac].Nickname)
    uploadParams := slack.FileUploadParameters{
      File: file.Path,
      Title: msg,
      InitialComment: msg,
      Filename: file.Path,
    }
    uploadedFile, err := userApi.UploadFile(uploadParams)
    if err != nil {
      fmt.Println(err)
    }

    publicFile,_,_,_ := userApi.ShareFilePublicURL(uploadedFile.ID)

    fmt.Printf("ID: %s, Title: %s\n", uploadedFile.ID, uploadedFile.Title)

    selectHeader := fmt.Sprintf("%s\n%s\n%s", time.UnixMilli(file.Timestamp).Format(time.RFC850), deviceMap[file.Mac].Nickname, publicFile.PermalinkPublic)
    textBlock := slack.NewTextBlockObject("mrkdwn", selectHeader, false, false)

    textSectionBlock := slack.NewSectionBlock(textBlock, nil, nil)

    _,_,_,msgErr := botApi.SendMessage(environment["SLACK_CHANNEL"], slack.MsgOptionBlocks(textSectionBlock, catNameSectionBlock, catActivitySectionBlock))

    if msgErr != nil {
      fmt.Println(msgErr)
    }
  }
}

func createConstantSectionBlocks() (*slack.SectionBlock,*slack.SectionBlock) {
  catNameTextBlock := slack.NewTextBlockObject("plain_text", "Cat Name:", false, false)
  catNameSydneyText := slack.NewTextBlockObject("plain_text", "Sydney", false, false)
  catNameSaviText := slack.NewTextBlockObject("plain_text", "Savi", false, false)
  catNameNoCatText := slack.NewTextBlockObject("plain_text", "Not a Cat", false, false)
  catNameSydneyOption := slack.NewOptionBlockObject("Sydney", catNameSydneyText, nil)
  catNameSaviOption := slack.NewOptionBlockObject("Savi", catNameSaviText, nil)
  catNameNoCatOption := slack.NewOptionBlockObject("NotACat", catNameNoCatText, nil)
  catNameOptionsBlock := slack.NewOptionsSelectBlockElement("static_select", catNameTextBlock, "cat_name", catNameSaviOption, catNameSydneyOption, catNameNoCatOption)
  catNameAccessory := slack.NewAccessory(catNameOptionsBlock)
  catNameSectionBlock := slack.NewSectionBlock(catNameTextBlock, nil, catNameAccessory)
  catNameSectionBlock.BlockID = "cat_name"

  catActivityTextBlock := slack.NewTextBlockObject("plain_text", "Cat Activity:", false, false)
  catActivityPeeText := slack.NewTextBlockObject("plain_text", "Pee", false, false)
  catActivityPoopText := slack.NewTextBlockObject("plain_text", "Poop", false, false)
  catActivityNoneText := slack.NewTextBlockObject("plain_text", "Neither", false, false)
  catActivityPeeOption := slack.NewOptionBlockObject("Pee", catActivityPeeText, nil)
  catActivityPoopOption := slack.NewOptionBlockObject("Poop", catActivityPoopText, nil)
  catActivityNoneOption := slack.NewOptionBlockObject("Neither", catActivityNoneText, nil)
  catActivityOptionsBlock := slack.NewOptionsSelectBlockElement("static_select", catActivityTextBlock, "cat_activity", catActivityPeeOption, catActivityPoopOption, catActivityNoneOption)
  catActivityAccessory := slack.NewAccessory(catActivityOptionsBlock)
  catActivitySectionBlock := slack.NewSectionBlock(catActivityTextBlock, nil, catActivityAccessory)
  catActivitySectionBlock.BlockID = "cat_activity"

  return catNameSectionBlock, catActivitySectionBlock
}

func parseCamList(camList string) []string {
  return strings.Split(camList, ":")
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

  fmt.Println(environment["TEST_VAR"])
  checkRequiredVar("WYZE_REFRESH_TOKEN", environment)
  checkRequiredVar("SLACK_OAUTH_BOT_TOKEN", environment)
  checkRequiredVar("SLACK_OAUTH_USER_TOKEN", environment)
  checkRequiredVar("WYZE_CAM_LIST", environment)
  checkRequiredVar("SLACK_CHANNEL", environment)
  getOptionalVar("WYZE_LOOKBACK_SECONDS", "330", &environment)
  getOptionalVar("WYZE_HOME", "./", &environment)

  return environment
}

func getTimeBounds(environment map[string]string) (time.Time, time.Time) {
  dirName := getOptionalVar("WYZE_HOME", "./", &environment)
  files, _ := ioutil.ReadDir(dirName)
  end_time := time.Now()
  lookback_seconds,_ := strconv.Atoi(environment["WYZE_LOOKBACK_SECONDS"])
  lookback_seconds *= -1

  defer deleteExpiredThumbnails(dirName, files)

  if len(files) == 0 {
    return end_time.Add(time.Second * time.Duration(lookback_seconds)), end_time
  } else {
    lastFileName := files[len(files)-1].Name()
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

