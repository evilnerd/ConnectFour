# ConnectFour API Swagger Documentation

This README explains how to use the Swagger documentation that has been integrated into the ConnectFour API.

## Overview

The ConnectFour API now includes automatically generated Swagger documentation using [Swag](https://github.com/swaggo/swag). This documentation:

- Provides an interactive UI to explore all API endpoints
- Shows request and response schemas
- Allows testing API endpoints directly from the browser
- Automatically updates when code annotations change

## Accessing the Documentation

Once the server is running, you can access the Swagger documentation at:

```
http://localhost:8443/swagger/index.html
```

This will open the Swagger UI interface where you can browse and interact with all available API endpoints.

## Authentication

The API uses JWT authentication. To use protected endpoints in the Swagger UI:

1. First, use the `/login` endpoint to obtain a JWT token
2. Click the "Authorize" button at the top of the Swagger UI
3. Enter your JWT token in the format: `Bearer {your-token-here}`
4. Click "Authorize" to apply the token to all subsequent requests

## Updating the Documentation

The Swagger documentation is generated from annotations in the Go code. If you make changes to the API:

1. Update the annotations in the relevant handler or model files
2. Run the generation script:
   ```
   # For Linux/Mac
   ./scripts/generate_swagger.sh
   
   # For Windows
   scripts\generate_swagger.bat
   ```
3. Restart the server to see the updated documentation

## Annotation Examples

Here are some examples of how the API is annotated:

### For Handlers
```go
// LoginHandler godoc
// @Summary User login
// @Description Authenticates a user and returns a JWT token
// @Tags authentication
// @Accept json
// @Produce plain
// @Param login body service.LoginRequest true "Login credentials"
// @Success 200 {string} string "JWT Token"
// @Failure 400 {object} service.ErrorResponse "Invalid request format"
// @Failure 401 {string} string "Invalid credentials"
// @Failure 500 {object} service.ErrorResponse "Internal server error"
// @Router /login [post]
func LoginHandler(w http.ResponseWriter, r *http.Request) {
    // Implementation...
}
```

### For Models
```go
// LoginRequest represents the user login credentials
// swagger:model
type LoginRequest struct {
    // User's email address
    // Required: true
    // example: player@example.com
    Email string `json:"email"`
    
    // User's password
    // Required: true
    // example: password123
    Password string `json:"password"`
}
```

## Further Reading

- [Swag Documentation](https://github.com/swaggo/swag)
- [OpenAPI Specification](https://swagger.io/specification/)
- [Swagger UI](https://swagger.io/tools/swagger-ui/)
