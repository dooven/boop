# Boop

## Features
* Generates RDS auth token
* Automatically copies generated token to clipboard



## Installing 
1. Run `go get github.com/dooven/boop `
    * Make sure you have the `$GOHOME/bin` added to your `PATH`
2. Run `boop` 
    * Optionally, you can pass in the `AWS_PROFILE` e.g. `AWS_PROFILE=personal boop`
    
## Demo
1. Pick the region  ![region](./screenshots/1.png)
2. Pick the RDS endpoint  ![endpoint](./screenshots/2.png)
3. Pick the RDS User  ![rds_user](./screenshots/3.png)
4. The token will be generated, and will be copied to your clipboard

