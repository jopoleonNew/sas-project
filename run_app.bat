@echo on
title golang SAS-project runner
echo SAS-project runner
echo Start MongoDB
title MongoDB
start cmd /k mongod --config F:\MongoDB\bin\mongodb.config"
title ngrok
start cmd /k ngrok http -subdomain=sas-project -region eu 8080
echo Waiting ~3 seconds
timeout 3 > NUL
echo Start Golang code
go run main.go
