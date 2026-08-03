import {
  Pressable,
  StyleSheet,
  Text,
  useWindowDimensions,
  View,
} from "react-native";
import { colors } from "@/constants";
import { type OnlineConversationFeedResponse } from "@/types/conversation";
import RenderHtml from "@native-html/render";

interface OnlineConversationSearchItemProps {
  conversation: OnlineConversationFeedResponse;
  onPress: () => void;
}

export default function OnlineConversationItem({
  conversation,
  onPress,
}: OnlineConversationSearchItemProps) {
  const { width } = useWindowDimensions();

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
        {conversation.novel && (
          <RenderHtml
            contentWidth={width}
            source={{ html: `<b>Novel:</b> ${conversation.novel}` }}
            baseStyle={styles.detail}
          />
        )}
        {conversation.shortStory && (
          <RenderHtml
            contentWidth={width}
            source={{ html: `<b>Short story:</b> ${conversation.shortStory}` }}
            baseStyle={styles.detail}
          />
        )}
        {conversation.poem && (
          <RenderHtml
            contentWidth={width}
            source={{ html: `<b>Poem:</b> ${conversation.poem}` }}
            baseStyle={styles.detail}
          />
        )}
        {conversation.play && (
          <RenderHtml
            contentWidth={width}
            source={{ html: `<b>Play:</b> ${conversation.play}` }}
            baseStyle={styles.detail}
          />
        )}
        {conversation.film && (
          <RenderHtml
            contentWidth={width}
            source={{ html: `<b>Film:</b> ${conversation.film}` }}
            baseStyle={styles.detail}
          />
        )}
        <RenderHtml
          contentWidth={width}
          source={{ html: `<b>Written by:</b> ${conversation.writtenBy}` }}
          baseStyle={styles.detail}
        />
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
