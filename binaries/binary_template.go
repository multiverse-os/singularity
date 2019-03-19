package binaries

// TODO: Create a template that can be used to store the binary data
// in such a way that it can be created during runtime. Ideally
// it could also be stored in a persistent data store too, and this
// could be a method of easily making it compilable for others
// outside this isntance

// Ideally we want a format that we can scan and build into a
// file and load. SO maybe turn json into Os.File and append
// that OS file. Then update something preallocated space
// in our primary program that lets us know where each
// appended program is
