# 🚫 go-sensitive

[![build](https://img.shields.io/badge/build-1.01-brightgreen)](https://github.com/sgoware/go-sensitive)[![go-version](https://img.shields.io/badge/go-~%3D1.19-30dff3?logo=go)](https://github.com/sgoware/go-sensitive)

English | [中文](README-zh_cn.md)

> Filter sensitive words, support multiple data sources, filter algorithms and functions

## 🌟 Feature

- support multiple functions
    - `Filter()` return filtered text
    - `Replace()` return text which sensitive words that is been replaced
    - `IsSensitive()` Check whether the text has sensitive word
    - `FindOne()` return first sensitive word that has been found in the text
    - `FindAll()` return all sensitive word that has been found in the text
    - `FindAllCount()` return all sensitive word with its count that has been found in the text
- support multiple data sources with dynamic modification
    - support memory storage
    - support mysql storage
    - support mongo storage
    - support multiple ways of add dict
    - support dynamic add/del sensitive word while running
- support multiple filter algorithms
    - **DFA** use `trie tree`  to filter sensitive words

## ⚙ Usage

```go
package main

import (
	"fmt"
	"github.com/sgoware/go-sensitive"
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
    
    // load dict
    
    err:=filterManager.GetStore().LoadDictPath("path-to-dict")
    if err != nil {
        fmt.Println(err)
        return
	}
    
    // dynamic add sensitive words
    
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
$ go get -u github.com/sgoware/go-sensitive
```

## 📂 Import

```go
import "github.com/sgoware/go-sensitive"
```

## 

## 📌 TODO

- [ ] add redis data source support
- [ ] add bloom algorithm
