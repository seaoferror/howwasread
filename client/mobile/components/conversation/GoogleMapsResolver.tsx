import { useState } from "react";
import { StyleSheet, View } from "react-native";
import { WebView, type WebViewNavigation } from "react-native-webview";

interface GoogleMapsResolverProps {
  shortUrl: string;
  onGeoInfoResolved: (coords: { lat: number; lng: number; placeName: string }) => void;
}

export default function GoogleMapsResolver({
  shortUrl,
  onGeoInfoResolved,
}: GoogleMapsResolverProps) {
  const [isResolving, setIsResolving] = useState(true);

  const handleNavigationStateChange = (navState: WebViewNavigation) => {
    const currentUrl = navState.url;
    const coordsMatch = currentUrl.match(/!3d(-?\d+\.\d+)!4d(-?\d+\.\d+)/);
    const placeNameMatch = currentUrl.match(/\/maps\/place\/([^/@]+)/);

    if (coordsMatch && placeNameMatch) {
      setIsResolving(false);
      const lat = parseFloat(coordsMatch[1]);
      const lng = parseFloat(coordsMatch[2]);

      const spaceCleaned = placeNameMatch[1].replace(/\+/g, " ");
      const decodedPlaceName = decodeURIComponent(spaceCleaned);
      const isCoords = /^[-]?\d+\.\d+,\s*[-]?\d+\.\d+$/.test(
        decodedPlaceName,
      );
      const placeName = isCoords ? "Dropped Pin" : decodedPlaceName;
      onGeoInfoResolved({ lat, lng, placeName });
    }
  };

  if (!isResolving) return null;

  return (
    <View style={styles.hiddenContainer}>
      <WebView
        source={{ uri: shortUrl }}
        onNavigationStateChange={handleNavigationStateChange}
        javaScriptEnabled={true}
        domStorageEnabled={true}
        incognito={true}
        startInLoadingState={false}
        scalesPageToFit={false}
        pointerEvents="none"
      />
    </View>
  );
}

const styles = StyleSheet.create({
  hiddenContainer: {
    width: 0,
    height: 0,
    opacity: 0,
    position: "absolute",
  },
});
