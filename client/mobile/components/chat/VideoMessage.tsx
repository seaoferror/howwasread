import { Pressable, StyleSheet } from "react-native";
import { useVideoPlayer, VideoView } from "expo-video";
import { useState } from "react";
import { colors } from "@/constants";

export default function VideoMessage({ url, onLongPress }: { url: string, onLongPress: () => void }) {
  const player = useVideoPlayer(url, (player) => {
    player.loop = true;
  });
  const [isPlaying, setIsPlaying] = useState<boolean>(false);

  return (
    <Pressable
      onPress={() => {
        setIsPlaying(!isPlaying);
        if (isPlaying) {
          player.pause();
          return;
        }
        player.seekBy(0);
        player.play();
      }}
      onLongPress={onLongPress}
    >
      <VideoView style={styles.media} player={player} />
    </Pressable>
  );
}

const styles = StyleSheet.create({
  media: {
    width: 220,
    height: 220,
    borderRadius: 4,
    backgroundColor: colors.GRAY_100,
  },
});
