import { KeyboardAvoidingView, Platform, StyleSheet } from "react-native";
import {
  SafeAreaView,
  useSafeAreaInsets,
} from "react-native-safe-area-context";
import { useLocalSearchParams, useNavigation } from "expo-router";
import { useGetChatRoomInfo } from "@/hooks/useChat";
import MessageList from "@/components/chat/MessageList";
import { colors } from "@/constants";
import MessageInput from "@/components/chat/MessageInput";
import useKeyboard from "@/hooks/useKeyboard";
import { useEffect } from "react";
import { getKVStore } from "@/db/storage";

export default function ChatScreen() {
  const { id: roomId } = useLocalSearchParams();
  const { data: roomInfo } = useGetChatRoomInfo(String(roomId));
  const { isKeyboardVisible } = useKeyboard();
  const navigation = useNavigation();
  const insets = useSafeAreaInsets();

  useEffect(() => {
    navigation.setOptions({
      title: roomInfo?.name ?? getKVStore(String(roomId)),
    });
  }, [roomInfo]);

  return (
    <SafeAreaView style={styles.container} edges={["bottom"]}>
      <KeyboardAvoidingView
        style={styles.keyboardAvoidingView}
        behavior={Platform.OS === "ios" ? "padding" : "height"}
        keyboardVerticalOffset={
          Platform.OS === "ios" || isKeyboardVisible
            ? Platform.OS === "ios"
              ? 100
              : 80
            : insets.bottom
        }
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
