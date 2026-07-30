import {
  FlatList,
  Keyboard,
  Platform,
  Pressable,
  StatusBar,
  StyleSheet,
  Text,
  View,
} from "react-native";
import { AppleMaps, type CameraMoveEvent, GoogleMaps } from "expo-maps";
import { useEffect, useMemo, useRef, useState } from "react";
import { gridDisk, latLngToCell } from "h3-js";
import { Feather } from "@expo/vector-icons";
import {
  useGetBlockedConversations,
  useInfiniteSearchOfflineConversations,
  useMapOfflineConversations,
} from "@/hooks/useConversation";
import {
  getCurrentPositionAsync,
  PermissionStatus,
  requestForegroundPermissionsAsync,
} from "expo-location";
import { router, useFocusEffect } from "expo-router";
import SearchInput from "@/components/conversation/SearchInput";
import { colors } from "@/constants";
import { TrueSheet } from "@lodev09/react-native-true-sheet";
import OfflineConversationSearchItem from "@/components/conversation/OfflineConversationSearchItem";

export default function OfflineConversationMap({isActive} : {isActive: boolean}) {
  const [h3Indexes, setH3Indexes] = useState<string[]>([]);
  const [resolution, setResolution] = useState<number>(7);
  const [initGeoInfo, setInitGeoInfo] = useState<{
    lat: number;
    lng: number;
    zoom: number;
  } | null>(null);
  const datas = useMapOfflineConversations({
    resolution: resolution,
    h3Indexes: h3Indexes,
  });
  const { data: blockedConversations } = useGetBlockedConversations();
  const [keyword, setKeyword] = useState("");
  const [submitKeyword, setSubmitKeyword] = useState("");
  const [showRetry, setShowRetry] = useState(false);
  const [submitH3Indexes, setSubmitH3Indexes] = useState<string[]>([]);
  const {
    data: searchData,
    fetchNextPage,
    hasNextPage,
    isFetchingNextPage,
  } = useInfiniteSearchOfflineConversations({
    input: submitKeyword,
    resolution: resolution,
    h3Indexes: submitH3Indexes,
  });
  const sheet = useRef<TrueSheet>(null);

  useEffect(() => {
    if (!isActive) {
      sheet.current?.dismiss();
    }
  }, [isActive]);

  useFocusEffect(() => {
    return () => {
      sheet.current?.dismiss();
    }
  })

  useEffect(() => {
    async function wrapper() {
      const { status } = await requestForegroundPermissionsAsync();
      if (status === PermissionStatus.DENIED) {
        return;
      }
      const location = await getCurrentPositionAsync({});
      setInitGeoInfo({
        lat: location.coords.latitude,
        lng: location.coords.longitude,
        zoom: Platform.OS === "ios" ? 15 : 17,
      });
    }
    wrapper();
  }, []);

  const markers = useMemo(() => {
    const blockedIds = new Set(
      blockedConversations?.map((block) => block.id) || [],
    );
    return searchData && searchData.pages.flat().length > 0
      ? searchData.pages
          .flat()
          .filter((data) => !blockedIds.has(data.id))
          .map((data) => ({
            id: data.id,
            coordinates: {
              latitude: data.lat,
              longitude: data.lng,
            },
            title: data.writtenBy,
          }))
      : datas
          .filter((data) => !blockedIds.has(data.id))
          .map((data) => ({
            id: data.id,
            coordinates: {
              latitude: data.lat,
              longitude: data.lng,
            },
            title: data.writtenBy,
          }));
  }, [blockedConversations, datas, searchData]);

  const handleEndReached = () => {
    if (hasNextPage && !isFetchingNextPage) {
      fetchNextPage();
    }
  };

  const handleCameraMove = (event: CameraMoveEvent) => {
    setShowRetry(true);
    const lat = event.coordinates.latitude;
    const lng = event.coordinates.longitude;
    if (event.zoom >= (Platform.OS === "ios" ? 12 : 14) && lat && lng) {
      setResolution(7);
      const h3Index = latLngToCell(lat, lng, 7);
      const h3Indexes = gridDisk(h3Index, 1);
      setH3Indexes(h3Indexes);
      return;
    }
    if (event.zoom >= (Platform.OS === "ios" ? 10 : 12) && lat && lng) {
      setResolution(5);
      const h3Index = latLngToCell(lat, lng, 5);
      const h3Indexes = gridDisk(h3Index, 1);
      setH3Indexes(h3Indexes);
    }
  };

  return (
    <View style={{ flex: 1 }}>
      <View style={styles.inputContainer}>
        <SearchInput
          placeholder="Search conversation"
          value={keyword}
          onChangeText={(text) => setKeyword(text)}
          onSubmit={() => {
            Keyboard.dismiss();
            setSubmitH3Indexes(h3Indexes);
            setSubmitKeyword(keyword);
            setResolution(7);
            sheet.current?.present();
          }}
          submitKeyWord={submitKeyword}
          onCancel={() => {
            setSubmitKeyword("");
            setKeyword("");
            sheet.current?.dismiss();
          }}
        />
      </View>

      {showRetry && submitKeyword && (
        <View style={styles.retrySearchContainer}>
          <Pressable
            style={styles.retrySearchButton}
            onPress={() => {
              sheet.current?.present();
              setShowRetry(false);
              setSubmitH3Indexes(h3Indexes);
            }}
          >
            <Feather name="rotate-cw" size={14} color="black" />
            <Text style={styles.retrySearchText}>Search this area</Text>
          </Pressable>
        </View>
      )}

      {Platform.OS === "ios" ? (
        <AppleMaps.View
          style={StyleSheet.absoluteFill}
          cameraPosition={{
            coordinates: {
              latitude: initGeoInfo?.lat,
              longitude: initGeoInfo?.lng,
            },
            zoom: initGeoInfo?.zoom,
          }}
          onCameraMove={(event) => {
            handleCameraMove(event);
          }}
          markers={markers}
          onMapClick={() => Keyboard.dismiss()}
          onMarkerClick={(event) => {
            router.push({
              pathname: "/offline/[id]",
              params: { id: String(event.id) },
            });
          }}
        />
      ) : (
        <GoogleMaps.View
          style={StyleSheet.absoluteFill}
          cameraPosition={{
            coordinates: {
              latitude: initGeoInfo?.lat,
              longitude: initGeoInfo?.lng,
            },
            zoom: initGeoInfo?.zoom ?? 1,
          }}
          onMapClick={() => Keyboard.dismiss()}
          onCameraMove={(event) => {
            handleCameraMove(event);
          }}
          uiSettings={{
            compassEnabled: true,
            myLocationButtonEnabled: true,
            zoomControlsEnabled: false,
            rotationGesturesEnabled: true,
          }}
          markers={markers}
          onMarkerClick={(event) => {
            router.push({
              pathname: "/offline/[id]",
              params: { id: String(event.id) },
            });
          }}
        />
      )}
      <TrueSheet
        ref={sheet}
        detents={["auto", 0.72]}
        dismissible={false}
        dimmed={false}
        backgroundColor={colors.SAND_100}
      >
        <FlatList
          data={searchData?.pages.flat() || []}
          ListEmptyComponent={
            <Pressable
              style={styles.emptyContainer}
              onPress={() => {
                sheet.current?.dismiss();
              }}
            >
              <Text style={styles.emptyText}>No Result</Text>
            </Pressable>
          }
          renderItem={({ item }) => {
            if (blockedConversations) {
              for (const c of blockedConversations) {
                if (c.id === String(item.id)) {
                  return null;
                }
              }
            }
            return <OfflineConversationSearchItem conversation={item} />;
          }}
          keyExtractor={(item) => String(item.id)}
          contentContainerStyle={styles.contentContainer}
          onEndReached={handleEndReached}
          onEndReachedThreshold={0.5}
        />
      </TrueSheet>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {},
  inputContainer: {
    zIndex: 1,
    marginTop: 17,
    marginRight: 75,
    paddingLeft: 10,
    gap: 8,
    backgroundColor: "transparent",
    flexDirection: "row",
    paddingTop: Platform.OS === "android" ? StatusBar.currentHeight : 0,
    height: 44,
  },
  retrySearchContainer: {
    marginTop: 20,
    alignItems: "center",
    zIndex: 10,
    marginHorizontal: 100,
  },
  retrySearchButton: {
    flexDirection: "row",
    alignItems: "center",
    backgroundColor: colors.WHITE,
    paddingVertical: 10,
    paddingHorizontal: 16,
    borderRadius: 24,
  },
  retrySearchText: {
    marginLeft: 8,
    fontSize: 14,
    fontWeight: "600",
    color: colors.BLACK,
  },
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
  ruleHeader: {
    marginTop: 6,
    fontSize: 18,
    fontWeight: 400,
  },
  contentContainer: {
    paddingVertical: 12,
    backgroundColor: colors.SAND_150,
    gap: 12,
  },
  emptyContainer: {
    paddingVertical: 40,
    alignItems: "center",
    justifyContent: "center",
  },
  emptyText: {
    fontSize: 16,
    color: colors.BLACK,
    fontWeight: "500",
  },
});
