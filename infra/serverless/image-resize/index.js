import { GetObjectCommand, PutObjectCommand, S3Client } from "@aws-sdk/client-s3";
import sharp from "sharp";
const S3 = new S3Client({});
const DESC_BUCKET = process.env.DESC_BUCKET;
const RESIZED_WIDTH = 800;
const ALLOWED_CONTENT_TYPE = {
    "image/jpg": true,
    "image/jpeg": true,
    "image/heic": true,
    "image/png": true,
    "image/webp": true,
    "image/gif": true,
    "image/avif": true
};
export const handler = async (event) => {
    const record = event.Records?.[0];
    if (!record) {
        console.log("No records found in the event.");
        return {
            statusCode: 400,
            body: "No S3 record provided"
        };
    }
    const { s3 } = record;
    const src = s3.bucket.name;
    const key = s3.object.key;
    const destKey = key.replace("image/", "resized-image/");
    try {
        const { Body, ContentType } = await S3.send(new GetObjectCommand({
            Bucket: src,
            Key: key
        }));
        if (!Body) {
            throw new Error("Empty body received from S3");
        }
        if (!ALLOWED_CONTENT_TYPE[ContentType ?? ""]) {
            throw new Error("not allowed content type");
        }
        let image = await Body.transformToByteArray();
        if (ContentType !== "image/heic") {
            image = await sharp(image).resize({
                width: RESIZED_WIDTH,
                withoutEnlargement: true
            }).toBuffer();
        }
        await S3.send(new PutObjectCommand({
            Bucket: DESC_BUCKET,
            Key: destKey,
            Body: image,
            ContentType: ContentType
        }));
        console.log(`Successfully resized and uploaded to ${destKey}`);
        return {
            statusCode: 200,
            body: "Success",
            desc: DESC_BUCKET
        };
    }
    catch (error) {
        console.error("Error processing image:", error);
        return {
            statusCode: 500,
            body: "Failed",
            desc: DESC_BUCKET
        };
    }
};
//# sourceMappingURL=index.js.map