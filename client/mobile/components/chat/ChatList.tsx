import { FlatList, StyleSheet, View } from "react-native";
import { useSQLiteContext } from "expo-sqlite";
import * as SQLite from "expo-sqlite";
import { useFocusEffect } from "expo-router";
import { stringify as uuidStringify } from "uuid";
import { useState } from "react";
import { colors } from "@/constants";
import { ChatPreview, MessageEntity } from "@/types/chat";
import ChatPreviewItem from "@/components/chat/ChatPreviewItem";

export default function ChatList() {
  const [preview, setPreview] = useState<ChatPreview[]>([]);
  const db = useSQLiteContext();
  SQLite.addDatabaseChangeListener(async (event) => {
    const newMessageRaw = await db.getFirstAsync<MessageEntity>(
      `SELECT * FROM message WHERE rowid = ?`,
      event.rowId,
    );

    if (!newMessageRaw) return;

    const newPreviewItem: ChatPreview = {
      id: uuidStringify(newMessageRaw.id),
      roomId: uuidStringify(newMessageRaw.room_id),
      fromId: uuidStringify(newMessageRaw.from_id),
      contentType: newMessageRaw.content_type,
      content: newMessageRaw.content,
      createdAt: newMessageRaw.created_at,
    };

    setPreview(prev => {
      const filtered = prev.filter(
        item => item.roomId !== newPreviewItem.roomId
      );
      return [newPreviewItem, ...filtered];
    });
  });
  useFocusEffect(() => {
    const wrapper = async () => {
      const previewRaw = await db.getAllAsync<MessageEntity>(
        `SELECT m.*
         FROM message m
                INNER JOIN
              (SELECT room_id, MAX(rowid) AS max_rowid
               FROM message
               GROUP BY room_id) latest
              ON m.room_id = latest.room_id AND m.rowid = latest.max_rowid
         ORDER BY m.rowid DESC;`
      );

      const p = previewRaw.map((row) => ({
        id: uuidStringify(row.id),
        roomId: uuidStringify(row.room_id),
        fromId: uuidStringify(row.from_id),
        contentType: row.content_type,
        content: row.content,
        createdAt: row.created_at,
      }));
      setPreview(p);
    };
    wrapper();
  });
  return <FlatList
    data={preview}
    renderItem={({item}) => <ChatPreviewItem preview={item}/>}
    keyExtractor={(item) => item.id}
    contentContainerStyle={styles.contentContainer}
  />;
}


const styles = StyleSheet.create({
  contentContainer: {
    paddingVertical: 12,
    backgroundColor: colors.SAND_150,
    gap: 12,
  },
});