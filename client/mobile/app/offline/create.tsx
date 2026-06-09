import {
  KeyboardAvoidingView,
  Platform,
  ScrollView,
  StyleSheet,
  Text,
  View,
} from "react-native";
import { FormProvider, useForm } from "react-hook-form";
import { colors } from "@/constants";
import FixedBottomCTA from "@/components/FixedBottomCTA";
import NovelInput from "@/components/conversation/NovelInput";
import ShortStoryInput from "@/components/conversation/ShortStoryInput";
import PoemInput from "@/components/conversation/PoemInput";
import PlayInput from "@/components/conversation/PlayInput";
import FilmInput from "@/components/conversation/FilmInput";
import ByInput from "@/components/conversation/ByInput";
import RuleInput from "@/components/conversation/RuleInput";
import YearInput from "@/components/conversation/YearInput";
import MonthDayInput from "@/components/conversation/MonthDayInput";
import HourInput from "@/components/conversation/HourInput";
import MinuteInput from "@/components/conversation/MinuteInput";
import LengthInput from "@/components/conversation/LengthInput";
import { useSafeAreaInsets } from "react-native-safe-area-context";
import useKeyboard from "@/hooks/useKeyboard";
import { openBrowserAsync } from "expo-web-browser";
import { useCreateOfflineConversation } from "@/hooks/useConversation";
import CustomButton from "@/components/CustomButton";
import MapsLinkInput from "@/components/conversation/MapsLinkInput";
import LocationInput from "@/components/conversation/LocationInput";
import { reverseGeocodeAsync } from "expo-location";
import { makeTime } from "@/util/time";
import Toast from "react-native-toast-message";
import { router } from "expo-router";
import { latLngToCell } from "h3-js";
import { extractAppleMapInfo } from "@/util/url";
import GoogleMapsResolver from "@/components/conversation/GoogleMapsResolver";
import { useEffect, useState } from "react";

interface FormValue {
  novel?: string;
  shortStory?: string;
  poem?: string;
  play?: string;
  film?: string;
  by: string;
  rule?: string;
  year: string;
  monthDay: string;
  hour: string;
  minute: string;
  length: string;
  mapsLink: string;
  location?: string;
}

export default function OfflineConversationScreen() {
  const createOfflineConversationMutation = useCreateOfflineConversation();
  const { isKeyboardVisible } = useKeyboard();
  const insets = useSafeAreaInsets();
  const now = new Date();
  const [resolvedCoords, setResolvedCoords] = useState<{
    lat: number;
    lng: number;
  } | null>(null);
  const offlineConversationForm = useForm<FormValue>({
    defaultValues: {
      novel: "",
      shortStory: "",
      poem: "",
      play: "",
      film: "",
      by: "",
      rule: "",
      year: String(now.getFullYear()),
      monthDay: `${now.getMonth() + 1}.${now.getDate()}`,
      hour: String(now.getHours()),
      minute: String(now.getMinutes()),
      length: "100",
      mapsLink: "",
      location: "",
    },
  });
  const mapsLinkValue = offlineConversationForm.watch("mapsLink");
  useEffect(() => {
    setResolvedCoords(null);
  }, [mapsLinkValue]);

  const onSubmit = async (formValues: FormValue) => {
    let {
      novel,
      shortStory,
      poem,
      play,
      film,
      by,
      rule,
      year,
      monthDay,
      hour,
      minute,
      length,
      mapsLink,
      location,
    } = formValues;
    const time = makeTime(now, monthDay, year, hour, minute);
    let lat = 0;
    let lng = 0;
    let placeName = "";
    if (Platform.OS === "ios") {
      const data = extractAppleMapInfo(mapsLink);
      lat = data.lat;
      lng = data.lng;
      placeName = data.placeName;
    }
    if (Platform.OS === "android" && resolvedCoords !== null) {
      lat = resolvedCoords.lat
      lng = resolvedCoords.lng
    }
    if (!location) {
      location = placeName;
    }
    const geoInfo = await reverseGeocodeAsync({
      latitude: lat,
      longitude: lng,
    });
    console.log(geoInfo);
    console.log(lat);
    console.log(lng);
    console.log(location);

    createOfflineConversationMutation.mutate(
      {
        novel: novel,
        shortStory: shortStory,
        poem: poem,
        play: play,
        film: film,
        writtenBy: by,
        rule: rule,
        time: time,
        length: Number(length),
        mapsLink: mapsLink,
        location: location,
        city: geoInfo[0].city ?? "",
        lat: lat,
        lng: lng,
        h3Res5: latLngToCell(lat, lng, 5),
        h3Res7: latLngToCell(lat, lng, 7),
      },
      {
        onSuccess: () => {
          router.replace("/conversations");
        },
        onError: (error) => {
          console.log(error);
          Toast.show({
            type: "error",
            text1: "Invalid Google Maps URL",
          });
        },
      },
    );
  };

  return (
    <FormProvider {...offlineConversationForm}>
      <View style={styles.container}>
        <KeyboardAvoidingView
          contentContainerStyle={styles.awareScrollViewContainer}
          behavior="height"
          keyboardVerticalOffset={
            Platform.OS === "ios" || isKeyboardVisible ? 100 : insets.bottom
          }
        >
          <ScrollView style={{ marginBottom: 100 }}>
            <View style={styles.content}>
              <NovelInput />
              <ShortStoryInput />
              <PoemInput />
              <PlayInput />
              <FilmInput />
              <ByInput />
              <RuleInput />
              <Text style={styles.whenLabel}>When</Text>
              <YearInput />
              <MonthDayInput />
              <View style={styles.timeRow}>
                <HourInput />
                <Text>:</Text>
                <MinuteInput />
              </View>
              <LengthInput />
              {Platform.OS === "ios" ? (
                <CustomButton
                  label="Go and get Apple Maps share link"
                  onPress={() => openBrowserAsync("https://maps.apple.com")}
                />
              ) : (
                <CustomButton
                  label="Go and get Google Maps share link"
                  onPress={() =>
                    openBrowserAsync("https://www.google.com/maps")
                  }
                />
              )}
              <MapsLinkInput />
              {mapsLinkValue &&
                !resolvedCoords &&
                Platform.OS === "android" && (
                  <GoogleMapsResolver
                    shortUrl={mapsLinkValue}
                    onCoordinatesResolved={(coords) => {
                      setResolvedCoords(coords);
                      console.log("Headless Resolution Success:", coords);
                    }}
                  />
                )}
              <LocationInput />
            </View>
          </ScrollView>
          <FixedBottomCTA
            label="Create"
            onPress={offlineConversationForm.handleSubmit(onSubmit)}
            disabled={Platform.OS === "android" && resolvedCoords === null}
          />
        </KeyboardAvoidingView>
      </View>
    </FormProvider>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: colors.SAND_110,
    borderTopWidth: StyleSheet.hairlineWidth,
    borderColor: colors.GRAY_700,
  },
  content: {
    flex: 1,
    margin: 16,
    gap: 16,
    backgroundColor: colors.SAND_110,
  },
  timeRow: {
    flex: 1,
    flexDirection: "row",
    alignItems: "center",
    gap: 10,
  },
  awareScrollViewContainer: {
    flex: 1,
  },
  whenLabel: {
    fontSize: 12,
    color: colors.GRAY_700,
    marginBottom: -10,
  },
});
