package store

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"github.com/imroc/req/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io"
	"net/http"
	"os"
)

const (
	defaultCollection = "dirties"
)

type MongoConfig struct {
	Address    string
	Port       string
	Username   string
	Password   string
	Database   string
	Collection string
}

type doc struct {
	Id   string `bson:"_id"`
	Word string `bson:"word"`
}

type MongoModel struct {
	store   *mongo.Collection
	addChan chan string
	delChan chan string
}

func NewMongoModel(config *MongoConfig) *MongoModel {
	clientOptions := options.Client().ApplyURI(
		fmt.Sprintf("mongodb://%s:%s",
			config.Address,
			config.Port,
		),
	)

	if config.Username != "" {
		clientOptions.SetAuth(options.Credential{
			Username: config.Username,
			Password: config.Password,
		})
	}

	mdb, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil
	}

	err = mdb.Ping(context.TODO(), nil)
	if err != nil {
		return nil
	}

	if config.Database == "" {
		return nil
	}

	if config.Collection == "" {
		config.Collection = defaultCollection
	}

	collection := mdb.Database(config.Database).Collection(config.Collection)

	_, err = collection.Indexes().CreateOne(context.Background(),
		mongo.IndexModel{
			Keys:    bson.D{{"word", 1}},
			Options: options.Index().SetUnique(true),
		},
	)
	if err != nil {
		return nil
	}

	return &MongoModel{
		store:   collection,
		addChan: make(chan string),
		delChan: make(chan string),
	}
}

func (m *MongoModel) LoadDictPath(paths ...string) error {
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

func (m *MongoModel) LoadDictHttp(urls ...string) error {
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

func (m *MongoModel) LoadDict(reader io.Reader) error {
	buf := bufio.NewReader(reader)
	var words []interface{}

	for {
		line, _, err := buf.ReadLine()
		if err != nil {
			if err != io.EOF {
				return err
			}
			break
		}

		word := string(line)

		words = append(words, bson.D{
			{"word", word},
		})

		m.addChan <- word
	}

	ctx := context.Background()

	_, _ = m.store.InsertMany(ctx, words, options.InsertMany().SetOrdered(false))

	return nil
}

func (m *MongoModel) ReadChan() <-chan string {
	ch := make(chan string)

	go func() {
		ctx := context.Background()
		cur, _ := m.store.Find(ctx,
			bson.D{},
			options.Find().SetProjection(
				bson.D{
					{"_id", 0},
					{"word", 1},
				},
			),
		)

		for cur.Next(ctx) {
			var word doc

			_ = cur.Decode(&word)

			ch <- word.Word
		}

		close(ch)
	}()

	return ch
}

func (m *MongoModel) ReadString() []string {
	ctx := context.Background()
	cur, err := m.store.Find(ctx,
		bson.D{},
		options.Find().SetProjection(
			bson.D{
				{"_id", 0},
				{"word", 1},
			},
		),
	)
	if err != nil {

	}

	var words []*doc

	err = cur.All(ctx, &words)

	res := make([]string, 0, len(words))

	for _, word := range words {
		res = append(res, word.Word)
	}

	return res
}

func (m *MongoModel) GetAddChan() <-chan string {
	return m.addChan
}

func (m *MongoModel) GetDelChan() <-chan string {
	return m.delChan
}

func (m *MongoModel) AddWord(words ...string) error {
	for _, word := range words {
		_, err := m.store.UpdateOne(context.Background(),
			bson.D{
				{"word", word},
			},
			bson.D{
				{"$set", bson.D{
					{"word", word},
				}},
			},
			options.Update().SetUpsert(true),
		)
		if err != nil {
			return err
		}

		m.addChan <- word
	}

	return nil
}

func (m *MongoModel) DelWord(words ...string) error {
	for _, word := range words {
		_, err := m.store.DeleteOne(context.Background(),
			bson.D{
				{"word", word},
			},
		)
		if err != nil {
			return err
		}

		m.delChan <- word
	}

	return nil
}
