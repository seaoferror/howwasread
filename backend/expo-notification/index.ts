import express from "express";
import Expo from "expo-server-sdk";

const app = express();
const expo = new Expo();
app.use(express.json());

app.get("/test", async (req, res) => {
  const token = req.body.token;
  const ticket = await expo.sendPushNotificationsAsync([
    {
      to: token,
      title: "test",
      body: "test",
      badge: 0,
    },
  ]);
  const receiptIds = ticket.filter((t) => t.status === "ok").map((t) => t.id);
  const chunk = expo.chunkPushNotificationReceiptIds(receiptIds);
  res.send(await expo.getPushNotificationReceiptsAsync(chunk[0]));
});

app.listen(3000, () => {
  console.log("express server running...");
});
