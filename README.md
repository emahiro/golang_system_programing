# golang_system_programing
「Goならわかるシステムプログラミング」の学習用リポジトリ

# Desc
archiveに作り終わったスクリプトは追加します。scriptを実行したい時はsrc配下にmain.goと登録します。  
image src `▶︎`  https://en.wikipedia.org/wiki/File:Lenna.png

# png_reader.go
pngファイルを読み込みます。サンプルファイルは以下のようなチャンクの構成になっています。

```bash
chunk 'IHDR' (13 byte)
chunk 'sRGB' (1 byte)
chunk 'IDAT' (473761 byte)
chunk 'IEND' (0 byte)
```
# write_text_chunk.go
サンプル画像を一度読み込んでから、新しく任意のテキストを追加します。

以下のように読み込まれます。

```bash
chunk 'IHDR' (13 byte)
chunk 'tExt' (15 byte)
Test Text Chunk
chunk 'sRGB' (1 byte)
chunk 'IDAT' (473761 byte)
chunk 'IEND' (0 byte)
```