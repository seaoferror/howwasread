import { Pressable, StyleSheet, Text, View } from "react-native";
import { colors } from "@/constants";
import { type OnlineConversationFeedResponse } from "@/types/conversation";

interface OnlineConversationSearchItemProps {
  conversation: OnlineConversationFeedResponse;
  onPress: () => void;
}

export default function OnlineConversationItem({
  conversation,
  onPress,
}: OnlineConversationSearchItemProps) {
  return (
    <Pressable style={styles.container} onPress={onPress}>
      <View style={[styles.content]}>
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
        {!!conversation.novel && (
          <Text style={styles.detail}>Novel: {conversation.novel}</Text>
        )}
        {!!conversation.shortStory && (
          <Text style={styles.detail}>
            Short story: {conversation.shortStory}
          </Text>
        )}
        {!!conversation.poem && (
          <Text style={styles.detail}>Poem: {conversation.poem}</Text>
        )}
        {!!conversation.play && (
          <Text style={styles.detail}>Play: {conversation.play}</Text>
        )}
        {!!conversation.film && (
          <Text style={styles.detail}>Film: {conversation.film}</Text>
        )}
        {!!conversation.writtenBy && (
          <Text style={styles.detail}>
            Written by: {conversation.writtenBy}
          </Text>
        )}
      </View>
    </Pressable>
  );
}

const styles = StyleSheet.create({
  container: { backgroundColor: colors.SAND_110 },
  content: {
    padding: 16,
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
