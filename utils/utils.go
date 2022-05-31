package utils

import (
	"fmt"
	"os"
)

func HandleError(err error) {
	if err != nil {
		fmt.Errorf("%s", err)
		os.Exit(1)
	}
}
