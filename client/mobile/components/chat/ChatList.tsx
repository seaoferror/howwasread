import { FlatList, StyleSheet } from "react-native";
import { colors } from "@/constants";
import ChatPreviewItem from "@/components/chat/ChatPreviewItem";
import { usePreview } from "@/hooks/useChat";

export default function ChatList() {
  const preview = usePreview();

  return (
    <FlatList
      data={preview}
      renderItem={({ item }) => <ChatPreviewItem preview={item} />}
      keyExtractor={(item) => item.id}
      contentContainerStyle={styles.contentContainer}
    />
  );
}

const styles = StyleSheet.create({
  contentContainer: {
    paddingVertical: 12,
    backgroundColor: colors.SAND_150,
    gap: 12,
  },
});
