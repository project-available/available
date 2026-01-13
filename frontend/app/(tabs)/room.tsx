import React, { useEffect, useState, useCallback, useRef } from "react";
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
import * as Sentry from "@sentry/react-native";
import {
  trackScreenView,
  trackRoomListLoaded,
  trackRoomListError,
  trackRoomCardTapped,
  trackPaginationTriggered,
  trackApiRequest,
  createRoomListTransaction,
} from "../../utils/tracking/roomListTracking";
import { API_BASE_URL } from "../../utils/apiConfig";

export default function Room() {
  const router = useRouter();
  const [rooms, setRooms] = useState<any[]>([]);
  const [page, setPage] = useState(1);
  const [loading, setLoading] = useState(false);
  const [hasMore, setHasMore] = useState(true);
  const [error, setError] = useState("");

  // Tracking refs
  const screenLoadTime = useRef(Date.now());
  const roomsViewedCount = useRef(0);

  const page_size = 6;

  // Track screen view on mount
  useEffect(() => {
    const checkAuth = async () => {
      const token = await AsyncStorage.getItem("access_token");
      return !!token;
    };

    checkAuth().then((isAuthenticated) => {
      trackScreenView(isAuthenticated);
    });
  }, []);

  const fetchRooms = useCallback(async () => {
    if (!hasMore || loading) return;

    const fetchStartTime = Date.now();
    const url = `${API_BASE_URL}/rooms?page_id=${page}&page_size=${page_size}`;

    try {
      setLoading(true);
      setError("");

      // Track API request initiation
      trackApiRequest("GET", url, page, "initiated");

      const res = await fetch(url, {
        headers: {
          "Content-Type": "application/json",
        },
      });

      const fetchDuration = Date.now() - fetchStartTime;

      if (!res.ok) {
        const errText = await res.text();

        // Track API error
        trackRoomListError({
          type: "api",
          message: `Fetch failed: ${res.status}`,
          page,
          statusCode: res.status,
          duration: fetchDuration,
        });

        // Capture in Sentry with context
        Sentry.captureException(
          new Error(`API fetch failed: ${res.status} - ${errText}`),
          {
            tags: {
              feature: "room_list",
              error_type: "api_fetch",
              page_number: page.toString(),
            },
            contexts: {
              fetch: {
                url,
                page,
                page_size,
                duration_ms: fetchDuration,
                rooms_loaded: rooms.length,
              },
            },
          }
        );

        throw new Error(`❌ Fetch failed: ${res.status} - ${errText}`);
      }

      const data = await res.json();
      const roomCount = data.length;

      if (data.length < page_size) setHasMore(false);

      const mapped = data.map((room: any) => ({
        id: room.roomId,
        name: room.name,
        image: room.image,
        location: room.location,
        capacity: room.customFields?.[0] || "",
        isAvailable: room.status === "available",
        status: room.status,
        availableAt: room.availableAt,
      }));

      setRooms((prev) => [...prev, ...mapped]);

      // Track successful load
      trackRoomListLoaded({
        roomCount,
        page,
        loadTime: fetchDuration,
        hasMore: data.length >= page_size,
        totalRoomsLoaded: rooms.length + roomCount,
      });

      // Track API success
      trackApiRequest("GET", url, page, "success", res.status, fetchDuration);
    } catch (err: any) {
      const fetchDuration = Date.now() - fetchStartTime;
      console.error("Fetch rooms error:", err);

      setError(err.message);

      // Track API failure
      trackApiRequest("GET", url, page, "error", 0, fetchDuration);

      // Track error with full context
      trackRoomListError({
        type: "network",
        message: err.message,
        page,
        duration: fetchDuration,
        roomsLoadedSoFar: rooms.length,
      });
    } finally {
      setLoading(false);
    }
  }, [page, hasMore, loading, rooms.length]);

  useEffect(() => {
    fetchRooms();
  }, [page]);

  const handleLoadMore = () => {
    if (!loading && hasMore) {
      // Track pagination trigger
      trackPaginationTriggered({
        currentPage: page,
        nextPage: page + 1,
        totalRoomsLoaded: rooms.length,
        triggerMethod: "scroll",
      });

      setPage((prev) => prev + 1);
    }
  };

  const handleRoomTap = (room: any, index: number) => {
    const timeOnScreen = Date.now() - screenLoadTime.current;
    roomsViewedCount.current++;

    // Track room card tap with rich context
    trackRoomCardTapped({
      roomId: room.id,
      roomName: room.name,
      roomStatus: room.status,
      position: index + 1,
      pageNumber: Math.ceil((index + 1) / page_size),
      timeOnScreen,
      roomsViewedBefore: roomsViewedCount.current - 1,
      isAvailable: room.isAvailable,
    });

    // Navigate to detail
    router.push({
      pathname: "/detailroom",
      params: { id: room.id },
    });
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
        renderItem={({ item, index }) => (
          <TouchableOpacity onPress={() => handleRoomTap(item, index)}>
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
                  <View
                    className={`h-[16px] rounded-lg border-[0.5px] border-white mt-[4px] items-center justify-center px-2 ${
                      item.isAvailable ? "bg-[#69B76B]" : "bg-[#E74C3C]"
                    }`}
                  >
                    <Text className="text-[8px] font-thin text-white">
                      {item.isAvailable
                        ? "Available"
                        : `Occupied till ${new Date(item.availableAt)
                            .toLocaleTimeString("en-US", {
                              hour: "numeric",
                              minute: "2-digit",
                              hour12: true,
                            })
                            .toLowerCase()}`}
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
