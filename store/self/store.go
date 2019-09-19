package store

// NOTE: And what if we use a preexisting database with a file format. We load
// this thing into memory, work and work, and periodically patch the binary with
// updates (incase of crash).

// NOTE: Having a variable magic sequence allows it to be stealth if desireable
type Store struct {
	MagicSequence []byte
	Data          map[string][]byte
}

// TODO: The storage information is stored in a location with zeros so it can
// expand. Then the actual data is just appended to the end. Each time a new
// item is appended a record gets added to the store. Things are pulled out by
// knowing the start/end points and using a checksum to verify the data before
// using it.
// The benefit of this style, is we should be able to add files to our store
// during runtime. Instead of only being able to turn files ijnto byte code and
// saving them into *.go files.
type StorageItem struct {
	Key        string
	Name       string
	Offset     uint64
	Size       uint64
	Checksum   string
	Encrypted  bool
	Compressed bool
}
