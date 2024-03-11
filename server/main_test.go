package main

import (
	"context"
	"fmt"
	"log"
	"runtime"
	"runtime/debug"
	"strings"
	"testing"

	"github.com/tidwall/gjson"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type User struct {
	Name string
}

type StructParams struct {
	StructName string
	Fields     []struct {
		Name     string
		BsonType string
	}
}

var StructTemplateText = `
type {{.StructName}} struct {
	{{ range Field := .Fields }}
		{{Field.Name}} {{Field.BsonType}}
	{{ end }}
}
`

func toPascalCase(s string) string {
	words := strings.Split(s, "_")
	for i, word := range words {
		words[i] = cases.Title(language.English).String(strings.ToLower(word))
	}
	return strings.Join(words, "")
}

func GenStructFromMongoDB(mongodbURI string, Database string) {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(mongodbURI).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.Background(), opts)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(context.Background()); err != nil {
			panic(err)
		}
	}()

	infos, err := client.Database(Database).ListCollectionSpecifications(context.TODO(), bson.M{})
	if err != nil {
		panic(err)
	}

	for _, info := range infos {
		log.Println(info.Name)

		sparams := StructParams{
			StructName: toPascalCase(info.Name),
		}

		// log.Println(info.Options.String())
		result := gjson.Parse(info.Options.String())
		result.ForEach(func(key, value gjson.Result) bool {

			if key.String() == "validator" {

				jsonSchema := value.Get("$jsonSchema")
				properties := jsonSchema.Get("properties").Map()

				for _, r := range jsonSchema.Get("required").Array() {
					key := r.String()
					var btype = "String"
					if prop, ok := properties[key]; ok {
						btype = prop.Get("bsonType").String()
						if btype == "int" {
							btype = "int64"
						} else if btype == "object" {
							btype = "interface{}"
						} else if btype == "array" {
							btype = "[]interface{}"
						}

						log.Println(key, btype)
					}

					sparams.Fields = append(sparams.Fields, struct {
						Name     string
						BsonType string
					}{
						Name:     toPascalCase(key),
						BsonType: toPascalCase(btype),
					})
				}
			}
			return true
		})

	}
}

func TestMain(t *testing.T) {

	// Connect to MongoDB using the MongoDB driver
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	mongodbURI := "mongodb://fusen-dev:fusen-dev@localhost:27017/?retryWrites=true&w=majority"

	opts := options.Client().ApplyURI(mongodbURI).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.Background(), opts)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(context.Background()); err != nil {
			panic(err)
		}
	}()

	GenStructFromMongoDB(mongodbURI, "fusen")

	// Access the user collection in the fusen database
	userCollection := client.Database("fusen").Collection("user")

	userCollection.InsertOne(context.TODO(), User{Name: "213"})

	// Count documents in the collection
	count, err := userCollection.CountDocuments(context.Background(), bson.D{})
	if err != nil {
		panic(err)
	}

	log.Println(count)

	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")
}

func TestOther(t *testing.T) {

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	fmt.Printf("Alloc: %v MiB", memStats.Alloc/1024/1024)

	stats := debug.GCStats{}
	debug.ReadGCStats(&stats)

	numGoroutines := runtime.NumGoroutine()
	fmt.Println("Number of Goroutines:", numGoroutines)

	fmt.Println("Last GC Pause Duration:", stats.PauseTotal.Seconds())
	t.Error()
	// http.HandleFunc("/redirect/test1", func(w http.ResponseWriter, r *http.Request) {

	// })

	// fmt.Println("Server is listening on port 8080...")
	// err := http.ListenAndServe(":8080", nil)
	// if err != nil {
	// 	panic(err)
	// }
}
