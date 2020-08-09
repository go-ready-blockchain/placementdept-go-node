package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"crypto/tls"
	"os"
	

	"github.com/go-ready-blockchain/blockchain-go-core/blockchain"
	
	"github.com/go-ready-blockchain/blockchain-go-core/notification"
)

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("Make POST request to /send \tTo Send Email to Eligible Students based on Eligibility Criteria")
	fmt.Println("Make POST request to /verify-PlacementDept \tPlacementDept Verifies Student's data")
}
func newnotification(link string, companyName string, Backlog string, StarOffer string, Branch []string, Gender string, CgpaCond string, Cgpa string, Perc10thCond string, Perc10th string, Perc12thCond string, Perc12th string) bool {

	emailitems := notification.ApplyFilter(Backlog, StarOffer, Branch, Gender, CgpaCond, Cgpa, Perc10thCond, Perc10th, Perc12thCond, Perc12th)

	for _, emailitem := range emailitems {
		email := emailitem.Email
		usn := emailitem.Usn
		name := emailitem.Name
		fmt.Println(email, usn, name)
		acceptlink := link + "/handlerequest?" + "approval=" + "true" + "&company=" + companyName + "&name=" + usn
		rejectlink := link + "/handlerequest?" + "approval=" + "false" + "&company=" + companyName + "&name=" + usn
		sendEmail("Request for Student Data", name, companyName, acceptlink, rejectlink, email)

	}
	return true
}

func sendEmail(subject string, studentName string, companyName string, acceptlink string, rejectlink string, To string){

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	req, err := http.NewRequest("GET", "https://us-central1-go-ready-blockchain.cloudfunctions.net/sendMail", nil)
    if err != nil {
        fmt.Println(err)
        
    }

    q := req.URL.Query()
    q.Add("dest", To)
	q.Add("subject", subject)
	q.Add("studentName", studentName)
	q.Add("companyName",companyName)
	q.Add("acceptlink",acceptlink)
	q.Add("rejectlink", rejectlink)
	req.URL.RawQuery = q.Encode()
	

    fmt.Println(req.URL.String())

	Url := req.URL.String()
	resp, err := http.Get(Url)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(body))
	


}
func sendNotification(w http.ResponseWriter, r *http.Request) {
	//name := time.Now().String()
	//logger.FileName = "Placement Send Notification " + name + ".log"
	//logger.NodeName = "Placement Node"
	//logger.CreateFile()

	type jsonBody struct {
		Company      string   `json:"company"`
		Backlog      string   `json:"backlog"`
		StarOffer    string   `json:"starOffer"`
		Branch       []string `json:"branch"`
		Gender       string   `json:"gender"`
		CgpaCond     string   `json:"cgpaCond"`
		Cgpa         string   `json:"cgpa"`
		Perc10thCond string   `json:"perc10thCond"`
		Perc10th     string   `json:"perc10th"`
		Perc12thCond string   `json:"perc12thCond"`
		Perc12th     string   `json:"perc12th"`
	}
	decoder := json.NewDecoder(r.Body)
	var b jsonBody
	if err := decoder.Decode(&b); err != nil {
		log.Fatal(err)
	}

	message := ""
	Apiurl := os.Getenv("STUDENT_URL")
	flag := newnotification(Apiurl, b.Company, b.Backlog, b.StarOffer, b.Branch, b.Gender, b.CgpaCond, b.Cgpa, b.Perc10thCond, b.Perc10th, b.Perc12thCond, b.Perc12th)

	if flag == true {
		message = "Notification sent successfully to Students!"
	} else {
		message = "Sending Notification to Student Failed!"
	}

	//logger.UploadToS3Bucket(//logger.NodeName)

	//logger.DeleteFile()

	

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(message))
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
	//name := time.Now().String()
	//logger.FileName = "Placement Verify " + name + ".log"
	//logger.NodeName = "Placement Node"
	//logger.CreateFile()

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

	//logger.UploadToS3Bucket(//logger.NodeName)

	//logger.DeleteFile()

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
	Apiurl := os.Getenv("COMPANY_URL")
	Apiurl = Apiurl + "/companyRetrieveData"
	resp, err := http.Post(Apiurl,
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
	http.HandleFunc("/send", sendNotification)
	http.HandleFunc("/verify-PlacementDept", callverificationByPlacementDept)
	http.HandleFunc("/usage", callprintUsage)
	fmt.Printf("Server listening on localhost:%s\n", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
