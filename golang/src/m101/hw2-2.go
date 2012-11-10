package main

import (
	"fmt"
	"github.com/garyburd/go-mongo/mongo"
	"os"
)

func main() {
	conn, err := mongo.Dial("localhost")
	if err != nil {
		fmt.Println("Error trying to connect to database")
		os.Exit(1)
	}

	defer conn.Close()

	db := mongo.Database{conn, "students", mongo.DefaultLastErrorCmd}

	grades := db.C("grades")

	cursor, err := grades.Find(mongo.M{"type": "homework"}).Sort(mongo.D{{"student_id", 1}, {"score", 1}}).Cursor()
	if err != nil {
		fmt.Println("Error trying to read collection")
		os.Exit(1)
	}

	defer cursor.Close()

	initialStudentId := -1
	for cursor.HasNext() {
		var m map[string]interface{}
		cursor.Next(&m)

		studentId := m["student_id"]

		if initialStudentId != studentId {
			evict(m, &grades)
			initialStudentId = studentId.(int)
		}
	}

}

func evict(m map[string]interface{}, collection *mongo.Collection) {
	fmt.Println("Removing", m["_id"])
	collection.Remove(mongo.M{"_id": m["_id"]})
}
