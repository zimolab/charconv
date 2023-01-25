# charconv

封装golang.org/x/text库中encoding的包，实现字符串编码转换。封装github.com/saintfish/chardet包实现文件编码检测。

由于golang使用utf-8作为其原生编码方式，在这种语境下，编码和解码的含义如下：

编码：utf-8  => 字符集x

解码：字符集x => utf-8


