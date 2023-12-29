import { Api, StackContext, Function } from "sst/constructs";
import * as iam from "aws-cdk-lib/aws-iam";

export function ApiStack({ stack , app }: StackContext) {
    const bot_function = new Function(stack, "bot", {
        handler: "packages/functions/src/slack/bot.go",
        runtime: "go",
    });
    bot_function.attachPermissions([
        new iam.PolicyStatement({
            actions: ["secretsmanager:GetSecretValue", "secretsmanager:DescribeSecret"],
            effect: iam.Effect.ALLOW,
            resources: [`arn:aws:secretsmanager:${app.region}:${app.account}:secret:slack_bot-xf2YyW`],
        }),
        new iam.PolicyStatement({
            actions: ["dynamodb:PutItem"],
            effect: iam.Effect.ALLOW,
            resources: [
                `arn:aws:dynamodb:${app.region}:${app.account}:table/ahayes-slack-bot-slack_event_app_mention`,
                `arn:aws:dynamodb:${app.region}:${app.account}:table/ahayes-slack-bot-slack_event_message`,
                `arn:aws:dynamodb:${app.region}:${app.account}:table/ahayes-slack-bot-slack_event_app_not_found`,
            ],
        }),
      ])
    const api = new Api(stack, "API", {
        routes: {
            "POST /bot": {function: bot_function}
        }
    });
    stack.addOutputs({
        ApiEndpoint:api.url
    });
}