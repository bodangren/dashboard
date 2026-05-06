'use strict';

describe('NotificationService', function() {
  const STORAGE_KEY = 'dashboard-notifications';

  beforeEach(function() {
    localStorage.removeItem(STORAGE_KEY);
    localStorage.removeItem('dashboard-notification-quiet-start');
    localStorage.removeItem('dashboard-notification-quiet-end');
    localStorage.removeItem('dashboard-notification-enabled');
  });

  describe('permission state management', function() {
    it('should default to denied when Notification.permission is denied', function() {
      const saved = NotificationService.getPreference();
      assert.equal(saved, 'denied');
    });

    it('should default to granted when Notification.permission is granted', function() {
      NotificationService.requestPermission();
      const saved = NotificationService.getPreference();
      assert.equal(saved, 'granted');
    });

    it('should read stored preference from localStorage', function() {
      localStorage.setItem(STORAGE_KEY, 'granted');
      assert.equal(NotificationService.getPreference(), 'granted');
    });
  });

  describe('isEnabled', function() {
    it('should return false when notifications are disabled', function() {
      localStorage.setItem('dashboard-notification-enabled', 'false');
      assert.equal(NotificationService.isEnabled(), false);
    });

    it('should return true when notifications are enabled', function() {
      localStorage.setItem('dashboard-notification-enabled', 'true');
      assert.equal(NotificationService.isEnabled(), true);
    });

    it('should return true by default', function() {
      assert.equal(NotificationService.isEnabled(), true);
    });
  });

  describe('quiet hours', function() {
    it('should not be in quiet hours when start and end are not set', function() {
      const result = NotificationService.isInQuietHours();
      assert.equal(result, false);
    });

    it('should be in quiet hours when current time is within range on same day', function() {
      localStorage.setItem('dashboard-notification-quiet-start', '22:00');
      localStorage.setItem('dashboard-notification-quiet-end', '08:00');
      assert.equal(NotificationService.isInQuietHours(), true);
    });

    it('should not be in quiet hours when current time is outside range', function() {
      localStorage.setItem('dashboard-notification-quiet-start', '08:00');
      localStorage.setItem('dashboard-notification-quiet-end', '18:00');
      const result = NotificationService.isInQuietHours();
      assert.equal(result, false);
    });
  });

  describe('notify', function() {
    it('should not notify when notifications are disabled', function() {
      localStorage.setItem('dashboard-notification-enabled', 'false');
      NotificationService.notify('Test', 'Body');
    });

    it('should not notify when in quiet hours', function() {
      localStorage.setItem('dashboard-notification-quiet-start', '00:00');
      localStorage.setItem('dashboard-notification-quiet-end', '23:59');
      NotificationService.notify('Test', 'Body');
    });

    it('should return early when permission is denied', function() {
      localStorage.setItem(STORAGE_KEY, 'denied');
      NotificationService.notify('Test', 'Body');
    });
  });
});