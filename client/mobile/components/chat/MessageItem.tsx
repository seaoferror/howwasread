import { StyleSheet, Text, View } from "react-native";
import { Message } from "@/types/chat";
import { getHourMinute, getLongDate } from "@/util/time";
import { useMyProfile } from "@/hooks/useMyProfile";
import { colors } from "@/constants";
import { useGetSignedURLs } from "@/hooks/useChat";
import { Image } from "expo-image";
import VideoMessage from "@/components/chat/VideoMessage";
import VoiceMessage from "@/components/chat/VoiceMessage";

interface MessageItemProps {
  message: Omit<Message, "roomId">;
  isDayFirst: boolean;
}

export default function MessageItem({ message, isDayFirst }: MessageItemProps) {
  const { profile } = useMyProfile();
  const { urls } = useGetSignedURLs({
    contentType: message.contentType,
    contents: message.contents,
  });

  const isMine = profile.id === message.fromId;

  return (
    <View style={styles.container}>
      {isDayFirst && (
        <View style={styles.dateContainer}>
          <View style={styles.datePill}>
            <Text style={styles.dateText}>
              {getLongDate(message.createdAt)}
            </Text>
          </View>
        </View>
      )}

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
            {message.contentType === "text" ? (
              <Text style={styles.content}>{message.contents[0]}</Text>
            ) : message.contentType === "image" ? (
              urls.map((url, idx) => (
                <Image
                  key={idx}
                  style={styles.media}
                  source={url}
                  cachePolicy="memory"
                />
              ))
            ) : message.contentType === "audio" ? (
              urls.some((url) => url) && <VoiceMessage url={urls[0]} />
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
  dateContainer: {
    justifyContent: "center",
    alignItems: "center",
    marginVertical: 12,
  },
  datePill: {
    backgroundColor: "rgba(0, 0, 0, 0.05)",
    borderRadius: 999,
    paddingHorizontal: 12,
    paddingVertical: 4,
  },
  dateText: {
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
