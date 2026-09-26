import { useRef, useState } from "react";
import { StyleSheet, View } from "react-native";
import { WebView, type WebViewNavigation } from "react-native-webview";
import { parseGoogleMapsUrl } from "@/util/geo";
import { GeoInfo } from "@/types/conversation";

interface GoogleMapsResolverProps {
  shortUrl: string;
  onGeoInfoResolved: (coords: GeoInfo) => void;
}

export default function GoogleMapsResolver({
  shortUrl,
  onGeoInfoResolved,
}: GoogleMapsResolverProps) {
  const [isResolving, setIsResolving] = useState(true);
  const resolvedRef = useRef(false);

  const handleNavigationStateChange = (navState: WebViewNavigation) => {
    if (resolvedRef.current) return;
    const info = parseGoogleMapsUrl(navState.url);
    if (!info) return;

    resolvedRef.current = true;
    setIsResolving(false);
    onGeoInfoResolved(info);
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
