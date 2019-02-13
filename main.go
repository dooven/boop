package main

import (
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/dooven/boop/config"
	"log"
	"os"
	"path"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/dooven/boop/rdsHelper"
	"github.com/manifoldco/promptui"
)

var (
	PORT             = "3006"
	regionTemplates  *promptui.SelectTemplates
	addressTemplates *promptui.SelectTemplates
)

func init() {
	regionTemplates = &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   "- {{ .Name | cyan }} ({{ .Region | red }})",
		Inactive: "  {{ .Name | cyan }} ({{ .Region | red }})",
		Selected: "* {{ .Name | red | cyan }}",
	}

	addressTemplates = &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   "- {{ . | cyan }}",
		Inactive: "  {{ . | cyan }}",
		Selected: "* {{ . | red | cyan }}",
	}
}

func main() {

	storedConfigs, err := config.GetOrWriteDefaults()

	if err != nil {
		log.Fatal(err)
	}

	regionOptions := storedConfigs.Regions

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

	rdsClient := rdsHelper.NewRdsClient(rds.New(sess))

	var dbInstanceAddress []string

	instances, err := rdsClient.GetRDSInstances(tempPath)

	if err != nil {
		log.Fatal(err)
	}

	for _, dbInstance := range instances {
		address := dbInstance.Endpoint.Address
		dbInstanceAddress = append(dbInstanceAddress, *address)
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

	authToken, err := rdsClient.GenerateToken(endpoint, "test-user")

	if err != nil {
		log.Fatal(err)
	}

	if err := clipboard.WriteAll(authToken); err != nil {
		log.Fatal(err)
	}

	log.Println(authToken)
}
