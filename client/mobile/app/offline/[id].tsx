import { Pressable, StyleSheet, Text, View } from "react-native";
import { useLocalSearchParams } from "expo-router";
import OfflineConversationDetail from "@/components/conversation/OfflineConversationDetail";
import { colors } from "@/constants";
import { reportUser } from "@/api/chat";
import {
  useBlockConversation,
  useGetOfflineConversationDetail,
} from "@/hooks/useConversation";
import { reportOfflineConversation } from "@/api/conversation";
import Toast from "react-native-toast-message";
import { useActionSheet } from "@expo/react-native-action-sheet";

export default function OfflineConversationDetailScreen() {
  const { id } = useLocalSearchParams();
  const { data } = useGetOfflineConversationDetail(String(id));
  const blockConversationMutation = useBlockConversation();

  const { showActionSheetWithOptions } = useActionSheet();

  const handlePress = () => {
    showActionSheetWithOptions(
      {
        options: ["Delete from map", `Report and Delete from map`, "Cancel"],
        destructiveButtonIndex: 1,
        cancelButtonIndex: 2,
      },
      async (selectedIndex?: number) => {
        switch (selectedIndex) {
          case 0:
            blockConversationMutation.mutate({
              id: String(id),
            });
            break;
          case 1:
            blockConversationMutation.mutate({
              id: String(id),
            });
            if (data) {
              await Promise.all([
                data.moderatorIds.map((modId) => reportUser({ id: modId })),
              ]);
              await reportOfflineConversation({ id: String(id) });
              Toast.show({
                type: "info",
                text1: "Success report",
                text2:
                  "We will review this conversation, sorry for inconvenience.",
              });
            }
        }
      },
    );
  };

  return (
    <View style={styles.container}>
      <View style={styles.box}>
        <OfflineConversationDetail id={String(id)} />
      </View>
      <View style={styles.footer}>
        <Pressable
          onPress={async () => handlePress()}
          style={({ pressed }) => [pressed && styles.reportPressed]}
        >
          <Text style={styles.reportText}>Report conversation</Text>
        </Pressable>
      </View>
    </View>
  );
}
const styles = StyleSheet.create({
  box: {
    padding: 16,
    marginHorizontal: 16,
    backgroundColor: "#ffffff",
    borderRadius: 8,
    shadowColor: "#000",
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.1,
    shadowRadius: 4,
    elevation: 2,
  },
  container: {
    flex: 1,
    backgroundColor: colors.SAND_110,
  },
  footer: {
    position: "absolute",
    bottom: 40,
    left: 0,
    right: 0,
    alignItems: "center",
  },
  reportText: {
    color: colors.GRAY_400,
    fontSize: 12,
  },
  reportPressed: {
    opacity: 0.6,
  },
});
