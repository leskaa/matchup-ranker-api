const core = require("@aws-cdk/core");
const apigateway = require("@aws-cdk/aws-apigateway");
const lambda = require("@aws-cdk/aws-lambda");
const dynamodb = require("@aws-cdk/aws-dynamodb");
const path = require("path");

class PrestigeAPI extends core.Construct {
  constructor(scope, id) {
    super(scope, id);

    // Create dynamodb tables
    const companyTable = new dynamodb.Table(this, "Companies", {
      tableName: 'prestige-companies',
      partitionKey: { name: 'Company', type: dynamodb.AttributeType.STRING },
      removalPolicy: core.RemovalPolicy.RETAIN,
      billingMode: dynamodb.BillingMode.PROVISIONED,
      readCapacity: 3,
      writeCapacity: 3
    })

    const matchupTable = new dynamodb.Table(this, "Matchups", {
      tableName: 'prestige-matchups',
      partitionKey: { name: 'VerificationCode', type: dynamodb.AttributeType.STRING },
      removalPolicy: core.RemovalPolicy.RETAIN,
      billingMode: dynamodb.BillingMode.PROVISIONED,
      readCapacity: 3,
      writeCapacity: 3
    })

    // Create lambda functions
    const rankingsFunction = new lambda.Function(this, "Rankings", {
      runtime: lambda.Runtime.GO_1_X,
      code: lambda.Code.fromAsset(path.join(__dirname, "..", "..", "services", "rankings-api")),
      handler: "main"
    });

    const matchupFunction = new lambda.Function(this, "Matchup", {
      runtime: lambda.Runtime.GO_1_X,
      code: lambda.Code.fromAsset(path.join(__dirname, "..", "..", "services", "matchup-api")),
      handler: "main"
    });

    // Give functions access to dynamodb tables
    companyTable.grantReadWriteData(rankingsFunction);
    companyTable.grantReadWriteData(matchupFunction);
    matchupTable.grantReadWriteData(matchupFunction);

    // Create API Gateway REST API
    const api = new apigateway.RestApi(this, "Api", {
      restApiName: "Matchup Ranker API",
      description: "This API handles rankings through matchups.",
      defaultCorsPreflightOptions: {
        allowOrigins: apigateway.Cors.ALL_ORIGINS,
        allowMethods: apigateway.Cors.ALL_METHODS
      }
    });

    // Add integrations between the REST API and functions
    const rankingsIntegration = new apigateway.LambdaIntegration(rankingsFunction);
    const matchupIntegration = new apigateway.LambdaIntegration(matchupFunction);

    // Create routes in REST API for each service
    const rankings = api.root.addResource('rankings');
    const matchup = api.root.addResource('matchup');

    // Add HTTP methods for each integration
    rankings.addMethod("ANY", rankingsIntegration);
    matchup.addMethod("ANY", matchupIntegration);
  }
}

module.exports = { PrestigeAPI }