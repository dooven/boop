package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/atotto/clipboard"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/aws/aws-sdk-go/service/rds/rdsutils"
)

type rdsClient struct {
	*rds.RDS
}

func (r *rdsClient) getData(tempPath string) ([]*rds.DBInstance, error) {
	if _, err := os.Stat(tempPath); os.IsNotExist(err) {
		input := &rds.DescribeDBInstancesInput{}
		result, err := r.DescribeDBInstances(input)

		if err != nil {
			return nil, err
		}

		jsonData, err := json.MarshalIndent(result, "", "    ")

		if err != nil {
			return nil, err
		}

		if err := os.Mkdir(path.Dir(tempPath), os.ModePerm); err != nil {
			return nil, err
		}

		jsonFile, err := os.Create(tempPath)

		if err != nil {
			log.Fatal(err)
		}

		defer jsonFile.Close()

		jsonFile.Write(jsonData)
		jsonFile.Close()

		return result.DBInstances, nil
	}

	content, err := ioutil.ReadFile(tempPath)

	if err != nil {
		log.Fatal(err)
	}

	var tempData *rds.DescribeDBInstancesOutput

	if err := json.Unmarshal(content, &tempData); err != nil {
		return nil, err
	}

	return tempData.DBInstances, nil

}

func (r *rdsClient) generateToken(endpoint string, dbUser string) (string, error) {
	config := r.Client.Config

	if authToken, err := rdsutils.BuildAuthToken(
		endpoint,
		*config.Region,
		dbUser,
		config.Credentials); err == nil {

		if err := clipboard.WriteAll(authToken); err != nil {
			return "", err
		}

		return authToken, nil
	} else {
		return "", err
	}
}

func newRdsClient(svc *rds.RDS) *rdsClient {
	return &rdsClient{svc}
}
