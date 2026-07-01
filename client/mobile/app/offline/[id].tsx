import { StyleSheet, View } from "react-native";
import { useLocalSearchParams } from "expo-router";
import OfflineConversationDetail from "@/components/conversation/OfflineConversationDetail";
import { colors } from "@/constants";

export default function OfflineConversationDetailScreen() {
  const { id } = useLocalSearchParams();

  return (
    <View style={styles.container}>
      <View style={styles.box}>
        <OfflineConversationDetail id={String(id)} />
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
    backgroundColor: colors.SAND_110
  }
});
