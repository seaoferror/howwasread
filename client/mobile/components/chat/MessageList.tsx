import { FlatList, StyleSheet } from "react-native";
import { useGetChatRoomInfo, useGetInfiniteMessages } from "@/hooks/useChat";
import { addDatabaseChangeListener, useSQLiteContext } from "expo-sqlite";
import { useEffect, useMemo, useState } from "react";
import { Message } from "@/types/chat";
import { stringify as uuidStringify } from "uuid";
import MessageItem from "@/components/chat/MessageItem";
import { useLocalSearchParams } from "expo-router";
import { findNewMessage } from "@/db/message";
import { useGetMyProfile } from "@/hooks/useProfile";
import { getKVStore } from "@/db/storage";

export default function MessageList() {
  const db = useSQLiteContext();
  const { id: roomId } = useLocalSearchParams();
  const { data, fetchNextPage, hasNextPage, isFetchingNextPage } =
    useGetInfiniteMessages(db, String(roomId));
  const { data: myProfile } = useGetMyProfile();
  const { data: roomInfo } = useGetChatRoomInfo(String(roomId));

  const isGroup =
    (roomInfo?.type ?? getKVStore("type" + String(roomId))) === "group";

  const [liveMessages, setLiveMessages] = useState<Omit<Message, "roomId">[]>(
    [],
  );

  const allMessages = useMemo(() => {
    const fetchedMessages = data?.pages?.flat() || [];
    return [...liveMessages, ...fetchedMessages];
  }, [data?.pages, liveMessages]);

  useEffect(() => {
    const subscription = addDatabaseChangeListener(async (event) => {
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
      if (newMessage.contentType === "delete") {
        setLiveMessages((prev) =>
          prev.filter((m) => m.id !== newMessage.contents[0]),
        );
        return;
      }
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
        if (!olderMessage) {
          return (
            <MessageItem
              message={item}
              isDayFirst={true}
              showName={isGroup}
              successDelete={(id) => {
                setLiveMessages((prev) => prev.filter((m) => m.id !== id));
              }}
            />
          );
        }
        const currentMsgDate = new Date(item.createdAt).toLocaleDateString();
        const olderMsgDate = new Date(
          olderMessage.createdAt,
        ).toLocaleDateString();
        const olderMsgFromId = olderMessage.fromId;
        const isDayFirst = currentMsgDate !== olderMsgDate;
        const isFromChange = olderMsgFromId !== item.fromId;
        const isMyId = item.fromId === (myProfile?.id ?? getKVStore("myId"));
        const wasEvent =
          olderMessage.contentType === "participate" ||
          olderMessage.contentType === "create" ||
          olderMessage.contentType === "unblock";
        return (
          <MessageItem
            message={item}
            isDayFirst={isDayFirst}
            showName={
              (wasEvent || isFromChange || isDayFirst) && !isMyId && isGroup
            }
            successDelete={(id) => {
              setLiveMessages((prev) => prev.filter((m) => m.id !== id));
            }}
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
