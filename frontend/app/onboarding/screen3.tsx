import { View, Text, TouchableOpacity } from "react-native";
import { router } from "expo-router";
import AsyncStorage from "@react-native-async-storage/async-storage";
import Confirmed from '../../image/confirmed.svg';

export default function Screen3() {
  const finishOnboarding = async () => {
    await AsyncStorage.setItem("hasSeenOnboarding", "true");
    router.replace("/"); 
  };

  return (
    <View className="flex-1 justify-center items-center bg-white">
      <Confirmed width={280} height={230} style={{ marginTop: 92 }} />

      <Text className="text-2xl font-bold mb-10 text-center">
        Book your space{"\n"}
        Own your focus
      </Text>

      <Text className="text-gray-500 text-lg font-semibold text-center mb-10">
        Choose the perfect room for quiet {"\n"}
        study or teamwork — in just a few {"\n"}
        taps.
      </Text>

      <TouchableOpacity activeOpacity={0.8} onPress={finishOnboarding}>
        <View className="w-[100px] h-[100px] justify-center items-center relative">

          <View
            style={{
              position: "absolute",
              width: 90,
              height: 90,
              borderRadius: 50,
              borderWidth: 4,
              borderColor: "#FDBA29",
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
