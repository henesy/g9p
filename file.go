package g9p



// Define the types of interactions that can be had with a file
type File interface {
	Walk() (error)
	Open() (error)
	Read() (error)
	Write() (error)
	Stat() (error)
	Wstat() (error)
}


