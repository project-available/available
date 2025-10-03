import { View, Text, TouchableOpacity } from 'react-native';
import Ionicons from "@expo/vector-icons/Ionicons";

export default function Profile() {
  return (
    <View className="flex-1 bg-[#f8f5f2]">
      {/* Header với Profile Details + Edit icon */}
      <View className="flex-row justify-between items-center mt-[66px] mx-[24px]">
        <Text className="text-[18px] leading-[24px] font-semibold">
          Profile Details
        </Text>

        <TouchableOpacity className="w-10 h-10 bg-white rounded-full items-center justify-center">
          <Ionicons name="pencil-outline" size={20} color="black" />
        </TouchableOpacity>
      </View>
    </View>
  );
}
