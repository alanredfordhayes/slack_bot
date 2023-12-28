import { Api, StackContext } from "sst/constructs";

export function ApiStack({stack}: StackContext) {
    const api = new Api(stack, "API", {
        routes: {
            "POST /bot": "packages/functions/src/bot.go"
        }
    });
    stack.addOutputs({
        ApiEndpoint:api.url
    });
}