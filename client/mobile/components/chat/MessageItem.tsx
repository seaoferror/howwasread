import { StyleSheet, Text, View } from "react-native";
import { Message } from "@/types/chat";
import { getHourMinute, getLongDate } from "@/util/time";
import { useGetMyProfile, useGetProfile } from "@/hooks/useProfile";
import { colors } from "@/constants";
import { useGetChatRoomInfo, useGetSignedURLs } from "@/hooks/useChat";
import { Image } from "expo-image";
import VideoMessage from "@/components/chat/VideoMessage";
import VoiceMessage from "@/components/chat/VoiceMessage";
import { useLocalSearchParams } from "expo-router";

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
  const { urls } = useGetSignedURLs({
    contentType: message.contentType,
    contents: message.contents,
  });
  if (!myProfile || !fromProfile || !roomInfo) {
    return null;
  }

  const isMine = myProfile.id === message.fromId;
  const isEvent =
    message.contentType === "participate" ||
    message.contentType === "quit" ||
    message.contentType === "create";

  if (
    (message.contentType === "image" ||
      message.contentType === "audio" ||
      message.contentType === "video") &&
    !urls[0]
  ) {
    return null;
  }
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
        roomInfo.type === "group" && (
          <View style={styles.pillPosition}>
            <View style={styles.pill}>
              <Text style={styles.pillText}>
                {fromProfile.name} {message.contentType}s chat room
              </Text>
            </View>
          </View>
        )
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
              {isFromChange && <Text>{fromProfile.name}</Text>}
              {message.contentType === "text" ? (
                <Text style={styles.content}>{message.contents[0]}</Text>
              ) : message.contentType === "image" ? (
                urls.map((url, idx) => (
                  <Image
                    key={idx}
                    style={styles.media}
                    source={url}
                    cachePolicy="memory"
                    priority="low"
                    onError={(event) => {
                      console.log(event.error);
                    }}
                  />
                ))
              ) : message.contentType === "audio" ? (
                urls.some((url) => url) && <VoiceMessage url={urls[0] ?? ""} />
              ) : message.contentType === "video" ? (
                <>
                  {urls.map(
                    (url, idx) => url && <VideoMessage key={idx} url={url} />,
                  )}
                </>
              ) : null}
            </View>
            <Text style={styles.time}>{getHourMinute(message.createdAt)}</Text>
          </View>
        </View>
      )}
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
