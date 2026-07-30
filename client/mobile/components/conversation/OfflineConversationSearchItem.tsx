import { Pressable, StyleSheet, Text } from "react-native";
import { colors } from "@/constants";
import { OfflineConversationSearchResponse } from "@/types/conversation";
import { router } from "expo-router";

interface OfflineConversationSearchProps {
  conversation: OfflineConversationSearchResponse;
}

export default function OfflineConversationSearchItem({
  conversation,
}: OfflineConversationSearchProps) {
  return (
    <Pressable
      style={styles.content}
      onPress={() =>
        router.push({
          pathname: "/offline/[id]",
          params: { id: conversation.id },
        })
      }
    >
      <Text style={styles.when}>
        {new Intl.DateTimeFormat("en-US", {
          weekday: "short",
          year: "numeric",
          month: "short",
          day: "numeric",
          hour: "2-digit",
          minute: "2-digit",
          hourCycle: "h12",
        })
          .format(new Date(conversation.time))
          .replace(/\sat\s/, " ")}
      </Text>
      {conversation.novel && (
        <Text style={styles.detail}>Novel: {conversation.novel}</Text>
      )}
      {conversation.shortStory && (
        <Text style={styles.detail}>
          Short story: {conversation.shortStory}
        </Text>
      )}
      {conversation.poem && (
        <Text style={styles.detail}>Poem: {conversation.poem}</Text>
      )}
      {conversation.play && (
        <Text style={styles.detail}>Play: {conversation.play}</Text>
      )}
      {conversation.film && (
        <Text style={styles.detail}>Film: {conversation.film}</Text>
      )}
      <Text style={styles.detail}>Written by: {conversation.writtenBy}</Text>
    </Pressable>
  );
}

const styles = StyleSheet.create({
  container: {},
  content: {
    padding: 16,
    gap: 17,
  },
  when: {
    fontSize: 19,
    color: colors.BLACK,
    fontWeight: 500,
    marginVertical: 6,
  },
  detail: {
    fontSize: 17,
    fontWeight: 300,
  },
});
