import { View, Text, TouchableOpacity } from "react-native";
import { router } from "expo-router";
import RemoteWorker from '../../image/remote-worker.svg';


export default function WelcomeScreen() {
  return (
    <View className="flex-1 bg-white items-center px-6 relative">

      <RemoteWorker width={308} height={248} style={{ marginTop: 100 }} />

      <Text className="text-[20px] font-medium text-black mt-6">
        Welcome
      </Text>

      <Text className="text-[16px] text-[#A0A0A0] mb-10 mt-[20px]">
        Have a better sharing experience
      </Text>

      <View className="absolute bottom-28 w-full items-center">

        <TouchableOpacity
          className="w-[90%] py-4 rounded-xl mb-8 bg-[#FDBA29] items-center"
          onPress={() => router.push("/(auth)/signup")}
        >
          <Text className="text-white text-[18px] font-semibold">
            Create an account
          </Text>
        </TouchableOpacity>

        <TouchableOpacity
          className="w-[90%] py-4 rounded-xl bg-white border items-center"
          style={{ borderColor: "#FDBA29", borderWidth: 2 }}
          onPress={() => router.push("/(auth)/login")}
        >
          <Text className="text-[18px] font-semibold" style={{ color: "#FDBA29" }}>
            Log in
          </Text>
        </TouchableOpacity>

        <TouchableOpacity
          className="w-full items-center mt-6"
          onPress={() => router.push("/(tabs)")}
        >
          <Text className="text-[16px] text-[#A0A0A0] underline">
            Skip
          </Text>
        </TouchableOpacity>

      </View>

    </View>
  );
}
