import { Table, StackContext, Function } from "sst/constructs";
import * as iam from "aws-cdk-lib/aws-iam";

export function DataStack({ stack }: StackContext) {
    const consumer_message_function = new Function(stack, "consumer_message_function", {
        handler: "packages/functions/src/dynamodb/message/consumer.go",
        runtime: "go",
    });
    const consumer_app_mention_function = new Function(stack, "consumer_app_mention_function", {
        handler: "packages/functions/src/dynamodb/app_mention/consumer.go",
        runtime: "go",
    });
    const slack_event_message_table = new Table( stack, "slack_event_message", {
        fields: {
            Channel: "string",
            User: "string",
            Text: "string",
            Ts: "string",
            Event_id: "string",
            Event_time: "number",
        },
        primaryIndex: { partitionKey: "Event_id", sortKey: "User" },
        globalIndexes: {
            userEventTimeIndex: { partitionKey: "User", sortKey: "Event_time" },
            userTextIndex: { partitionKey: "User", sortKey: "Text" },
            tsChannelIndex: { partitionKey: "Ts", sortKey: "Channel" },
        },
        stream: "new_image",
        consumers: {
            function: consumer_message_function
        }
    })
    const slack_event_app_mention_table = new Table( stack, "slack_event_app_mention", {
        fields: {
            Channel: "string",
            User: "string",
            Text: "string",
            Ts: "string",
            Event_id: "string",
            Event_time: "number",
        },
        primaryIndex: { partitionKey: "Event_id", sortKey: "User" },
        globalIndexes: {
            userEventTimeIndex: { partitionKey: "User", sortKey: "Event_time" },
            userTextIndex: { partitionKey: "User", sortKey: "Text" },
            tsChannelIndex: { partitionKey: "Ts", sortKey: "Channel" },
        },
        stream: "new_image",
        consumers: {
            function: consumer_app_mention_function
        }
    })
    const slack_event_app_not_found= new Table( stack, "slack_event_app_not_found", {
        fields: {
            Event_id: "string",
        },
        primaryIndex: { partitionKey: "Event_id" },
    })
    stack.addOutputs({
        slack_event_message_table_arn:slack_event_message_table.tableArn,
        slack_event_app_mention_table_arn:slack_event_app_mention_table.tableArn,
        slack_event_app_not_found_table_arn:slack_event_app_not_found.tableArn,
    });
}