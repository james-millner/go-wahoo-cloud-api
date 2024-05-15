GoLang - L&D Wahoo Cloud API Project
----------------------------

This app serves as a backend layer to intercept Wahoo Cloud API requests and responses.

Endpoints: 

- Health (GET); `/healthz` - Simple health endpoint.
- Authorize (GET); `/authorize` - Kicks off the OAuth 2.0 flow with Wahoo.
- Root (GET); `/` - Handles the wahoo access token request.
- Callback (POST); `/callback` - Exposes an interface for wahoo to call when a ride is uploaded. The request will contain a [workout summary.](https://cloud-api.wahooligan.com/#workout-summary) 

Warning Beginner Gopher here.

The app is primarily configured by environment variables and is fairly simple to setup and run yourself.

The following environment variables are required:

PORT = "8080"
REDIRECT_URI = "MY_REDIRECT_URI"
WAHOO_CLIENT_ID = "MY_WAHOO_CLIENT_ID"
WAHOO_CLIENT_SECRET = "MY_WAHOO_CLIENT_SECRET"
WAHOO_AUTH_BASE_URL = "https://api.wahooligan.com/oauth/authorize"
WAHOO_TOKEN_BASE_URL = "https://api.wahooligan.com/oauth/token"
TIGRIS_ENABLED = "true" // Optional, and defaults to false
```