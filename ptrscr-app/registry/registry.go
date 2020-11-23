package registry

import (
	log "github.com/sirupsen/logrus"
	"golang.org/x/sys/windows/registry"
)

// GetStringFromLocalMachine returns string for given path and key
func GetStringFromLocalMachine(path string, key string) string {
	s := ""
	k, err := registry.OpenKey(registry.CURRENT_USER, path, registry.QUERY_VALUE)
	if err != nil {
		log.Println(err)
	}
	defer k.Close()
	s, _, err = k.GetStringValue(key)
	if err != nil {
		log.Println(err)
	}
	return s
}

func CreateRegistryKey(appPath string, tokenName string, tokenValue string) {

	softwareK, err := registry.OpenKey(registry.CURRENT_USER, "Software", registry.QUERY_VALUE)
	if err != nil {
		log.Fatal(err)
	}
	defer softwareK.Close()

	k, _, err := registry.CreateKey(softwareK, appPath, registry.CREATE_SUB_KEY|registry.QUERY_VALUE|registry.SET_VALUE)
	if err != nil {
		log.Fatal(err)
	}
	defer k.Close()

	if err := k.SetStringValue(tokenName, tokenValue); err != nil {
		log.Fatal(err)
	}
	if err := k.Close(); err != nil {
		log.Fatal(err)
	}
}
