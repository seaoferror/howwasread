import React from "react";
import {
  Modal,
  Pressable,
  StyleSheet,
  View,
  useWindowDimensions,
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
  useSharedValue,
  useAnimatedStyle,
  withTiming,
} from "react-native-reanimated";

interface ImageModalProps {
  imageContent: string | null;
  onClose: () => void;
}

const AnimatedImage = Animated.createAnimatedComponent(Image);

export default function ImageModal({ imageContent, onClose }: ImageModalProps) {
  // 1. Get window dimensions to calculate boundaries
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
    .onUpdate((e) => {
      scale.value = Math.min(Math.max(savedScale.value * e.scale, 0.5), 10);
    })
    .onEnd(() => {
      if (scale.value < 1) {
        scale.value = withTiming(1);
        translateX.value = withTiming(0);
        translateY.value = withTiming(0);
        savedScale.value = 1;
        savedTranslateX.value = 0;
        savedTranslateY.value = 0;
      } else {
        savedScale.value = scale.value;

        // 2. If the user pinches OUT (shrinks the image) while panned to the edge,
        // we need to recalculate the bounds and snap the image back into view.
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
        savedTranslateX.value = clampedX;
        savedTranslateY.value = clampedY;
      }
    });

  const panGesture = Gesture.Pan()
    .onUpdate((e) => {
      if (scale.value > 1) {
        // 3. Calculate maximum allowed panning distance
        const maxTranslateX = (width * (scale.value - 1)) / 2;
        const maxTranslateY = (height * (scale.value - 1)) / 2;

        const nextTranslateX = savedTranslateX.value + e.translationX;
        const nextTranslateY = savedTranslateY.value + e.translationY;

        // 4. Clamp the translation values so it can't go past the edges
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
    .onEnd(() => {
      savedTranslateX.value = translateX.value;
      savedTranslateY.value = translateY.value;
    });

  const doubleTapGesture = Gesture.Tap()
    .numberOfTaps(2)
    .onEnd(() => {
      if (scale.value > 1) {
        scale.value = withTiming(1);
        translateX.value = withTiming(0);
        translateY.value = withTiming(0);
        savedScale.value = 1;
        savedTranslateX.value = 0;
        savedTranslateY.value = 0;
      } else {
        scale.value = withTiming(2);
        savedScale.value = 2;
      }
    });

  const composedGestures = Gesture.Simultaneous(
    pinchGesture,
    panGesture,
    doubleTapGesture,
  );

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
