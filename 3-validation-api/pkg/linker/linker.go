package linker

import (
	"fmt"
)

func GetHashUrl(serverPath string, hashString string) string {
	return fmt.Sprintf("%s/{%s}", serverPath, hashString)
}
