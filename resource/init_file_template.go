package resource

var InitFileEnvContent = `#-------------------------------------
# project name
#-------------------------------------
PROJECT_NAME=dockser

#-------------------------------------
# project dir
#
# your code projects dir
#-------------------------------------
PROJECT_DIR=../projects

#-------------------------------------
# network name
#-------------------------------------
NETWORK=dockser
`

var InitFileEnvDemoContent = InitFileEnvContent

var InitFileEvnExampleContent = InitFileEnvContent

var InitFileGroupContent = `default:
  services:
    # - nginx
  template: "docker-compose.yml"
  output: "docker-compose.yml"
  override: "auto"`
var InitFileGroupDemoContent = InitFileGroupContent + `
demo:
  services:
    - nginx
  template: "docker-compose-demo.yml"
  output: "docker-compose-demo.yml"
  override: "auto"`

var InitFileTemplateContent = `version: "3"
networks:
  #@_NETWORK_@#:
      driver: bridge`

var InitFileTemplateDemoContent = InitFileTemplateContent

var InitFileServiceNginxContent = `#--------------------------------------------------------------------------
# nginx
#--------------------------------------------------------------------------
nginx:
  image: "nginx:alpine"
  container_name: "${PROJECT_NAME}_nginx"
  ports:
    - 80:80
    - 443:443
  tty: true
  networks:
    #@_NETWORK_@#:
`
