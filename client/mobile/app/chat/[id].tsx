import {
  KeyboardAvoidingView,
  Platform,
  Pressable,
  StyleSheet,
} from "react-native";
import {
  SafeAreaView,
  useSafeAreaInsets,
} from "react-native-safe-area-context";
import {
  router,
  useFocusEffect,
  useLocalSearchParams,
  useNavigation,
} from "expo-router";
import { useGetChatRoomInfo } from "@/hooks/useChat";
import MessageList from "@/components/chat/MessageList";
import { colors } from "@/constants";
import MessageInput from "@/components/chat/MessageInput";
import useKeyboard from "@/hooks/useKeyboard";
import { Ionicons } from "@expo/vector-icons";

export default function ChatScreen() {
  const { id: roomId } = useLocalSearchParams();
  const { data: roomInfo } = useGetChatRoomInfo(String(roomId));
  const { isKeyboardVisible } = useKeyboard();
  const navigation = useNavigation();
  const insets = useSafeAreaInsets();

  useFocusEffect(() => {
    navigation.setOptions({
      title: roomInfo?.name,
      headerLeft: () => (
        <Pressable onPress={() => router.back()}>
          <Ionicons name="chevron-back" size={28} color="black" />
        </Pressable>
      ),
    });
  });

  return (
    <SafeAreaView style={styles.container} edges={["bottom"]}>
      <KeyboardAvoidingView
        style={styles.keyboardAvoidingView}
        behavior={Platform.OS === "ios" ? "padding" : undefined}
        keyboardVerticalOffset={
          Platform.OS === "ios" || isKeyboardVisible ? 100 : insets.bottom //TODO: this need adjust in both platform
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
