package main

import (
  "fmt"
  "time"
  "github.com/slack-go/slack"
  "github.com/stevejones-softwaredev/wyze-api"
)

func createSlackMessageWithFile(file wyze.WyzeDownloadedFile, deviceMap map[string]wyze.WyzeDevice, channelId string, botApi *slack.Client, userApi *slack.Client) {
  catNameSectionBlock, catActivitySectionBlock := createConstantSectionBlocks()
  
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

  _,_,_,msgErr := botApi.SendMessage(channelId, slack.MsgOptionBlocks(textSectionBlock, catNameSectionBlock, catActivitySectionBlock))

  if msgErr != nil {
    fmt.Println(msgErr)
  }
}

func createConstantSectionBlocks() (*slack.SectionBlock, *slack.SectionBlock) {
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

