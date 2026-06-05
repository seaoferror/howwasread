import { Platform, StyleSheet, View } from "react-native";
import { AppleMaps, type CameraMoveEvent, GoogleMaps } from "expo-maps";
import { useMemo, useState } from "react";
import { latLngToCell, gridDisk } from "h3-js";
import { useMapOfflineConversations } from "@/hooks/useConversation";

export default function OfflineConversationMap() {
  const [h3Res5, setH3Res5] = useState<string[]>([]);
  const [h3Res7, setH3Res7] = useState<string[]>([]);
  const [zoom, setZoom] = useState<number>(12);
  const res5Datas = useMapOfflineConversations({
    resolution: 5,
    h3Indexes: h3Res5,
  });
  const res7Datas = useMapOfflineConversations({
    resolution: 7,
    h3Indexes: h3Res7,
  });

  const markers = useMemo(
    () =>
      zoom >= 12
        ? res7Datas.map((data) => ({
            id: data.id,
            coordinates: {
              latitude: data.lat,
              longitude: data.lng,
            },
            title: data.writtenBy,
          }))
        : zoom >= 10
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
    setZoom(event.zoom);
    const lat = event.coordinates.latitude;
    const lng = event.coordinates.longitude;
    if (event.zoom >= 12 && lat && lng) {
      const h3Index = latLngToCell(lat, lng, 7);
      const h3Indexes = gridDisk(h3Index, 1);
      setH3Res7(h3Indexes);
      return;
    }
    if (event.zoom >= 10 && lat && lng) {
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
          onCameraMove={(event) => {
            handleCameraMove(event);
          }}
          markers={markers}
          onMarkerClick={(event) => {
            event.id;
          }}
        />
      ) : (
        <GoogleMaps.View style={{ flex: 1 }} />
      )}
    </View>
  );
}

const styles = StyleSheet.create({
  container: {},
});
