package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/jugalw13/placementdept-go-node/blockchain"
)

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("verify-PlacementDept -student USN \tPlacementDept Verifies Student's data")
}

func verificationByPlacementDept(name string, company string) bool {
	fmt.Println("\nStarting Verification by Placement Dept\n")
	flag := blockchain.PlacementDeptVerification(name, company)
	if flag == true {
		fmt.Println("Verification by Placement Dept Successfully completed!")
		return true
	} else {
		fmt.Println("Verification by Placement Dept Failed!")
		return false
	}
}

func callverificationByPlacementDept(w http.ResponseWriter, r *http.Request) {
	type jsonBody struct {
		Name    string `json:"name"`
		Company string `json:"company"`
	}
	decoder := json.NewDecoder(r.Body)
	var b jsonBody
	if err := decoder.Decode(&b); err != nil {
		log.Fatal(err)
	}
	message := ""
	flag := verificationByPlacementDept(b.Name, b.Company)
	if flag == true {
		message = "Verification by Placement Dept Successfully completed!"
	} else {
		message = "Verification by Placement Dept Failed!"
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(message))

	fmt.Println("\n\nSending Notification to Company to retrieve Student's Data\n\n")
	callCompanyRetrieveData(b.Name, b.Company)
}

func callCompanyRetrieveData(name string, company string) {
	reqBody, err := json.Marshal(map[string]string{
		"name":    name,
		"company": company,
	})
	if err != nil {
		print(err)
	}
	resp, err := http.Post("http://company-cluster2-jen-ci.devtools-dev.ext.devshift.net/companyRetrieveData",
		"application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		print(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		print(err)
	}
	fmt.Println(string(body))
}

func callprintUsage(w http.ResponseWriter, r *http.Request) {

	printUsage()

	w.Header().Set("Content-Type", "application/json")
	message := "Printed Usage!!"
	w.Write([]byte(message))
}

func main() {
	port := "8080"
	http.HandleFunc("/verify-PlacementDept", callverificationByPlacementDept)
	http.HandleFunc("/usage", callprintUsage)
	fmt.Printf("Server listening on localhost:%s\n", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
