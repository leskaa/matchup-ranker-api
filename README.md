# prestige-api

This is a serverless API for the [prestige app](https://github.com/Brojaga/prestige-app-frontend). It was built with AWS Lambda and Golang. The project includes a CDK sub-project to define the infrastructure. Currently the API only supports ranking technology companies. In the future the API will support ranking computer science programs at American universities. The API ranks entities by winrate from individual 1v1 matchups.

## File Structure

* **/services** - Contains each of the Golang Lambda functions
* **/libs** - Contains code shared by multiple functions (nothing currently)
* **/infrastructure** - Contains an AWS CDK project that defines *infrastructure as code*

## Deploying to AWS

1. Ensure you have **NodeJS** and **NPM** installed!
2. Install the AWS CLI
   1. [Installing, updating, and uninstalling the AWS CLI version 2
](https://docs.aws.amazon.com/cli/latest/userguide/install-cliv2.html)
3. Configure the AWS CLI and ensure you have billing enabled
   1. [Quick configuration with `aws configure`](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-quickstart.html)
4. Navigate to the **/infrastructure** directory `cd infrastructure`
5. Install the cdk cli `npm install -g aws-cdk`
6. Run `npm install` to install dependencies
7. Run `cdk deploy` to create the CloudFormation template and deploy it

Check out the [INFRA.md](https://github.com/leskaa/matchup-ranker-api/infrastructure/INFRA.md) in the CDK sub-project for more information about the infrastructure and deployment.

## API Endpoints

* **GET /rankings** - Get a list of all companies sorted by winrate
* **GET /ranking?company=name** - Get a company by name (Not included in infrastructure)
* **GET /matchup** - Get a randomly generated matchup with verification code
* **POST /matchup** - Send the result of a matchup to update the rankings
  * Requires JSON body that includes `verificationCode` and `winner` with value `1` or `2`
