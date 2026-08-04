import {
  Modal,
  Pressable,
  StyleSheet,
  useWindowDimensions,
  View,
} from "react-native";
import { SafeAreaView } from "react-native-safe-area-context";
import { Ionicons } from "@expo/vector-icons";
import { Image } from "expo-image";
import { getKVStore } from "@/db/storage";

import {
  Gesture,
  GestureDetector,
  GestureHandlerRootView,
} from "react-native-gesture-handler";
import Animated, {
  useAnimatedStyle,
  useSharedValue,
  withTiming,
  withDecay, // 1. Imported withDecay
} from "react-native-reanimated";

interface ImageModalProps {
  imageContent: string | null;
  onClose: () => void;
}

const AnimatedImage = Animated.createAnimatedComponent(Image);

export default function ReanimatedImageModal({
  imageContent,
  onClose,
}: ImageModalProps) {
  const { width, height } = useWindowDimensions();

  const scale = useSharedValue(1);
  const savedScale = useSharedValue(1);
  const translateX = useSharedValue(0);
  const translateY = useSharedValue(0);
  const savedTranslateX = useSharedValue(0);
  const savedTranslateY = useSharedValue(0);

  const handleClose = () => {
    scale.value = 1;
    savedScale.value = 1;
    translateX.value = 0;
    translateY.value = 0;
    savedTranslateX.value = 0;
    savedTranslateY.value = 0;
    onClose();
  };

  const pinchGesture = Gesture.Pinch()
    .onStart(() => {
      savedScale.value = scale.value;
      savedTranslateX.value = translateX.value;
      savedTranslateY.value = translateY.value;
    })
    .onUpdate((e) => {
      scale.value = Math.min(Math.max(savedScale.value * e.scale, 0.5), 10);
    })
    .onEnd(() => {
      if (scale.value < 1) {
        scale.value = withTiming(1);
        translateX.value = withTiming(0);
        translateY.value = withTiming(0);
      } else {
        const maxTranslateX = (width * (scale.value - 1)) / 2;
        const maxTranslateY = (height * (scale.value - 1)) / 2;

        const clampedX = Math.max(
          -maxTranslateX,
          Math.min(translateX.value, maxTranslateX),
        );
        const clampedY = Math.max(
          -maxTranslateY,
          Math.min(translateY.value, maxTranslateY),
        );
        translateX.value = withTiming(clampedX);
        translateY.value = withTiming(clampedY);
      }
    });

  const panGesture = Gesture.Pan()
    .onStart(() => {
      savedTranslateX.value = translateX.value;
      savedTranslateY.value = translateY.value;
    })
    .onUpdate((e) => {
      if (scale.value > 1) {
        const maxTranslateX = (width * (scale.value - 1)) / 2;
        const maxTranslateY = (height * (scale.value - 1)) / 2;

        const nextTranslateX = savedTranslateX.value + e.translationX;
        const nextTranslateY = savedTranslateY.value + e.translationY;

        translateX.value = Math.max(
          -maxTranslateX,
          Math.min(nextTranslateX, maxTranslateX),
        );
        translateY.value = Math.max(
          -maxTranslateY,
          Math.min(nextTranslateY, maxTranslateY),
        );
      }
    })
    .onEnd((e) => {
      if (scale.value > 1) {
        const maxTranslateX = (width * (scale.value - 1)) / 2;
        const maxTranslateY = (height * (scale.value - 1)) / 2;
        translateX.value = withDecay({
          velocity: e.velocityX,
          clamp: [-maxTranslateX, maxTranslateX],
        });
        translateY.value = withDecay({
          velocity: e.velocityY,
          clamp: [-maxTranslateY, maxTranslateY],
        });
      }
    });

  const composedGestures = Gesture.Simultaneous(pinchGesture, panGesture);

  const animatedStyle = useAnimatedStyle(() => {
    return {
      transform: [
        { translateX: translateX.value },
        { translateY: translateY.value },
        { scale: scale.value },
      ],
    };
  });

  return (
    <Modal
      visible={imageContent !== null}
      transparent={true}
      animationType="fade"
      onRequestClose={handleClose}
    >
      <GestureHandlerRootView style={styles.fullScreenBackground}>
        <SafeAreaView style={styles.fullScreenWrapper}>
          <Pressable style={styles.closeButton} onPress={handleClose}>
            <Ionicons name="close" size={32} color="white" />
          </Pressable>

          {imageContent && (
            <GestureDetector gesture={composedGestures}>
              <View style={styles.gestureContainer}>
                <AnimatedImage
                  style={[styles.fullScreenImage, animatedStyle]}
                  source={getKVStore(imageContent)}
                  contentFit="contain"
                  cachePolicy="memory"
                />
              </View>
            </GestureDetector>
          )}
        </SafeAreaView>
      </GestureHandlerRootView>
    </Modal>
  );
}

const styles = StyleSheet.create({
  fullScreenBackground: {
    flex: 1,
    backgroundColor: "rgba(0, 0, 0, 0.9)",
  },
  fullScreenWrapper: {
    flex: 1,
    width: "100%",
    height: "100%",
  },
  gestureContainer: {
    flex: 1,
    justifyContent: "center",
    alignItems: "center",
    overflow: "hidden",
  },
  fullScreenImage: {
    width: "100%",
    height: "100%",
  },
  closeButton: {
    position: "absolute",
    top: 60,
    right: 16,
    zIndex: 10,
    padding: 8,
    backgroundColor: "rgba(0, 0, 0, 0.5)",
    borderRadius: 20,
  },
});
