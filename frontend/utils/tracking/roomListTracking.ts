/**
 * Room List Tracking Module
 *
 * Centralized tracking for room selection and booking funnel
 */

import * as Sentry from '@sentry/react-native';

// ============================================================================
// TYPE DEFINITIONS
// ============================================================================

interface RoomListLoadedData {
  roomCount: number;
  page: number;
  loadTime: number;
  hasMore: boolean;
  totalRoomsLoaded: number;
}

interface RoomListErrorData {
  type: 'network' | 'api' | 'parse' | 'state';
  message: string;
  page: number;
  statusCode?: number;
  duration?: number;
  roomsLoadedSoFar?: number;
}

interface RoomCardTappedData {
  roomId: number;
  roomName: string;
  roomStatus: string;
  position: number;
  pageNumber: number;
  timeOnScreen: number;
  roomsViewedBefore: number;
  isAvailable: boolean;
}

interface PaginationTriggeredData {
  currentPage: number;
  nextPage: number;
  totalRoomsLoaded: number;
  triggerMethod: 'scroll' | 'manual';
}

// ============================================================================
// SCREEN TRACKING
// ============================================================================

/**
 * Track room list screen view
 */
export function trackScreenView(isAuthenticated: boolean) {
  // Add breadcrumb
  Sentry.addBreadcrumb({
    category: 'navigation',
    message: 'Navigated to Room List screen',
    level: 'info',
    data: {
      screen: 'room-list',
      authenticated: isAuthenticated,
      timestamp: Date.now()
    }
  });

  // Set context
  Sentry.setContext('screen', {
    name: 'Room List',
    feature: 'booking',
    authenticated: isAuthenticated
  });

  // Set tag for filtering
  Sentry.setTag('current_screen', 'room-list');
}

/**
 * Track successful room list load
 */
export function trackRoomListLoaded(data: RoomListLoadedData) {
  Sentry.addBreadcrumb({
    category: 'data',
    message: `Loaded ${data.roomCount} rooms (page ${data.page})`,
    level: 'info',
    data: {
      room_count: data.roomCount,
      page_number: data.page,
      load_time_ms: data.loadTime,
      has_more: data.hasMore,
      total_rooms: data.totalRoomsLoaded,
      timestamp: Date.now()
    }
  });

  // Track performance metric
  Sentry.setMeasurement('room_list.load_time', data.loadTime, 'millisecond');

  // Track business metric
  Sentry.setMeasurement('room_list.rooms_loaded', data.roomCount, 'none');
}

/**
 * Track room list error
 */
export function trackRoomListError(data: RoomListErrorData) {
  Sentry.addBreadcrumb({
    category: 'error',
    message: `Room list error: ${data.message}`,
    level: 'error',
    data: {
      error_type: data.type,
      page_number: data.page,
      status_code: data.statusCode,
      duration_ms: data.duration,
      rooms_loaded: data.roomsLoadedSoFar,
      timestamp: Date.now()
    }
  });

  // Capture as warning if recoverable
  if (data.type === 'network' && (data.roomsLoadedSoFar || 0) > 0) {
    Sentry.captureMessage(`Room list ${data.type} error on page ${data.page}`, {
      level: 'warning',
      tags: {
        error_category: data.type,
        page_number: data.page.toString()
      }
    });
  }
}

// ============================================================================
// INTERACTION TRACKING
// ============================================================================

/**
 * Track room card tap
 */
export function trackRoomCardTapped(data: RoomCardTappedData) {
  // Add detailed breadcrumb
  Sentry.addBreadcrumb({
    category: 'ui.interaction',
    message: `User tapped room: ${data.roomName}`,
    level: 'info',
    data: {
      action: 'tap',
      target: 'room-card',
      room_id: data.roomId,
      room_name: data.roomName,
      room_status: data.roomStatus,
      position: data.position,
      page_number: data.pageNumber,
      time_on_screen_ms: data.timeOnScreen,
      rooms_viewed_before: data.roomsViewedBefore,
      is_available: data.isAvailable,
      timestamp: Date.now()
    }
  });

  // Track conversion funnel progression
  Sentry.addBreadcrumb({
    category: 'conversion',
    message: 'Room selection - booking funnel step 1',
    level: 'info',
    data: {
      funnel_step: 'room_selected',
      room_id: data.roomId,
      time_to_selection_ms: data.timeOnScreen,
      exploration_count: data.roomsViewedBefore + 1
    }
  });

  // Track as custom event for analytics
  Sentry.setContext('room_selection', {
    room_id: data.roomId,
    room_name: data.roomName,
    position_in_list: data.position,
    is_available: data.isAvailable,
    decision_time_ms: data.timeOnScreen
  });

  // Performance metric
  Sentry.setMeasurement('room_list.time_to_selection', data.timeOnScreen, 'millisecond');
}

/**
 * Track pagination trigger
 */
export function trackPaginationTriggered(data: PaginationTriggeredData) {
  Sentry.addBreadcrumb({
    category: 'ui.interaction',
    message: `Pagination triggered: page ${data.nextPage}`,
    level: 'info',
    data: {
      action: 'pagination',
      current_page: data.currentPage,
      next_page: data.nextPage,
      total_rooms_loaded: data.totalRoomsLoaded,
      trigger_method: data.triggerMethod,
      timestamp: Date.now()
    }
  });

  // Track engagement metric
  Sentry.setMeasurement('room_list.pages_viewed', data.nextPage, 'none');
}

// ============================================================================
// API TRACKING
// ============================================================================

/**
 * Track API request with enhanced context
 */
export function trackApiRequest(
  method: string,
  url: string,
  page: number,
  status: 'initiated' | 'success' | 'error',
  statusCode?: number,
  duration?: number
) {
  const message = status === 'initiated'
    ? `API request: ${method} ${url}`
    : `API ${status}: ${method} ${url}`;

  Sentry.addBreadcrumb({
    category: 'http',
    type: 'http',
    message,
    level: status === 'error' ? 'error' : 'info',
    data: {
      method,
      url,
      page_number: page,
      status,
      status_code: statusCode,
      duration_ms: duration,
      timestamp: Date.now()
    }
  });

  // Track performance if completed
  if (duration && status === 'success') {
    Sentry.setMeasurement(`room_list.api_response_time_p${page}`, duration, 'millisecond');
  }
}

// ============================================================================
// PERFORMANCE TRANSACTIONS
// ============================================================================

/**
 * Create performance transaction for room list operations
 * Note: Using startSpan instead of deprecated startTransaction
 */
export function createRoomListTransaction(operation: 'initial_load' | 'pagination' | 'refresh') {
  // Use startSpan for newer Sentry SDK versions
  return Sentry.startSpan(
    {
      name: `RoomList.${operation}`,
      op: operation === 'initial_load' ? 'screen.load' : 'ui.action',
      attributes: {
        screen: 'room-list',
        feature: 'booking',
        operation
      }
    },
    () => {
      // Return a span-like object for compatibility
      return {
        startChild: () => ({ finish: () => {} }),
        finish: () => {}
      };
    }
  );
}