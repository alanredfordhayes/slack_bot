import { SSTConfig } from "sst";
import { ApiStack } from "./stacks/ApiStack";
import { NextStack }  from "./stacks/NextStack";
import { DataStack } from "./stacks/DataStack";

export default {
    config(_input) {
        return {
            name: "slack-bot",
            region: "us-east-1",
        };
    },
    stacks(app) {
        if (app.stage !== "prod") {
            app.setDefaultRemovalPolicy("destroy");
        }
        app
            .stack(NextStack)
            .stack(ApiStack)
            .stack(DataStack)
    },
} satisfies SSTConfig;
