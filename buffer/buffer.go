package buffer

// Buffer contains the text to edit.
type Buffer interface {
	Filename() string
	SetFilename(filename string)
	Data() []byte
	SetData(data []byte)
}

// Bytes is a simple implenetation of the buffer that stores data in a []byte.
type Bytes struct {
	filename string
	data     []byte
}

// New constructs a new ByteBuffer object containing data.
func New() Bytes {
	return Bytes{}
}

// Filename returns the buffers filename as a byte slice.
func (buffer *Bytes) Filename() string {
	return buffer.filename
}

// SetFilename ses the filename of that will be used when the buffer is saved.
func (buffer *Bytes) SetFilename(filename string) {
	buffer.filename = filename
}

// Data returns the buffers data as a byte slice.
func (buffer *Bytes) Data() []byte {
	return buffer.data
}

// SetData sets the data of the buffer.
func (buffer *Bytes) SetData(data []byte) {
	buffer.data = data
}
