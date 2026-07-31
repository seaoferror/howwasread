import { Pressable, StyleSheet, Text } from "react-native";
import { colors } from "@/constants";
import { OfflineConversationSearchResponse } from "@/types/conversation";
import RenderHtml from "@native-html/render";

interface OfflineConversationSearchProps {
  conversation: OfflineConversationSearchResponse;
  onPress: () => void;
}

export default function OfflineConversationSearchItem({
  conversation,
  onPress,
}: OfflineConversationSearchProps) {
  return (
    <Pressable style={styles.content} onPress={onPress}>
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
        <RenderHtml
          source={{ html: `<b>Novel:</b> ${conversation.novel}` }}
          baseStyle={styles.detail}
        />
      )}
      {conversation.shortStory && (
        <RenderHtml
          source={{ html: `<b>Short story:</b> ${conversation.shortStory}` }}
          baseStyle={styles.detail}
        />
      )}
      {conversation.poem && (
        <RenderHtml
          source={{ html: `<b>Poem:</b> ${conversation.poem}` }}
          baseStyle={styles.detail}
        />
      )}
      {conversation.play && (
        <RenderHtml
          source={{ html: `<b>Play:</b> ${conversation.play}` }}
          baseStyle={styles.detail}
        />
      )}
      {conversation.film && (
        <RenderHtml
          source={{ html: `<b>Film:</b> ${conversation.film}` }}
          baseStyle={styles.detail}
        />
      )}
      <RenderHtml
        source={{ html: `<b>Written by:</b> ${conversation.writtenBy}` }}
        baseStyle={styles.detail}
      />
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
