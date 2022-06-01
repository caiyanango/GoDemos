package mongodb_example

import (
	"bufio"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io"
	"os"
	"strings"
)

const (
	insert = iota + 1
	delete
	update
	retrieve
	disconnect
)

var client *mongo.Client
var ctx context.Context

type student struct {
	Name string
	Age  int
}

func Connect() {
	ctx = context.TODO()
	clientOptions := options.Client().ApplyURI("mongodb://192.168.9.53:27017")
	var err error
	client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		client.Disconnect(ctx)
		fmt.Println("Disconnect from MongoDB!")
	}()
	err = client.Ping(ctx, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Connected to MongoDB!")
	c := client.Database("go_db").Collection("student")
	var cmd int
	for {
		fmt.Println("Please select your operation:")
		fmt.Println("1、Insert  2、Delete  3、Update  4、Retrieve  5、Disconnect")
		fmt.Scan(&cmd)
		switch cmd {
		case insert:
			{
				stuInfo := []interface{}{}
				var s student
				fmt.Println("Enter student info[name=x age=x]:")
				for {
					input := readLine(os.Stdin)
					if strings.EqualFold(input, "over") {
						break
					}
					fmt.Sscanf(input, "name=%s age=%d", &s.Name, &s.Age)
					stuInfo = append(stuInfo, s)
				}
				err := insertData(c, stuInfo)
				if err != nil {
					fmt.Printf("insert data failed, %v\n", err)
					continue
				}
			}
		case delete:
			{
				filter := bson.D{}
				element := bson.E{Key: "name"}
				array := bson.A{}
				fmt.Println("Enter name you want to delete[name name name]:")
				input := readLine(os.Stdin)
				if strings.EqualFold(input, "all") {
					fmt.Println("delete all data")
				} else {
					for _, name := range strings.Split(input, " ") {
						array = append(array, name)
					}
					element.Value = bson.D{{Key: "$in", Value: array}}
					filter = append(filter, element)
					fmt.Println(filter)
				}
				err := deleteData(c, filter)
				if err != nil {
					fmt.Printf("delete data failed, %v\n", err)
					continue
				}
			}
		case update:
			{
				filterUpdate := bson.D{}
				fmt.Println("Enter name you want to modify:")
				name := readLine(os.Stdin)
				filterFind := bson.D{{Key: "name", Value: name}}
				element := bson.E{Key: "$set"}
				var age int
				fmt.Println("Enter new name and age[name=x age=x]:")
				stuInfo := readLine(os.Stdin)
				fmt.Sscanf(stuInfo, "name=%s age=%d", &name, &age)
				element.Value = bson.D{{Key: "name", Value: name}, {Key: "age", Value: age}}
				filterUpdate = append(filterUpdate, element)
				err := updateData(c, filterFind, filterUpdate)
				if err != nil {
					fmt.Printf("update data failed, %v\n", err)
					continue
				}
			}
		case retrieve:
			{
				err := retrieveData(c, bson.D{})
				if err != nil {
					fmt.Printf("retrieve data failed, %v\n", err)
					continue
				}
			}
		case disconnect:
			return
		default:
			fmt.Println("Unsupport option.")
		}
	}

}

func insertData(collection *mongo.Collection, info []interface{}) error {
	result, err := collection.InsertMany(ctx, info)
	if err != nil {
		return err
	}
	fmt.Println("Inserted documents:", result.InsertedIDs)
	return nil
}

func deleteData(collection *mongo.Collection, filter interface{}) error {
	result, err := collection.DeleteMany(ctx, filter)
	if err != nil {
		return err
	}
	fmt.Println("DeleteCount:", result.DeletedCount)
	return nil
}

func updateData(collection *mongo.Collection, filter interface{}, update interface{}) error {
	result, err := collection.UpdateMany(ctx, filter, update)
	if err != nil {
		return err
	}
	fmt.Println("ModifiedCount:", result.ModifiedCount)
	return nil
}
func retrieveData(collection *mongo.Collection, filter interface{}) error {
	cur, err := collection.Find(ctx, filter)
	if err != nil {
		return err
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var result bson.D
		err := cur.Decode(&result)
		if err != nil {
			return err
		}
		fmt.Printf("%v\n", result)
	}
	if err := cur.Err(); err != nil {
		return err
	}
	return nil
}

func readLine(reader io.Reader) string {
	var inputStr string
	rd := bufio.NewReader(reader)
	for {
		inputStr, _ = rd.ReadString('\n')
		if inputStr == "\r\n" || inputStr == "\n" {
			continue
		}
		break
	}
	if strings.Contains(inputStr, "\r\n") {
		inputStr = inputStr[0 : len(inputStr)-2]
	} else if strings.Contains(inputStr, "\n") {
		inputStr = inputStr[0 : len(inputStr)-1]
	}
	return inputStr
}
