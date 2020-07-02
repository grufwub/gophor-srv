package core

var (
	// ReadFromClient is the globally set function to read data from a client
	ReadFromClient func(*Client) ([]byte, Error)
)
