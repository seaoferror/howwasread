import { Pressable, StyleSheet, View } from "react-native";
import { Ionicons } from "@expo/vector-icons";
import { colors } from "@/constants";
import InputField from "@/components/InputField";
import {
  RecordingOptions,
  RecordingPresets,
  useAudioPlayer,
  useAudioPlayerStatus,
  useAudioRecorder,
  useAudioRecorderState,
} from "expo-audio";
import { type Dispatch, type SetStateAction, useEffect } from "react";
import { formatToMinuteSecond } from "@/util/time";

interface VoiceInputProps {
  recordingPresets?: RecordingOptions
  setIsVoice: Dispatch<SetStateAction<boolean>>;
  handleFileMessage: (
    contentType: string,
    mimeType: string,
    content: string[],
  ) => Promise<void>;
}

export default function VoiceInput({
  recordingPresets,
  setIsVoice,
  handleFileMessage,
}: VoiceInputProps) {
  const audioRecorder = useAudioRecorder(recordingPresets ?? RecordingPresets.HIGH_QUALITY);
  const recorderState = useAudioRecorderState(audioRecorder);
  const audioPlayer = useAudioPlayer();
  const audioStatus = useAudioPlayerStatus(audioPlayer);

  useEffect(() => {
    const wrapper = async () => {
      await audioRecorder.prepareToRecordAsync();
      audioRecorder.record();
    };

    wrapper();
  }, []);

  return (
    <InputField
      value={
        recorderState.isRecording
          ? `${formatToMinuteSecond(recorderState.durationMillis / 1000)} recording...`
          : audioPlayer.playing
            ? `${formatToMinuteSecond(audioStatus.currentTime)} playing...`
            : `${formatToMinuteSecond(audioStatus.duration)}`
      }
      editable={false}
      leftChild={
        <View style={styles.buttons}>
          <Pressable
            style={styles.buttonContainer}
            onPress={async () => {
              await audioRecorder.stop();
              setIsVoice(false);
            }}
          >
            <Ionicons name="trash-bin" size={20} color={colors.RED_500} />
          </Pressable>
          {!recorderState.isRecording &&
            (audioStatus.playing ? (
              <Pressable
                style={styles.buttonContainer}
                onPress={async () => {
                  audioPlayer.pause();
                  await audioPlayer.seekTo(0);
                }}
              >
                <Ionicons name="stop" size={20} color="black" />
              </Pressable>
            ) : (
              <Pressable
                style={styles.buttonContainer}
                onPress={async () => {
                  audioPlayer.play();
                }}
              >
                <Ionicons name="play" size={20} color="black" />
              </Pressable>
            ))}
        </View>
      }
      rightChild={
        recorderState.isRecording ? (
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
              if (!audioRecorder.uri) return;
              await handleFileMessage("audio", "audio/mp4", [
                audioRecorder.uri,
              ]);
              setIsVoice(false);
            }}
          >
            <Ionicons name="send-outline" size={20} color={colors.WHITE} />
          </Pressable>
        )
      }
    />
  );
}

const styles = StyleSheet.create({
  buttonContainer: {
    backgroundColor: colors.SAND_300,
    padding: 8,
    borderRadius: 5,
  },
  buttons: {
    flexDirection: "row",
    gap: 20,
  },
});
