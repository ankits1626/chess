# Chess Coach API - Quick Start

## ✅ Setup Complete!

Your Chess Coach API now has Swagger documentation enabled!

---

## 🚀 Quick Commands

### Generate Swagger Docs
```bash
cd backend
swag init -g cmd/server/main.go
```

### Run the Server
```bash
# With hot reload (recommended)
air

# Or without hot reload
go run cmd/server/main.go
```

### Access Swagger UI
Open in browser:
```
http://localhost:8080/swagger/index.html
```

---

## 📝 Current API Endpoints

### 1. Root Endpoint
- **URL:** `GET /`
- **Description:** Welcome message
- **Response:** `"Chess Coach API - Server Running!"`

### 2. Health Check
- **URL:** `GET /api/v1/health`
- **Description:** API health status
- **Response:**
  ```json
  {
    "status": "healthy"
  }
  ```

### 3. Swagger Documentation
- **URL:** `GET /swagger/index.html`
- **Description:** Interactive API documentation

---

## 🎯 Next Steps

### 1. Start the Server
```bash
cd backend
air
```

### 2. Test the Endpoints

**Using curl:**
```bash
# Test root endpoint
curl http://localhost:8080/

# Test health check
curl http://localhost:8080/api/v1/health
```

**Using browser:**
- Visit: http://localhost:8080/swagger/index.html
- Try out the endpoints interactively!

### 3. Add More Endpoints

When you add new endpoints:

1. **Add Swagger annotations:**
```go
// @Summary Your endpoint summary
// @Description Detailed description
// @Tags category
// @Accept json
// @Produce json
// @Success 200 {object} YourResponse
// @Router /api/v1/your-path [get]
func yourHandler(c *gin.Context) {
    // Your code
}
```

2. **Regenerate docs:**
```bash
swag init -g cmd/server/main.go
```

3. **Refresh Swagger UI** in browser

---

## 📖 Documentation Files

- **[swagger-setup.md](./swagger-setup.md)** - Complete Swagger setup guide
- **[setup-guide.md](./setup-guide.md)** - Go backend setup guide
- **docs/swagger.json** - OpenAPI JSON specification
- **docs/swagger.yaml** - OpenAPI YAML specification
- **docs/docs.go** - Generated Go documentation code

---

## 🔥 Pro Tips

### Auto-regenerate on File Changes

The error you'll see initially is normal - the `docs` package won't exist until you run `swag init`.

**Workflow:**
1. Write your handler with annotations
2. Run: `swag init -g cmd/server/main.go`
3. Start server: `air`
4. Test in Swagger UI

### Testing API in Swagger UI

1. Visit http://localhost:8080/swagger/index.html
2. Click on an endpoint to expand it
3. Click "Try it out"
4. Fill in parameters (if any)
5. Click "Execute"
6. See the response!

---

## 🐛 Troubleshooting

### Error: "could not import docs package"

This is expected before first `swag init`. Just run:
```bash
swag init -g cmd/server/main.go
```

### Swagger UI not showing new endpoints

1. Regenerate docs: `swag init -g cmd/server/main.go`
2. Restart server
3. Hard refresh browser (Cmd+Shift+R on Mac)

---

## 📚 Learn More

- Full Swagger guide: [swagger-setup.md](./swagger-setup.md)
- Backend setup: [setup-guide.md](./setup-guide.md)
- Swaggo Docs: https://github.com/swaggo/swag

---

**Happy Coding! 🎉**
