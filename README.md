# slack-module-svc

## Description

Module for access control for slack.

## Slack Access Tokens generation rules

To run Slack ACS Module, a pair of tokens is needed - *User OAuth Token* and *User Bot OAuth token*. Follow (this guide from Slack)[https://api.slack.com/tutorials/tracks/getting-a-token] to generate them for your Slack workspace.

Keep in mind that each token has its own scopes that you have to specify once you created a Slack app. See the list of required scopes below.

### User OAuth token scopes

Always starts with `xoxp`

Required scopes: 
- admin

### User Bot OAuth token scopes

Always starts with `xoxb`

Required scopes: 
- channels:read
- groups:read
- team:read
- usergroups:read
- users.profile:read
- users:read

## Install

  ```
  git clone github.com/acs-dl/slack-module-svc
  cd slack-module-svc
  go build main.go
  export KV_VIPER_FILE=./config.yaml
  export USER_TOKEN=<user_token>
  export BOT_TOKEN=<user_bot_token>
  ./main migrate up
  ./main run service
  ```

## Documentation

We do use openapi:json standard for API. We use swagger for documenting our API.

To open online documentation, go to [swagger editor](http://localhost:8080/swagger-editor/) here is how you can start it
```
  cd docs
  npm install
  npm start
```
To build documentation use `npm run build` command,
that will create open-api documentation in `web_deploy` folder.

To generate resources for Go models run `./generate.sh` script in root folder.
use `./generate.sh --help` to see all available options.

Note: if you are using Gitlab for building project `docs/spec/paths` folder must not be
empty, otherwise only `Build and Publish` job will be passed.  

## Running from docker 
  
Make sure that docker installed.

use `docker run ` with `-p 8080:80` to expose port 80 to 8080

  ```
  docker build -t github.com/acs-dl/slack-module-svc .
  docker run -e KV_VIPER_FILE=/config.yaml github.com/acs-dl/slack-module-svc
  ```

## Running from Source

* Set up environment value with config file path `KV_VIPER_FILE=./config.yaml`
* Provide valid config file
* Launch the service with `migrate up` command to create database schema
* Launch the service with `run service` command


### Database
For services, we do use ***PostgresSQL*** database. 
You can [install it locally](https://www.postgresql.org/download/) or use [docker image](https://hub.docker.com/_/postgres/).


### Third-party services


## Contact

Responsible Katheryn Hryhorieva
The primary contact for this project is  t.me/wiyenya
