package credentials

import "os"

// ServiceAccount is the credentials to use with outside services
type ServiceAccount interface {
	GenerateAccount()
}

// JSONCredentials are credentials stored in a .json file
type JSONCredentials interface {
	createJSON(path string) (*os.File, error)
	writeToJSONFile(file *os.File, key []byte) error
}
