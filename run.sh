# export BASE_GUID=$(uuidgen)
export BASE_GUID=2A0A5573-5AB7-470A-9010-958150FA6710

export CREDENTIALS='{"port": "4000", "host": "1.2.3.4"}'
export SERVICE_NAME=myservice
export SERVICE_PLAN_NAME=shared
export TAGS=simple,shared
export AUTH_USER=broker
export AUTH_PASSWORD=broker
go run cmd/worlds-simplest-service-broker/main.go