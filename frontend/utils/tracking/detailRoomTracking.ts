/**
 * Detail Room Tracking Module
 *
 * Tracking for room detail view and booking calendar
 */

import * as Sentry from '@sentry/react-native';

// ============================================================================
// TYPE DEFINITIONS
// ============================================================================

interface RoomDetailViewedData {
  roomId: string;
  roomName: string;
  fromScreen: string;
}

interface RoomDetailLoadedData {
  roomId: string;
  roomName: string;
  capacity: string;
  location: string;
  loadTime: number;
}

interface BookingsLoadedData {
  roomId: string;
  date: string;
  bookingCount: number;
  loadTime: number;
  hasRetry: boolean;
  retryCount?: number;
}

interface DateChangedData {
  roomId: string;
  fromDate: string;
  toDate: string;
  daysFromToday: number;
}

interface BookNowTappedData {
  roomId: string;
  roomName: string;
  selectedDate: string;
  timeOnScreen: number;
  datesViewedCount: number;
}

interface CalendarErrorData {
  roomId: string;
  date: string;
  errorType: 'network' | 'api' | 'retry_failed';
  statusCode?: number;
  retryCount?: number;
  duration?: number;
}

// ============================================================================
// SCREEN TRACKING
// ============================================================================

/**
 * Track room detail screen view
 */
export function trackRoomDetailViewed(data: RoomDetailViewedData) {
  Sentry.addBreadcrumb({
    category: 'navigation',
    message: `Viewed room detail: ${data.roomName}`,
    level: 'info',
    data: {
      screen: 'room-detail',
      room_id: data.roomId,
      room_name: data.roomName,
      from_screen: data.fromScreen,
      timestamp: Date.now()
    }
  });

  // Set context
  Sentry.setContext('room_detail', {
    room_id: data.roomId,
    room_name: data.roomName,
    from_screen: data.fromScreen
  });

  // Set tag for filtering
  Sentry.setTag('current_screen', 'room-detail');
  Sentry.setTag('room_id', data.roomId);
}

/**
 * Track successful room detail load
 */
export function trackRoomDetailLoaded(data: RoomDetailLoadedData) {
  Sentry.addBreadcrumb({
    category: 'data',
    message: `Loaded room details: ${data.roomName}`,
    level: 'info',
    data: {
      room_id: data.roomId,
      room_name: data.roomName,
      capacity: data.capacity,
      location: data.location,
      load_time_ms: data.loadTime,
      timestamp: Date.now()
    }
  });

  // Track performance metric
  Sentry.setMeasurement('room_detail.load_time', data.loadTime, 'millisecond');
}

/**
 * Track bookings loaded for calendar
 */
export function trackBookingsLoaded(data: BookingsLoadedData) {
  Sentry.addBreadcrumb({
    category: 'data',
    message: `Loaded ${data.bookingCount} bookings for ${data.date}`,
    level: 'info',
    data: {
      room_id: data.roomId,
      date: data.date,
      booking_count: data.bookingCount,
      load_time_ms: data.loadTime,
      has_retry: data.hasRetry,
      retry_count: data.retryCount,
      timestamp: Date.now()
    }
  });

  // Track performance
  Sentry.setMeasurement('calendar.booking_load_time', data.loadTime, 'millisecond');

  // Track if retry was needed
  if (data.hasRetry) {
    Sentry.addBreadcrumb({
      category: 'retry',
      message: `Booking fetch succeeded after ${data.retryCount} retries`,
      level: 'warning',
      data: {
        room_id: data.roomId,
        retry_count: data.retryCount
      }
    });
  }
}

// ============================================================================
// INTERACTION TRACKING
// ============================================================================

/**
 * Track date selection in calendar
 */
export function trackDateChanged(data: DateChangedData) {
  Sentry.addBreadcrumb({
    category: 'ui.interaction',
    message: `Date changed: ${data.fromDate} → ${data.toDate}`,
    level: 'info',
    data: {
      action: 'date_change',
      room_id: data.roomId,
      from_date: data.fromDate,
      to_date: data.toDate,
      days_from_today: data.daysFromToday,
      timestamp: Date.now()
    }
  });
}

/**
 * Track "Book Now" button tap
 */
export function trackBookNowTapped(data: BookNowTappedData) {
  Sentry.addBreadcrumb({
    category: 'ui.interaction',
    message: `Book Now tapped for ${data.roomName}`,
    level: 'info',
    data: {
      action: 'tap',
      target: 'book-now-button',
      room_id: data.roomId,
      room_name: data.roomName,
      selected_date: data.selectedDate,
      time_on_screen_ms: data.timeOnScreen,
      dates_viewed: data.datesViewedCount,
      timestamp: Date.now()
    }
  });

  // Track conversion funnel progression
  Sentry.addBreadcrumb({
    category: 'conversion',
    message: 'Book Now initiated - booking funnel step 2',
    level: 'info',
    data: {
      funnel_step: 'book_now_tapped',
      room_id: data.roomId,
      selected_date: data.selectedDate,
      time_to_decision_ms: data.timeOnScreen,
      exploration_count: data.datesViewedCount
    }
  });

  // Performance metric
  Sentry.setMeasurement('room_detail.time_to_book', data.timeOnScreen, 'millisecond');
}

// ============================================================================
// ERROR TRACKING
// ============================================================================

/**
 * Track calendar/booking fetch errors
 */
export function trackCalendarError(data: CalendarErrorData) {
  Sentry.addBreadcrumb({
    category: 'error',
    message: `Calendar error: ${data.errorType}`,
    level: 'error',
    data: {
      error_type: data.errorType,
      room_id: data.roomId,
      date: data.date,
      status_code: data.statusCode,
      retry_count: data.retryCount,
      duration_ms: data.duration,
      timestamp: Date.now()
    }
  });

  // Capture as warning if it's a retry scenario
  if (data.errorType === 'retry_failed' || (data.retryCount && data.retryCount > 0)) {
    Sentry.captureMessage(`Calendar fetch failed after retries for room ${data.roomId}`, {
      level: 'warning',
      tags: {
        error_category: data.errorType,
        room_id: data.roomId,
        retry_count: data.retryCount?.toString()
      },
      contexts: {
        calendar_fetch: {
          room_id: data.roomId,
          date: data.date,
          retry_count: data.retryCount,
          duration_ms: data.duration
        }
      }
    });
  }
}

/**
 * Track room detail fetch error
 */
export function trackRoomDetailError(roomId: string, errorMessage: string, statusCode?: number) {
  Sentry.addBreadcrumb({
    category: 'error',
    message: `Room detail fetch error: ${errorMessage}`,
    level: 'error',
    data: {
      room_id: roomId,
      error_message: errorMessage,
      status_code: statusCode,
      timestamp: Date.now()
    }
  });
}

// ============================================================================
// API TRACKING
// ============================================================================

/**
 * Track API request for room details
 */
export function trackRoomDetailApiRequest(
  roomId: string,
  status: 'initiated' | 'success' | 'error',
  statusCode?: number,
  duration?: number
) {
  const message = status === 'initiated'
    ? `API request: GET /rooms/${roomId}`
    : `API ${status}: GET /rooms/${roomId}`;

  Sentry.addBreadcrumb({
    category: 'http',
    type: 'http',
    message,
    level: status === 'error' ? 'error' : 'info',
    data: {
      method: 'GET',
      url: `/rooms/${roomId}`,
      room_id: roomId,
      status,
      status_code: statusCode,
      duration_ms: duration,
      timestamp: Date.now()
    }
  });
}

/**
 * Track API request for bookings
 */
export function trackBookingsApiRequest(
  roomId: string,
  date: string,
  status: 'initiated' | 'success' | 'error' | 'retry',
  statusCode?: number,
  duration?: number,
  retryCount?: number
) {
  const message = status === 'initiated'
    ? `API request: GET /rooms/${roomId}/bookings`
    : `API ${status}: GET /rooms/${roomId}/bookings`;

  Sentry.addBreadcrumb({
    category: 'http',
    type: 'http',
    message,
    level: status === 'error' ? 'error' : status === 'retry' ? 'warning' : 'info',
    data: {
      method: 'GET',
      url: `/rooms/${roomId}/bookings`,
      room_id: roomId,
      date,
      status,
      status_code: statusCode,
      duration_ms: duration,
      retry_count: retryCount,
      timestamp: Date.now()
    }
  });
}