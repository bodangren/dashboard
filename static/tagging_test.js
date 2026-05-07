'use strict';

describe('Tag UI', function() {
  beforeEach(function() {
    TagManager.clearAllTags();
    TagManager.setTagsForRepo('/repo1', ['work', 'client']);
    TagManager.setTagsForRepo('/repo2', ['personal']);
    TagManager.setTagsForRepo('/repo3', []);
  });

  describe('renderTagChips', function() {
    it('should render tags as chips for a project', function() {
      const tags = TagManager.getTagsForRepo('/repo1');
      assert.deepEqual(tags, ['work', 'client']);
    });

    it('should return empty for project with no tags', function() {
      const tags = TagManager.getTagsForRepo('/repo3');
      assert.deepEqual(tags, []);
    });
  });

  describe('add tag via UI', function() {
    it('should add a new tag to repo', function() {
      const result = TagManager.addTag('/repo3', 'newtag');
      assert.deepEqual(result, ['newtag']);
    });

    it('should not add duplicate tag', function() {
      TagManager.addTag('/repo1', 'work');
      const tags = TagManager.getTagsForRepo('/repo1');
      assert.deepEqual(tags, ['work', 'client']);
    });
  });

  describe('remove tag via UI', function() {
    it('should remove a tag from repo', function() {
      TagManager.removeTag('/repo1', 'client');
      const tags = TagManager.getTagsForRepo('/repo1');
      assert.deepEqual(tags, ['work']);
    });
  });

  describe('tag filter', function() {
    it('should get repos with specific tag', function() {
      const repos = TagManager.getReposWithTag('work');
      assert.deepEqual(repos, ['/repo1']);
    });

    it('should return multiple repos for same tag', function() {
      TagManager.setTagsForRepo('/repo4', ['work']);
      const repos = TagManager.getReposWithTag('work');
      assert.deepEqual(repos, ['/repo1', '/repo4']);
    });

    it('should return empty for non-existent tag', function() {
      const repos = TagManager.getReposWithTag('nonexistent');
      assert.deepEqual(repos, []);
    });
  });

  describe('getAllTags', function() {
    it('should return sorted unique tags across all repos', function() {
      const allTags = TagManager.getAllTags();
      assert.deepEqual(allTags, ['client', 'personal', 'work']);
    });
  });
});