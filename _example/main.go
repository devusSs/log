package main

import (
	"context"
	"os"

	"github.com/devusSs/log"
)

func main() {
	l := log.NewLogger()
	// l.Debug() will not print if we do not set this
	// since the default level is "info".
	l.SetLevel(log.LevelDebug)
	// Defaults to os.Stderr.
	l.SetOut(os.Stdout)

	l.Debug("TEST", "KEY", "VALUE")
	l.Info("TEST", "key", 42, "value", "meaning")
	l.Warn("TEST", "key", 42, "value", "meaning")
	l.Error("TEST", "key", 42, "value", "meaning")
	// l.Fatal("TEST", "key", 42, "value", "meaning")

	sampleMap := make(map[string]interface{})
	sampleMap["mapKey"] = "mapValue"

	l.Info("TEST", "map", sampleMap)

	sampleStruct := struct {
		Name string
		Age  int
	}{Name: "Anton", Age: 45}

	l.Info("TEST", "struct", sampleStruct)

	// This will also work despite only having a value and not a key.
	// This will be printed as no_key=sampleStruct.
	l.Info("TEST", sampleStruct)

	// We can change the handler
	// and therefor change the output format.
	l.SetHandler(log.JSONHandler)

	l.Info("TEST", "TEST")
	l.Info("TEST", "KEY", "VALUE")

	type ctxName string
	ctx := context.WithValue(context.Background(), ctxName("name"), "anton")
	l.Info("TEST WITH CTX", "CTX", ctx)                              // will print Context = {}
	l.Info("TEST WITH CTX VALUE", "CTX", ctx.Value(ctxName("name"))) // will print CTX = anton
}
