import { NextjsSite, StackContext } from "sst/constructs";

export function NextStack({ stack }: StackContext) {
    const site = new NextjsSite(stack, "site");
    stack.addOutputs({
        SiteUrl: site.url,
    });
    stack.addOutputs({
        NextStackEndpoint:site.url
    });
}
