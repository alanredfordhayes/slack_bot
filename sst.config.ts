import { SSTConfig } from "sst";

export default {
  config(_input) {
    return {
      name: "slack-bot",
      region: "us-east-1",
    };
  },
  stacks(app) {}
} satisfies SSTConfig;
