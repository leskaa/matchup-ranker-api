#!/usr/bin/env node

const cdk = require('@aws-cdk/core');
const { PrestigeAPIStack } = require('../lib/prestige_api-stack');

const app = new cdk.App();
new PrestigeAPIStack(app, 'PrestigeAPIStack');
