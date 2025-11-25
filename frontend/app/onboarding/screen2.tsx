import { View, Text, TouchableOpacity } from "react-native";
import { router } from "expo-router";
import { Ionicons } from "@expo/vector-icons";
import Programming from '../../image/programming.svg';

export default function Screen2() {
  return (
    <View className="flex-1 justify-center items-center">
      <Programming width={280} height={230} style={{ marginTop: 48 }} />
      <Text className="text-3xl font-bold mb-10">At anytime</Text>
      <Text className="text-3xl font-bold mb-10">Your Study space awaits.</Text>

      <Text className="text-gray-500 text-lg font-semibold">Reserve a room instantly, whether it’s day or night — your learning never stops.</Text>

      <TouchableOpacity
        activeOpacity={0.8}
        onPress={() => router.push("/onboarding/screen3")}
      >
        <View className="w-[100px] h-[100px] justify-center items-center relative">

          <View
            style={{
              position: "absolute",
              width: 100,
              height: 100,
              borderRadius: 50,
              borderWidth: 4,
              borderColor: "#FBBF24",
              borderBottomColor: "transparent",
            }}
          />

          <View className="w-[75px] h-[75px] rounded-full bg-[#FDBA29] justify-center items-center">
            <Ionicons name="arrow-forward" size={28} color="white" />
          </View>
        </View>
      </TouchableOpacity>
    </View>
  );
}
