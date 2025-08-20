package utils


/*
this file will manage the config for
everything

imports:
os
encoding/json OR "github.com/spf13/viper"

functions:
LoadConfig() Config
type config struct {DBPath string, CacheSize int, Port int, MODELURL string}

this will store all the logs of where the req came from, current DB and model
paths and request management
*/