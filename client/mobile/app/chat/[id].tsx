import { KeyboardAvoidingView, Platform, StyleSheet } from "react-native";
import { SafeAreaView } from "react-native-safe-area-context";
import {
  useFocusEffect,
  useLocalSearchParams,
  useNavigation,
} from "expo-router";
import { useGetChatRoomInfo } from "@/hooks/useChat";
import MessageList from "@/components/chat/MessageList";
import { colors } from "@/constants";
import MessageInput from "@/components/chat/MessageInput";

export default function ChatScreen() {
  const { id: roomId } = useLocalSearchParams();
  const { data: roomInfo } = useGetChatRoomInfo(String(roomId));
  const navigation = useNavigation();

  useFocusEffect(() => {
    navigation.setOptions({
      title: roomInfo?.name,
    });
  });

  return (
    <SafeAreaView style={styles.container} edges={["bottom"]}>
      <KeyboardAvoidingView
        style={styles.keyboardAvoidingView}
        behavior={Platform.OS === "ios" ? "padding" : undefined}
        keyboardVerticalOffset={Platform.OS === "ios" ? 160 : 0}
      >
        <MessageList />
        <MessageInput />
      </KeyboardAvoidingView>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: colors.WHITE,
  },
  keyboardAvoidingView: {
    flex: 1,
  },
});
