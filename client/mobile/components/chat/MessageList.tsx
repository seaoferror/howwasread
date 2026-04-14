import { FlatList, StyleSheet } from "react-native";
import { useGetInfiniteMessages } from "@/hooks/useChat";
import { useSQLiteContext } from "expo-sqlite";
import { useEffect, useState } from "react";
import { Message, MessageEntity } from "@/types/chat";
import * as SQLite from "expo-sqlite";
import { stringify as uuidStringify } from "uuid";
import MessageItem from "@/components/chat/MessageItem";
import { useLocalSearchParams } from "expo-router";

export default function MessageList() {
  const db = useSQLiteContext();
  const { id: roomId } = useLocalSearchParams();
  const { data, fetchNextPage, hasNextPage, isFetchingNextPage } =
    useGetInfiniteMessages(db, String(roomId));
  const [messages, setMessages] = useState<Omit<Message, "roomId">[]>([]);

  useEffect(() => {
    if (data?.pages) {
      setMessages(data.pages.flat());
    }
  }, [data]);

  SQLite.addDatabaseChangeListener(async (event) => {
    const newMessageRaw = await db.getFirstAsync<MessageEntity>(
      `SELECT * FROM message WHERE rowid = ?`,
      event.rowId,
    );
    if (!newMessageRaw) return;
    if (roomId !== uuidStringify(newMessageRaw.room_id)) {
      return;
    }
    setMessages([
      {
        id: uuidStringify(newMessageRaw.id),
        fromId: uuidStringify(newMessageRaw.from_id),
        contentType: newMessageRaw.content_type,
        content: newMessageRaw.content,
        createdAt: newMessageRaw.created_at,
        isDayFirst: Boolean(newMessageRaw.is_day_first),
      },
      ...messages,
    ]);
  });

  const handleEndReached = () => {
    if (hasNextPage && !isFetchingNextPage) {
      fetchNextPage();
    }
  };

  return (
    <FlatList
      data={messages}
      renderItem={({ item }) => <MessageItem message={item} />}
      keyExtractor={(item) => String(item.id)}
      contentContainerStyle={{ ...styles.contentContainer }}
      onEndReached={handleEndReached}
      onEndReachedThreshold={0.5}
      inverted={true}
    />
  );
}

const styles = StyleSheet.create({
  contentContainer: {
    paddingVertical: 12,
    backgroundColor: "transparent",
    gap: 12,
  },
});
