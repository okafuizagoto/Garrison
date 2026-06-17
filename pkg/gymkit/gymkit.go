// Package gymkit provides shared helpers for Gold Gym microservices.
// It is a public replacement for the private mataharibiz/sange library.
package gymkit

import "os"

const Version = "1"
const OAUTH = "Bearer "

var SYSTEM_TOKEN = os.Getenv("SYSTEM_TOKEN")
var TOKEN = OAUTH + SYSTEM_TOKEN
