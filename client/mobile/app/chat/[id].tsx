import {
  KeyboardAvoidingView,
  Platform,
  Pressable,
  StyleSheet,
  Text,
} from "react-native";
import {
  SafeAreaView,
  useSafeAreaInsets,
} from "react-native-safe-area-context";
import { router, useLocalSearchParams, useNavigation } from "expo-router";
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
  const roomName = roomInfo?.name ?? getKVStore(String(roomId));
  const isPersonal =
    (roomInfo?.type ?? getKVStore("type" + String(roomId))) === "personal";

  useEffect(() => {
    navigation.setOptions({
      headerTitle: () => (
        <Pressable
          onPress={() =>
            router.push({
              pathname: "/chat/detail/[id]",
              params: { id: String(roomId) },
            })
          }
          disabled={isPersonal}
        >
          <Text
            style={[
              { fontSize: 17 },
              !isPersonal && {
                fontWeight: "bold",
                textDecorationLine: "underline",
              },
            ]}
          >
            {roomName.length > 15 ? roomName.slice(0, 15) + "..." : roomName}
          </Text>
        </Pressable>
      ),
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
