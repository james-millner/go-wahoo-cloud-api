GoLang - L&D Wahoo Cloud API Project
----------------------------

This app serves as a backend layer to intercept Wahoo Cloud API requests and responses.

Endpoints: 

- Health (GET); `/healthz` - Simple health endpoint.
- Authorize (GET); `/authorize` - Kicks off the OAuth 2.0 flow with Wahoo.
- Root (GET); `/` - Handles the wahoo access token request.
- Callback (POST); `/callback` - Exposes an interface for wahoo to call when a ride is uploaded. The request will contain a [workout summary.](https://cloud-api.wahooligan.com/#workout-summary) 

Warning Beginner Gopher here.