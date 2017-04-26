# go_dynamodb_search
Simple web application for searching dynamoDB database

## Provides

* form to provide search criteria DynamoDB table
* the result is the number of db records satisfying the criteria

## Technologies

* Golang (1.8)
* AWS-SDK library
        <ul>
        <li>go get -u github.com/aws/aws-sdk-go/aws...</li>
        <li>go get -u github.com/aws/aws-sdk-go/service/...</li>
        </ul>

* UI uses Bootstrap
        <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/css/bootstrap.min.css">

        <script src="https://ajax.googleapis.com/ajax/libs/jquery/1.12.4/jquery.min.js"></script>
        <script src="https://code.jquery.com/ui/1.12.1/jquery-ui.js"></script>

        <script src = "https://cdnjs.cloudflare.com/ajax/libs/angular-ui-bootstrap/2.5.0/ui-bootstrap.min.js"></script>
        <script src = "https://cdnjs.cloudflare.com/ajax/libs/bootstrap-datepicker/1.6.4/js/bootstrap-datepicker.min.js"></script>

## Creating DynamoDB tables using the aws CLI

sudo apt-get install aws-cli

aws dynamodb create-table --table-name loyalty-service-member --attribute-definitions AttributeName=memberID,AttributeType=N AttributeName=namespace,AttributeType=S AttributeName=date,AttributeType=S --key-schema AttributeName=memberId,KeyType=HASH AttributeName=namespace,KeyType=RANGE --provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5 

aws dynamodb create-table --table-name loyalty-service-order --attribute-definitions AttributeName=orderID,AttributeType=N AttributeName=namespace,AttributeType=S AttributeName=createDate,AttributeType=S --key-schema AttributeName=orderId,KeyType=HASH AttributeName=namespace,KeyType=RANGE --provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5 

## Note

If the DynamoDB endpoint is not provided as a command line parameter, a locally running DynamoDB is considered by default.

The AWS credentials and the AWS region must be provided even for the local DynamoDB. The AWS_ACCESS_KEY_ID and the region are used for naming the database. The AWS_SECRET_ACCESS_KEY is ignored but must be provided.
  
