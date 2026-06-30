import { StyleSheet, View } from "react-native";
import { useLocalSearchParams } from "expo-router";
import MemberList from "@/components/chat/MemberList";

export default function ChatRoomDetailScreen() {
  const { id: roomId } = useLocalSearchParams();

  return (
    <View>
      <MemberList roomId={String(roomId)} />
    </View>
  );
}

const styles = StyleSheet.create({
  container: {},
});
