'use strict';

const NotificationService = (function() {
  const STORAGE_KEY = 'dashboard-notifications';
  const ENABLED_KEY = 'dashboard-notification-enabled';
  const QUIET_START_KEY = 'dashboard-notification-quiet-start';
  const QUIET_END_KEY = 'dashboard-notification-quiet-end';

  function getStoredPermission() {
    return localStorage.getItem(STORAGE_KEY) || 'denied';
  }

  function getBrowserPermission() {
    if (typeof Notification === 'undefined') return 'denied';
    return Notification.permission;
  }

  function getEffectivePermission() {
    const stored = getStoredPermission();
    if (stored === 'default') return getBrowserPermission();
    return stored;
  }

  function requestPermission() {
    if (typeof Notification === 'undefined') return;
    if (!('Notification' in window)) return;
    Notification.requestPermission().then(function(permission) {
      localStorage.setItem(STORAGE_KEY, permission);
    });
  }

  function getPreference() {
    const stored = getStoredPermission();
    if (stored === 'denied' || stored === 'granted') return stored;
    return getBrowserPermission();
  }

  function isEnabled() {
    const stored = localStorage.getItem(ENABLED_KEY);
    if (stored === null) return true;
    return stored === 'true';
  }

  function isInQuietHours() {
    const start = localStorage.getItem(QUIET_START_KEY);
    const end = localStorage.getItem(QUIET_END_KEY);
    if (!start || !end) return false;

    const now = new Date();
    const currentMinutes = now.getHours() * 60 + now.getMinutes();

    const [startH, startM] = start.split(':').map(Number);
    const [endH, endM] = end.split(':').map(Number);
    const startMinutes = startH * 60 + startM;
    const endMinutes = endH * 60 + endM;

    if (startMinutes < endMinutes) {
      return currentMinutes >= startMinutes && currentMinutes < endMinutes;
    } else {
      return currentMinutes >= startMinutes || currentMinutes < endMinutes;
    }
  }

  function notify(title, body, options) {
    if (!isEnabled()) return;
    if (isInQuietHours()) return;

    const permission = getEffectivePermission();
    if (permission !== 'granted') return;

    if (typeof Notification === 'undefined') return;
    if (!('Notification' in window)) return;

    const notification = new Notification(title, Object.assign({ body: body }, options));
    setTimeout(function() { notification.close(); }, 5000);
  }

  function setEnabled(enabled) {
    localStorage.setItem(ENABLED_KEY, enabled ? 'true' : 'false');
  }

  function setQuietHours(start, end) {
    if (start) localStorage.setItem(QUIET_START_KEY, start);
    else localStorage.removeItem(QUIET_START_KEY);
    if (end) localStorage.setItem(QUIET_END_KEY, end);
    else localStorage.removeItem(QUIET_END_KEY);
  }

  function getQuietHours() {
    return {
      start: localStorage.getItem(QUIET_START_KEY) || '',
      end: localStorage.getItem(QUIET_END_KEY) || ''
    };
  }

  return {
    requestPermission: requestPermission,
    getPreference: getPreference,
    isEnabled: isEnabled,
    setEnabled: setEnabled,
    isInQuietHours: isInQuietHours,
    setQuietHours: setQuietHours,
    getQuietHours: getQuietHours,
    notify: notify
  };
})();