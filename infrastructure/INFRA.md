# PrestigeAPI CDK JavaScript Project

This is a CDK Javascript project to define *infrastructure as code* for the **PrestigeAPI**.

## Usage Steps

1. Ensure you have **NodeJS** and **NPM** installed!
2. Install the AWS CLI
   1. [Installing, updating, and uninstalling the AWS CLI version 2
](https://docs.aws.amazon.com/cli/latest/userguide/install-cliv2.html)
3. Configure the AWS CLI and ensure you have billing enabled
   1. [Quick configuration with `aws configure`](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-quickstart.html)
4. Install the cdk cli `npm install -g aws-cdk`
5. Run `npm install` to install dependencies
6. Run `cdk deploy` to create the CloudFormation template and deploy it

## Useful commands

* `cdk deploy`           deploy this stack to your default AWS account/region
* `cdk diff`             compare deployed stack with current state
* `cdk synth`            emits the synthesized CloudFormation template

## Notes from CDK CLI Tool

The `cdk.json` file tells the CDK Toolkit how to execute your app. The build step is not required when using JavaScript.
