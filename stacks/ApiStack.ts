import { Api, StackContext } from "sst/constructs";

export function ApiStack({stack}: StackContext) {
    const api = new Api(stack, "API", {
        routes: {
            "POST /bot": {
                function: {
                    handler: "packages/functions/src/bot.go",
                    runtime: "go",
                    copyFiles: [{ from: "packages/functions/src/bin/http.go" }]
                }
            }
        }
    });
    stack.addOutputs({
        ApiEndpoint:api.url
    });
}