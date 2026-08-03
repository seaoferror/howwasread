import { StyleSheet } from "react-native";
import { SafeAreaView } from "react-native-safe-area-context";
import InitRoute from "@/components/auth/InitRoute";
import { Image } from "expo-image";

const AppIcon = require("@/assets/images/icon_transparent_background.png");

export default function InitScreen() {
  return (
    <InitRoute>
      <SafeAreaView style={styles.container}>
        <Image style={styles.icon} source={AppIcon} />
      </SafeAreaView>
    </InitRoute>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    justifyContent: "center",
    alignItems: "center",
    backgroundColor: "#082247",
  },
  icon: {
    alignContent: "center",
    justifyContent: "center",
    width: 160,
    height: 120,
  },
});
