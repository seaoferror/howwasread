import { Pressable, StyleSheet, View } from "react-native";
import { ChatPreview } from "@/types/chat";
import { colors } from "@/constants";
import { router } from "expo-router";
import { use } from "react";

interface ChatPreviewItemProps {
  preview: ChatPreview;
}

export default function ChatPreviewItem({ preview }: ChatPreviewItemProps) {
  useGetRoomName()
  return <Pressable
    style={styles.container}
    onPress={() => router.push({
      pathname: `/chat`,
      params: {
        roomId: preview.roomId
      },
    })}
  >
    <View style={styles.content}>

    </View>

  </Pressable>;
}

const styles = StyleSheet.create({
  container: { backgroundColor: colors.SAND_110 },
  content: {
    padding: 16,
  },
});
