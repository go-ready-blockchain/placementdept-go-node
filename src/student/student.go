package student

import (
	"encoding/json"
	"fmt"

	"github.com/jugalw13/placementdept-go-node/security"
)

type Student struct {
	Usn       string `json:"Usn"`
	Branch    string `json:"Branch"`
	Name      string `json:"Name"`
	Gender    string `json:"Gender"`
	Dob       string `json:"Dob"`
	Perc10th  string `json:"Perc10th"`
	Perc12th  string `json:"Perc12th"`
	Cgpa      string `json:"Cgpa"`
	Backlog   bool   `json:"Backlog"`
	Email     string `json:"Email"`
	Mobile    string `json:"Mobile"`
	StarOffer bool   `json:"StarOffer"`
}

func EncodeToBytes(s *Student) []byte {
	sbytes, _ := json.Marshal(s)
	return sbytes
}

func DecodeToStruct(sbytes []byte) *Student {
	result := &Student{}
	json.Unmarshal(sbytes, &result)
	return result
}

func EnterStudentData(usn string, branch string, name string, gender string, dob string, perc10th string, perc12th string, cgpa string, backlog bool, email string, mobile string, staroffer bool) *Student {
	s := &Student{usn, branch, name, gender, dob, perc10th, perc12th, cgpa, backlog, email, mobile, staroffer}
	return s
}

func PrintStudentData(s *Student) {
	fmt.Println("USN: ", s.Usn)
	fmt.Println("Branch: ", s.Branch)
	fmt.Println("Name: ", s.Name)
	fmt.Println("Gender: ", s.Gender)
	fmt.Println("DOB: ", s.Dob)
	fmt.Println("10th Percentage: ", s.Perc10th)
	fmt.Println("12th Percentage: ", s.Perc12th)
	fmt.Println("CGPA: ", s.Cgpa)
	fmt.Println("Backlog: ", s.Backlog)
	fmt.Println("Email: ", s.Email)
	fmt.Println("Mobile: ", s.Mobile)
	fmt.Println("Star Offer: ", s.StarOffer)
}

func GenerateStudentSignature(name string, data []byte) []byte {

	student := security.GetUserFromDB(name)

	signature := security.PSSSignature(data, student.PrivateKey)

	return signature

}

func EncryptStudentData(receiver string, data []byte) []byte {

	receiverPublicKey := security.GetPublicKeyFromDB(receiver)
	studentdata := security.EncryptMessage(data, receiverPublicKey)

	return studentdata

}

// func TestStudent() {

// 	//Initialize the Student struct for a new student
// 	s := EnterStudentData("usn", "branch", "name", "gender", "dob", "perc10th", "perc12th", "cgpa", false, "email", "mobile", true)

// 	fmt.Println(s) //Print Struct

// 	//Encode it to Bytes
// 	b := EncodeToBytes(s)
// 	//fmt.Println(b)

// 	//Store it in the block
// 	block := &TestBlock{b}

// 	//fmt.Println(block)

// 	//Fetch it from Block and Decode it
// 	nb := block.Test
// 	nv := DecodeToStruct(nb)

// 	fmt.Println(nv) //Print Struct
// 	PrintStudentData(nv)

// }

// func main() {
// 	TestStudent()
// }
