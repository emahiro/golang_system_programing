package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"hash/crc32"
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
	if bytes.Equal(buf, []byte("tExt")) {
		rawText := make([]byte, l)
		chunk.Read(rawText)
		fmt.Printf("%s\n", string(rawText))
	}
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

		// 最初のchunk(4byte)から12(3 * 4byte)byte含めた16byteが1chunk
		chunks = append(chunks, io.NewSectionReader(file, ofs, int64(l)+12))

		// 次のchunkの先頭に移動
		// 現在位置は長さを読み終わった箇所なのでchunk名（4byte） + data長 + CRC(4byte)先に移動する
		ofs, _ = file.Seek(int64(l+8), 1)
	}

	return chunks
}

func textChunk(txt string) io.Reader {
	byteData := []byte(txt)
	var buf bytes.Buffer
	binary.Write(&buf, binary.BigEndian, int32(len(byteData))) // 書き込むためのbufferを確保。長さは入力文字長。
	buf.WriteString("tExt")                                    //bufferに書き込む
	buf.Write(byteData)                                        // 入力対象の文字をbufferに書き込む
	// CRCを計算して追加
	crc := crc32.NewIEEE()
	io.WriteString(crc, "tExt")
	binary.Write(&buf, binary.BigEndian, crc.Sum32()) // crc分を新しくbufferに書き込む
	return &buf                                       // 書き込み終わったbufferを返す
}

func WriteNewTextChunk() {
	// pngをチャンクに分けて読み込み
	file, err := os.Open("./imgs/Lenna.png")
	if err != nil {
		fmt.Printf("error! err: %v", err)
		os.Exit(-1)
	}
	defer file.Close()

	newFile, _ := os.Create("./imgs/result.png")
	defer newFile.Close()

	chunks := readChunks(file)
	// pngのシグニチャー書き込み
	io.WriteString(newFile, "\x89PNG\r\n\x1a\n")

	// 先頭にIHDR chunkを書き込む
	io.Copy(newFile, chunks[0])

	// TextChunk を追加する
	io.Copy(newFile, textChunk("Test Text Chunk"))

	// 残りのチャンクを新しく追加する
	for _, c := range chunks[1:] {
		io.Copy(newFile, c)
	}
}

func main() {
	// 私いchunkを書き込んだfileを生成する
	WriteNewTextChunk()

	file, err := os.Open("./imgs/result.png")
	if err != nil {
		fmt.Printf("error! err: %v", err)
		os.Exit(-1)
	}

	chunks := readChunks(file)
	for _, c := range chunks {
		dumpChunk(c)
	}
}
