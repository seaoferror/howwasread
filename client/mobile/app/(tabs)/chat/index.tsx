import { Text } from "react-native";
import { SafeAreaView } from "react-native-safe-area-context";
import * as SQLite from "expo-sqlite";

export default function ChatScreen() {
  SQLite.addDatabaseChangeListener((event) => {
    event.rowId
  })
  return (
    <SafeAreaView>
      <Text>first</Text>
    </SafeAreaView>
  );
}
//
// const styles = StyleSheet.create({
// });
