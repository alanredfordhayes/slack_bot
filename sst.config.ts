import { SSTConfig } from "sst";
import { ApiStack } from "./stacks/ApiStack";
import { NextStack }  from "./stacks/NextStack";

export default {
    config(_input) {
        return {
            name: "slack-bot",
            region: "us-east-1",
        };
    },
    stacks(app) {
        app.setDefaultFunctionProps({
            runtime: "go",
        });
        app
            .stack(NextStack)
            .stack(ApiStack);
    },
} satisfies SSTConfig;
