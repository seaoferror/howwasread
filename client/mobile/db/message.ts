import { Message, MessageEntity } from "@/types/chat";
import { parse as uuidParse, stringify as uuidStringify } from "uuid";
import { type SQLiteDatabase } from "expo-sqlite";

export async function initDB(db: SQLiteDatabase) {
  await db.execAsync(`
    PRAGMA journal_mode = 'wal';
    DROP TABLE message;
    CREATE TABLE IF NOT EXISTS message (
        id BLOB PRIMARY KEY NOT NULL,
        room_id BLOB NOT NULL,
        from_id BLOB NOT NULL,
        content_type TEXT NOT NULL,
        contents TEXT NOT NULL,
        created_at TEXT NOT NULL,
        is_day_first INTEGER DEFAULT 0
    );`);
}

export async function checkIfFirstOfDay(
  db: SQLiteDatabase,
  roomId: Uint8Array,
  timestamp: string,
): Promise<number> {
  const datePrefix = timestamp.substring(0, 10) + "%";

  const existing = await db.getFirstAsync(
    `SELECT 1 FROM message WHERE room_id = ? AND created_at LIKE ? LIMIT 1`,
    roomId,
    datePrefix,
  );

  return existing ? 0 : 1;
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
            created_at,
            is_day_first
     FROM message
     WHERE room_id = ?
     ORDER BY rowid DESC LIMIT ?
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
    isDayFirst: Boolean(row.is_day_first),
  }));
}
