package main

import (
	"github.com/powerman/structlog"
)

func init()  {
	structlog.DefaultLogger.
		SetPrefixKeys(
			//structlog.KeyApp,
			structlog.KeyPID, structlog.KeyLevel, structlog.KeyUnit, structlog.KeyTime,
		).
		SetDefaultKeyvals(
			//structlog.KeyApp, filepath.Base(os.Args[0]),
			structlog.KeySource, structlog.Auto,
			structlog.KeyStack, structlog.Auto,
		).
		SetSuffixKeys(
			//structlog.KeyStack,
			//structlog.KeySource,
		).
		SetKeysFormat(map[string]string{
			structlog.KeyTime:   " %[2]s",
			structlog.KeySource: " %6[2]s",
			structlog.KeyUnit:   " %6[2]s",
		})
}

var log = structlog.New()