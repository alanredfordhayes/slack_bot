import { StackContext, Api, Function } from "sst/constructs";
export function ApiStack( { stack }: StackContext){

    const events_function = new Function(stack, "events_function", {
        runtime: "go",
        handler: "packages/functions/lambda/events.go",
        timeout: 30,
        memorySize: 1536,
    })

    const interactivity_function = new Function(stack, "interactivity_function", {
        runtime: "go",
        handler: "packages/functions/lambda/interactivity.go",
        timeout: 30,
        memorySize: 1536,
    })

    const api = new Api(stack, "api", {
        routes: {
            "POST /events": events_function,
            "POST /interactivity": interactivity_function,
        },
    });

    stack.addOutputs({
        ApiEndpoint: api.url
    })
}