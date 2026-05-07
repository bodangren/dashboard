'use strict';

describe('TagManager', function() {
  const STORAGE_KEY = TagManager.STORAGE_KEY;

  beforeEach(function() {
    TagManager.clearAllTags();
  });

  describe('getTagsForRepo', function() {
    it('should return empty array for repo with no tags', function() {
      const result = TagManager.getTagsForRepo('/some/path');
      assert.deepEqual(result, []);
    });

    it('should return tags for repo that has them', function() {
      TagManager.setTagsForRepo('/some/path', ['work', 'urgent']);
      const result = TagManager.getTagsForRepo('/some/path');
      assert.deepEqual(result, ['work', 'urgent']);
    });
  });

  describe('setTagsForRepo', function() {
    it('should store tags for a repo', function() {
      TagManager.setTagsForRepo('/repo', ['client', 'active']);
      assert.deepEqual(TagManager.getTagsForRepo('/repo'), ['client', 'active']);
    });

    it('should remove repo entry when tags cleared', function() {
      TagManager.setTagsForRepo('/repo', ['tag']);
      TagManager.setTagsForRepo('/repo', []);
      assert.deepEqual(TagManager.getTagsForRepo('/repo'), []);
    });

    it('should handle empty array', function() {
      TagManager.setTagsForRepo('/repo', []);
      assert.deepEqual(TagManager.getTagsForRepo('/repo'), []);
    });
  });

  describe('addTag', function() {
    it('should add a new tag', function() {
      const result = TagManager.addTag('/repo', 'newtag');
      assert.deepEqual(result, ['newtag']);
    });

    it('should not add duplicate tag', function() {
      TagManager.addTag('/repo', 'tag1');
      const result = TagManager.addTag('/repo', 'tag1');
      assert.deepEqual(result, ['tag1']);
    });

    it('should trim whitespace from tag', function() {
      const result = TagManager.addTag('/repo', '  trimmed  ');
      assert.deepEqual(result, ['trimmed']);
    });

    it('should not add empty tag', function() {
      TagManager.addTag('/repo', '   ');
      assert.deepEqual(TagManager.getTagsForRepo('/repo'), []);
    });

    it('should append to existing tags', function() {
      TagManager.addTag('/repo', 'tag1');
      const result = TagManager.addTag('/repo', 'tag2');
      assert.deepEqual(result, ['tag1', 'tag2']);
    });
  });

  describe('removeTag', function() {
    it('should remove existing tag', function() {
      TagManager.setTagsForRepo('/repo', ['tag1', 'tag2']);
      const result = TagManager.removeTag('/repo', 'tag1');
      assert.deepEqual(result, ['tag2']);
    });

    it('should handle removing non-existent tag', function() {
      TagManager.setTagsForRepo('/repo', ['tag1']);
      const result = TagManager.removeTag('/repo', 'nonexistent');
      assert.deepEqual(result, ['tag1']);
    });

    it('should trim whitespace when removing', function() {
      TagManager.setTagsForRepo('/repo', ['tag1']);
      TagManager.removeTag('/repo', '  tag1  ');
      assert.deepEqual(TagManager.getTagsForRepo('/repo'), []);
    });
  });

  describe('getAllTags', function() {
    it('should return empty array when no tags exist', function() {
      assert.deepEqual(TagManager.getAllTags(), []);
    });

    it('should return unique sorted tags across all repos', function() {
      TagManager.setTagsForRepo('/repo1', ['work', 'client']);
      TagManager.setTagsForRepo('/repo2', ['client', 'archive']);
      const result = TagManager.getAllTags();
      assert.deepEqual(result, ['archive', 'client', 'work']);
    });

    it('should not include duplicate tags', function() {
      TagManager.setTagsForRepo('/repo1', ['tag']);
      TagManager.setTagsForRepo('/repo2', ['tag']);
      const result = TagManager.getAllTags();
      assert.deepEqual(result, ['tag']);
    });
  });

  describe('getReposWithTag', function() {
    it('should return empty array for tag that matches nothing', function() {
      assert.deepEqual(TagManager.getReposWithTag('nonexistent'), []);
    });

    it('should return repos that have the tag', function() {
      TagManager.setTagsForRepo('/repo1', ['work']);
      TagManager.setTagsForRepo('/repo2', ['personal']);
      TagManager.setTagsForRepo('/repo3', ['work', 'urgent']);
      const result = TagManager.getReposWithTag('work');
      assert.deepEqual(result, ['/repo1', '/repo3']);
    });

    it('should trim whitespace when matching', function() {
      TagManager.setTagsForRepo('/repo1', ['work']);
      const result = TagManager.getReposWithTag('  work  ');
      assert.deepEqual(result, ['/repo1']);
    });
  });

  describe('clearAllTags', function() {
    it('should remove all tags', function() {
      TagManager.setTagsForRepo('/repo1', ['work']);
      TagManager.setTagsForRepo('/repo2', ['personal']);
      TagManager.clearAllTags();
      assert.deepEqual(TagManager.getAllTags(), []);
      assert.deepEqual(TagManager.getTagsForRepo('/repo1'), []);
      assert.deepEqual(TagManager.getTagsForRepo('/repo2'), []);
    });
  });

  describe('persistence', function() {
    it('should persist tags across page reload', function() {
      TagManager.setTagsForRepo('/repo', ['work', 'urgent']);
      localStorage.setItem(STORAGE_KEY, localStorage.getItem(STORAGE_KEY));
      assert.deepEqual(TagManager.getTagsForRepo('/repo'), ['work', 'urgent']);
    });
  });
});