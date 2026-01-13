/**
 * Auth Tracking Module
 *
 * Tracking for authentication related events
 */

import * as Sentry from '@sentry/react-native';

// ============================================================================
// TYPE DEFINITIONS
// ============================================================================

interface LoginAttemptData {
  method: 'email';
}

interface LoginSuccessData {
  method: 'email';
  duration: number;
}

interface LoginFailureData {
  method: 'email';
  duration: number;
  error: string;
}

interface LoginApiRequestData {
  duration: number;
  status: 'success' | 'failure';
  statusCode?: number;
}

// ============================================================================
// AUTH TRACKING
// ============================================================================

/**
 * Track login attempt start
 */
export function trackLoginAttempt(data: LoginAttemptData) {
  Sentry.addBreadcrumb({
    category: 'auth',
    message: `Login attempt started via ${data.method}`,
    level: 'info',
    data: {
      action: 'login_attempt',
      method: data.method,
      timestamp: Date.now()
    }
  });

  console.log(`[Tracking] Login attempt started: ${data.method}`);
}

/**
 * Track login success
 */
export function trackLoginSuccess(data: LoginSuccessData) {
  Sentry.addBreadcrumb({
    category: 'auth',
    message: `Login successful via ${data.method}`,
    level: 'info',
    data: {
      action: 'login_success',
      method: data.method,
      duration_ms: data.duration,
      timestamp: Date.now()
    }
  });

  Sentry.setMeasurement('login.duration', data.duration, 'millisecond');

  console.log(`[Tracking] Login success: ${data.duration}ms`);
}

/**
 * Track login failure
 */
export function trackLoginFailure(data: LoginFailureData) {
  Sentry.addBreadcrumb({
    category: 'auth',
    message: `Login failed via ${data.method}: ${data.error}`,
    level: 'warning',
    data: {
      action: 'login_failure',
      method: data.method,
      duration_ms: data.duration,
      error_message: data.error,
      timestamp: Date.now()
    }
  });

  console.log(`[Tracking] Login failure: ${data.error} (${data.duration}ms)`);
}

/**
 * Track Login API Request performance specifically
 */
export function trackLoginApiRequest(data: LoginApiRequestData) {
    Sentry.addBreadcrumb({
        category: 'http',
        message: `Login API request completed`,
        level: data.status === 'success' ? 'info' : 'error',
        data: {
            url: '/accounts/login',
            method: 'POST',
            duration_ms: data.duration,
            status: data.status,
            status_code: data.statusCode
        }
    });

    Sentry.setMeasurement('login.api_duration', data.duration, 'millisecond');

    console.log(`[Tracking] Login API Request: ${data.status} (${data.duration}ms)`);
}
