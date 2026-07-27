import {
  Keyboard,
  Platform,
  Pressable,
  StatusBar,
  StyleSheet,
  Text,
  View,
} from "react-native";
import { AppleMaps, type CameraMoveEvent, GoogleMaps } from "expo-maps";
import { useEffect, useMemo, useState } from "react";
import { gridDisk, latLngToCell } from "h3-js";
import { Feather } from "@expo/vector-icons";
import {
  useGetBlockedConversations,
  useMapOfflineConversations,
  useSearchOfflineConversations,
} from "@/hooks/useConversation";
import {
  getCurrentPositionAsync,
  PermissionStatus,
  requestForegroundPermissionsAsync,
} from "expo-location";
import { router } from "expo-router";
import SearchInput from "@/components/conversation/SearchInput";
import { colors } from "@/constants";

export default function OfflineConversationMap() {
  const [h3Res5, setH3Res5] = useState<string[]>([]);
  const [h3Res7, setH3Res7] = useState<string[]>([]);
  const [zoom, setZoom] = useState<number>(1);
  const [initGeoInfo, setInitGeoInfo] = useState<{
    lat: number;
    lng: number;
    zoom: number;
  } | null>(null);
  const res5Datas = useMapOfflineConversations({
    resolution: 5,
    h3Indexes: h3Res5,
  });
  const res7Datas = useMapOfflineConversations({
    resolution: 7,
    h3Indexes: h3Res7,
  });
  const { data: blockedConversations } = useGetBlockedConversations();
  const [keyword, setKeyword] = useState("");
  const [submitKeyword, setSubmitKeyword] = useState("");
  const [showRetry, setShowRetry] = useState(false);
  const [submitH3Res5, setSubmitH3Res5] = useState<string[]>([]);
  const [submitH3Res7, setSubmitH3Res7] = useState<string[]>([]);
  const searchRes5Datas = useSearchOfflineConversations({
    input: submitKeyword,
    resolution: 5,
    h3Indexes: submitH3Res5,
  });
  const searchRes7Datas = useSearchOfflineConversations({
    input: submitKeyword,
    resolution: 7,
    h3Indexes: submitH3Res7,
  });

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
    return searchRes7Datas.length > 0
      ? searchRes7Datas
          .filter((data) => !blockedIds.has(data.id))
          .map((data) => ({
            id: data.id,
            coordinates: {
              latitude: data.lat,
              longitude: data.lng,
            },
            title: data.writtenBy,
          }))
      : searchRes5Datas.length > 0
        ? searchRes5Datas
            .filter((data) => !blockedIds.has(data.id))
            .map((data) => ({
              id: data.id,
              coordinates: {
                latitude: data.lat,
                longitude: data.lng,
              },
              title: data.writtenBy,
            }))
        : zoom >= (Platform.OS === "ios" ? 12 : 14)
          ? res7Datas
              .filter((data) => !blockedIds.has(data.id))
              .map((data) => ({
                id: data.id,
                coordinates: {
                  latitude: data.lat,
                  longitude: data.lng,
                },
                title: data.writtenBy,
              }))
          : zoom >= (Platform.OS === "ios" ? 10 : 12)
            ? res5Datas
                .filter((data) => !blockedIds.has(data.id))
                .map((data) => ({
                  id: data.id,
                  coordinates: {
                    latitude: data.lat,
                    longitude: data.lng,
                  },
                  title: data.writtenBy,
                }))
            : [];
  }, [
    blockedConversations,
    searchRes7Datas,
    searchRes5Datas,
    zoom,
    res7Datas,
    res5Datas,
  ]);

  const handleCameraMove = (event: CameraMoveEvent) => {
    setShowRetry(true);

    setZoom(event.zoom);
    const lat = event.coordinates.latitude;
    const lng = event.coordinates.longitude;
    if (zoom >= (Platform.OS === "ios" ? 12 : 14) && lat && lng) {
      const h3Index = latLngToCell(lat, lng, 7);
      const h3Indexes = gridDisk(h3Index, 1);
      setH3Res7(h3Indexes);
      return;
    }
    if (zoom >= (Platform.OS === "ios" ? 10 : 12) && lat && lng) {
      const h3Index = latLngToCell(lat, lng, 5);
      const h3Indexes = gridDisk(h3Index, 1);
      setH3Res5(h3Indexes);
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
            setSubmitH3Res5([]);
            setSubmitH3Res7(h3Res7);
            setSubmitKeyword(keyword);
          }}
          submitKeyWord={submitKeyword}
          onCancel={() => {
            setSubmitKeyword("");
            setKeyword("");
          }}
        />
      </View>

      {showRetry && submitKeyword && (
        <View style={styles.retrySearchContainer}>
          <Pressable
            style={styles.retrySearchButton}
            onPress={() => {
              setShowRetry(false);
              if (zoom >= (Platform.OS === "ios" ? 12 : 14)) {
                setSubmitH3Res7(h3Res7);
                setSubmitH3Res5([]);
              }
              if (zoom >= (Platform.OS === "ios" ? 10 : 12)) {
                setSubmitH3Res5(h3Res5);
                setSubmitH3Res7([]);
              }
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
      <View>
        {searchRes7Datas
          ? searchRes7Datas.map((data, idx) => (
              <View key={idx} style={styles.content}>
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
                    .format(new Date(data.time))
                    .replace(/\sat\s/, " ")}
                </Text>
                {data.novel && (
                  <Text style={styles.detail}>Novel: {data.novel}</Text>
                )}
                {data.shortStory && (
                  <Text style={styles.detail}>
                    Short story: {data.shortStory}
                  </Text>
                )}
                {data.poem && (
                  <Text style={styles.detail}>Poem: {data.poem}</Text>
                )}
                {data.play && (
                  <Text style={styles.detail}>Play: {data.play}</Text>
                )}
                {data.film && (
                  <Text style={styles.detail}>Film: {data.film}</Text>
                )}
                <Text style={styles.detail}>Written by: {data.writtenBy}</Text>
              </View>
            ))
          : searchRes5Datas
            ? searchRes5Datas.map((data, idx) => (
                <View key={idx} style={styles.content}>
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
                      .format(new Date(data.time))
                      .replace(/\sat\s/, " ")}
                  </Text>
                  {data.novel && (
                    <Text style={styles.detail}>Novel: {data.novel}</Text>
                  )}
                  {data.shortStory && (
                    <Text style={styles.detail}>
                      Short story: {data.shortStory}
                    </Text>
                  )}
                  {data.poem && (
                    <Text style={styles.detail}>Poem: {data.poem}</Text>
                  )}
                  {data.play && (
                    <Text style={styles.detail}>Play: {data.play}</Text>
                  )}
                  {data.film && (
                    <Text style={styles.detail}>Film: {data.film}</Text>
                  )}
                  <Text style={styles.detail}>
                    Written by: {data.writtenBy}
                  </Text>
                </View>
              ))
            : null}
      </View>
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
});
