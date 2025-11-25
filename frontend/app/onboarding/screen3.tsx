import { View, Text, TouchableOpacity } from "react-native";
import { router } from "expo-router";
import AsyncStorage from "@react-native-async-storage/async-storage";
import Confirmed from '../../image/confirmed.svg';

export default function Screen3() {
  const finishOnboarding = async () => {
    await AsyncStorage.setItem("hasSeenOnboarding", "true");
    router.replace("/"); 
  };

  const skipOnboarding = async () => {
    await AsyncStorage.setItem("hasSeenOnboarding", "true");
    router.replace("/");
  };

  return (
    <View className="flex-1 justify-center items-center">
      <Confirmed width={280} height={230} style={{ marginTop: 48 }} />

      <TouchableOpacity
        onPress={skipOnboarding}
        style={{
          position: "absolute",
          top: 50,
          right: 20,
          padding: 10,
        }}
      >
        <Text className="text-gray-500 text-lg font-semibold">Skip</Text>
      </TouchableOpacity>

      <Text className="text-1xl font-bold mb-10">Book your space</Text>
      <Text className="text-1xl font-bold mb-10">Own your focus.</Text>

      <Text className="text-gray-500 text-lg font-semibold">Reserve a room instantly, whether it’s day or night — your learning never stops.</Text>

      <TouchableOpacity activeOpacity={0.8} onPress={finishOnboarding}>
        <View className="w-[100px] h-[100px] justify-center items-center relative">

          <View
            style={{
              position: "absolute",
              width: 100,
              height: 100,
              borderRadius: 50,
              borderWidth: 4,
              borderColor: "#FBBF24",
            }}
          />

          <View className="w-[75px] h-[75px] rounded-full bg-[#FDBA29] justify-center items-center">
            <Text className="text-black font-bold text-lg">Go</Text>
          </View>
        </View>
      </TouchableOpacity>
    </View>
  );
}
