package vmio

type VMKubeElementsStream interface {
	Read()		error
	Write() 	error
	Export(prettify bool) 	([]byte, error)
	Import(file string, format string) 	error
}

