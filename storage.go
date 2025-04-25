package main

type storage interface {
	validateClientId(clientId string) error
}
