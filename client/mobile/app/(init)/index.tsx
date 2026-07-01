import { Text } from "react-native";
import { SafeAreaView } from "react-native-safe-area-context";
import InitRoute from "@/components/auth/InitRoute";

export default function InitScreen() {
  return (
    <InitRoute>
      <SafeAreaView>
        <Text>loading</Text>
      </SafeAreaView>
    </InitRoute>
  );
}
//
// const styles = StyleSheet.create({
// });
