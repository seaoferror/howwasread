import { Keyboard, Pressable, StyleSheet, View } from "react-native";
import InputField from "@/components/InputField";
import { Feather, Ionicons } from "@expo/vector-icons";
import { colors } from "@/constants";
import { useState } from "react";
import { useLocalSearchParams } from "expo-router";
import { useGetChatRoomInfo, useSendMessaging } from "@/hooks/useChat";
import {
  RecordingPresets,
  useAudioRecorder,
  useAudioRecorderState,
} from "expo-audio";

export default function MessageInput() {
  const { id: roomId } = useLocalSearchParams();
  const { data: roomInfo } = useGetChatRoomInfo(String(roomId));
  const [textContent, setTextContent] = useState("");
  const sendMessagingMutation = useSendMessaging();
  const audioRecorder = useAudioRecorder(RecordingPresets.HIGH_QUALITY);
  const recorderState = useAudioRecorderState(audioRecorder);

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
        placeholder={
          recorderState.isRecording
            ? `${recorderState.durationMillis / 1000}s`
            : ""
        }
        submitBehavior="newline"
        leftChild={
          recorderState.durationMillis === 0 ? (
            <Pressable
              style={styles.buttonContainer}
              onPress={() => handleMoreButton()}
            >
              <Feather name="plus" size={20} color={colors.WHITE} />
            </Pressable>
          ) : (
            <Pressable
              style={styles.buttonContainer}
              onPress={async () => {
                await audioRecorder.prepareToRecordAsync();
              }}
            >
              <Ionicons name="trash-bin" size={20} color={colors.RED_500} />
            </Pressable>
          )
        }
        rightChild={
          textContent.trim() ? (
            <Pressable
              style={styles.buttonContainer}
              onPress={() => handleSendMessage("text")}
            >
              <Ionicons name="send-sharp" size={20} color={colors.WHITE} />
            </Pressable>
          ) : recorderState.durationMillis !== 0 ? (
            <Pressable style={styles.buttonContainer}>
              <Ionicons name="send-outline" size={20} color={colors.WHITE} />
            </Pressable>
          ) : recorderState.isRecording ? (
            <Pressable
              style={styles.buttonContainer}
              onPress={async () => {
                await audioRecorder.stop();
              }}
            >
              <Ionicons name="stop" size={20} color="black" />
            </Pressable>
          ) : (
            <Pressable
              style={styles.buttonContainer}
              onPress={async () => {
                await audioRecorder.prepareToRecordAsync();
                audioRecorder.record();
              }}
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
