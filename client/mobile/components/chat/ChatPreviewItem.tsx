import { Pressable, StyleSheet, Text, View } from "react-native";
import { Message } from "@/types/chat";
import { colors } from "@/constants";
import { router } from "expo-router";
import { useGetChatRoomInfo } from "@/hooks/useChat";
import { useGetProfile, useMyProfile } from "@/hooks/useMyProfile";

interface ChatPreviewItemProps {
  preview: Omit<Message, "isDayFirst">;
}

function formatPreviewDate(createdAt: string) {
  const date = new Date(createdAt);
  const now = new Date();

  const isToday =
    date.getFullYear() === now.getFullYear() &&
    date.getMonth() === now.getMonth() &&
    date.getDate() === now.getDate();

  const yesterday = new Date(now);
  yesterday.setDate(now.getDate() - 1);

  const isYesterday =
    date.getFullYear() === yesterday.getFullYear() &&
    date.getMonth() === yesterday.getMonth() &&
    date.getDate() === yesterday.getDate();

  if (isToday) {
    return date.toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" });
  }
  if (isYesterday) {
    return "yesterday";
  }
  return date.toLocaleDateString(undefined, {
    month: "numeric",
    day: "numeric",
  });
}

export default function ChatPreviewItem({ preview }: ChatPreviewItemProps) {
  const { data: roomInfo } = useGetChatRoomInfo(preview.roomId);
  const { data: fromProfile } = useGetProfile(preview.fromId);
  const { profile: myProfile } = useMyProfile();

  return (
    <Pressable
      style={({ pressed }) => [styles.container, pressed && styles.pressed]}
      onPress={() =>
        router.push({
          pathname: `/chat/[id]`,
          params: {
            id: preview.roomId,
          },
        })
      }
    >
      <View style={styles.avatar}>
        <Text style={styles.avatarText}>
          {(roomInfo?.name ?? "?").charAt(0)}
        </Text>
      </View>

      <View style={styles.textContainer}>
        <Text style={styles.roomName} numberOfLines={1}>
          {roomInfo?.name ?? ""}
        </Text>
        <Text style={styles.messagePreview} numberOfLines={1}>
          {preview.fromId === preview.roomId || preview.fromId === myProfile.id
            ? ""
            : `${fromProfile?.name}: `}
          {preview.contentType === "text"
            ? preview.content
            : preview.contentType}
        </Text>
      </View>

      <Text style={styles.timestamp}>
        {formatPreviewDate(preview.createdAt)}
      </Text>
    </Pressable>
  );
}

const styles = StyleSheet.create({
  container: {
    flexDirection: "row",
    alignItems: "center",
    backgroundColor: colors.SAND_110,
    paddingHorizontal: 16,
    paddingVertical: 14,
    borderBottomWidth: StyleSheet.hairlineWidth,
    borderBottomColor: colors.SAND_200,
  },
  pressed: {
    backgroundColor: colors.SAND_150,
  },
  avatar: {
    width: 48,
    height: 48,
    borderRadius: 24,
    backgroundColor: colors.ORANGE_150,
    justifyContent: "center",
    alignItems: "center",
    marginRight: 12,
  },
  avatarText: {
    fontSize: 18,
    fontWeight: "600",
    color: colors.GRAY_700,
  },
  textContainer: {
    flex: 1,
    justifyContent: "center",
    gap: 4,
  },
  roomName: {
    fontSize: 16,
    fontWeight: "600",
    color: colors.GRAY_900,
  },
  messagePreview: {
    fontSize: 14,
    color: colors.GRAY_500,
  },
  timestamp: {
    fontSize: 12,
    color: colors.GRAY_500,
    marginLeft: 8,
    alignSelf: "flex-start",
    marginTop: 2,
  },
});
