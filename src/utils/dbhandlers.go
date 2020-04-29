package utils

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Item struct {
	Usn       string
	Branch    string
	Name      string
	Gender    string
	Dob       string
	Perc10th  string
	Perc12th  string
	Cgpa      string
	Backlog   bool
	Email     string
	Mobile    string
	StarOffer bool
}

func StoreInBLOCKCHAIN(firstblockHASH []byte, firstblockbytes []byte) bool {

	type Item struct {
		Hash  []byte
		Block []byte
	}

	sess, _ := session.NewSession(&aws.Config{

		Region: aws.String("us-east-1"),
	})
	svc := dynamodb.New(sess)

	item := Item{
		Hash:  firstblockHASH,
		Block: firstblockbytes,
	}
	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		fmt.Println("Got error marshalling new item:")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// Create item in table Movies
	tableName := "BLockchain"

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		fmt.Println("Got error calling PutItem:")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return true

}

func StoreLastHash(hash []byte) bool {
	type Item struct {
		Hash     []byte
		LastHash string
	}

	sess, _ := session.NewSession(&aws.Config{

		Region: aws.String("us-east-1"),
	})
	svc := dynamodb.New(sess)

	item := Item{
		Hash:     hash,
		LastHash: "LastHash",
	}
	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		fmt.Println("Got error marshalling new item:")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// Create item in table Movies
	tableName := "Lasthash"

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		fmt.Println("Got error calling PutItem:")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return true

}

func GetLastHash() []byte {
	type Item struct {
		Hash     []byte
		LastHash string
	}

	sess, _ := session.NewSession(&aws.Config{

		Region: aws.String("us-east-1"),
	})
	svc := dynamodb.New(sess)

	params := &dynamodb.ScanInput{
		TableName: aws.String("Lasthash"),
	}
	// Make the DynamoDB Query API call
	result, err := svc.Scan(params)
	if err != nil {
		fmt.Println("Query API call failed:")
		fmt.Println((err.Error()))
		os.Exit(1)
	}
	if result.Items == nil {
		fmt.Println("Could not find Last Hash\nExit")
		os.Exit(1)
	}
	item := Item{}

	for _, i := range result.Items {

		err = dynamodbattribute.UnmarshalMap(i, &item)

		if err != nil {
			fmt.Println("Got error unmarshalling:")
			fmt.Println(err.Error())
			os.Exit(1)
		}

	}

	return item.Hash

}

func GetFromBLOCKCHAIN(hash []byte) []byte {
	type Item struct {
		Hash  []byte
		Block []byte
	}
	sess, _ := session.NewSession(&aws.Config{

		Region: aws.String("us-east-1"),
	})
	svc := dynamodb.New(sess)
	//tableName = "BLockchain"
	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("BLockchain"),
		Key: map[string]*dynamodb.AttributeValue{
			"Hash": {
				B: hash,
			},
		},
	})
	if err != nil {
		fmt.Println(err.Error())

	}
	if result.Item == nil {
		fmt.Println("Could not Block in BlockChain!\nExit")
		os.Exit(1)
	}
	item := Item{}

	err = dynamodbattribute.UnmarshalMap(result.Item, &item)
	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	}

	return item.Block

}

func GetStudentData(usn string) []byte {

	sess, _ := session.NewSession(&aws.Config{

		Region: aws.String("us-east-1"),
	})
	svc := dynamodb.New(sess)
	//tableName = "BLockchain"
	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("Student"),
		Key: map[string]*dynamodb.AttributeValue{
			"Usn": {
				S: aws.String(usn),
			},
		},
	})
	if err != nil {
		fmt.Println(err.Error())

	}
	if result.Item == nil {
		fmt.Println("Student not Found!\nExiting")
		os.Exit(1)

	}
	item := Item{}

	err = dynamodbattribute.UnmarshalMap(result.Item, &item)
	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	}

	resp, _ := json.Marshal(item)
	return resp

}

func StoreInBuffer(block []byte, name string) {

	type Item struct {
		Block []byte
		Name  string
	}

	sess, _ := session.NewSession(&aws.Config{

		Region: aws.String("us-east-1"),
	})
	svc := dynamodb.New(sess)

	item := Item{
		Block: block,
		Name:  name,
	}
	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		fmt.Println("Got error marshalling new item:")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// Create item in table Movies
	tableName := "Buffer"

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		fmt.Println("Got error calling PutItem:")
		fmt.Println(err.Error())
		os.Exit(1)
	}

}

func FetchBlockFromBuffer(name string) []byte {
	type Item struct {
		Block []byte
		Name  string
	}

	sess, _ := session.NewSession(&aws.Config{

		Region: aws.String("us-east-1"),
	})
	svc := dynamodb.New(sess)

	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("Buffer"),
		Key: map[string]*dynamodb.AttributeValue{
			"Name": {
				S: aws.String(name),
			},
		},
	})
	if result.Item == nil {
		fmt.Println("Can not find Block in Buffer\nExiting")
		os.Exit(1)
	}
	if err != nil {
		fmt.Println(err.Error())

	}
	item := Item{}

	err = dynamodbattribute.UnmarshalMap(result.Item, &item)
	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	}

	return item.Block

}

func DeleteBlockFromBuffer(name string, company string) bool {
	namecompany := name + "/" + company
	fmt.Println("hello i am deleting")
	sess, _ := session.NewSession(&aws.Config{

		Region: aws.String("us-east-1"),
	})
	svc := dynamodb.New(sess)

	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"Name": {
				S: aws.String(namecompany),
			},
		},
		TableName: aws.String("Buffer"),
	}

	_, err := svc.DeleteItem(input)
	if err != nil {
		fmt.Println("Got error calling DeleteItem")
		fmt.Println(err.Error())
		return false
	} else {
		return true
	}

	return false
}

func PutUserBytesIntoDB(ubytes []byte) bool {
	type Item struct {
		Name       string
		PrivateKey []byte
		PublicKey  []byte
	}
	sess, _ := session.NewSession(&aws.Config{

		Region: aws.String("us-east-1"),
	})
	svc := dynamodb.New(sess)
	item := Item{}
	json.Unmarshal(ubytes, &item)
	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		fmt.Println("Got error marshalling new item:")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// Create item in table Movies
	tableName := "Encryption"

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		fmt.Println("Got error calling PutItem:")
		fmt.Println(err.Error())
		return false

	}

	return true

}

func GetUserBytesFromDB(name string) []byte {
	type Item struct {
		Name       string
		PrivateKey []byte
		PublicKey  []byte
	}
	sess, _ := session.NewSession(&aws.Config{

		Region: aws.String("us-east-1"),
	})
	svc := dynamodb.New(sess)

	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("Encryption"),
		Key: map[string]*dynamodb.AttributeValue{
			"Name": {
				S: aws.String(name),
			},
		},
	})
	if err != nil {
		fmt.Println(err.Error())

	}
	if result.Item == nil {
		fmt.Println("Could not Find Private and Public Keys!\nExit")
		os.Exit(1)
	}
	item := Item{}

	err = dynamodbattribute.UnmarshalMap(result.Item, &item)
	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	}

	resp, _ := json.Marshal(item)
	return resp

}

func PutPublickeyIntoDB(publickeybytes []byte, name string) bool {
	type Item struct {
		Name      string
		PublicKey []byte
	}
	sess, _ := session.NewSession(&aws.Config{

		Region: aws.String("us-east-1"),
	})
	svc := dynamodb.New(sess)
	item := Item{
		Name:      name,
		PublicKey: publickeybytes,
	}
	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		fmt.Println("Got error marshalling new item:")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// Create item in table Movies
	tableName := "Publickeys"

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		fmt.Println("Got error calling PutItem:")
		fmt.Println(err.Error())
		return false
	}

	return true

}

func GetPublicKeyFromDB(name string) []byte {
	type Item struct {
		Name      string
		PublicKey []byte
	}
	sess, _ := session.NewSession(&aws.Config{

		Region: aws.String("us-east-1"),
	})
	svc := dynamodb.New(sess)

	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("Publickeys"),
		Key: map[string]*dynamodb.AttributeValue{
			"Name": {
				S: aws.String(name),
			},
		},
	})
	if err != nil {
		fmt.Println(err.Error())

	}
	if result.Item == nil {
		fmt.Println("Could not Find Public Keys!\nExit")
		os.Exit(1)
	}
	item := Item{}

	err = dynamodbattribute.UnmarshalMap(result.Item, &item)
	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	}

	return item.PublicKey

}
