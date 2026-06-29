import { Pressable, StyleSheet, Text, View } from "react-native";
import { Message } from "@/types/chat";
import { getHourMinute, getLongDate } from "@/util/time";
import { useGetMyProfile, useGetProfile } from "@/hooks/useProfile";
import { colors } from "@/constants";
import { useGetChatRoomInfo } from "@/hooks/useChat";
import { Image } from "expo-image";
import VideoMessage from "@/components/chat/VideoMessage";
import VoiceMessage from "@/components/chat/VoiceMessage";
import { useLocalSearchParams } from "expo-router";
import { getKVStore, getSecure, setKVStore } from "@/db/storage";
import { useEffect, useState } from "react";
import ImageModal from "@/components/ImageModal";

interface MessageItemProps {
  message: Omit<Message, "roomId">;
  isDayFirst: boolean;
  isFromChange: boolean;
}

export default function MessageItem({
  message,
  isDayFirst,
  isFromChange,
}: MessageItemProps) {
  const { data: myProfile } = useGetMyProfile();
  const { data: fromProfile } = useGetProfile(message.fromId);
  const { id: roomId } = useLocalSearchParams();
  const { data: roomInfo } = useGetChatRoomInfo(String(roomId));
  const [pressedImageContent, setPressedImageContent] = useState<string | null>(
    null,
  );

  useEffect(() => {
    if (fromProfile) {
      setKVStore(message.fromId, fromProfile.name);
    }
  }, [fromProfile, message.fromId]);

  const isMine = (myProfile?.id ?? getSecure("myId")) === message.fromId;
  const isEvent =
    message.contentType === "participate" ||
    message.contentType === "create" ||
    message.contentType === "block" ||
    message.contentType === "unblock";

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

      {isEvent ? (
        <View style={styles.pillPosition}>
          <View style={styles.pill}>
            <Text style={styles.pillText}>
              {message.contentType === "participate" &&
                (fromProfile?.name ??
                  getKVStore(message.fromId) + "participates chatroom")}
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
              {isFromChange && (
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
                <VoiceMessage url={getKVStore(message.contents[0])} />
              ) : message.contentType === "video" ? (
                message.contents.map((content, idx) => (
                  <VideoMessage key={idx} url={getKVStore(content)} />
                ))
              ) : null}
            </View>
            <Text style={styles.time}>{getHourMinute(message.createdAt)}</Text>
          </View>
        </View>
      )}
      <ImageModal
        imageContent={pressedImageContent}
        onClose={() => setPressedImageContent(null)}
      />
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    marginVertical: 4,
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
    gap: 6,
    maxWidth: "85%",
  },
  otherName: {
    paddingBottom: 12,
    fontWeight: "bold",
  },
  bubbleWrapperReverse: {
    flexDirection: "row-reverse",
  },
  messageContainer: {
    borderRadius: 16,
    paddingVertical: 10,
    paddingHorizontal: 14,
    overflow: "hidden",
  },
  mine: {
    backgroundColor: colors.SAND_150,
    borderBottomRightRadius: 4,
  },
  theirs: {
    backgroundColor: colors.WHITE,
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
