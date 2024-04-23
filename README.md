# Wyze Integration

This project is a bespoke solution in Go intended to facilitate a fairly simple flow of capturing thumbnails from Wyze camera events and posting them to a slack channel. Its focus is extremely narrow to this purpose. It was done as a learning project and has a number of quirks that are not Go best practices.

## How to build the app
The project can be built as a standard Go application using Go v19 or later, or can be built using the included `Dockerfile_build`. The file is coded to default to building for linux with amd64 architecture. However, build arguments for the docker file are available to specify different operating systems and architectures. For instance, to build for an M-series Mac, you would use:
  `docker build -o bin --no-cache --build-arg "GOOS=darwin" --build-arg "GOARCH=arm64" -f Dockerfile_build .`
  
  A list of the valid GOOS and GOARCH combinations can be found in the Variables section of [the documentation for the Go internal/platform package](https://pkg.go.dev/internal/platform).
  
  In the above example, the executable would be written to the bin subfolder; however any desired subfolder can be passed as the `-o` argument to the `docker build` command.

## How to use the app
The app receives all sensitive credentials and specific configuration from environment variables. Some are required and some are optional. Here is the list of them:
  * WYZE_USERNAME **Required** - this is the email address associated with the account that will be used to access WYZE API data

  * WYZE_PASSWORD_HASH **Required** - this is the triple MD5-hashed value of the password associated with the username above (that is, *MD5(MD5(MD5(password)))*)

  * WYZE_KEY_ID **Required** - this is a UUID associated with the API key below [this page from the Wyze support pages](https://support.wyze.com/hc/en-us/articles/16129834216731-Creating-an-API-Key). **Wyze does have a restriction that refresh tokens are only valid for 30 days, so it will be necessary to rotate this credential.**

  * WYZE_API_KEY **Required** - this is an API key used to access Wyze APIs. One can be generated as documented at [this page from the Wyze support pages](https://support.wyze.com/hc/en-us/articles/16129834216731-Creating-an-API-Key). **Wyze does have a restriction that API keys are valid for one year only, so it will be necessary to rotate this credential.**

  * SLACK_OAUTH_BOT_TOKEN **Required** -

  * SLACK_OAUTH_USER_TOKEN **Required** -

  * WYZE_CAM_LIST **Required** - Colon (:) separated list of Wyze device MACs of cameras for whom we should capture all events

  * SLACK_CHANNEL **Required** - Slack channel to which images will be posted

  * WYZE_FILTERED_CAM_LIST *Optional (Default [])* - Colon (:) separated list of Wyze device MACs of cameras for whom we should capture events matching the filters in WYZE_FILTER_VALUES

  * WYZE_FILTER_VALUES *Optional (Default [])* -  Colon (:) separated list of event types we should capture from cameras listed in WYZE_FILTERED_CAM_LIST. Allowed values are a combination of *Pet*, *Person*, *Vehicle* and *Package*

  * WYZE_LOOKBACK_SECONDS *Optional (Default 330)* -

  * WYZE_HOME *Optional (Default "./")* - the folder where thumbnails will be downloaded and Wyze API refresh token will be stored. Images will be purged on next run after they are 24 hours old.

## How to create a container of the app
A container image can be built after doing a local build using the included `Dockerfile`. The file expects the wyxe-go executable to be in the same folder as the `Dockerfile`, so if the Docker-based build method notied above was used, you may need to copy the executable to the folder containing `Dockerfile`.

## How to deploy a kubernetes scheduled job of the app
helm install sj-test --set volume.wyzeDownloads=$WYZE_HOME --set wyze.camList=$WYZE_CAM_LIST --set wyze.filteredCamList=$WYZE_FILTERED_CAM_LIST --set wyze.filterValues=$WYZE_FILTER_VALUES --set wyze.lookbackInSeconds=$WYZE_LOOKBACK_SECONDS --set wyze.userName=$WYZE_USERNAME --set wyze.passwordHash=$WYZE_PASSWORD_HASH --set wyze.apiKeyId=$WYZE_KEY_ID --set wyze.apiKey=$WYZE_API_KEY --set slack.botToken=$SLACK_OAUTH_BOT_TOKEN --set slack.userToken=$SLACK_OAUTH_USER_TOKEN --set slack.channel=$SLACK_CHANNEL .
