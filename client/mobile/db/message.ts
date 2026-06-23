import { Message, MessageEntity, MessagingResponse } from "@/types/chat";
import { parse as uuidParse, stringify as uuidStringify } from "uuid";
import { type SQLiteDatabase } from "expo-sqlite";
import { getTimestamp } from "@/util/time";

export async function initDB(db: SQLiteDatabase) {
  console.log("init db");
  await db.execAsync(`
    PRAGMA journal_mode = 'wal';
    CREATE TABLE IF NOT EXISTS message (
        id BLOB PRIMARY KEY NOT NULL,
        room_id BLOB NOT NULL,
        from_id BLOB NOT NULL,
        content_type TEXT NOT NULL,
        contents TEXT NOT NULL,
        created_at TEXT NOT NULL
    );`);
}

export async function deleteAllMessages(db: SQLiteDatabase) {
  await db.execAsync(`
    DELETE FROM message;
  `);
}

export async function findPreview(db: SQLiteDatabase) {
  return db.getAllAsync<MessageEntity>(
    `SELECT m.id,
            m.room_id,
            m.from_id,
            m.content_type,
            m.contents,
            m.created_at
     FROM message m
            INNER JOIN
          (SELECT room_id, MAX(id) AS id
           FROM message
           GROUP BY room_id) latest
          ON m.room_id = latest.room_id AND m.id = latest.id
     ORDER BY m.id DESC;`,
  );
}

export async function findMessagesByRoomId(
  db: SQLiteDatabase,
  roomId: string,
  page = 1,
): Promise<Omit<Message, "roomId">[]> {
  const size = 20;
  const offset = (page - 1) * size;
  const messagesRaw = await db.getAllAsync<Omit<MessageEntity, "room_id">>(
    `SELECT id,
            from_id,
            content_type,
            contents,
            created_at
     FROM message
     WHERE room_id = ?
     ORDER BY id DESC LIMIT ?
     OFFSET ?`,
    uuidParse(String(roomId)),
    size,
    offset,
  );

  return messagesRaw.map((row) => ({
    id: uuidStringify(row.id),
    fromId: uuidStringify(row.from_id),
    contentType: row.content_type,
    contents: JSON.parse(row.contents),
    createdAt: row.created_at,
  }));
}

export async function findNewMessage(db: SQLiteDatabase, rowId: number) {
  return db.getFirstAsync<MessageEntity>(
    `SELECT *
         FROM message
         WHERE rowid = ?`,
    rowId,
  );
}

export async function saveRecentMessage(
  db: SQLiteDatabase,
  m: MessagingResponse,
) {
  return db.runAsync(
    `INSERT
        OR IGNORE INTO message (id, room_id, from_id, content_type, contents, created_at)
               VALUES (?, ?, ?, ?, ?, ?);`,
    uuidParse(m.id),
    uuidParse(m.roomId),
    uuidParse(m.fromId),
    m.contentType,
    JSON.stringify(m.contents),
    getTimestamp(m.id),
  );
}

export async function deleteMessagesBeforeQuit(
  db: SQLiteDatabase,
  m: MessagingResponse,
) {
  return db.runAsync(
    `DELETE FROM message
     WHERE room_id = ? AND id <= ?;`,
    uuidParse(m.roomId),
    uuidParse(m.id),
  );
}
