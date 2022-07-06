/**
 *	from github.com/dtylman/scp
 */

package scp

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"golang.org/x/crypto/ssh"
)

const (
	fileMode = "0644"
	buffSize = 1024 * 256
)

//CopyTo copy from local to remote
func CopyTo(sshClient *ssh.Client, local string, remote string, progressCall func(cur, total int64)) (int64, error) {
	session, err := sshClient.NewSession()
	if err != nil {
		return 0, err
	}
	defer session.Close()
	stderr := &bytes.Buffer{}
	session.Stderr = stderr
	stdout := &bytes.Buffer{}
	session.Stdout = stdout
	writer, err := session.StdinPipe()
	if err != nil {
		return 0, err
	}
	defer writer.Close()

	err = session.Start("scp -t " + filepath.ToSlash(filepath.Dir(remote)))
	if err != nil {
		return 0, err
	}

	localFile, err := os.Open(local)
	if err != nil {
		return 0, err
	}
	defer localFile.Close()
	fileInfo, err := localFile.Stat()
	if err != nil {
		return 0, err
	}
	_, err = fmt.Fprintf(writer, "C%s %d %s\n", fileMode, fileInfo.Size(), filepath.ToSlash(filepath.Base(remote)))
	if err != nil {
		return 0, err
	}
	n, err := copyN(writer, localFile, fileInfo.Size(), func(cur, total int64) {
		progressCall(cur, fileInfo.Size())
	})
	if err != nil {
		return 0, err
	}
	err = ack(writer)
	if err != nil {
		return 0, err
	}

	session.Wait()
	//NOTE: Process exited with status 1 is not an error, it just how scp work. (waiting for the next control message and we send EOF)
	return n, nil
}

//CopyFrom copy from remote to local
func CopyFrom(sshClient *ssh.Client, remote string, local string) (int64, error) {
	session, err := sshClient.NewSession()
	if err != nil {
		return 0, err
	}
	defer session.Close()
	stderr := &bytes.Buffer{}
	session.Stderr = stderr
	writer, err := session.StdinPipe()
	if err != nil {
		return 0, err
	}
	defer writer.Close()
	reader, err := session.StdoutPipe()
	if err != nil {
		return 0, err
	}
	err = session.Start("scp -f " + filepath.ToSlash(remote))
	if err != nil {
		return 0, err
	}
	err = ack(writer)
	if err != nil {
		return 0, err
	}
	msg, err := NewMessageFromReader(reader)
	if err != nil {
		return 0, err
	}
	if msg.Type == ErrorMessage || msg.Type == WarnMessage {
		return 0, msg.Error
	}
	log.Printf("Receiving %v", msg)

	err = ack(writer)
	if err != nil {
		return 0, err
	}
	outFile, err := os.Create(local)
	if err != nil {
		return 0, err
	}
	defer outFile.Close()
	n, err := copyN(outFile, reader, msg.Size, nil)
	if err != nil {
		return 0, err
	}
	err = outFile.Sync()
	if err != nil {
		return 0, err
	}
	err = outFile.Close()
	if err != nil {
		return 0, err
	}
	session.Wait()
	return n, nil
}

func ack(writer io.Writer) error {
	var msg = []byte{0, 0, 10, 13}
	n, err := writer.Write(msg)
	if err != nil {
		return err
	}
	if n < len(msg) {
		return errors.New("failed to write ack buffer")
	}
	return nil
}

func copyN(writer io.Writer, src io.Reader, size int64, progressCall func(cur, total int64)) (int64, error) {
	reader := io.LimitReader(src, size)
	var total int64
	for total < size {
		n, err := CopyBuffer(writer, reader, make([]byte, buffSize), progressCall)
		if err != nil {
			return 0, err
		}
		total += n

	}
	return total, nil
}

func CopyBuffer(dst io.Writer, src io.Reader, buf []byte, progressCall func(cur, total int64)) (written int64, err error) {
	if buf != nil && len(buf) == 0 {
		panic("empty buffer in CopyBuffer")
	}
	return copyBuffer(dst, src, buf, progressCall)
}

var errInvalidWrite = errors.New("invalid write result")

// copyBuffer is the actual implementation of Copy and CopyBuffer.
// if buf is nil, one is allocated.
func copyBuffer(dst io.Writer, src io.Reader, buf []byte, progressCall func(cur, total int64)) (written int64, err error) {
	// If the reader has a WriteTo method, use it to do the copy.
	// Avoids an allocation and a copy.
	if wt, ok := src.(io.WriterTo); ok {
		return wt.WriteTo(dst)
	}
	// Similarly, if the writer has a ReadFrom method, use it to do the copy.
	if rt, ok := dst.(io.ReaderFrom); ok {
		return rt.ReadFrom(src)
	}

	if buf == nil {
		size := 32 * 1024
		if l, ok := src.(*io.LimitedReader); ok && int64(size) > l.N {
			if l.N < 1 {
				size = 1
			} else {
				size = int(l.N)
			}
		}
		buf = make([]byte, size)
	}
	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			nw, ew := dst.Write(buf[0:nr])
			if nw < 0 || nr < nw {
				nw = 0
				if ew == nil {
					ew = errInvalidWrite
				}
			}
			written += int64(nw)

			progressCall(written, int64(nr))

			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er != nil {
			if er != io.EOF {
				err = er
			}
			break
		}
	}
	return written, err
}
