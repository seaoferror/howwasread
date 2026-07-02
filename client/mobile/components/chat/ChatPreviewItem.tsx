import { Pressable, StyleSheet, Text, View } from "react-native";
import { Message } from "@/types/chat";
import { colors } from "@/constants";
import { router } from "expo-router";
import {
  useCheckBlock,
  useGetChatRoomInfo,
  useSendMessage,
} from "@/hooks/useChat";
import { useGetProfile } from "@/hooks/useProfile";
import { formatPreviewDate } from "@/util/time";
import { useEffect } from "react";
import { getKVStore, setKVStore } from "@/db/storage";
import { useActionSheet } from "@expo/react-native-action-sheet";

interface ChatPreviewItemProps {
  preview: Message;
  successQuit: (roomId: string) => void;
}

export default function ChatPreviewItem({
  preview,
  successQuit,
}: ChatPreviewItemProps) {
  const { data: roomInfo } = useGetChatRoomInfo(preview.roomId);
  const { data: fromProfile } = useGetProfile(preview.fromId);
  const { data } = useCheckBlock(preview.roomId);
  const { showActionSheetWithOptions } = useActionSheet();
  const sendMessageMutation = useSendMessage();

  useEffect(() => {
    if (roomInfo) {
      setKVStore(preview.roomId, roomInfo.name);
      setKVStore("type" + preview.roomId, roomInfo.type);
    }
    if (fromProfile) {
      setKVStore(preview.fromId, fromProfile.name);
    }
  }, [roomInfo, fromProfile, preview.roomId, preview.fromId]);

  function handleLongPress() {
    if (!roomInfo || !fromProfile) {
      return;
    }
    const roomType = roomInfo.type;
    const opponentName = fromProfile.name;
    showActionSheetWithOptions(
      {
        title:
          roomType === "personal"
            ? opponentName
            : `Quit from the group chat room of ${roomInfo.name}`,
        options:
          roomType === "personal"
            ? ["Quit", data?.didBlock ? "Unblock" : "Block", "Cancel"]
            : [`Quit`, "Cancel"],
        destructiveButtonIndex: roomType === "personal" ? [0, 1] : 0,
        cancelButtonIndex: 2,
      },
      (selectedIndex?: number) => {
        console.log(selectedIndex);
        switch (selectedIndex) {
          case 0:
            sendMessageMutation.mutate(
              {
                toIdType: roomType,
                toId: preview.roomId,
                contentType: "quit",
                contents: [],
              },
              {
                onSuccess: () => {
                  successQuit(preview.roomId);
                },
              },
            );
            break;
          case 1:
            roomType === "personal" &&
              sendMessageMutation.mutate({
                toIdType: roomType,
                toId: preview.roomId,
                contentType: data?.didBlock ? "unblock" : "block",
                contents: [],
              });
            break;
          case 2:
            break;
        }
      },
    );
  }

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
      onLongPress={() => handleLongPress()}
    >
      <View style={styles.avatar}>
        <Text style={styles.avatarText}>
          {(roomInfo?.name ?? getKVStore(preview.roomId)).charAt(0)}
        </Text>
      </View>

      <View style={styles.textContainer}>
        <Text style={styles.roomName} numberOfLines={1}>
          {roomInfo?.name ?? getKVStore(preview.roomId)}
        </Text>
        <Text style={styles.messagePreview} numberOfLines={1}>
          {preview.fromId === preview.roomId ||
          preview.fromId === getKVStore("myId")
            ? ""
            : `${fromProfile?.name ?? getKVStore(preview.fromId)}: `}
          {preview.contentType === "text"
            ? preview.contents[0]
            : "(" + preview.contentType + ")"}
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
