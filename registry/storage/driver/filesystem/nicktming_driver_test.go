package filesystem

import (
	"fmt"
	"os"
	"testing"
)

func TestDriverFileWrite(t *testing.T)  {
	file, err := os.OpenFile("test.txt", os.O_CREATE | os.O_RDWR, 0644)
	if err != nil {
		fmt.Printf("could not open file and err:%v\n", err)
		return
	}

	fw := newFileWriter(file, 0)
	fw.Write([]byte("test"))
	fmt.Printf("size:%v\n", fw.size)
	defer file.Close()
}
