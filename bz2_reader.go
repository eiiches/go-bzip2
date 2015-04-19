package bzip2

// #cgo LDFLAGS: -lbz2
// #include <bzlib.h>
import "C"
import "fmt"
import "unsafe"
import "io"

type Bz2Reader struct {
	context C.bz_stream
	reader  io.Reader
	buf     []byte
	eof     bool
}

func NewReader(reader io.Reader) (*Bz2Reader, error) {
	result := &Bz2Reader{
		reader: reader,
		buf:    make([]byte, 1024),
	}
	if err := C.BZ2_bzDecompressInit(&result.context, 0, 0); err != C.BZ_OK {
		return nil, fmt.Errorf("%v", err)
	}
	return result, nil
}

func (self *Bz2Reader) Read(outbuf []byte) (int, error) {
	if self.context.avail_in == 0 && !self.eof {
		n, err := self.reader.Read(self.buf[:])
		if err != nil {
			if err == io.EOF {
				self.eof = true
			} else {
				return 0, err
			}
		}
		self.context.next_in = (*C.char)(unsafe.Pointer(&self.buf[0]))
		self.context.avail_in = C.uint(n)
	}
	self.context.next_out = (*C.char)(unsafe.Pointer(&outbuf[0]))
	self.context.avail_out = C.uint(len(outbuf))

	ret := C.BZ2_bzDecompress(&self.context)
	if ret != C.BZ_OK && ret != C.BZ_STREAM_END {
		return 0, fmt.Errorf("%v", ret)
	}
	if ret == C.BZ_OK && self.eof && self.context.avail_in == 0 && self.context.avail_out > 0 {
		return 0, io.ErrUnexpectedEOF
	}
	outn := len(outbuf) - int(self.context.avail_out)
	if ret == C.BZ_STREAM_END {
		return outn, io.EOF
	}
	return outn, nil
}

func (self *Bz2Reader) Close() error {
	r := C.BZ2_bzDecompressEnd(&self.context)
	if r != C.BZ_OK {
		return fmt.Errorf("%v", r)
	}
	return nil
}
