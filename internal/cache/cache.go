package cache

/*
this file will manage the caching
using ristretto cacher

imports required:
dgraph-io/ristretto
the utils package

functions:
func initCache() *ristretto.Cache
func get(key string) (interface{}, bool) get from cache
func set(key string, value interface{}) set in cache

here interface{} means that it can hold any type of data
*/
