import { SafeAreaView } from "react-native-safe-area-context";
import { StyleSheet } from "react-native";

import ChatList from "@/components/chat/ChatList";
import { colors } from "@/constants";

export default function ChatsScreen() {
  return (
    <SafeAreaView style={styles.container}>
      <ChatList />
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: colors.SAND_110,
  },
});
