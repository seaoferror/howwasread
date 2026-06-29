import { Alert, Keyboard, Pressable, StyleSheet, View } from "react-native";
import InputField from "@/components/InputField";
import { Ionicons } from "@expo/vector-icons";
import { colors } from "@/constants";
import { useState } from "react";
import { useLocalSearchParams } from "expo-router";
import {
  useCheckBlock,
  useGeneratePresignedURL,
  useGetChatRoomInfo,
  useSendMessage,
} from "@/hooks/useChat";
import {
  launchImageLibraryAsync,
  requestMediaLibraryPermissionsAsync,
  VideoExportPreset,
} from "expo-image-picker";
import { uploadToS3 } from "@/api/chat";
import Toast from "react-native-toast-message";
import VoiceInput from "@/components/chat/VoiceInput";

export default function MessageInput() {
  const { id: roomId } = useLocalSearchParams();
  const { data: roomInfo } = useGetChatRoomInfo(String(roomId));
  const { data } = useCheckBlock(String(roomId));
  const sendMessageMutation = useSendMessage();
  const presignedURLMutation = useGeneratePresignedURL();
  const [textContent, setTextContent] = useState("");
  const [isVoice, setIsVoice] = useState(false);

  const handleSendMessage = (contentType: string, contents: string[]) => {
    const message = {
      toIdType: String(roomInfo?.type),
      toId: String(roomId),
      contentType: contentType,
      contents: contents,
    };
    console.log(message);

    if (contentType === "text") {
      setTextContent("");
    }

    sendMessageMutation.mutate(message);
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
      quality: 0.9,
      videoExportPreset: VideoExportPreset.H264_1920x1080,
      videoMaxDuration: 30,
      shouldDownloadFromNetwork: true,
      allowsMultipleSelection: true,
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
        handleFileMessage(
          "video",
          videoAssets[0].mimeType ?? "video/mp4",
          videoAssets.map((asset) => {
            return asset.uri;
          }),
        ),
      );
    }
    if (imageAssets.length > 0) {
      uploadTasks.push(
        handleFileMessage(
          "image",
          imageAssets[0].mimeType ?? "image/jpg",
          imageAssets.map((asset) => {
            return asset.uri;
          }),
        ),
      );
      console.log(imageAssets[0].mimeType);
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
            const tasks = data.map((res, idx) => {
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
  if (data?.isBlocked) {
    return null;
  }

  return (
    <View style={styles.container}>
      {isVoice ? (
        <VoiceInput
          setIsVoice={setIsVoice}
          handleFileMessage={handleFileMessage}
        />
      ) : (
        <InputField
          value={textContent}
          onChangeText={(text) => setTextContent(text)}
          submitBehavior="newline"
          leftChild={
            <Pressable
              style={styles.buttonContainer}
              onPress={() => handleImagePickerButton()}
            >
              <Ionicons name="images" size={20} color={colors.WHITE} />
            </Pressable>
          }
          rightChild={
            textContent.trim() ? (
              <Pressable
                style={styles.buttonContainer}
                onPress={() => handleSendMessage("text", [textContent])}
                disabled={sendMessageMutation.isPending}
              >
                <Ionicons name="send-sharp" size={20} color={colors.WHITE} />
              </Pressable>
            ) : (
              <Pressable
                style={styles.buttonContainer}
                onPress={async () => {
                  setIsVoice(true);
                }}
              >
                <Ionicons name="mic" size={20} color="black" />
              </Pressable>
            )
          }
        />
      )}
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
