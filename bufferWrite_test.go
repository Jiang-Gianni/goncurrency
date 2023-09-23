package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"testing"

	"github.com/Jiang-Gianni/goncurrency/chapt4"
)

func BenchmarkUnbufferedWrite(b *testing.B) {
	performWrite(b, tmpFileOrFatal())
}

// BufferWrite is faster because bufio.Writer writes in chunks, while bytes.Buffer must grow its allocated memory
func BenchmarkBufferedWrite(b *testing.B) {
	bufferredFile := bufio.NewWriter(tmpFileOrFatal())
	performWrite(b, bufio.NewWriter(bufferredFile))
}

func tmpFileOrFatal() *os.File {
	file, err := os.CreateTemp("", "tmp")
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	return file
}

func performWrite(b *testing.B, writer io.Writer) {
	done := make(chan interface{})
	defer close(done)
	b.ResetTimer()
	for bt := range chapt4.Take(done, chapt4.Repeat(done, byte(0)), b.N) {
		writer.Write([]byte{bt.(byte)})
	}
}
