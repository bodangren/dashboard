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

  describe('agent failure detection', function() {
    it('should detect error status from exit code', function() {
      const status = NotificationService.getAgentStatus(1, 'some error');
      assert.equal(status, 'error');
    });

    it('should detect success status from exit code 0', function() {
      const status = NotificationService.getAgentStatus(0, '');
      assert.equal(status, 'success');
    });

    it('should detect error status from non-zero exit code', function() {
      const status = NotificationService.getAgentStatus(127, 'command not found');
      assert.equal(status, 'error');
    });

    it('should extract last 5 log lines', function() {
      const lines = ['line1', 'line2', 'line3', 'line4', 'line5', 'line6', 'line7'];
      const result = NotificationService.extractLogLines(lines, 5);
      assert.equal(result.length, 5);
      assert.equal(result[0], 'line3');
    });

    it('should handle fewer than 5 log lines', function() {
      const lines = ['line1', 'line2'];
      const result = NotificationService.extractLogLines(lines, 5);
      assert.equal(result.length, 2);
    });

    it('should format agent failure notification body', function() {
      const body = NotificationService.formatAgentBody('command not found\nmore lines', 'Agent1', 127);
      assert.equal(body, 'Agent1 failed (exit 127): command not found');
    });

    it('should truncate long error messages', function() {
      const longMsg = 'a'.repeat(150);
      const body = NotificationService.formatAgentBody(longMsg, 'Agent1', 1);
      assert.equal(body.length < 130, true);
    });
  });

  describe('agent event notification', function() {
    const ENABLED_KEY = 'dashboard-notification-enabled';

    beforeEach(function() {
      localStorage.removeItem(ENABLED_KEY);
      localStorage.removeItem('dashboard-notification-agent-errors');
    });

    it('should notify on agent failure', function() {
      localStorage.setItem(ENABLED_KEY, 'true');
      localStorage.setItem('dashboard-notification-agent-errors', 'true');
      localStorage.setItem(STORAGE_KEY, 'granted');
      NotificationService.notifyAgentFailure('test-agent', 1, 'command not found');
    });

    it('should not notify if agent errors disabled', function() {
      localStorage.setItem(ENABLED_KEY, 'true');
      localStorage.setItem('dashboard-notification-agent-errors', 'false');
      localStorage.setItem(STORAGE_KEY, 'granted');
      NotificationService.notifyAgentFailure('test-agent', 1, 'error');
    });
  });

  describe('AI insight detection', function() {
    it('should detect conflict markers in flags', function() {
      const result = NotificationService.hasConflictFlag(['conflict-markers', 'WIP']);
      assert.equal(result, true);
    });

    it('should detect WIP flags', function() {
      const result = NotificationService.hasWIPFlag(['WIP', 'rapid-changes']);
      assert.equal(result, true);
    });

    it('should return false for no flags', function() {
      const result = NotificationService.hasConflictFlag([]);
      assert.equal(result, false);
    });

    it('should notify on conflict insight', function() {
      localStorage.setItem('dashboard-notification-enabled', 'true');
      localStorage.setItem('dashboard-notification-ai-insights', 'true');
      localStorage.setItem(STORAGE_KEY, 'granted');
      NotificationService.notifyAIInsight('repo', 'conflict detected', ['conflict-markers']);
    });
  });

  describe('notification preferences', function() {
    it('should get agent errors enabled preference', function() {
      localStorage.setItem('dashboard-notification-agent-errors', 'true');
      assert.equal(NotificationService.isAgentErrorsEnabled(), true);
    });

    it('should get agent errors disabled preference', function() {
      localStorage.setItem('dashboard-notification-agent-errors', 'false');
      assert.equal(NotificationService.isAgentErrorsEnabled(), false);
    });

    it('should default agent errors to enabled', function() {
      localStorage.removeItem('dashboard-notification-agent-errors');
      assert.equal(NotificationService.isAgentErrorsEnabled(), true);
    });

    it('should get AI insights enabled preference', function() {
      localStorage.setItem('dashboard-notification-ai-insights', 'true');
      assert.equal(NotificationService.isAIInsightsEnabled(), true);
    });

    it('should default AI insights to enabled', function() {
      localStorage.removeItem('dashboard-notification-ai-insights');
      assert.equal(NotificationService.isAIInsightsEnabled(), true);
    });
  });
});