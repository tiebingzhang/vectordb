package database

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"os"
)

const (
	ErrorFileNotFound = "file not found"
)

type DatabaseInterface interface {
	Add(doc Document)
	Save() error
	Close() error
	GetTextByIndex(i int) string
	GetEmbeddingByIndex(i int) []float64
	GetEmbeddings() [][]float64
	GetTexts() []string
}

type Document struct {
	Text      string    `json:"text"`
	Embedding []float64 `json:"embedding"`
}
type Database struct {
	filePath  string
	Documents []Document
}

func (db *Database) Add(doc Document) {
	db.Documents = append(db.Documents, doc)
}

func (db *Database) RemoveByIndex(i int) {
	db.Documents = append(db.Documents[:i], db.Documents[i+1:]...)
}

func (db *Database) RemoveByText(text string) {
	for i, doc := range db.Documents {
		if doc.Text == text {
			db.RemoveByIndex(i)
			return
		}
	}
}

func (db *Database) Save() error {
	return Store(db.filePath, db.Documents)
}

func (db *Database) GetTextByIndex(i int) string {
	return db.Documents[i].Text
}

func (db *Database) GetEmbeddingByIndex(i int) []float64 {
	return db.Documents[i].Embedding
}

func (db *Database) GetEmbeddings() [][]float64 {
	embeddings := make([][]float64, len(db.Documents))
	for i, doc := range db.Documents {
		embeddings[i] = doc.Embedding
	}
	return embeddings
}

func (db *Database) GetTexts() []string {
	texts := make([]string, len(db.Documents))
	for i, doc := range db.Documents {
		texts[i] = doc.Text
	}
	return texts
}

func New(filePath string) (*Database, error) {
	documents, err := Load(filePath)
	if err == nil {
		return &Database{
		filePath:  filePath,
		Documents: documents,
		}, nil
	}
	
	if err.Error() == ErrorFileNotFound {
		return &Database{
			filePath:  filePath,
			Documents: []Document{},
		}, nil
	}
	return nil, err
}

func Load(filePath string) ([]Document, error) {
	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf(ErrorFileNotFound)
	}
	// Read from a file
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Println("error reading from file:", err)
		return nil, err
	}

	buf := *bytes.NewBuffer(data)
	dec := gob.NewDecoder(&buf)

	var doc []Document
	err = dec.Decode(&doc)
	if err != nil {
		log.Println("decode error:", err)
		return nil, err
	}

	return doc, nil
}

func Store(filePath string, embeddings []Document) error {

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	err := enc.Encode(embeddings)
	if err != nil {
		log.Println("encode error:", err)
		return err
	}

	// Save to a file
	err = os.WriteFile(filePath, buf.Bytes(), 0644)
	if err != nil {
		log.Println("error writing to file:", err)
		return err
	}

	return nil
}
