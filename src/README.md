# Blockchain Implementation in GoLang For Placement System

## Proof of Work is the Consensus Algorithm used in our BlockChain Implementation

### Run `go run main.go` to run the app, run `go build main.go` to build an executable file.

### Usage :

#### To Create a New BlockChain    
####    `go run main.go createBlockChain`

#### To Add a New Student
####    `go run main.go student -usn "USN" -branch "BRANCH" -name "NAME" -gender "GENDER" -dob "DOB" -cgpa "CGPA" -perc10th "PERC10TH" -perc12th "PERC12TH"  -backlog=false -email "EMAIL" -mobile "MOBILE" -staroffer=true`

#### To Add a new Companu    
####    `go run main.go company -name COMPANY`

#### To Company to Request Student Data
####    `go run main.go request -company "COMPANY" -student "USN"`

#### To Run Verification by Academic Department
####    `go run main.go verify-AcademicDept -student "USN" -company "COMPANY"`

#### To Run Verification by Placement Department
####    `go run main.go verify-PlacementDept -student "USN" -company "COMPANY"`

#### To Retrieve the data for the Company
####    `go run main.go companyRetrieveData -student "USN" -company "COMPANY"`

#### To Print the entire BlockChain
####    `go run main.go print`


