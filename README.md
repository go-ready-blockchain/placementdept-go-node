# PLACEMENT DEPT NODE

## Blockchain Implementation in GoLang For Placement System

## The Consensus Algorithm implemented in Blockchain System is a combination of Proof Of Work and Proof Of Elapsed Time


### Run `go run src/main.go` to Start the Server and listen on localhost:8084

### Usage :

#### To Print Usage
####    Make POST request to `/usage`

#### Part of the Pipeline - 

#### To Send Email to Eligible Students based on Eligibility Criteria
####    Make POST request to `/send` with body -
```json
{
	"company" : "JPMC",
	"backlog" : "",
	"starOffer" : "",
	"branch" : ["CSE","ISE"],
	"gender" : "",
	"cgpaCond" : "GreaterThan",
	"cgpa" : "2",
	"perc10thCond" : "GreaterThan",
	"perc10th" : "10",
	"perc12thCond" : "GreaterThan",
	"perc12th" : "10"
}
```

#### To Run Verification by Placement Department
####    Make POST request to `/verify-PlacementDept` with body -
```json
{
	"name":"1MS16CS034",
    "company": "JPMC"
  
}
```





