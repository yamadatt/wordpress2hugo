# wordpressからhugoに変換する

## 概要

WordPressのエクスポートファイルをhugoで使用するファイルにエクスポートする。

画像ファイルも含めてエクスポートする。



## 使い方

以下のエクスポートのファイルを書き換えて実行する。

```golang
	// The path to the exported XML file containing all posts
	WordPressXMLFile = "./test.WordPress.2023-11-03.xml"
```

## 出力するファイル構成

以下の構成で出力される

```
BaseDir/
  content/
    posts/
      Some-post-title/
        index.md
        ImagesDir/
          post-image-1.jpg
          post-image-2.png
      A-different-post-title/
        index.md
        ImagesDir/
          post-image-for-this-other-arcicle.jpg
```

## 参考

以下のリポジトリを参考にした。

自分用に改造したら、汎用性がなくなってしまった。なので自分用のリポジトリとして使用することにした。

[wjessop/wordpress\_to\_hugo: A small program for converting a wordpress XML dump to Hugo files, including images\.](https://github.com/wjessop/wordpress_to_hugo)