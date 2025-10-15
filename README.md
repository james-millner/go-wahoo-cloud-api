# Go - L&D Wahoo Cloud API Project

This app serves as a backend layer to intercept Wahoo Cloud API requests and responses.

## Endpoints

- **Health** (GET): `/healthz` - Simple health endpoint.
- **Authorize** (GET): `/authorize` - Kicks off the OAuth 2.0 flow with Wahoo.
- **Root** (GET): `/` - Handles the Wahoo access token request.
- **Callback** (POST): `/callback` - Exposes an interface for Wahoo to call when a ride is uploaded. The request will contain a [workout summary](https://cloud-api.wahooligan.com/#workout-summary).

> **Warning**: Beginner Gopher here.

## Configuration

The app is primarily configured by environment variables and is fairly simple to set up and run yourself.

The following environment variables are required:

```
PORT = "8080"
REDIRECT_URI = "MY_REDIRECT_URI"
WAHOO_CLIENT_ID = "MY_WAHOO_CLIENT_ID"
WAHOO_CLIENT_SECRET = "MY_WAHOO_CLIENT_SECRET"
WAHOO_AUTH_BASE_URL = "https://api.wahooligan.com/oauth/authorize"
WAHOO_TOKEN_BASE_URL = "https://api.wahooligan.com/oauth/token"
TIGRIS_ENABLED = "true" // Optional, and defaults to false
FITFILE_SERVICE_URL = "https://fit-file-backend-billowing-cloud-731.fly.dev/api/v1/fitfiles" // Optional, if set will POST FIT files to this service
```

## Deployment

This project is deployed using [Fly.io](https://fly.io). Enjoyed using Fly to be honest, its been quite user friendly to setup and run, and has cost my nothig so far! Added bonus!

### Fly.io Features Used

1. **Fly Storage**: Used for storing FIT files. This provides persistent storage for workout data files.

2. **Fly Secrets**: Utilized for securely storing sensitive information, particularly the Wahoo Client Secret. This ensures that confidential credentials are not exposed in the codebase or environment variables.

### Deployment Process

To deploy this application on Fly.io:

1. Ensure you have the Fly CLI installed and are logged in.
2. Set up your `fly.toml` file with the necessary configuration.
3. Use Fly Secrets to securely store your Wahoo Client Secret:
   ```
   fly secrets set WAHOO_CLIENT_SECRET=your_secret_here
   ```
4. Deploy the application using:
   ```
   fly deploy
   ```
5. (Optional) If you wish to store the FIT files with Storage, see [Fly.io Global Object Storage Docs](https://fly.io/docs/tigris/) this requires some additional configuration.

For more detailed information on deploying Go applications with Fly.io, refer to the [Fly.io documentation](https://fly.io/docs/languages-and-frameworks/golang/).

## Contributing

As a beginner Gopher, contributions, suggestions, and feedback are welcome to improve this.