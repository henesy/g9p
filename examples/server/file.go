package main

/* File constructs */

// 9p file
type File9 struct {
	FileName string
}

// Walk
func (f *File9) Walk() (error){
	// TODO
	return nil
}

// Open
func (f *File9) Open() (error){
	// TODO
	return nil
}

// Read
func (f *File9) Read() (error){
	// TODO
	return nil
}

// Write
func (f *File9) Write() (error){
	// TODO
	return nil
}

// Stat
func (f *File9) Stat() (error){
	// TODO
	return nil
}

// Wstat
func (f *File9) Wstat() (error){
	// TODO
	return nil
}


// Name returns the name of the file
func (f *File9) Name() string {
	return f.FileName
}

// 9p directory
type Dir9 struct {
	FileName string
}

// Walk
func (f *Dir9) Walk() (error){
	// TODO
	return nil
}

// Open
func (f *Dir9) Open() (error){
	// TODO
	return nil
}

// Read
func (f *Dir9) Read() (error){
	// TODO
	return nil
}

// Write
func (f *Dir9) Write() (error){
	// TODO
	return nil
}

// Stat
func (f *Dir9) Stat() (error){
	// TODO
	return nil
}

// Wstat
func (f *Dir9) Wstat() (error){
	// TODO
	return nil
}

// Name returns the name of the file
func (f *Dir9) Name() string {
	return f.FileName
}
