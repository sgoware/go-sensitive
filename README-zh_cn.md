# 🚫 go-sensitive

[![build](https://img.shields.io/badge/build-1.01-brightgreen)](https://github.com/StellarisW/go-sensitive)[![go-version](https://img.shields.io/badge/go-~%3D1.19-30dff3?logo=go)](https://github.com/StellarisW/go-sensitive)

[English](README.md) | 中文

> 敏感词过滤, 支持多种数据源加载, 多种过滤算法, 多种操作功能

## 🌟 Feature

- 支持多种操作功能
    - `Filter()` 返回过滤后的文本
    - `Replace()` 返回替换了敏感词后的文本
    - `IsSensitive()` 返回文本是否含有敏感词
    - `FindOne()` 返回匹配到的第一个敏感词
    - `FindAll()` 返回匹配到的所有敏感词
    - `FindAllCount()` 返回匹配到的所有敏感词及出现次数
- 支持多种数据源加载, 动态修改数据源
    - 支持内存存储
    - 支持mysql存储
    - 支持mongo存储
    - 支持多种字典加载方式
    - 支持运行过程中动态修改数据源
- 支持多种过滤算法
    - **DFA** 使用 `trie tree` 数据结构匹配敏感词

## ⚙ Usage

```go
package main

import (
	"fmt"
	"github.com/StellarisW/go-sensitive"
)

func main() {
    filterManager := sensitive.NewFilter(
        sensitive.StoreOption{
            Type: sensitive.StoreMemory
        },
        sensitive.FilterOption{
            Type: sensitive.FilterDfa
        }
    )
    
    // 加载字典
    
    err:=filterManager.GetStore().LoadDictPath("path-to-dict")
    if err != nil {
        fmt.Println(err)
        return
	}
    
    // 动态增加词汇
    
    err=filterManager.GetStore().AddWord("这是敏感词1", "这是敏感词2", "这是敏感词3")
    if err != nil {
        fmt.Println(err)
        return
	}
    
    fmt.Println(filterManager.GetFilter().IsSensitive("这是敏感词1,这是敏感词2,这是敏感词3,这是敏感词1,这里没有敏感词"))
    
    fmt.Println(filterManager.GetFilter().Filter("这是敏感词1,这是敏感词2,这是敏感词3,这是敏感词1,这里没有敏感词"))
    
    fmt.Println(filterManager.GetFilter().Replace("这是敏感词1,这是敏感词2,这是敏感词3,这是敏感词1,这里没有敏感词", '*'))
    
    fmt.Println(filterManager.GetFilter().FindOne("这是敏感词1,这是敏感词2,这是敏感词3,这是敏感词1,这里没有敏感词"))

    fmt.Println(filterManager.GetFilter().FindAll("这是敏感词1,这是敏感词2,这是敏感词3,这是敏感词1,这里没有敏感词"))

    fmt.Println(filterManager.GetFilter().FindAllCount("这是敏感词1,这是敏感词2,这是敏感词3,这是敏感词1,这里没有敏感词"))
}
```

## ✔ Get

```
$ go get -u github.com/StellarisW/go-sensitive
```

## 📂 Import

```go
import "github.com/StellarisW/go-sensitive"
```

## 

## 📌 TODO

- [ ] add mongo data source support
- [ ] add  bloom algorithm