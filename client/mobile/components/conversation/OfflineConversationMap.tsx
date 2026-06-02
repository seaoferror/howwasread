import { Platform, StyleSheet, Text, View } from "react-native";
import { AppleMaps, CameraMoveEvent, GoogleMaps } from "expo-maps";
import { useEffect, useMemo, useRef, useState } from "react";
import {
  getCurrentPositionAsync,
  LocationObject,
  PermissionStatus,
  requestForegroundPermissionsAsync,
} from "expo-location";
import { latLngToCell, gridDisk } from "h3-js";
import { useMapOfflineConversations } from "@/hooks/useConversation";

export default function OfflineConversationMap() {
  const [initialLocation, setInitialLocation] = useState<LocationObject | null>(
    null,
  );
  const [zoom, setZoom] = useState<number>(0);
  const centerH3Res5 = useRef<string | null>(null);
  const centerH3Res7 = useRef<string | null>(null);
  const [h3Res5, setH3Res5] = useState<string[]>([]);
  const [h3Res7, setH3Res7] = useState<string[]>([]);
  const res5Datas = useMapOfflineConversations({
    resolution: 5,
    h3Indexes: h3Res5,
  });
  const res7Datas = useMapOfflineConversations({
    resolution: 7,
    h3Indexes: h3Res7,
  });
  const markers = useMemo(() => {
    if (zoom >= 12) {
      return res7Datas.map((data) => ({
        id: data.id,
        coordinates: {
          latitude: data.lat,
          longitude: data.lng,
        },
        title: data.writtenBy,
      }));
    }
    if (zoom >= 10) {
      return res5Datas.map((data) => ({
        id: data.id,
        coordinates: {
          latitude: data.lat,
          longitude: data.lng,
        },
        title: data.writtenBy,
      }));
    }
    return [];
  }, [zoom, res5Datas, res7Datas]);

  useEffect(() => {
    async function getCurrentLocation() {
      const { status } = await requestForegroundPermissionsAsync();
      if (status === PermissionStatus.DENIED) {
        return;
      }

      const location = await getCurrentPositionAsync({});
      setInitialLocation(location);
      setZoom(12.6);
      console.log(location);
    }

    getCurrentLocation();
  }, []);

  const handleCameraMove = (event: CameraMoveEvent) => {
    setZoom(event.zoom);
    const lat = event.coordinates.latitude;
    const lng = event.coordinates.longitude;
    if (event.zoom >= 12 && lat && lng) {
      const h3Index = latLngToCell(lat, lng, 7);
      if (centerH3Res7.current !== h3Index) {
        centerH3Res7.current = h3Index;
        const h3Indexes = gridDisk(h3Index, 1);
        setH3Res7(h3Indexes);
      }
      return;
    }
    if (event.zoom >= 10 && lat && lng) {
      const h3Index = latLngToCell(lat, lng, 5);
      if (centerH3Res5.current !== h3Index) {
        centerH3Res5.current = h3Index;
        const h3Indexes = gridDisk(h3Index, 1);
        setH3Res5(h3Indexes);
      }
    }
  };

  return (
    <View style={{ flex: 1 }}>
      {Platform.OS === "ios" ? (
        <AppleMaps.View
          style={{ flex: 1 }}
          cameraPosition={{
            coordinates: {
              latitude: initialLocation?.coords.latitude,
              longitude: initialLocation?.coords.longitude,
            },
            zoom: zoom,
          }}
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
      <Text>current zoom: {zoom.toFixed(2)}</Text>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {},
});
