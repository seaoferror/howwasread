import { FlatList, StyleSheet } from "react-native";
import { useGetInfiniteMessages } from "@/hooks/useChat";
import * as SQLite from "expo-sqlite";
import { useSQLiteContext } from "expo-sqlite";
import { useEffect, useMemo, useState } from "react";
import { Message } from "@/types/chat";
import { stringify as uuidStringify } from "uuid";
import MessageItem from "@/components/chat/MessageItem";
import { useLocalSearchParams } from "expo-router";
import { findNewMessage } from "@/db/message";
import { useGetMyProfile } from "@/hooks/useProfile";
import { getSecure } from "@/db/storage";

export default function MessageList() {
  const db = useSQLiteContext();
  const { id: roomId } = useLocalSearchParams();
  const { data, fetchNextPage, hasNextPage, isFetchingNextPage } =
    useGetInfiniteMessages(db, String(roomId));
  const { data: myProfile } = useGetMyProfile();

  const [liveMessages, setLiveMessages] = useState<Omit<Message, "roomId">[]>(
    [],
  );

  const allMessages = useMemo(() => {
    const fetchedMessages = data?.pages?.flat() || [];
    return [...liveMessages, ...fetchedMessages];
  }, [data?.pages, liveMessages]);

  useEffect(() => {
    const subscription = SQLite.addDatabaseChangeListener(async (event) => {
      const newMessageRaw = await findNewMessage(db, event.rowId);
      if (!newMessageRaw) return;
      if (roomId !== uuidStringify(newMessageRaw.room_id)) return;
      const newMessage: Omit<Message, "roomId"> = {
        id: uuidStringify(newMessageRaw.id),
        fromId: uuidStringify(newMessageRaw.from_id),
        contentType: newMessageRaw.content_type,
        contents: JSON.parse(newMessageRaw.contents),
        createdAt: newMessageRaw.created_at,
      };
      setLiveMessages((prev) => [newMessage, ...prev]);
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

  return (
    <FlatList
      data={allMessages}
      renderItem={({ item, index }) => {
        const olderMessage = allMessages[index + 1];
        const currentMsgDate = new Date(item.createdAt).toLocaleDateString();
        const olderMsgDate = olderMessage
          ? new Date(olderMessage.createdAt).toLocaleDateString()
          : null;
        const olderMsgFromId = olderMessage ? olderMessage.fromId : null;
        const isDayFirst = currentMsgDate !== olderMsgDate;
        const isFromChange =
          olderMsgFromId !== item.fromId ||
          item.fromId !== (myProfile?.id ?? getSecure("myId"));
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
