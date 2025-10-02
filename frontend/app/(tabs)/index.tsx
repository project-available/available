import { View, Text, TouchableOpacity } from "react-native";
import { useRouter } from "expo-router";

export default function HomeTab() {
  const router = useRouter();

  return (
    <View className="flex-1 bg-[#f8f5f2] p-4">
      {/* Header */}
      <View className="flex-row justify-between items-center mb-6">
        <Text className="text-lg font-bold">Hello!!! Nhi</Text>
        <TouchableOpacity>
          <Text>🔔</Text>
        </TouchableOpacity>
      </View>

      {/* Title */}
      <Text className="text-xl font-bold mb-4">Services</Text>

      {/* Services Grid */}
      <View className="flex-row flex-wrap justify-between gap-4">
        <TouchableOpacity
          className="bg-blue-500 rounded-xl w-[48%] h-32 justify-center items-center"
          onPress={() => router.push("/room")}
        >
          <Text className="text-white font-bold">Book</Text>
          <Text className="text-white">Making your booking</Text>
        </TouchableOpacity>

        <TouchableOpacity
          className="bg-white rounded-xl w-[48%] h-32 justify-center items-center"
          onPress={() => router.push("/room")}
        >
          <Text className="font-bold">Cancel</Text>
          <Text>Cancel your booking</Text>
        </TouchableOpacity>

        <TouchableOpacity
          className="bg-white rounded-xl w-[48%] h-32 justify-center items-center"
          onPress={() => router.push("/notification")}
        >
          <Text className="font-bold">View History</Text>
          <Text>View booking history</Text>
        </TouchableOpacity>

        <TouchableOpacity
          className="bg-white rounded-xl w-[48%] h-32 justify-center items-center"
          onPress={() => router.push("/profile")}
        >
          <Text className="font-bold">Available</Text>
          <Text>View available rooms</Text>
        </TouchableOpacity>
      </View>
    </View>
  );
}
