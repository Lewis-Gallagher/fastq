package main

import (
	"os"
)

// https://stackoverflow.com/questions/6059302/how-can-i-check-if-a-file-is-gzip-compressed#comment7013798_60593420
// testBytes, err := bReader.Peek(2) //read 2 bytes
// ....
// if testBytes[0] == 31 && testBytes[1] == 139 {
//     //gzip
// }else{
//    ...
// }

// func IsGzipCompressed(r io.Reader) (bool, error) {

// 	// peak at the first 1 MiB of file
// 	br := bufio.NewReaderSize(r, 1>>20)

// 	// check first two bytes should be 0x1 0x8b the 'magic number'
// 	// https://stackoverflow.com/questions/6059302/how-can-i-check-if-a-file-is-gzip-compressed
// 	magic, err := br.Peek(2)

// 	if err != nil && err != io.EOF {
// 		return false, err
// 	}

// 	if len(magic) == 2 && magic[0] == 0x1f && magic[1] == 0x8b {
// 		return true, err
// 	}

// 	return false, err

// }

// func OpenFastq(path string) (io.ReadCloser, error) {

// 	f, err := os.Open(path)

// 	defer f.Close()

// 	if err != nil {
// 		return nil, err
// 	}

// 	if IsGzipCompressed(f) {
// 		r, err := gzip.NewReader(f)
// 	} else {
// 		r, err := bufio.Reader()
// 	}

// 	// r, ok, err := autogzip.Open(f)
// 	// // ok is true if the file is compressed
// 	// // r contains the decompressed stream
// 	// if err != nil {
// 	// 	return nil, err
// 	// }

// 	// // input is compressed
// 	// if ok == true {
// 	// 	return r, nil
// 	// }

// 	// // input is not compressed
// 	return f, nil
// }

func main() {
	OpenFastq(os.Args[1])
}
