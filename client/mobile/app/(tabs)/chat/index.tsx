import { Text } from "react-native";
import { SafeAreaView } from "react-native-safe-area-context";
import * as SQLite from "expo-sqlite";
import { useFocusEffect } from "expo-router";
import { useSQLiteContext } from "expo-sqlite";

export default function ChatScreen() {
  const db = useSQLiteContext();
  SQLite.addDatabaseChangeListener((event) => {
    event.rowId;
  });
  useFocusEffect(() => {
    const wrapper = async () => {
      db.getAllSync(`SELECT * FROM message`)
    };
    wrapper();
  });
  return (
    <SafeAreaView>
      <Text>first</Text>
    </SafeAreaView>
  );
}
//
// const styles = StyleSheet.create({
// });
