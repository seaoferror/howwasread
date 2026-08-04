import { Modal, Platform, Pressable, ScrollView, StyleSheet, View } from "react-native";
import { SafeAreaView } from "react-native-safe-area-context";
import { Ionicons } from "@expo/vector-icons";
import { Image } from "expo-image";
import { getKVStore } from "@/db/storage";
import Zoom from "react-native-zoom-reanimated";

interface ImageModalProps {
  imageContent: string | null;
  onClose: () => void;
}

export default function ImageModal({ imageContent, onClose }: ImageModalProps) {
  return (
    <Modal
      visible={imageContent !== null}
      transparent={true}
      animationType="fade"
      onRequestClose={onClose}
    >
      <View style={styles.fullScreenBackground}>
        <SafeAreaView style={styles.fullScreenWrapper}>
          <Pressable style={styles.closeButton} onPress={onClose}>
            <Ionicons name="close" size={32} color="white" />
          </Pressable>

          {imageContent &&
            (Platform.OS === "android" ? (
              <Zoom>
                <Image
                  style={styles.fullScreenImage}
                  source={getKVStore(imageContent)}
                  contentFit="contain"
                  cachePolicy="memory"
                />
              </Zoom>
            ) : (
              <ScrollView
                contentContainerStyle={{ flex: 1 }}
                maximumZoomScale={10}
                minimumZoomScale={1}
                bouncesZoom={true}
                centerContent={true}
                showsHorizontalScrollIndicator={false}
                showsVerticalScrollIndicator={false}
              >
                <Image
                  style={styles.fullScreenImage}
                  source={getKVStore(imageContent)}
                  contentFit="contain"
                  cachePolicy="memory"
                />
              </ScrollView>
            ))}
        </SafeAreaView>
      </View>
    </Modal>
  );
}

const styles = StyleSheet.create({
  fullScreenBackground: {
    flex: 1,
    backgroundColor: "rgba(0, 0, 0, 0.9)",
    justifyContent: "center",
    alignItems: "center",
  },
  fullScreenWrapper: {
    width: "100%",
    height: "100%",
  },
  fullScreenImage: {
    flex: 1,
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
