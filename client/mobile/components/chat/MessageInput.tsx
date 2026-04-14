import { Keyboard, Pressable, StyleSheet, View } from "react-native";
import InputField from "@/components/InputField";
import { Feather, Ionicons } from "@expo/vector-icons";
import { colors } from "@/constants";
import { useState } from "react";
import { useLocalSearchParams } from "expo-router";
import { useGetChatRoomInfo, useSendMessaging } from "@/hooks/useChat";

export default function MessageInput() {
  const { id: roomId } = useLocalSearchParams();
  const { data: roomInfo } = useGetChatRoomInfo(String(roomId));
  const [textContent, setTextContent] = useState("");
  const sendMessagingMutation = useSendMessaging();

  const handleSendMessage = (contentType: string) => {
    const message = {
      toIdType: String(roomInfo?.type),
      toId: String(roomId),
      contentType: String(contentType),
      content: textContent,
    };

    sendMessagingMutation.mutate(message, {
      onSuccess: () => {
        if (contentType === "text") {
          setTextContent("");
        }
      },
    });
  };

  function handleMoreButton() {
    Keyboard.dismiss();
  }

  return (
    <View style={styles.container}>
      <InputField
        value={textContent}
        onChangeText={(text) => setTextContent(text)}
        placeholder={"recording"}
        submitBehavior="newline"
        leftChild={
          <Pressable
            style={styles.buttonContainer}
            onPress={() => handleMoreButton()}
          >
            <Feather name="plus" size={20} color={colors.WHITE} />
          </Pressable>
        }
        rightChild={
          textContent.trim() ? (
            <Pressable
              style={styles.buttonContainer}
              onPress={() => handleSendMessage("text")}
            >
              <Ionicons name="send-sharp" size={20} color={colors.WHITE} />
            </Pressable>
          ) : (
            <Pressable
              style={styles.buttonContainer}
              onLongPress={() => {}}
              onPressOut={() => {}}
            >
              <Ionicons name="mic" size={20} color="black" />
            </Pressable>
          )
        }
      />
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    backgroundColor: colors.WHITE,
    padding: 16,
    borderTopColor: colors.GRAY_200,
    borderTopWidth: StyleSheet.hairlineWidth,
    width: "100%",
  },
  buttonContainer: {
    backgroundColor: colors.SAND_300,
    padding: 8,
    borderRadius: 5,
  },
});
