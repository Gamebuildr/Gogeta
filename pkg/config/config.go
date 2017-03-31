package config

// Region is the world region messages are coming from
const Region string = "REGION"

// QueueURL is the URL endpoint to recieve messages
const QueueURL string = "QUEUE_URL"

// MrrobotNotifications is the URL endpoint to send MrRobot messages
const MrrobotNotifications string = "MRROBOT_NOTIFICATIONS"

// GamebuildrNotifications is the URL endpoint to send Gamebuildr messages
const GamebuildrNotifications string = "GAMEBUILDR_NOTIFICATIONS"

// GcloudServiceKey is the base64 generated key from the gcloud .json
const GcloudServiceKey string = "GCLOUD_SERVICE_KEY"

// GcloudServiceAccount is the path to the gcloud service account .json
// The .json is automatically generated and set when the client is ran
const GcloudServiceAccount string = "GOOGLE_APPLICATION_CREDENTIALS"

// CodeRepoStorage is the location to save source code
const CodeRepoStorage string = "CODE_REPO_STORAGE"

// GoEnv is the environment the current system is operating in
const GoEnv string = "GO_ENV"

// LogEndpoint is the endpoint for sending log data
const LogEndpoint string = "PAPERTRAIL_ENDPOINT"
