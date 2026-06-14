import { Platform, StyleSheet, View } from "react-native";
import { AppleMaps, type CameraMoveEvent, GoogleMaps } from "expo-maps";
import { useEffect, useMemo, useState } from "react";
import { gridDisk, latLngToCell } from "h3-js";
import { useMapOfflineConversations } from "@/hooks/useConversation";
import {
  getCurrentPositionAsync,
  PermissionStatus,
  requestForegroundPermissionsAsync,
} from "expo-location";
import OfflineConversationDetail from "@/components/conversation/OfflineConversationDetail";

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
  const [clickedMarkerId, setClickedMarkerId] = useState<string>("");
  console.log(clickedMarkerId);
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
  const markers = useMemo(
    () =>
      zoom >= (Platform.OS === "ios" ? 12 : 14)
        ? res7Datas.map((data) => ({
            id: data.id,
            coordinates: {
              latitude: data.lat,
              longitude: data.lng,
            },
            title: data.writtenBy,
          }))
        : zoom >= (Platform.OS === "ios" ? 10 : 12)
          ? res5Datas.map((data) => ({
              id: data.id,
              coordinates: {
                latitude: data.lat,
                longitude: data.lng,
              },
              title: data.writtenBy,
            }))
          : [],
    [zoom, res7Datas, res5Datas],
  );

  const handleCameraMove = (event: CameraMoveEvent) => {
    // console.log(event.zoom);
    setZoom(event.zoom);
    const lat = event.coordinates.latitude;
    const lng = event.coordinates.longitude;
    if (event.zoom >= (Platform.OS === "ios" ? 12 : 14) && lat && lng) {
      const h3Index = latLngToCell(lat, lng, 7);
      const h3Indexes = gridDisk(h3Index, 1);
      setH3Res7(h3Indexes);
      return;
    }
    if (event.zoom >= (Platform.OS === "ios" ? 10 : 12) && lat && lng) {
      const h3Index = latLngToCell(lat, lng, 5);
      const h3Indexes = gridDisk(h3Index, 1);
      setH3Res5(h3Indexes);
    }
  };

  return (
    <View style={{ flex: 1 }}>
      {Platform.OS === "ios" ? (
        <AppleMaps.View
          style={{ flex: 1 }}
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
          onMarkerClick={(event) => {
            console.log("marker clicked");
            setClickedMarkerId(event.id ?? "");
          }}
          onMapClick={() => {
            setClickedMarkerId("");
          }}
        />
      ) : (
        <GoogleMaps.View
          style={{ flex: 1 }}
          cameraPosition={{
            coordinates: {
              latitude: initGeoInfo?.lat,
              longitude: initGeoInfo?.lng,
            },
            zoom: initGeoInfo?.zoom ?? 1,
          }}
          onCameraMove={(event) => {
            handleCameraMove(event);
          }}
          markers={markers}
          onMarkerClick={(event) => {
            setClickedMarkerId(event.id ?? "");
          }}
          onMapClick={() => {
            setClickedMarkerId("");
          }}
        />
      )}
      {clickedMarkerId && <OfflineConversationDetail id={clickedMarkerId} />}
    </View>
  );
}

const styles = StyleSheet.create({
  container: {},
});
