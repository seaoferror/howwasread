import { useEffect, useRef, useState } from "react";
import {
  FlatList,
  Keyboard,
  Platform,
  Pressable,
  StatusBar,
  StyleSheet,
  View,
} from "react-native";
import { colors } from "@/constants";
import { useFocusEffect, useScrollToTop } from "expo-router";
import {
  useGetBlockedConversations,
  useGetInfiniteOnlineConversations,
} from "@/hooks/useConversation";
import SearchInput from "@/components/conversation/SearchInput";
import OnlineConversationItem from "@/components/conversation/OnlineConversationItem";
import { TrueSheet } from "@lodev09/react-native-true-sheet";
import { Ionicons } from "@expo/vector-icons";
import OnlineConversationDetail from "@/components/conversation/OnlineConversationDetail";

export default function OnlineConversationList({
  isActive,
}: {
  isActive: boolean;
}) {
  const [keyword, setKeyword] = useState("");
  const [submitKeyword, setSubmitKeyword] = useState("");
  const { data, fetchNextPage, hasNextPage, isFetchingNextPage, refetch } =
    useGetInfiniteOnlineConversations(submitKeyword);
  const { data: blockedConversations } = useGetBlockedConversations();
  const [isRefreshing, setIsRefreshing] = useState(false);
  const ref = useRef<FlatList | null>(null);
  useScrollToTop(ref);
  const sheet = useRef<TrueSheet>(null);
  const [detailId, setDetailId] = useState("");

  useEffect(() => {
    if (!isActive) {
      sheet.current?.dismiss();
    }
  }, [isActive]);

  useFocusEffect(() => {
    return () => {
      sheet.current?.dismiss();
    };
  });

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

  const handleSearch = () => {
    Keyboard.dismiss();
    setSubmitKeyword(keyword);
  };

  const handleSearchCancel = () => {
    setSubmitKeyword("");
    setKeyword("");
  };

  const handleSheetCloseButtonPress = () => {
    setDetailId("");
    sheet.current?.dismiss();
  };

  return (
    <View style={{ flex: 1 }}>
      <View style={styles.inputContainer}>
        <SearchInput
          placeholder="Search conversation"
          value={keyword}
          onChangeText={(text) => setKeyword(text)}
          onSubmit={() => {
            handleSearch();
          }}
          submitKeyword={submitKeyword}
          keyword={keyword}
          onCancel={() => {
            handleSearchCancel();
          }}
        />
      </View>
      <FlatList
        ref={ref}
        data={data?.pages.flat()}
        renderItem={({ item }) => {
          if (blockedConversations) {
            for (const c of blockedConversations) {
              if (c.id === String(item.id)) {
                return null;
              }
            }
          }
          return (
            <OnlineConversationItem
              conversation={item}
              onPress={() => {
                setDetailId(item.id);
                sheet.current?.present();
              }}
            />
          );
        }}
        keyExtractor={(item) => String(item.id)}
        contentContainerStyle={styles.contentContainer}
        onEndReached={handleEndReached}
        onEndReachedThreshold={0.5}
        refreshing={isRefreshing}
        onRefresh={handleRefresh}
      />
      <TrueSheet
        ref={sheet}
        detents={[0.123, 0.7]}
        dismissible={false}
        dimmed={false}
        backgroundColor={colors.SAND_100}
      >
        <Pressable
          style={styles.closeButton}
          onPress={() => {
            handleSheetCloseButtonPress();
          }}
        >
          <Ionicons name="close" size={20} color="white" />
        </Pressable>
        <OnlineConversationDetail id={detailId} />
      </TrueSheet>
    </View>
  );
}

const styles = StyleSheet.create({
  contentContainer: {
    paddingVertical: 12,
    backgroundColor: colors.SAND_150,
    gap: 12,
  },
  inputContainer: {
    zIndex: 1,
    marginTop: 17,
    paddingHorizontal: 10,
    gap: 8,
    backgroundColor: "transparent",
    flexDirection: "row",
    paddingTop: Platform.OS === "android" ? StatusBar.currentHeight : 0,
    height: 44,
  },
  closeButton: {
    position: "absolute",
    top: 10,
    right: 16,
    zIndex: 1,
    padding: 8,
    backgroundColor: "rgba(0, 0, 0, 0.5)",
    borderRadius: 100,
  },
});
