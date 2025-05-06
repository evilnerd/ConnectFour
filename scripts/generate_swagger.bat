@echo off
setlocal

REM Check if swag is installed
where swag >nul 2>&1
if %ERRORLEVEL% neq 0 (
    echo Swag not found, installing...
    go install github.com/swaggo/swag/cmd/swag@latest
    go get -u github.com/swaggo/swag/cmd/swag
    go get -u github.com/swaggo/http-swagger
)

REM Generate Swagger documentation
echo Generating Swagger documentation...
swag init -g ./cmd/server/main.go -o ./docs

echo Swagger documentation generated successfully!
echo You can access the Swagger UI at http://localhost:8443/swagger/index.html when the server is running.

endlocal

