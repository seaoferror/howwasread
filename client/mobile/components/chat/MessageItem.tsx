import { StyleSheet, View, Text } from "react-native";
import { Message } from "@/types/chat";
import { getHourMinute, getLongDate } from "@/util/time";
import { useMyProfile } from "@/hooks/useMyProfile";
import { colors } from "@/constants";

interface MessageItemProps {
  message: Omit<Message, "roomId">;
}

export default function MessageItem({ message }: MessageItemProps) {
  const { profile } = useMyProfile();
  console.log(message.contentType)

  return (
    <View style={styles.container}>
      {message.isDayFirst && (
        <View style={styles.dateContainer}>
          <View style={styles.datePill}>
            <Text style={styles.dateText}>
              {getLongDate(message.createdAt)}
            </Text>
          </View>
        </View>
      )}

      <View
        style={[
          styles.row,
          profile.id === message.fromId
            ? { alignItems: "flex-end" }
            : { alignItems: "flex-start" },
        ]}
      >
        <View style={styles.messageContainer}>
          <Text style={styles.content}>
            {message.contentType === "text" ? message.content : ""}
          </Text>
          <Text style={styles.time}>{getHourMinute(message.createdAt)}</Text>
        </View>
      </View>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    marginVertical: 2,
  },
  row: {
    width: "100%",
  },
  messageContainer: {
    gap: 6,
    borderRadius: 12,
    paddingVertical: 8,
    paddingHorizontal: 10,
    maxWidth: "80%",
  },
  mine: {
    backgroundColor: colors.SAND_150,
    borderBottomRightRadius: 4,
  },
  content: {
    fontSize: 14,
    lineHeight: 20,
    color: colors.GRAY_900,
  },
  time: {
    fontSize: 11,
    alignSelf: "flex-end",
    color: colors.GRAY_500,
  },
  dateContainer: {
    justifyContent: "center",
    alignItems: "center",
    marginVertical: 8,
  },
  datePill: {
    backgroundColor: colors.GRAY_100,
    borderRadius: 999,
    paddingHorizontal: 10,
    paddingVertical: 4,
  },
  dateText: {
    color: colors.GRAY_700,
    fontSize: 12,
  },
});
