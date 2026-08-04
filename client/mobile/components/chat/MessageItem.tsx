import { Pressable, StyleSheet, Text, View } from "react-native";
import { Message } from "@/types/chat";
import { getHourMinute, getLongDate } from "@/util/time";
import { useGetMyProfile, useGetProfile } from "@/hooks/useProfile";
import { colors } from "@/constants";
import { useGetChatRoomInfo, useSendMessage } from "@/hooks/useChat";
import { Image } from "expo-image";
import VideoMessage from "@/components/chat/VideoMessage";
import VoiceMessage from "@/components/chat/VoiceMessage";
import { useLocalSearchParams } from "expo-router";
import { getKVStore, setKVStore } from "@/db/storage";
import { useEffect, useState } from "react";
import ImageModal from "@/components/ImageModal";
import { useActionSheet } from "@expo/react-native-action-sheet";
import { useSQLiteContext } from "expo-sqlite";
import { deleteMessage } from "@/db/message";
import { reportUser } from "@/api/chat";
import Toast from "react-native-toast-message";

interface MessageItemProps {
  message: Omit<Message, "roomId">;
  isDayFirst: boolean;
  showName: boolean;
  successDelete: (id: string) => void;
}

export default function MessageItem({
  message,
  isDayFirst,
  showName,
  successDelete,
}: MessageItemProps) {
  const { data: myProfile } = useGetMyProfile();
  const { data: fromProfile } = useGetProfile(message.fromId);
  const { id: roomId } = useLocalSearchParams();
  const { data: roomInfo } = useGetChatRoomInfo(String(roomId));
  const [pressedImageContent, setPressedImageContent] = useState<string | null>(
    null,
  );
  const { showActionSheetWithOptions } = useActionSheet();
  const db = useSQLiteContext();
  const sendMessageMutation = useSendMessage();

  useEffect(() => {
    if (fromProfile) {
      setKVStore(message.fromId, fromProfile.name);
    }
  }, [fromProfile, message.fromId]);

  const isMine = (myProfile?.id ?? getKVStore("myId")) === message.fromId;
  const isEvent =
    message.contentType === "participate" ||
    message.contentType === "create" ||
    message.contentType === "block" ||
    message.contentType === "unblock";

  const handleLongPress = () => {
    showActionSheetWithOptions(
      {
        options:
          message.fromId === (myProfile?.id ?? getKVStore("myId"))
            ? [`Delete`, "Cancel"]
            : ["Delete", "Report and Delete", "Cancel"],
        cancelButtonIndex: 2,
        destructiveButtonIndex: 0,
      },
      async (selectedIndex?: number) => {
        console.log(selectedIndex);
        switch (selectedIndex) {
          case 0:
            await deleteMessage(db, message.id);
            successDelete(message.id);
            if (roomInfo) {
              sendMessageMutation.mutate({
                toId: String(roomId),
                toIdType: roomInfo.type,
                contentType: "delete",
                contents: [message.id],
              });
            }
            break;
          case 1:
            if (message.fromId === (myProfile?.id ?? getKVStore("myId"))) {
              return;
            }
            await deleteMessage(db, message.id);
            successDelete(message.id);
            console.log(roomInfo);
            if (roomInfo) {
              sendMessageMutation.mutate({
                toId: String(roomId),
                toIdType: roomInfo.type,
                contentType: "delete",
                contents: [message.id],
              });
            }
            console.log("start report...");
            await reportUser({ id: message.fromId });
            console.log("success report");
            Toast.show({
              type: "info",
              text1: "Success report",
              text2: "We will review this message, sorry for inconvenience.",
            });
        }
      },
    );
  };

  return (
    <View style={styles.container}>
      {isDayFirst && (
        <View style={styles.pillPosition}>
          <View style={styles.pill}>
            <Text style={styles.pillText}>
              {getLongDate(message.createdAt)}
            </Text>
          </View>
        </View>
      )}

      <Pressable onLongPress={() => handleLongPress()} disabled={isEvent}>
        {isEvent ? (
          <View style={styles.pillPosition}>
            <View style={styles.pill}>
              <Text style={styles.pillText}>
                {message.contentType === "participate" &&
                  (fromProfile?.name ?? getKVStore(message.fromId)) +
                    " participates chatroom"}
                {message.contentType === "create" && "You create chat room"}
                {message.contentType === "block" &&
                  `You block ${roomInfo?.name ?? getKVStore(String(roomId))}`}
                {message.contentType === "unblock" &&
                  `You unblock ${roomInfo?.name ?? getKVStore(String(roomId))}`}
              </Text>
            </View>
          </View>
        ) : (
          <View style={[styles.row, isMine ? styles.rowRight : styles.rowLeft]}>
            <View
              style={[
                styles.bubbleWrapper,
                isMine ? styles.bubbleWrapperReverse : null,
              ]}
            >
              <View
                style={[
                  styles.messageContainer,
                  isMine ? styles.mine : styles.theirs,
                ]}
              >
                {showName && (
                  <Text style={styles.otherName}>
                    {fromProfile?.name ?? getKVStore(message.fromId)}
                  </Text>
                )}
                {message.contentType === "text" ? (
                  <Text style={styles.content}>{message.contents[0]}</Text>
                ) : message.contentType === "image" ? (
                  message.contents.map((content, idx) => (
                    <Pressable
                      key={idx}
                      onPress={() => setPressedImageContent(content)}
                      onLongPress={() => handleLongPress()}
                    >
                      <Image
                        style={styles.media}
                        source={getKVStore(content)}
                        cachePolicy="memory"
                        priority="low"
                        onError={(event) => {
                          console.log(event.error);
                        }}
                      />
                    </Pressable>
                  ))
                ) : message.contentType === "audio" ? (
                  <VoiceMessage
                    url={getKVStore(message.contents[0])}
                    onLongPress={() => handleLongPress()}
                  />
                ) : message.contentType === "video" ? (
                  message.contents.map((content, idx) => (
                    <VideoMessage
                      key={idx}
                      url={getKVStore(content)}
                      onLongPress={() => handleLongPress()}
                    />
                  ))
                ) : null}
              </View>
              <Text style={styles.time}>
                {getHourMinute(message.createdAt)}
              </Text>
            </View>
          </View>
        )}
        {
          <ImageModal
            imageContent={pressedImageContent}
            onClose={() => setPressedImageContent(null)}
          />
        }
      </Pressable>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    marginVertical: 3.7,
    paddingHorizontal: 16,
  },
  row: {
    flexDirection: "row",
    width: "100%",
  },
  rowRight: {
    justifyContent: "flex-end",
  },
  rowLeft: {
    justifyContent: "flex-start",
  },
  bubbleWrapper: {
    flexDirection: "row",
    alignItems: "flex-end",
    gap: 3.7,
    maxWidth: "85%",
  },
  otherName: {
    paddingBottom: 4,
    fontWeight: "600",
  },
  bubbleWrapperReverse: {
    flexDirection: "row-reverse",
  },
  messageContainer: {
    borderRadius: 7,
    paddingVertical: 7,
    paddingHorizontal: 12,
    overflow: "hidden",
  },
  mine: {
    backgroundColor: colors.SAND_110,
    borderBottomRightRadius: 4,
  },
  theirs: {
    backgroundColor: colors.GRAY_50,
    borderBottomLeftRadius: 4,
    borderWidth: StyleSheet.hairlineWidth,
    borderColor: colors.GRAY_200 || "#E5E7EB",
  },
  content: {
    fontSize: 15,
    lineHeight: 22,
    color: colors.GRAY_900,
  },
  time: {
    fontSize: 11,
    color: colors.GRAY_500,
    marginBottom: 2,
  },
  pillPosition: {
    justifyContent: "center",
    alignItems: "center",
    marginVertical: 12,
  },
  pill: {
    backgroundColor: "rgba(0, 0, 0, 0.05)",
    borderRadius: 999,
    paddingHorizontal: 12,
    paddingVertical: 4,
  },
  pillText: {
    color: colors.GRAY_700,
    fontSize: 12,
    fontWeight: "500",
  },
  media: {
    width: 220,
    height: 220,
    borderRadius: 8,
    backgroundColor: colors.GRAY_100,
  },
});
