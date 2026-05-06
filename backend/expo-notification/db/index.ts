import { type Client } from "cassandra-driver";

export async function removePushTokenByIdAndDeviceId(
  cassandra: Client,
  id: string,
  deviceId: string,
) {
  await cassandra.execute(
    `DELETE
         FROM device_push_token_by_id
         WHERE id = ?
           AND device_id = ?`,
    [id, deviceId],
    {
      prepare: true,
    },
  );
}
