import { Keyboard, Pressable, StyleSheet, View } from "react-native";
import InputField from "@/components/InputField";
import { Feather, Ionicons } from "@expo/vector-icons";
import { colors } from "@/constants";
import { use, useState } from "react";
import { useLocalSearchParams } from "expo-router";
import {
  useGeneratePresignedURL,
  useGetChatRoomInfo,
  useSendMessaging,
  useUploadToS3,
} from "@/hooks/useChat";
import {
  RecordingPresets,
  useAudioPlayer,
  useAudioPlayerStatus,
  useAudioRecorder,
  useAudioRecorderState,
} from "expo-audio";

export default function MessageInput() {
  const { id: roomId } = useLocalSearchParams();
  const { data: roomInfo } = useGetChatRoomInfo(String(roomId));
  const sendMessagingMutation = useSendMessaging();
  const presignedURLMutation = useGeneratePresignedURL();
  const uploadToS3Mutation = useUploadToS3();

  const [textContent, setTextContent] = useState("");
  const audioRecorder = useAudioRecorder(RecordingPresets.HIGH_QUALITY);
  const recorderState = useAudioRecorderState(audioRecorder);
  const audioPlayer = useAudioPlayer();
  const playerState = useAudioPlayerStatus(audioPlayer);

  const handleSendMessage = (contentType: string, content: string) => {
    const message = {
      toIdType: String(roomInfo?.type),
      toId: String(roomId),
      contentType: contentType,
      content: content,
    };

    sendMessagingMutation.mutate(message, {
      onSuccess: async () => {
        if (contentType === "text") {
          setTextContent("");
          return;
        }
        if (contentType === "voice") {
          await audioRecorder.prepareToRecordAsync();
          audioPlayer.remove();
          return;
        }
      },
    });
  };

  function handleMoreButton() {
    Keyboard.dismiss();
  }

  const handleFileMessage = async (
    contentType: string,
    mimeType: string,
    content: string,
  ) => {
    presignedURLMutation.mutate(
      { contentType },
      {
        onSuccess: (data) => {
          uploadToS3Mutation.mutate(
            {
              awsFields: data.fields,
              awsPresignedURL: data.url,
              localFileURI: content,
              mimeType: mimeType,
            },
            {
              onSuccess: () => {
                handleSendMessage(contentType, data.filename);
              },
            },
          );
        },
      },
    );
  };
  return (
    <View style={styles.container}>
      <InputField
        value={textContent}
        onChangeText={(text) => setTextContent(text)}
        placeholder={
          recorderState.isRecording
            ? `recording... ${recorderState.durationMillis / 1000}s`
            : playerState.isLoaded
              ? playerState.playing
                ? `□ stop ${playerState.currentTime}`
                : `▷ play ${playerState.duration}`
              : ``
        }
        onPress={
          playerState.isLoaded
            ? playerState.playing
              ? async () => {
                  audioPlayer.pause();
                  await audioPlayer.seekTo(0);
                }
              : audioPlayer.play
            : () => {}
        }
        submitBehavior="newline"
        leftChild={
          recorderState.durationMillis !== 0 && !audioPlayer.isLoaded ? (
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
                audioPlayer.remove();
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
              onPress={() => handleSendMessage("text", textContent)}
            >
              <Ionicons name="send-sharp" size={20} color={colors.WHITE} />
            </Pressable>
          ) : recorderState.durationMillis !== 0 ? (
            <Pressable
              style={styles.buttonContainer}
              onPress={async () => {
                if (!audioRecorder.uri) return;
                await handleFileMessage(
                  "voice",
                  "audio/mp4",
                  audioRecorder.uri,
                );
              }}
            >
              <Ionicons name="send-outline" size={20} color={colors.WHITE} />
            </Pressable>
          ) : recorderState.isRecording ? (
            <Pressable
              style={styles.buttonContainer}
              onPress={async () => {
                await audioRecorder.stop();
                audioPlayer.replace(audioRecorder.uri);
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
