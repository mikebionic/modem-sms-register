GOOS=windows GOARCH=amd64 go build -o bin/modem-sms-register-amd64.exe .
GOOS=windows GOARCH=386 go build -o bin/modem-sms-register-386.exe .

GOOS=windows GOARCH=amd64 go build -ldflags -H=windowsgui -o bin/modem-sms-register-amd64.exe .
GOOS=windows GOARCH=386 go build -ldflags -H=windowsgui -o bin/modem-sms-register-386.exe