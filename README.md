go-bzip2
--------

This package provides a bzip2 Reader interface over libbzip2, similar to the one in the standard `compress/bzip2` package.

The standard package denies to decompress bzip2 files which contain *randomised* blocks.  Though the feature is deprecated, some libraries still seem to be using it to compress data, generating files which `compress/bzip2` cannot decompress.  I needed to handle such files and this is a library for that.  This package exists just as a temporal workaround to the problem.


Installation &amp; Usage
--------------------

```sh
go get github.com/eiiches/go-bzip2
```

```go
import "github.com/eiiches/go-bzip2"
```

Example
-------

```go
fp, err := os.Open("foo.bz2")
if err != nil {
	panic(err)
}
defer fp.Close()

reader, err := bzip2.NewReader(fp)
if err != nil {
	panic(err)
}
defer reader.Close() // Bz2Reader must be Close()ed to free memory allocated internally

var buf [1024]byte
for {
	n, err := reader.Read(buf[:])
	if err != nil && err != io.EOF {
		panic(err)
	}
	os.Stdout.Write(buf[:n])
	if err == io.EOF {
		break
	}
}
```

License
-------

This project is licensed under the terms of the 3-clause BSD license. See [LICENSE.txt](LICENSE.txt) for the full license text.
