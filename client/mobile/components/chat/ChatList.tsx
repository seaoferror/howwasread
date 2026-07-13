import { FlatList, StyleSheet } from "react-native";
import { colors } from "@/constants";
import ChatPreviewItem from "@/components/chat/ChatPreviewItem";
import { addDatabaseChangeListener, useSQLiteContext } from "expo-sqlite";
import { useEffect, useState } from "react";
import { Message } from "@/types/chat";
import {
  deleteMessage,
  findNewMessage,
  findPreview,
  getBeforeMessage,
  quitLive,
} from "@/db/message";
import { stringify as uuidStringify } from "uuid";

export default function ChatList() {
  const db = useSQLiteContext();
  const [preview, setPreview] = useState<Message[]>([]);

  useEffect(() => {
    const wrapper = async () => {
      const previewRaw = await findPreview(db);

      const p = previewRaw.map((row) => ({
        id: uuidStringify(row.id),
        roomId: uuidStringify(row.room_id),
        fromId: uuidStringify(row.from_id),
        contentType: row.content_type,
        contents: JSON.parse(row.contents),
        createdAt: row.created_at,
      }));
      setPreview(p);

      addDatabaseChangeListener(async (event) => {
        const newMessageRaw = await findNewMessage(db, event.rowId);

        if (!newMessageRaw) return;

        const newPreviewItem: Message = {
          id: uuidStringify(newMessageRaw.id),
          roomId: uuidStringify(newMessageRaw.room_id),
          fromId: uuidStringify(newMessageRaw.from_id),
          contentType: newMessageRaw.content_type,
          contents: JSON.parse(newMessageRaw.contents),
          createdAt: newMessageRaw.created_at,
        };
        if (newPreviewItem.contentType === "quit") {
          await quitLive(db, newPreviewItem.roomId);
          setPreview((prev) =>
            prev.filter((item) => item.roomId !== newPreviewItem.roomId),
          );
          return;
        }
        if (newPreviewItem.contentType === "delete") {
          await deleteMessage(db, newPreviewItem.id);
          await deleteMessage(db, newPreviewItem.contents[0]);
          const found = preview.find(
            (p) => p.id === newPreviewItem.contents[0],
          );
          if (found) {
            const beforeMsgRaw = await getBeforeMessage(
              db,
              newPreviewItem.contents[0],
              found.roomId,
            );
            if (beforeMsgRaw) {
              const beforeMsg: Message = {
                id: uuidStringify(beforeMsgRaw.id),
                roomId: found.roomId,
                fromId: uuidStringify(beforeMsgRaw.room_id),
                contentType: beforeMsgRaw.content_type,
                contents: JSON.parse(beforeMsgRaw.contents),
                createdAt: beforeMsgRaw.created_at,
              };
              setPreview((prev) => [
                beforeMsg,
                ...prev.filter(
                  (item) => item.id !== newPreviewItem.contents[0],
                ),
              ]);
              return;
            }
            setPreview((prev) =>
              prev.filter((item) => item.id !== newPreviewItem.contents[0]),
            );
          }
          return;
        }
        setPreview((prev) => {
          const filtered = prev.filter(
            (item) => item.roomId !== newPreviewItem.roomId,
          );
          return [newPreviewItem, ...filtered];
        });
      });
    };
    wrapper();
  }, []);

  return (
    <FlatList
      data={preview}
      renderItem={({ item }) => (
        <ChatPreviewItem
          preview={item}
          successQuit={async (roomId) => {
            await quitLive(db, roomId);
            setPreview((prev) => prev.filter((item) => item.roomId !== roomId));
          }}
        />
      )}
      keyExtractor={(item) => item.id}
      contentContainerStyle={styles.contentContainer}
    />
  );
}

const styles = StyleSheet.create({
  contentContainer: {
    paddingVertical: 12,
    backgroundColor: colors.SAND_150,
    gap: 12,
  },
});
