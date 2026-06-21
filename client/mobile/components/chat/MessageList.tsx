import { FlatList, StyleSheet } from "react-native";
import { useGetInfiniteMessages } from "@/hooks/useChat";
import * as SQLite from "expo-sqlite";
import { useSQLiteContext } from "expo-sqlite";
import { useEffect, useState } from "react";
import { Message } from "@/types/chat";
import { stringify as uuidStringify } from "uuid";
import MessageItem from "@/components/chat/MessageItem";
import { useLocalSearchParams } from "expo-router";
import { findNewMessage } from "@/db/message";
import { useGetMyProfile } from "@/hooks/useProfile";

export default function MessageList() {
  const db = useSQLiteContext();
  const { id: roomId } = useLocalSearchParams();
  const { data, fetchNextPage, hasNextPage, isFetchingNextPage } =
    useGetInfiniteMessages(db, String(roomId));
  const [messages, setMessages] = useState<Omit<Message, "roomId">[]>([]);
  const { data: myProfile } = useGetMyProfile();

  useEffect(() => {
    if (data?.pages) {
      setMessages(data.pages.flat());
    }
  }, [data?.pages]);

  useEffect(() => {
    const subscription = SQLite.addDatabaseChangeListener(async (event) => {
      const newMessageRaw = await findNewMessage(db, event.rowId);
      if (!newMessageRaw) return;
      if (roomId !== uuidStringify(newMessageRaw.room_id)) return;
      const newMessage = {
        id: uuidStringify(newMessageRaw.id),
        fromId: uuidStringify(newMessageRaw.from_id),
        contentType: newMessageRaw.content_type,
        contents: JSON.parse(newMessageRaw.contents),
        createdAt: newMessageRaw.created_at,
      };
      if (
        newMessage.contentType === "image" ||
        newMessage.contentType === "video"
      ) {
        setTimeout(() => {
          setMessages((prev) => [newMessage, ...prev]);
        }, 1000);
        return;
      }
      setMessages((prev) => [newMessage, ...prev]);
    });

    return () => {
      subscription.remove();
    };
  }, [db, roomId]);

  const handleEndReached = () => {
    if (hasNextPage && !isFetchingNextPage) {
      fetchNextPage();
    }
  };

  if (!myProfile) {
    return null;
  }

  return (
    <FlatList
      data={messages}
      renderItem={({ item, index }) => {
        const olderMessage = messages[index + 1];
        const currentMsgDate = new Date(item.createdAt).toLocaleDateString();
        const olderMsgDate = olderMessage
          ? new Date(olderMessage.createdAt).toLocaleDateString()
          : null;
        const olderMsgFromId = olderMessage ? olderMessage.fromId : null;
        const isDayFirst = currentMsgDate !== olderMsgDate;
        const isFromChange =
          olderMsgFromId !== item.fromId && item.fromId !== myProfile.id;
        return (
          <MessageItem
            message={item}
            isDayFirst={isDayFirst}
            isFromChange={isFromChange}
          />
        );
      }}
      keyExtractor={(item, index) => `${item.id}-${index}`}
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
