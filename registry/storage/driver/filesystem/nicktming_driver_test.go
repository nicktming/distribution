package filesystem

import (
	"fmt"
	"os"
	"testing"
)

// go test -v -test.run TestDriverFileWrite

func TestDriverFileWrite(t *testing.T)  {
	file, err := os.OpenFile("test.txt", os.O_CREATE | os.O_RDWR, 0644)
	if err != nil {
		fmt.Printf("could not open file and err:%v\n", err)
		return
	}

	fw := newFileWriter(file, 0)
	fw.Write([]byte("test"))
	// 此时发现test1.txt中并没有任何内容  因为在缓存中
	// 需要通过commit才可以把缓存中的内容放到文件中
	fmt.Printf("size:%v\n", fw.size)
	defer file.Close()
}

func TestDriverFileWriteCommit(t *testing.T)  {
	file, err := os.OpenFile("test1.txt", os.O_CREATE | os.O_RDWR, 0644)
	if err != nil {
		fmt.Printf("could not open file and err:%v\n", err)
		return
	}
	content1 := "test1"
	file.Write([]byte(content1))
	// 到文件末尾
	n, _ := file.Seek(0, os.SEEK_END)
	fmt.Printf("file end : %v\n", n)
	if n != int64(len(content1)) {
		t.Fatalf("end of file %v != %v\n", n, len(content1))
	}
	fw := newFileWriter(file, n)
	content2 := "test"
	fw.Write([]byte(content2))
	// 此时发现test1.txt中并没有任何内容  因为在缓存中
	// 需要通过commit才可以把缓存中的内容放到文件中
	if fw.Size() != int64(len(content1) + len(content2)) {
		t.Fatalf("after filewrite write %v != %v\n", fw.Size(), len(content1) + len(content2))
	}
	fw.Commit()
	defer file.Close()
}

func TestDriverFileWriteCancel(t *testing.T)  {
	file, err := os.OpenFile("test2.txt", os.O_CREATE | os.O_RDWR, 0644)
	if err != nil {
		fmt.Printf("could not open file and err:%v\n", err)
		return
	}
	fw := newFileWriter(file, 0)
	fw.Write([]byte("test"))
	fw.Cancel()
}




