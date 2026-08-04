import { useState } from "react";
import { Pressable, StyleSheet, Text, View } from "react-native";
import { useAudioPlayer, useAudioPlayerStatus } from "expo-audio";
import { Ionicons } from "@expo/vector-icons";
import { colors } from "@/constants";
import { formatToMinuteSecond } from "@/util/time";

export default function VoiceMessage({
  url,
  onLongPress,
}: {
  url: string;
  onLongPress: () => void;
}) {
  const player = useAudioPlayer(url);
  const status = useAudioPlayerStatus(player);
  const [trackWidth, setTrackWidth] = useState(0);
  const duration = Number.isFinite(status.duration) ? status.duration : 0;
  const currentTime = Number.isFinite(status.currentTime)
    ? status.currentTime
    : 0;
  const progress = duration > 0 ? Math.min(currentTime / duration, 1) : 0;

  const togglePlay = async () => {
    if (status.playing) {
      player.pause();
      await player.seekTo(0);
    } else {
      if (duration > 0 && currentTime >= duration - 0.1) {
        await player.seekTo(0);
      }
      player.play();
    }
  };

  const handleSeek = async (event: any) => {
    if (duration === 0 || trackWidth === 0) return;
    const touchX = event.nativeEvent.locationX;
    const percentage = Math.max(0, Math.min(touchX / trackWidth, 1));
    const seekTime = percentage * duration;
    await player.seekTo(seekTime);
  };

  return (
    <Pressable
      style={styles.audioContainer}
      onPress={togglePlay}
      onLongPress={onLongPress}
    >
      <View style={styles.iconBox}>
        <Ionicons
          name={status.playing ? "stop" : "play"}
          size={18}
          color={colors.GRAY_900}
        />
      </View>

      <View style={styles.contentBox}>
        <Pressable
          style={styles.progressTrack}
          onLayout={(e) => setTrackWidth(e.nativeEvent.layout.width)}
          onPress={handleSeek}
          hitSlop={{ top: 10, bottom: 10 }}
        >
          <View
            style={[styles.progressFill, { width: `${progress * 100}%` }]}
            pointerEvents="none"
          />
        </Pressable>

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
    paddingHorizontal: 8,
    borderRadius: 7,
    backgroundColor: colors.SAND_150,
    borderWidth: 1,
    borderColor: colors.ORANGE_100,
  },
  iconBox: {
    width: 36,
    height: 36,
    borderRadius: 18,
    alignItems: "center",
    justifyContent: "center",
    backgroundColor: colors.SAND_200,
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
