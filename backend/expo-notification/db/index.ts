import { type Client } from "cassandra-driver";

export async function removeNotificationInfoByIdAndToken(
  cassandra: Client,
  id: string,
  token: string,
) {
  await cassandra.execute(
    `DELETE
         FROM notification_info_by_id
         WHERE id = ?
           AND device_push_token = ?`,
    [id, token],
    {
      prepare: true,
    },
  );
}
