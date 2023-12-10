package utils

import "flag"

var (
	ConfigFilePath string
)

func InitFlag() {
	flag.StringVar(&ConfigFilePath, "c", ".", "set config file path")

	flag.Parse()
}
