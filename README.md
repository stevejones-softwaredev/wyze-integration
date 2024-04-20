#Wyze Integration

This project is a bespoke solution in Go intended to facilitate a fairly simple flow of capturing thumbnails from Wyze camera events and posting them to a slack channel. Its focus is extremely narrow to this purpose. It was done as a learning project and has a number of quirks that are not Go best practices.

##How to build the app
The project can be built as a standard Go application using Go v19 or later, or can be built using the included `Dockerfile_build`. The file is coded to default to building for linux with amd64 architecture. However, build arguments for the docker file are available to specify different operating systems and architectures. For instance, to build for an M-series Mac, you would use:
  `docker build -o bin --no-cache --build-arg "GOOS=darwin" --build-arg "GOARCH=arm64" -f Dockerfile_build .`
  
  A list of the valid GOOS and GOARCH combinations can be found in the Variables section of [the documentation for the Go internal/platform package](https://pkg.go.dev/internal/platform).
  
  In the above example, the executable would be written to the bin subfolder; however any desired subfolder can be passed as the `-o` argument to the `docker build` command.

##How to use the app
The app receives all sensitive credentials and specific configuration from environment variables. Some are required and some are optional. Here is the list of them:
  * WYZE_ACCESS_TOKEN **Required** - this is a refresh API token used to access Wyze APIs. One can be generated as documented at [this page from the Wyze support pages](https://support.wyze.com/hc/en-us/articles/16129834216731-Creating-an-API-Key). **Wyze does have a restriction that refresh tokens are only valid for 30 days, so it will be necessary to rotate this credential.**

  * SLACK_OAUTH_TOKEN **Required** - 

  * WYZE_CAM_LIST **Required** - 

  * SLACK_CHANNEL **Required** - 

  * WYZE_LOOKBACK_SECONDS *Optional (Default 330)* -

  * WYZE_HOME *Optional (Default "./")* - the folder where thumbnails will be downloaded 

##How to create a container of the app
A container image can be built after doing a local build using the included `Dockerfile`. The file expects the wyxe-go executable to be in the same folder as the `Dockerfile`, so if the Docker-based build method notied above was used, you may need to copy the executable to the folder containing `Dockerfile`.


helm install sj-test --set volume.wyzeDownloads=$WYZE_HOME --set wyze.camList=$WYZE_CAM_LIST --set wyze.lookbackInSeconds=$WYZE_LOOKBACK_SECONDS=330 --set wyze.refreshToken=$WYZE_ACCESS_TOKEN --set slack.token=$SLACK_OAUTH_TOKEN --set slack.channel=$SLACK_CHANNEL .
