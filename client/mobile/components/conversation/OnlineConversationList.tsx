import { useRef, useState } from "react";
import { FlatList, StyleSheet } from "react-native";
import { colors } from "@/constants";
import { useScrollToTop } from "expo-router";
import OnlineConversationItem from "@/components/conversation/OnlineConversationItem";
import { useGetInfiniteOnlineConversations } from "@/hooks/useConversation";

export default function OnlineConversationList() {
  const { data, fetchNextPage, hasNextPage, isFetchingNextPage, refetch } =
    useGetInfiniteOnlineConversations();
  const [isRefreshing, setIsRefreshing] = useState(false);
  const ref = useRef<FlatList | null>(null);
  useScrollToTop(ref);

  const handleRefresh = async () => {
    setIsRefreshing(true);
    await refetch();
    setIsRefreshing(false);
  };

  const handleEndReached = () => {
    if (hasNextPage && !isFetchingNextPage) {
      fetchNextPage();
    }
  };

  return (
    <FlatList
      ref={ref}
      data={data?.pages.flat()}
      renderItem={({ item }) => <OnlineConversationItem conversation={item} />}
      keyExtractor={(item) => String(item.id)}
      contentContainerStyle={styles.contentContainer}
      onEndReached={handleEndReached}
      onEndReachedThreshold={0.5}
      refreshing={isRefreshing}
      onRefresh={handleRefresh}
    />
  );
}

const styles = StyleSheet.create({
  contentContainer: {
    paddingVertical: 12,
    backgroundColor: colors.SAND_150,
    gap: 12,
  },
});
