package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/manifoldco/promptui"
)

var (
	PORT = "3006"
)

type regionOption struct {
	Name   string
	Region string
}

func main() {

	regionOptions := []regionOption{
		{Name: "test-region", Region: endpoints.EuWest1RegionID},
	}

	regionTemplates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   "- {{ .Name | cyan }} ({{ .Region | red }})",
		Inactive: "  {{ .Name | cyan }} ({{ .Region | red }})",
		Selected: "* {{ .Name | red | cyan }}",
	}

	regionSearcher := func(input string, index int) bool {
		region := regionOptions[index]
		name := region.Name
		input = strings.Replace(strings.ToLower(input), " ", "", -1)

		return strings.Contains(name, input)
	}

	regionPrompt := promptui.Select{
		Label:     "Region",
		Items:     regionOptions,
		Templates: regionTemplates,
		Size:      4,
		Searcher:  regionSearcher,
	}

	selectionRegionIndex, _, selectRegionErr := regionPrompt.Run()

	if selectRegionErr != nil {
		log.Fatal(selectRegionErr)
	}

	selectedRegion := regionOptions[selectionRegionIndex].Region

	fmt.Printf("Region: %s\n\n", selectedRegion)

	tempPath := path.Join(
		os.TempDir(),
		"boop-cache",
		fmt.Sprintf("%s.json", selectedRegion),
	)

	sess, err := session.NewSessionWithOptions(
		session.Options{
			Config: aws.Config{
				Region: aws.String(selectedRegion),
			},
			AssumeRoleTokenProvider: stscreds.StdinTokenProvider,
			SharedConfigState:       session.SharedConfigEnable,
		},
	)

	if err != nil {
		log.Fatal(err)
	}

	rdsClient := newRdsClient(rds.New(sess))

	var dbInstanceAddress []string

	instances, err := rdsClient.getData(tempPath)

	if err != nil {
		log.Fatal(err)
	}

	for _, dbInstance := range instances {
		address := dbInstance.Endpoint.Address
		dbInstanceAddress = append(dbInstanceAddress, *address)
	}

	addressTemplates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   "- {{ . | cyan }}",
		Inactive: "  {{ . | cyan }}",
		Selected: "* {{ . | red | cyan }}",
	}

	addressSearcher := func(input string, index int) bool {
		address := dbInstanceAddress[index]
		input = strings.Replace(strings.ToLower(input), " ", "", -1)

		return strings.Contains(address, input)
	}

	addressPrompt := promptui.Select{
		Label:     "Endpoint",
		Items:     dbInstanceAddress,
		Templates: addressTemplates,
		Size:      10,
		Searcher:  addressSearcher,
	}

	selectAddressIndex, _, selectAddressErr := addressPrompt.Run()

	if selectAddressErr != nil {
		log.Fatal(selectAddressErr)
	}

	endpoint := fmt.Sprintf("%s:%s/", dbInstanceAddress[selectAddressIndex], PORT)

	if authToken, err := rdsClient.generateToken(endpoint, "test-user"); err == nil {
		log.Println(authToken)
	} else {
		log.Fatal(err)
	}
}
