const cdk = require('@aws-cdk/core');
const prestige_api = require('./prestige_api');

class PrestigeAPIStack extends cdk.Stack {
  /**
   *
   * @param {cdk.Construct} scope
   * @param {string} id
   * @param {cdk.StackProps=} props
   */
  constructor(scope, id, props) {
    super(scope, id, props);

    // The code that defines your stack goes here
    new prestige_api.PrestigeAPI(this, 'Prestige');
  }
}

module.exports = { PrestigeAPIStack }
