package store

import (
	"bufio"
	"errors"
	"github.com/imroc/req/v3"
	cmap "github.com/orcaman/concurrent-map/v2"
	"io"
	"net/http"
	"os"
)

type MemoryModel struct {
	store   cmap.ConcurrentMap[string, struct{}]
	AddChan chan string
	DelChan chan string
}

func NewMemoryModel() *MemoryModel {
	return &MemoryModel{
		store:   cmap.New[struct{}](),
		AddChan: make(chan string),
		DelChan: make(chan string),
	}
}

func (m *MemoryModel) LoadDictPath(paths ...string) error {
	for _, path := range paths {
		err := func(path string) error {
			f, err := os.Open(path)
			defer func(f *os.File) {
				_ = f.Close()
			}(f)
			if err != nil {
				return err
			}

			return m.LoadDict(f)
		}(path)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *MemoryModel) LoadDictHttp(urls ...string) error {
	for _, url := range urls {
		err := func(url string) error {
			httpRes, err := req.Get(url)
			if err != nil {
				return err
			}
			if httpRes == nil {
				return errors.New("nil http response")
			}
			if httpRes.StatusCode != http.StatusOK {
				return errors.New(httpRes.GetStatus())
			}

			defer func(Body io.ReadCloser) {
				_ = Body.Close()
			}(httpRes.Body)

			return m.LoadDict(httpRes.Body)
		}(url)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *MemoryModel) LoadDict(reader io.Reader) error {
	buf := bufio.NewReader(reader)
	for {
		line, _, err := buf.ReadLine()
		if err != nil {
			if err != io.EOF {
				return err
			}
			break
		}

		m.store.Set(string(line), struct{}{})
		m.AddChan <- string(line)
	}

	return nil
}

func (m *MemoryModel) ReadChan() <-chan string {
	ch := make(chan string)

	go func() {
		for key := range m.store.Items() {
			ch <- key
		}
		close(ch)
	}()

	return ch
}

func (m *MemoryModel) ReadString() []string {
	res := make([]string, 0, m.store.Count())

	for key := range m.store.Items() {
		res = append(res, key)
	}

	return res
}

func (m *MemoryModel) GetAddChan() <-chan string {
	return m.AddChan
}

func (m *MemoryModel) GetDelChan() <-chan string {
	return m.DelChan
}

func (m *MemoryModel) AddWord(words ...string) error {
	for _, word := range words {
		m.store.Set(word, struct{}{})
		m.AddChan <- word
	}

	return nil
}

func (m *MemoryModel) DelWord(words ...string) error {
	for _, word := range words {
		m.store.Remove(word)
		m.DelChan <- word
	}

	return nil
}
