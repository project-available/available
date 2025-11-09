import React, { useEffect, useState, useCallback } from "react";
import {
  View,
  Text,
  FlatList,
  Image,
  TouchableOpacity,
  ActivityIndicator,
} from "react-native";
import { useRouter } from "expo-router";
import AsyncStorage from "@react-native-async-storage/async-storage";

export default function Room() {
  const router = useRouter();
  const [rooms, setRooms] = useState<any[]>([]);
  const [page, setPage] = useState(1);
  const [loading, setLoading] = useState(false);
  const [hasMore, setHasMore] = useState(true);
  const [error, setError] = useState("");

  const page_size = 6;

  const fetchRooms = useCallback(async () => {
    if (!hasMore || loading) return;
    try {
      setLoading(true);
      setError("");

      const access_token = await AsyncStorage.getItem("access_token");
      if (!access_token) throw new Error("⚠️ Not logged in");

      const url = `https://querulous-valerie-quanghia-967df8a0.koyeb.app/rooms?page_id=${page}&page_size=${page_size}`;
      console.log("📡 Fetching:", url);

      const res = await fetch(url, {
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${access_token}`,
        },
      });

      if (!res.ok) {
        const errText = await res.text();
        throw new Error(`❌ Fetch failed: ${res.status} - ${errText}`);
      }

      const data = await res.json();

      if (data.length < page_size) setHasMore(false);

      const mapped = data.map((room: any) => ({
        id: room.roomId,
        name: room.name,
        image: room.image,
        location: room.location,
        capacity: room.customFields?.[0] || "",
        isAvailable: true,
      }));

      setRooms((prev) => [...prev, ...mapped]);
    } catch (err: any) {
      console.error("Fetch rooms error:", err);
      setError(err.message);
    } finally {
      setLoading(false);
    }
  }, [page, hasMore, loading]);

  useEffect(() => {
    fetchRooms();
  }, [page]);

  const handleLoadMore = () => {
    if (!loading && hasMore) {
      console.log("📥 Loading next page:", page + 1);
      setPage((prev) => prev + 1);
    }
  };

  if (error && rooms.length === 0) {
    return (
      <View className="flex-1 items-center justify-center">
        <Text className="text-red-500 text-center px-6">
          Error loading rooms: {error}
        </Text>
      </View>
    );
  }

  return (
    <View className="flex-1 bg-[#f8f5f2]">
      <Text className="mt-[66px] ml-[24px] text-[18px] leading-[24px] font-semibold">
        Select Room
      </Text>
      <View className="mt-[40px] h-[1px] bg-gray-300 mx-[24px]" />

      <FlatList
        data={rooms}
        keyExtractor={(item) => item.id.toString()}
        renderItem={({ item }) => (
          <TouchableOpacity
            onPress={() =>
              router.push({
                pathname: "/detailroom",
                params: { id: item.id },
              })
            }
          >
            <View className="h-[115px] bg-[#EEEBe5] mx-[24px] mt-[12px] rounded-3xl flex-row items-center">
              <View className="ml-[18px] items-center justify-center">
                <Image
                  source={{ uri: item.image }}
                  className="w-[126px] h-[72px] rounded-[12px] border-[1px] border-white"
                />
              </View>
              <View className="ml-[8px] flex-1 flex-col justify-between h-[72px]">
                <Text className="text-[16px] font-semibold text-[#000]">
                  {item.name}
                </Text>
                <Text className="text-[14px] font-thin text-[#333]">
                  {item.location}
                </Text>
                <View className="flex-row">
                  <View className="w-[80px] h-[16px] rounded-lg border-[0.5px] border-white mt-[4px] bg-[#69B76B] items-center justify-center">
                    <Text className="text-[8px] font-thin text-white">
                      Available
                    </Text>
                  </View>
                  <View className="w-[56px] h-[16px] bg-[#8C8C8C] rounded-lg border-[0.5px] border-white mt-[4px] ml-1 items-center justify-center">
                    <Text className="text-[8px] font-thin text-[#fff]">
                      {item.capacity}
                    </Text>
                  </View>
                </View>
              </View>
            </View>
          </TouchableOpacity>
        )}
        onEndReached={handleLoadMore}
        onEndReachedThreshold={0.5}
        ListFooterComponent={
          loading ? (
            <View className="py-4">
              <ActivityIndicator size="small" color="#607FBA" />
            </View>
          ) : null
        }
      />
    </View>
  );
}
