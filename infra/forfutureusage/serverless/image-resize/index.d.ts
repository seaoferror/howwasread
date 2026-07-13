import type { S3Event } from "aws-lambda";
export declare const handler: (event: S3Event) => Promise<{
    statusCode: number;
    body: string;
    desc?: never;
} | {
    statusCode: number;
    body: string;
    desc: string | undefined;
}>;
//# sourceMappingURL=index.d.ts.map