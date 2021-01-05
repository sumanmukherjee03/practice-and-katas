package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

var bufPool = sync.Pool{
	New: func() interface{} {
		fmt.Println("Allocating a new bytes.Buffer")
		return new(bytes.Buffer) // The new builtin function allocates memory for an empty value of the type and returns a pointer to that
	},
}

func main() {
	log(os.Stdout, "debug-message-1")
	log(os.Stdout, "debug-message-2")
}

// If the log function is getting called from hundreds of places,
// it could lead to lots of stale memory and degrade performance because the GC would have to collect more objects
// That's why it is useful to use a sync pool for the buffer.
func log(w io.Writer, msg string) {
	b := bufPool.Get().(*bytes.Buffer) // type cast the empty interface back to the bytes.Buffer pointer
	b.Reset()                          // reset the buffer so that previous values are flushed - just a safety net
	b.WriteString(time.Now().Format("15:04:05"))
	b.WriteString(" : ")
	b.WriteString(msg)
	b.WriteString("\n")
	w.Write(b.Bytes()) // write the bytes to the io writer
	bufPool.Put(b)     // put the resource back into the pool since the work is done
}
