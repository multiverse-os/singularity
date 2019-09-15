package binaries

//https://github.com/mewspring/binary/blob/master/binary.go
type Binary struct {
	Name     string
	Filename string
	Size     int
	Data     []byte
}

// TODO: Intended to be stored in a binary for later insertion of data.
type File struct {
	Format   Format
	Arch     Arch
	Entry    uint64
	Segments []*Segment
	Sections []*Section
	Imports  map[uint64]string
	Exports  map[string]uint64
}

type Arch uint8

const (
	ArchX86 Arch = 1 + iota
	ArchX86_64
)

type Format uint8

const (
	FormatELF Format = 1 + iota
	FormatPE
)

type Segment struct {
	Address     uint64
	Permissions Permissions
	Data        []byte
}

type Section struct {
	Name        string
	Address     uint64
	Permissions Permissions
	Data        []byte
}

type Permissions uint8

// TODO: Set these to the right int value so they can be convereted easily
const (
	PermissionExecute Permissions = 1 << iota
	PermissionWrite
	PermissionRead
)
