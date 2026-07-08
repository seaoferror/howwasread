import React from 'react';
import {
  View,
  Text,
  ScrollView,
  StyleSheet,
} from 'react-native';
import { SafeAreaView } from "react-native-safe-area-context";
import CustomButton from "@/components/CustomButton";
import { router } from "expo-router";
import { setKVStore } from "@/db/storage";

interface TermsSection {
  id: string;
  body: string;
}

const SECTIONS: TermsSection[] = [
  {
    id: '1',
    body:
      'Defamatory, discriminatory, or mean-spirited content, including references or ' +
      'commentary about religion, race, sexual orientation, gender, national/ethnic ' +
      'origin, or other targeted groups, particularly if the app is likely to humiliate, ' +
      'intimidate, or harm a targeted individual or group. Professional political ' +
      'satirists and humorists are generally exempt from this requirement.',
  },
  {
    id: '2',
    body:
      'Realistic portrayals of people or animals being killed, maimed, tortured, or ' +
      'abused, or content that encourages violence. "Enemies" within the context of a ' +
      'game cannot solely target a specific race, culture, real government, corporation, ' +
      'or any other real entity.',
  },
  {
    id: '3',
    body:
      'Depictions that encourage illegal or reckless use of weapons and dangerous ' +
      'objects, or facilitate the purchase of firearms or ammunition.',
  },
  {
    id: '4',
    body:
      'Overtly sexual or pornographic material, defined as "explicit descriptions or ' +
      'displays of sexual organs or activities intended to stimulate erotic rather than ' +
      'aesthetic or emotional feelings." This includes "hookup" apps and other apps that ' +
      'may include pornography or be used to facilitate prostitution, or human ' +
      'trafficking and exploitation.',
  },
  {
    id: '5',
    body:
      'Inflammatory religious commentary or inaccurate or misleading quotations of ' +
      'religious texts.',
  },
  {
    id: '6',
    body:
      'False information and features, including inaccurate device data or trick/joke ' +
      'functionality, such as fake location trackers. Stating that the app is "for ' +
      'entertainment purposes" won\'t overcome this guideline. Apps that enable ' +
      'anonymous or prank phone calls or SMS/MMS messaging will be rejected.',
  },
  {
    id: '7',
    body:
      'Harmful concepts which capitalize or seek to profit on recent or current events, ' +
      'such as violent conflicts, terrorist attacks, and epidemics.',
  },
];

export default function TermsOfUseScreen() {
  return (
    <SafeAreaView style={styles.safeArea}>
      <ScrollView
        contentContainerStyle={styles.scrollContent}
        showsVerticalScrollIndicator={true}
      >
        <Text style={styles.title}>Terms of Use</Text>
        <Text style={styles.lastUpdated}>Last updated: July 8, 2026</Text>

        <View style={styles.section}>
          <Text style={styles.sectionHeading}>Objectionable Content</Text>
          <Text style={styles.paragraph}>
            By using this app, you agree not to submit, upload, or share content that
            is offensive, insensitive, upsetting, intended to disgust, in exceptionally
            poor taste, or otherwise objectionable. Examples of prohibited content
            include, but are not limited to, the following:
          </Text>

          {SECTIONS.map((section) => (
            <View key={section.id} style={styles.rule}>
              <Text style={styles.ruleId}>{section.id}</Text>
              <Text style={styles.ruleBody}>{section.body}</Text>
            </View>
          ))}

          <Text style={styles.paragraph}>
            Content that violates any of the above may be removed without notice, and
            accounts responsible for repeated or severe violations may be suspended or
            terminated at our discretion.
          </Text>
          <CustomButton label="Agree" onPress={() => {
            setKVStore("didAgree", "1");
            router.replace("/auth")
          }}/>
        </View>
      </ScrollView>
    </SafeAreaView>
  );
};

const styles = StyleSheet.create({
  safeArea: {
    flex: 1,
    backgroundColor: '#FFFFFF',
  },
  scrollContent: {
    paddingHorizontal: 20,
    paddingVertical: 24,
  },
  title: {
    fontSize: 26,
    fontWeight: '700',
    color: '#111111',
    marginBottom: 4,
  },
  lastUpdated: {
    fontSize: 13,
    color: '#767676',
    marginBottom: 20,
  },
  section: {
    marginBottom: 24,
  },
  sectionHeading: {
    fontSize: 18,
    fontWeight: '600',
    color: '#111111',
    marginBottom: 10,
  },
  paragraph: {
    fontSize: 14,
    lineHeight: 21,
    color: '#333333',
    marginBottom: 14,
  },
  rule: {
    flexDirection: 'row',
    marginBottom: 12,
    paddingLeft: 4,
  },
  ruleId: {
    fontSize: 14,
    fontWeight: '700',
    color: '#111111',
    width: 52,
  },
  ruleBody: {
    flex: 1,
    fontSize: 14,
    lineHeight: 21,
    color: '#333333',
  },
});

