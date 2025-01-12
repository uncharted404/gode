# Gode  #

#### 项目简介

Gode 基于 Go 语言实现，它让你能够在 Go 应用程序里无缝调用 JavaScript 代码片段。提供了一套简单直观的 API，借助
Node.js 的强大功能与 JavaScript 代码交互。

#### 安装使用

##### 需要 node.js 环境
##### go version >= 1.17

#### 用法示例

##### 创建 `Gode` 实例

```go
// 创建一个 Gode 实例
gode, err := New()
if err != nil {
    retrun err
}
// 创建一个带有上下文的 Gode 实例
gode, err := NewWithContext(context.Background())
if err != nil {
    return err
}
```

##### 执行 JavaScript 代码
```go
// 创建一个 Gode 实例
gode, err := New()
if err != nil {
    retrun err
}

res, err := gode.Eval("1 + 2")
if err != nil {
	return err
}

fmt.Printf("%+v\n", res) // 3
```
##### 执行指定函数
```go
// 创建一个 Gode 实例
gode, err := New("function add(a, b) {return a + b}")
if err != nil {
    retrun err
}

res, err := gode.Call("add", 1, 2)
if err != nil {
	return err
}

fmt.Printf("%+v\n", res) // 3
```


#### 参考项目

* [doloopwhile/PyExecJS](https://github.com/doloopwhile/PyExecJS)

