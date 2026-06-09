import { useState } from "react";
import { StyleSheet, View } from "react-native";
import { WebView, type WebViewNavigation } from "react-native-webview";

interface GoogleMapsResolverProps {
  shortUrl: string;
  onCoordinatesResolved: (coords: { lat: number; lng: number }) => void;
}

export default function GoogleMapsResolver({
  shortUrl,
  onCoordinatesResolved,
}: GoogleMapsResolverProps) {
  const [isResolving, setIsResolving] = useState(true);

  const handleNavigationStateChange = (navState: WebViewNavigation) => {
    const currentUrl = navState.url;
    const pinPattern = /!3d(-?\d+\.\d+)!4d(-?\d+\.\d+)/;
    const match = currentUrl.match(pinPattern);

    if (match) {
      setIsResolving(false); // Stop loading the WebView
      const lat = parseFloat(match[1]);
      const lng = parseFloat(match[2]);
      onCoordinatesResolved({ lat, lng });
    }
  };

  if (!isResolving) return <></>;

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
