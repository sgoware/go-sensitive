package store

import (
	"bufio"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/imroc/req/v3"
	"github.com/jmoiron/sqlx"
	"io"
	"net/http"
	"os"
)

const (
	defaultTable = "dirties"
)

type MysqlConfig struct {
	Dsn       string
	Database  string
	TableName string
}

type Subject struct {
	Id   int64  `db:"id"`
	Word string `db:"word"`
}

type MysqlModel struct {
	store     *sqlx.DB
	TableName string
	addChan   chan string
	delChan   chan string
}

func NewMysqlModel(config *MysqlConfig) *MysqlModel {
	db, err := sqlx.Connect("mysql", config.Dsn)
	if err != nil {
		return nil
	}

	if config.TableName == "" {
		config.TableName = defaultTable
	}

	var tableName string

	err = db.Get(
		&tableName,
		fmt.Sprintf(
			"SELECT `TABLE_NAME` FROM information_schema.tables WHERE table_schema = '%s' AND table_name = '%s' LIMIT 1",
			config.Database,
			config.TableName),
	)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil
		}
	}

	if tableName == "" {
		_, err = db.Exec(fmt.Sprintf("CREATE TABLE `%s` "+
			"(`id` bigint(20) NOT NULL AUTO_INCREMENT, "+
			"`word` varchar(255) NOT NULL, "+
			"PRIMARY KEY (`id`) USING BTREE)",
			config.TableName),
		)
		if err != nil {
			return nil
		}
	}

	return &MysqlModel{
		store:     db,
		TableName: config.TableName,
		addChan:   make(chan string),
		delChan:   make(chan string),
	}
}

func (m *MysqlModel) LoadDictPath(paths ...string) error {
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

func (m *MysqlModel) LoadDictHttp(urls ...string) error {
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

func (m *MysqlModel) LoadDict(reader io.Reader) error {
	buf := bufio.NewReader(reader)
	var words []*Subject
	set := make(map[string]struct{})

	for {
		line, _, err := buf.ReadLine()
		if err != nil {
			if err != io.EOF {
				return err
			}
			break
		}

		word := string(line)

		if _, ok := set[word]; !ok {
			words = append(words, &Subject{
				Id:   0,
				Word: word,
			})
			set[word] = struct{}{}
		}

		m.addChan <- word
	}

	_, err := m.store.NamedExec(fmt.Sprintf("INSERT INTO `%s` (`word`) VALUES (:word)", m.TableName), words)
	if err != nil {
		return err
	}

	_, err = m.store.Exec(fmt.Sprintf("DELETE FROM `%s` AS t1 "+
		"WHERE t1.`id` <> "+
		"(SELECT t.minid FROM "+
		"(SELECT MIN(t2.`id`) AS minid FROM `%s` AS t2 WHERE t1.`word` = t2.`word`) t )",
		m.TableName,
		m.TableName),
	)
	if err != nil {
		return err
	}

	return nil
}

func (m *MysqlModel) ReadChan() <-chan string {
	ch := make(chan string)
	var words []string

	go func() {
		defer close(ch)

		err := m.store.Select(&words, fmt.Sprintf("SELECT `word` FROM `%s`", m.TableName))
		if err != nil {
			return
		}

		for _, word := range words {
			ch <- word
		}
	}()

	return ch
}

func (m *MysqlModel) ReadString() []string {
	var words []string

	err := m.store.Select(&words, fmt.Sprintf("SELECT `word` FROM `%s`", m.TableName))
	if err != nil {
		return nil
	}

	return words
}

func (m *MysqlModel) GetAddChan() <-chan string {
	return m.addChan
}

func (m *MysqlModel) GetDelChan() <-chan string {
	return m.delChan
}

func (m *MysqlModel) AddWord(words ...string) error {
	insertedWords := make([]*Subject, 0, len(words))
	set := make(map[string]struct{})

	for _, word := range words {
		if _, ok := set[word]; !ok {
			insertedWords = append(insertedWords, &Subject{
				Id:   0,
				Word: word,
			})
			set[word] = struct{}{}
		}
	}

	_, err := m.store.NamedExec(fmt.Sprintf("INSERT INTO `%s` (`word`) VALUES (:word)", m.TableName), insertedWords)
	if err != nil {
		return err
	}
	_, err = m.store.Exec(fmt.Sprintf("DELETE FROM `%s` AS t1 "+
		"WHERE t1.`id` <> "+
		"(SELECT t.minid FROM "+
		"(SELECT MIN(t2.`id`) AS minid FROM `%s` AS t2 WHERE t1.`word` = t2.`word`) t )",
		m.TableName,
		m.TableName),
	)
	if err != nil {
		return err
	}

	return nil
}

func (m *MysqlModel) DelWord(words ...string) error {
	query, args, _ := sqlx.In(fmt.Sprintf("DELETE FROM `%s` WHERE `word` IN (?)", m.TableName), words)
	_, err := m.store.Exec(query, args)
	if err != nil {
		return err
	}

	return nil
}
