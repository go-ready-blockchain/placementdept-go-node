package blockchain

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/jugalw13/placementdept-go-node/security"
	"github.com/jugalw13/placementdept-go-node/utils"
)

type Verification struct {
	Verified      string               `json:"Verified"`
	Timestamps    map[string]time.Time `json:Timestamps"`
	Verifications map[string]string    `json:Verifications"`
}

func AcademicDeptVerification(name string, company string) bool {

	block := GetBlockFromBuffer(name, company)

	studentdata, dflag := security.DecryptMessage(block.StudentData, security.GetUserFromDB("AcademicDept").PrivateKey)
	if dflag == false {
		fmt.Println("Decrytion of Message Failed")
		utils.DeleteBlockFromBuffer(name, company)
		return false
	}

	sflag := security.VerifyPSSSignature(security.GetPublicKeyFromDB(name), block.Signature, studentdata)
	if sflag == false {
		fmt.Println("Signature Verification Failed, Authentication Failed")
		utils.DeleteBlockFromBuffer(name, company)
		return false
	}

	block.StudentData = studentdata
	v, vflag := ValidationByAcademicDept(DecodeToStruct(block.Verification), block)
	if vflag == false {
		fmt.Println("Validation By Academic Dept Failed")
		utils.DeleteBlockFromBuffer(name, company)
		return false
	}

	//fmt.Println("Verified: ", v.Verified, "\n", "Verifications: ", v.Verifications, "\n", "Timestamps: ", v.Timestamps, "\n") //Print Struct

	block.Verification = EncodeToBytes(v)
	block.StudentData = security.EncryptMessage(studentdata, security.GetPublicKeyFromDB("PlacementDept"))
	block.Signature = security.PSSSignature(studentdata, security.GetUserFromDB("AcademicDept").PrivateKey)

	PutBlockIntoBuffer(block, name, company)
	return true
}

func PlacementDeptVerification(name string, company string) bool {

	block := GetBlockFromBuffer(name, company)

	studentdata, dflag := security.DecryptMessage(block.StudentData, security.GetUserFromDB("PlacementDept").PrivateKey)
	if dflag == false {
		fmt.Println("Decrytion of Message Failed")
		utils.DeleteBlockFromBuffer(name, company)
		return false
	}

	sflag := security.VerifyPSSSignature(security.GetPublicKeyFromDB("AcademicDept"), block.Signature, studentdata)
	if sflag == false {
		fmt.Println("Signature Verification Failed, Authentication Failed")
		utils.DeleteBlockFromBuffer(name, company)
		return false
	}

	block.StudentData = studentdata
	v, vflag := ValidationByPlacementDept(DecodeToStruct(block.Verification), block)
	if vflag == false {
		fmt.Println("Validation By Placement Dept Failed")
		utils.DeleteBlockFromBuffer(name, company)
		return false
	}

	//fmt.Println("Verified: ", v.Verified, "\n", "Verifications: ", v.Verifications, "\n", "Timestamps: ", v.Timestamps, "\n") //Print Struct

	block.Verification = EncodeToBytes(v)
	block.StudentData = security.EncryptMessage(studentdata, security.GetPublicKeyFromDB(company))
	block.Signature = security.PSSSignature(studentdata, security.GetUserFromDB("PlacementDept").PrivateKey)

	PutBlockIntoBuffer(block, name, company)

	//Add block to blockchain as a transaction

	AddBlock(block)
	fmt.Println("Validation successfully completed.\nCompany can retrieve the data")
	return true
}

func InitVerification() *Verification {
	v := &Verification{"", make(map[string]time.Time), make(map[string]string)}

	v.Verified = "Not Done Yet"

	v.Verifications["Academic Dept"] = "Not Done Yet"
	v.Verifications["Placement Dept"] = "Not Done Yet"

	v.Timestamps["Created At"] = time.Now()
	v.Timestamps["Academic Dept"] = time.Time{}  //zero value
	v.Timestamps["Placement Dept"] = time.Time{} //zero value

	return v
}

func CheckIfVerifiedByAll(v *Verification) bool {

	if v.Verified == "True" {
		return true
	} else if v.Verified == "Not Done Yet" {
		return false
	} else {
		fmt.Println("Error")
	}
	return false
}

func CheckIfVerifiedByAcademicDept(v *Verification) bool {

	if v.Verifications["Academic Dept"] == "True" {
		return true
	} else if v.Verifications["Academic Dept"] == "Not Done Yet" {
		return false
	} else {
		fmt.Println("Error")
	}
	return false
}

func ValidationByAcademicDept(v *Verification, block *Block) (*Verification, bool) {

	if CheckIfVerifiedByAll(v) {
		fmt.Println("Already Verified")
		return v, true
	}

	//TODO: validateBlockchain()
	pow := NewProof(block)
	vflag := pow.Validate()
	if vflag == false {
		fmt.Println("Placement Dept Validation of Proof Of Work Failed")
		return v, false
	}
	fmt.Println("Academic Dept Successfully completed Validation of Proof Of Work!")

	tflag := ProofOfElapsedTime(v.Timestamps["Created At"])
	if tflag == false {
		fmt.Println("Proof Of Elapsed Time failed")
		return v, false
	}

	v.Verifications["Academic Dept"] = "True"
	v.Timestamps["Academic Dept"] = time.Now()

	return v, true
}

func ValidationByPlacementDept(v *Verification, block *Block) (*Verification, bool) {

	if CheckIfVerifiedByAll(v) {
		fmt.Println("Already Verified")
		return v, true
	}

	if CheckIfVerifiedByAcademicDept(v) == false {
		fmt.Println("Verification Not Yet Done by Academic Dept")
		return v, false
	}

	//TODO: validateBlockchain()
	pow := NewProof(block)
	vflag := pow.Validate()
	if vflag == false {
		fmt.Println("Placement Dept Validation of Proof Of Work Failed")
		return v, false
	}
	fmt.Println("Placement Dept Successfully completed Validation Proof Of Work!")

	tflag := ProofOfElapsedTime(v.Timestamps["Created At"])
	if tflag == false {
		fmt.Println("Proof Of Elapsed Time failed")
		return v, false
	}

	v.Verifications["Placement Dept"] = "True"
	v.Timestamps["Placement Dept"] = time.Now()
	v.Verified = "True"

	return v, true
}

func EncodeToBytes(v *Verification) []byte {
	vbytes, _ := json.Marshal(v)
	return vbytes
}

func DecodeToStruct(vbytes []byte) *Verification {
	result := &Verification{}
	json.Unmarshal(vbytes, &result)
	return result
}

func ProofOfElapsedTime(creation time.Time) bool {
	limit := creation.Add(24 * time.Hour) //24 Hr Limit
	now := time.Now()
	return now.After(creation) && now.Before(limit)
}

// func TestVerify() {

// 	//Initialize the Verification struct for a new student
// 	v := InitVerification()

// 	fmt.Println("Verified: ", v.Verified, "\n", "Verifications: ", v.Verifications, "\n", "Timestamps: ", v.Timestamps, "\n") //Print Struct

// 	//Function Call For Academic Dept to Verify
// 	v, flag := ValidationByAcademicDept(v)

// 	fmt.Println("Verified: ", v.Verified, "\n", "Verifications: ", v.Verifications, "\n", "Timestamps: ", v.Timestamps, "\n") //Print Struct

// 	//Function Call For Placement Dept to Verify
// 	v, flag = ValidationByPlacementDept(v)

// 	fmt.Println("Verified: ", v.Verified, "\n", "Verifications: ", v.Verifications, "\n", "Timestamps: ", v.Timestamps, "\n") //Print Struct

// 	fmt.Println(flag)
// 	//Encode it to Bytes
// 	b := EncodeToBytes(v)
// 	//fmt.Println(b)

// 	//Store it in the block
// 	block := &TestBlock{b}

// 	//fmt.Println(block)

// 	//Fetch it from Block and Decode it
// 	nb := block.Verify
// 	nv := DecodeToStruct(nb)

// 	fmt.Println("Verified: ", nv.Verified, "\n", nv.Timestamps, "\n", nv.Verifications) //Print Struct

// }

// func main() {

// }
