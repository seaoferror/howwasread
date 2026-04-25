import {
  Alert,
  Keyboard,
  Platform,
  Pressable,
  StyleSheet,
  View,
} from "react-native";
import InputField from "@/components/InputField";
import { Feather, Ionicons } from "@expo/vector-icons";
import { colors } from "@/constants";
import { useState } from "react";
import { useLocalSearchParams } from "expo-router";
import {
  useGeneratePresignedURL,
  useGetChatRoomInfo,
  useSendMessage,
} from "@/hooks/useChat";
import {
  RecordingPresets,
  useAudioPlayer,
  useAudioPlayerStatus,
  useAudioRecorder,
  useAudioRecorderState,
} from "expo-audio";
import {
  launchImageLibraryAsync,
  requestMediaLibraryPermissionsAsync,
} from "expo-image-picker";
import { uploadToS3 } from "@/api/chat";
import Toast from "react-native-toast-message";

export default function MessageInput() {
  const { id: roomId } = useLocalSearchParams();
  const { data: roomInfo } = useGetChatRoomInfo(String(roomId));
  const sendMessageMutation = useSendMessage();
  const presignedURLMutation = useGeneratePresignedURL();

  const [textContent, setTextContent] = useState("");
  const audioRecorder = useAudioRecorder(RecordingPresets.HIGH_QUALITY);
  const recorderState = useAudioRecorderState(audioRecorder);
  const audioPlayer = useAudioPlayer();
  const playerState = useAudioPlayerStatus(audioPlayer);

  const handleSendMessage = (contentType: string, contents: string[]) => {
    const message = {
      toIdType: String(roomInfo?.type),
      toId: String(roomId),
      contentType: contentType,
      contents: contents,
    };

    sendMessageMutation.mutate(message, {
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

  const handleImagePickerButton = async () => {
    Keyboard.dismiss();
    const permissionResult = await requestMediaLibraryPermissionsAsync();

    if (!permissionResult.granted) {
      Alert.alert(
        "Permission required",
        "Permission to access the media library is required.",
      );
      return;
    }

    const result = await launchImageLibraryAsync({
      mediaTypes: ["images", "videos"],
      allowsEditing: false,
      allowsMultipleSelection: true,
      aspect: [4, 3],
      quality: 1,
    });

    console.log(result);

    if (result.canceled) {
      return;
    }
    const videoAssets = result.assets.filter((asset) => asset.type === "video");
    const imageAssets = result.assets.filter((asset) => asset.type === "image");
    const uploadTasks = [];
    if (videoAssets.length > 0) {
      uploadTasks.push(
        videoAssets.map((asset) => {
          handleFileMessage("video", videoAssets[0].mimeType ?? "video/mp4", [asset.uri])
        })
      );
    }
    if (imageAssets.length > 0) {
      uploadTasks.push(
        handleFileMessage(
          "image",
          imageAssets[0].mimeType ??
            (Platform.OS === "android" ? "image/jpg" : "image/heic"),
          imageAssets.map((asset) => {
            return asset.uri;
          }),
        ),
      );
    }
    await Promise.all(uploadTasks);
  };

  const handleFileMessage = async (
    contentType: string,
    mimeType: string,
    content: string[],
  ) => {
    presignedURLMutation.mutate(
      { contentType: contentType, n: content.length },
      {
        onSuccess: async (data) => {
          try {
            const tasks: Promise<void>[] = data.map((res, idx) => {
              const task = uploadToS3({
                awsFields: res.fields,
                awsPresignedURL: res.url,
                localFileURI: content[idx],
                mimeType: mimeType,
                filename: res.filename,
              });
              return task;
            });
            await Promise.all(tasks);
            handleSendMessage(
              contentType,
              data.map((res) => res.filename),
            );
          } catch (error) {
            console.log(error);
            Toast.show({
              type: "error",
              text1: String(error),
            });
            //TODO: remove all file via redis by X-User-Id, so it is transaction rollback
          }
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
          !audioPlayer.isLoaded ? (
            <Pressable
              style={styles.buttonContainer}
              onPress={() => handleImagePickerButton()}
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
              onPress={() => handleSendMessage("text", [textContent])}
            >
              <Ionicons name="send-sharp" size={20} color={colors.WHITE} />
            </Pressable>
          ) : recorderState.durationMillis !== 0 ? (
            <Pressable
              style={styles.buttonContainer}
              onPress={async () => {
                if (!audioRecorder.uri) return;
                await handleFileMessage("voice", "audio/mp4", [
                  audioRecorder.uri,
                ]);
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
