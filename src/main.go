package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

const (
	OFFSET = 8
)

func dumpChunk(chunk io.Reader) {
	var l int32
	binary.Read(chunk, binary.BigEndian, &l)
	// png を4byteごとに分けるときのいれもの
	buf := make([]byte, 4)
	chunk.Read(buf)
	fmt.Printf("chunk '%v' (%d byte)\n", string(buf), l)
}

func readChunks(file *os.File) []io.Reader {
	// pngをchunkに分けたときのchunkを格納する入れ物
	var chunks []io.Reader

	// pngの最初の8byteを飛ばして9byte目から読み込む
	file.Seek(OFFSET, 0)
	var ofs int64 = OFFSET

	for {
		var l int32
		err := binary.Read(file, binary.BigEndian, &l)
		if err == io.EOF {
			break
		}

		// 最初のchunkから12(3 * 4byte)byteずらす
		chunks = append(chunks, io.NewSectionReader(file, ofs, int64(l)+12))

		// 次のchunkの先頭に移動
		// 現在位置は長さを読み終わった箇所なのでchunk名（4byte） + data長 + CRC(4byte)先に移動する
		ofs, _ = file.Seek(int64(l+8), 1)
	}

	return chunks
}

func main() {
	// pngをチャンクに分けて読み込み
	file, err := os.Open("./imgs/Lenna.png")
	if err != nil {
		fmt.Printf("error! err: %v", err)
		os.Exit(-1)
	}
	defer file.Close()

	chunks := readChunks(file)
	for _, c := range chunks {
		dumpChunk(c)
	}
}
