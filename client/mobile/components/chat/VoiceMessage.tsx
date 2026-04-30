import { Pressable, StyleSheet, Text, View } from "react-native";
import { useAudioPlayer, useAudioPlayerStatus } from "expo-audio";
import { Ionicons } from "@expo/vector-icons";
import { colors } from "@/constants";
import { formatToMinuteSecond } from "@/util/time";

export default function VoiceMessage({ url }: { url: string }) {
  const player = useAudioPlayer(url);
  const status = useAudioPlayerStatus(player);
  const duration = Number.isFinite(status.duration) ? status.duration : 0;
  const currentTime = Number.isFinite(status.currentTime)
    ? status.currentTime
    : 0;
  const progress = duration > 0 ? Math.min(currentTime / duration, 1) : 0;

  return (
    <Pressable
      style={styles.audioContainer}
      onPress={
        status.playing
          ? async () => {
              player.pause();
              await player.seekTo(0);
            }
          : () => {
              player.play();
            }
      }
    >
      <View style={styles.iconBox}>
        <Ionicons
          name={status.playing ? "stop" : "play"}
          size={18}
          color={colors.GRAY_900}
        />
      </View>

      <View style={styles.contentBox}>
        <Text style={styles.audioTitle}>Voice message</Text>

        <View style={styles.progressTrack}>
          <View
            style={[styles.progressFill, { width: `${progress * 100}%` }]}
          />
        </View>

        <Text style={styles.audioMeta}>
          {formatToMinuteSecond(currentTime)} / {formatToMinuteSecond(duration)}
        </Text>
      </View>
    </Pressable>
  );
}

const styles = StyleSheet.create({
  audioContainer: {
    flexDirection: "row",
    alignItems: "center",
    gap: 12,
    minWidth: 220,
    minHeight: 64,
    paddingVertical: 10,
    paddingHorizontal: 12,
    borderRadius: 14,
    backgroundColor: colors.SAND_100,
    borderWidth: 1,
    borderColor: colors.ORANGE_100,
  },
  iconBox: {
    width: 36,
    height: 36,
    borderRadius: 18,
    alignItems: "center",
    justifyContent: "center",
    backgroundColor: colors.ORANGE_100,
  },
  contentBox: {
    flex: 1,
    gap: 6,
  },
  audioTitle: {
    fontSize: 14,
    color: colors.GRAY_900,
    fontWeight: "500",
  },
  progressTrack: {
    width: "100%",
    height: 6,
    borderRadius: 999,
    backgroundColor: colors.ORANGE_100,
    overflow: "hidden",
  },
  progressFill: {
    height: "100%",
    borderRadius: 999,
    backgroundColor: colors.ORANGE_200,
  },
  audioMeta: {
    fontSize: 12,
    color: colors.GRAY_700,
    fontWeight: "500",
  },
});
