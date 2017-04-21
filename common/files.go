package common

var GetFile func(string) ([]byte, error)

var WriteFile func(string, []byte, int) error

var Mkdir func(string, int) error

var Exists func(string) (bool, error)
